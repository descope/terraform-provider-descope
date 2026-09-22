package read

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type Result struct {
	Instance discover.Instance
	Object   types.Object
}

func Singletons(loaded map[string]resources.ExportableResource, projectID string) []discover.Instance {
	var instances []discover.Instance
	for name, exportable := range loaded {
		if exportable.ExportSingleton() {
			instances = append(instances, discover.Instance{Resource: name, ID: projectID, Name: name})
		}
	}
	return instances
}

const UnreadableWarning = "the entity was discovered but could not be read"

func All(ctx context.Context, client *infra.Client, loaded map[string]resources.ExportableResource, projectID string, instances []discover.Instance) (results []Result, warnings []warn.Warning) {
	for _, instance := range instances {
		exportable, ok := loaded[instance.Resource]
		if !ok {
			warnings = append(warnings, warn.Lost("Skipped %s %s: unknown resource type", instance.Resource, instance.ID))
			continue
		}
		model, found, diags := exportable.ExportRead(ctx, client, projectID, instance.Scope, instance.ID)
		if !found {
			// a singleton or uncustomized placeholder has nothing to export, but any other unreadable entity is a silent loss and must be reported
			if !exportable.ExportSingleton() && !instance.MayBeAbsent {
				warnings = append(warnings, warn.Lost("Skipped %s %s: %s", instance.Resource, instance.ID, UnreadableWarning))
			}
			continue
		}
		object, objectDiags := stateObject(ctx, exportable.ExportSchema(), model, projectID, instance.Scope)
		diags.Append(objectDiags...)
		if diags.HasError() {
			for _, diagnostic := range diags.Errors() {
				warnings = append(warnings, warn.Lost("Failed reading %s %s: %s: %s", instance.Resource, instance.ID, diagnostic.Summary(), diagnostic.Detail()))
			}
			continue
		}
		results = append(results, Result{Instance: instance, Object: object})
	}
	return results, warnings
}

func stateObject(ctx context.Context, sc schema.Schema, model any, projectID, scope string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	state := tfsdk.State{Schema: sc, Raw: tftypes.NewValue(sc.Type().TerraformType(ctx), nil)}
	diags.Append(state.Set(ctx, model)...)
	if _, ok := sc.Attributes["project_id"]; ok {
		diags.Append(state.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	}
	if _, ok := sc.Attributes["method"]; ok && scope != "" {
		diags.Append(state.SetAttribute(ctx, path.Root("method"), scope)...)
	}
	if _, ok := sc.Attributes["app_id"]; ok && scope != "" {
		diags.Append(state.SetAttribute(ctx, path.Root("app_id"), scope)...)
	}
	if diags.HasError() {
		return basetypes.NewObjectNull(nil), diags
	}

	value, err := sc.Type().ValueFromTerraform(ctx, state.Raw)
	if err != nil {
		diags.AddError("Error converting state", err.Error())
		return basetypes.NewObjectNull(nil), diags
	}
	object, ok := value.(types.Object)
	if !ok {
		diags.AddError("Error converting state", "Unexpected value type")
		return basetypes.NewObjectNull(nil), diags
	}
	return object, diags
}
