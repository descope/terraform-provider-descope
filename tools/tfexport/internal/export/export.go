package export

import (
	"context"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/emit"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/read"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ReadProject(ctx context.Context, client *infra.Client, projectID, only string) ([]read.Result, []warn.Warning, error) {
	loaded := registry.Load(ctx)

	files, err := discover.Snapshot(ctx, client, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to export project snapshot: %w", err)
	}
	instances, warnings := discover.Instances(files, registry.ConnectorTypes(ctx, loaded))
	instances = append(instances, read.Singletons(loaded, projectID)...)
	if _, ok := loaded["descope_project"]; ok {
		instances = append(instances, discover.Instance{Resource: "descope_project", ID: projectID, Name: "imported_project"})
	}
	inboundApps, err := discover.InboundApps(ctx, client, projectID)
	if err != nil {
		warnings = append(warnings, warn.Lost("Failed to list inbound apps: %s", err.Error()))
	}
	instances = append(instances, inboundApps...)
	instances, resolveWarnings := discover.ResolveAuthorizationIDs(ctx, client, projectID, instances)
	warnings = append(warnings, resolveWarnings...)

	if only != "" {
		var filtered []discover.Instance
		for _, instance := range instances {
			if strings.Contains(instance.Resource, only) {
				filtered = append(filtered, instance)
			}
		}
		instances = filtered
	}

	results, readWarnings := read.All(ctx, client, loaded, projectID, instances)
	return results, append(warnings, readWarnings...), nil
}

func projectNameOf(results []read.Result) string {
	for _, result := range results {
		if result.Instance.Resource != "descope_project" {
			continue
		}
		if value, ok := result.Object.Attributes()["name"].(basetypes.StringValue); ok {
			return value.ValueString()
		}
	}
	return ""
}

func Run(ctx context.Context, client *infra.Client, projectID, outDir, only string) (int, []warn.Warning, error) {
	results, warnings, err := ReadProject(ctx, client, projectID, only)
	if err != nil {
		return 0, warnings, err
	}

	loaded := registry.Load(ctx)
	plan := &emit.Plan{ProjectID: projectID}
	labels := emit.NewLabels()
	projectName := projectNameOf(results)
	for _, result := range results {
		exportable := loaded[result.Instance.Resource]
		schema := exportable.ExportSchema()
		attrs, secrets := prune.Object(ctx, schema.Attributes, result.Object, true)
		if !prune.Meaningful(attrs) && exportable.ExportSingleton() {
			continue // nothing but defaults, so the resource is left out entirely
		}
		attrs, secrets = ensureValidConfig(ctx, exportable, result.Instance.Name, result.Object, attrs, secrets, func(format string, args ...any) {
			warnings = append(warnings, warn.Note(format, args...))
		})

		_, hasProjectID := schema.Attributes["project_id"]
		_, hasMethod := schema.Attributes["method"]
		_, hasAppID := schema.Attributes["app_id"]
		label, fallback := result.Instance.Name, result.Instance.ID
		if exportable.ExportSingleton() || result.Instance.Resource == "descope_project" {
			label, fallback = projectName, "main"
		}

		resource := emit.Resource{
			Type:         result.Instance.Resource,
			EntityID:     result.Instance.ID,
			Label:        labels.Assign(result.Instance.Resource, label, fallback),
			Scope:        result.Instance.Scope,
			HasProjectID: hasProjectID,
			HasMethod:    hasMethod,
			HasAppID:     hasAppID,
			Attrs:        attrs,
		}
		// the ID formats match what baseResource.ImportState expects
		switch {
		case !hasProjectID:
			resource.ImportID = result.Instance.ID
		case exportable.ExportSingleton():
			resource.ImportID = projectID
		case (hasMethod || hasAppID) && result.Instance.Scope != "":
			resource.ImportID = projectID + "/" + result.Instance.Scope + "/" + result.Instance.ID
		default:
			resource.ImportID = projectID + "/" + result.Instance.ID
		}

		for _, secret := range secrets {
			switch {
			case secret.Required:
				warnings = append(warnings, warn.Note("Secret %s of %s.%s must be assigned via its generated variable", secret.Path, resource.Type, resource.Label))
			case secret.Dropped:
				warnings = append(warnings, warn.Note("Secret %s of %s.%s has a stored value that the generated configuration cannot carry: set it manually before applying, or the apply clears it", secret.Path, resource.Type, resource.Label))
			}
		}
		plan.Resources = append(plan.Resources, resource)
	}

	if err := emit.Write(ctx, outDir, plan); err != nil {
		return 0, warnings, fmt.Errorf("failed to write output files: %w", err)
	}
	return len(plan.Resources), warnings, nil
}
