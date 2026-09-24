package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var legacyProjectOwnershipAttributes = [...]string{
	"project_settings",
	"invite_settings",
	"authentication",
	"authorization",
	"attributes",
	"connectors",
	"applications",
	"jwt_templates",
	"styles",
	"flows",
	"widgets",
	"lists",
	"admin_portal",
}

var _ resource.ResourceWithUpgradeState = &projectResource{}

type projectResource struct {
	*baseResource[project.ProjectModel, *project.ProjectModel]
}

// The project container resource on its dedicated /v1/mgmt/project route. The resource id is the project
// id: reads, updates and deletes scope the bearer token with it, while create runs with a bare key and
// gets the new id from the response.
func NewProjectResource() resource.Resource {
	const projectPath = "/v1/mgmt/project"
	ops := operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, projectPath, data)
			if err != nil {
				return "", nil, err
			}
			id, _ := body["id"].(string)
			return id, body, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, _ string) (map[string]any, error) {
			return c.Get(ctx, projectID, projectPath, nil)
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, _ string, data map[string]any) (map[string]any, error) {
			return c.PutData(ctx, projectID, projectPath, data)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, _ string) error {
			return c.Del(ctx, projectID, projectPath, nil)
		},
	}
	return &projectResource{baseResource: &baseResource[project.ProjectModel, *project.ProjectModel]{
		name: "project", schema: project.Schema, ops: ops,
	}}
}

func (r *projectResource) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			// UpgradeState cannot fan one legacy project into standalone resources. Only clean state may pass.
			PriorSchema:   nil,
			StateUpgrader: upgradeProjectState,
		},
	}
}

func upgradeProjectState(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var rawState map[string]json.RawMessage
	if err := json.Unmarshal(req.RawState.JSON, &rawState); err != nil {
		resp.Diagnostics.AddError("Invalid Descope Project State", "The existing descope_project state could not be decoded: "+err.Error())
		return
	}

	populated := make([]string, 0, len(legacyProjectOwnershipAttributes))
	known := map[string]bool{
		"id": true, "name": true, "environment": true, "tags": true, "deletion_protection": true,
	}
	for _, attribute := range legacyProjectOwnershipAttributes {
		known[attribute] = true
		if raw, ok := rawState[attribute]; ok && legacyOwnershipPopulated(raw) {
			populated = append(populated, attribute)
		}
	}
	for attribute, raw := range rawState {
		if !known[attribute] && legacyOwnershipPopulated(raw) {
			populated = append(populated, attribute)
		}
	}
	if len(populated) > 0 {
		slices.Sort(populated)
		resp.Diagnostics.AddError(
			"Descope Project State Requires Migration",
			fmt.Sprintf("Legacy descope_project state contains populated nested ownership attributes: %s. These nested blocks are now standalone resources, and Terraform cannot split one resource's state into many automatically. With the old provider version still pinned, run `go run github.com/descope/terraform-provider-descope/tools/tfmigrate plan -state <state.json> -address <legacy address> -source-provider-version <X.Y.Z> -target-provider-version <A.B.C> -out <dir>`, review the generated artifacts, and follow the generated README to detach with a `removed` block and adopt with `import` blocks. No state was modified. Downgrading the provider pin restores the previous behavior.", strings.Join(populated, ", ")),
		)
		return
	}

	stringValues := make(map[string]types.String, 3)
	for _, attribute := range [...]string{"id", "name", "environment"} {
		value, err := decodeProjectString(rawState, attribute)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Descope Project State", err.Error())
			return
		}
		stringValues[attribute] = value
	}

	var tagsValue any = types.SetNull(types.StringType)
	if raw, ok := rawState["tags"]; ok && !jsonNull(raw) {
		var tags []string
		if err := json.Unmarshal(raw, &tags); err != nil {
			resp.Diagnostics.AddError("Invalid Descope Project State", "The existing descope_project tags could not be decoded: "+err.Error())
			return
		}
		tagsValue = tags
	}
	deletionProtection := types.BoolNull()
	if raw, ok := rawState["deletion_protection"]; ok && !jsonNull(raw) {
		var enabled bool
		if err := json.Unmarshal(raw, &enabled); err != nil {
			resp.Diagnostics.AddError("Invalid Descope Project State", "The existing descope_project deletion_protection could not be decoded: "+err.Error())
			return
		}
		deletionProtection = types.BoolValue(enabled)
	}

	resp.State.Raw = tftypes.NewValue(resp.State.Schema.Type().TerraformType(ctx), nil)
	for _, attribute := range [...]string{"id", "name", "environment"} {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attribute), stringValues[attribute])...)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tags"), tagsValue)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("deletion_protection"), deletionProtection)...)
}

func legacyOwnershipPopulated(raw json.RawMessage) bool {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return true
	}
	switch value := value.(type) {
	case nil:
		return false
	case map[string]any:
		return len(value) > 0
	case []any:
		return len(value) > 0
	default:
		return true
	}
}

func decodeProjectString(rawState map[string]json.RawMessage, attribute string) (types.String, error) {
	raw, ok := rawState[attribute]
	if !ok || jsonNull(raw) {
		return types.StringNull(), nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return types.StringNull(), fmt.Errorf("the existing descope_project %s could not be decoded: %w", attribute, err)
	}
	return types.StringValue(value), nil
}

func jsonNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}
