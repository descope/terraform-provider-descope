package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/policyrule"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &policyRuleResource{}
	_ resource.ResourceWithConfigure   = &policyRuleResource{}
	_ resource.ResourceWithImportState = &policyRuleResource{}
)

type policyRuleResource struct {
	client *infra.Client
}

func (r *policyRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *policyRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_rule"
}

func (r *policyRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = policyrule.Schema
}

func (r *policyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model policyrule.PolicyRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, diags := model.ToPolicyRule(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreatePolicyRule(ctx, model.ProjectID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Policy Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(model.SetPolicyRule(ctx, created)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *policyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	var model policyrule.PolicyRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.LoadPolicyRule(ctx, model.ProjectID.ValueString(), model.ID.ValueString())
	if descope.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Policy Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(model.SetPolicyRule(ctx, rule)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *policyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model policyrule.PolicyRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, diags := model.ToPolicyRule(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdatePolicyRule(ctx, model.ProjectID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Policy Rule", err.Error())
		return
	}

	resp.Diagnostics.Append(model.SetPolicyRule(ctx, updated)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *policyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model policyrule.PolicyRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeletePolicyRule(ctx, model.ProjectID.ValueString(), model.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error Deleting Policy Rule", err.Error())
	}
}

func (r *policyRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.MarkImportState(ctx, resp)
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Import ID must be in the format 'project_id/policy_rule_id', got %q", req.ID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
