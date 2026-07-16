package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/resourcepolicy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &resourcePolicyResource{}
	_ resource.ResourceWithConfigure   = &resourcePolicyResource{}
	_ resource.ResourceWithImportState = &resourcePolicyResource{}
)

type resourcePolicyResource struct {
	client *infra.Client
}

type resourcePolicyImportIdentity struct {
	ProjectID     string
	ApplicationID string
	ResourceID    string
}

func NewResourcePolicyResource() resource.Resource {
	return &resourcePolicyResource{}
}

func (r *resourcePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *resourcePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_policy"
}

func (r *resourcePolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourcepolicy.Schema
}

func (r *resourcePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating resource policy")
	model := &resourcepolicy.Model{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	policy, err := r.client.CreateResourcePolicy(ctx, model.ProjectID.ValueString(), model.Policy(handler))
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource policy", err.Error())
		return
	}
	model.SetPolicy(policy)
	model.ID = stringattr.Value(resourcePolicyID(policy.ApplicationID, policy.ResourceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *resourcePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading resource policy")
	model := &resourcepolicy.Model{}
	resp.Diagnostics.Append(req.State.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.ReadResourcePolicy(ctx, model.ProjectID.ValueString(), model.Identity())
	if errors.Is(err, infra.ErrResourcePolicyNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource policy", err.Error())
		return
	}

	model.SetPolicy(policy)
	model.ID = stringattr.Value(resourcePolicyID(policy.ApplicationID, policy.ResourceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *resourcePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating resource policy")
	model := &resourcepolicy.Model{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	policy, err := r.client.UpdateResourcePolicy(ctx, model.ProjectID.ValueString(), model.Policy(handler))
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource policy", err.Error())
		return
	}
	model.SetPolicy(policy)
	model.ID = stringattr.Value(resourcePolicyID(policy.ApplicationID, policy.ResourceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *resourcePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting resource policy")
	model := &resourcepolicy.Model{}
	resp.Diagnostics.Append(req.State.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteResourcePolicy(ctx, model.ProjectID.ValueString(), model.Identity()); err != nil {
		resp.Diagnostics.AddError("Error deleting resource policy", err.Error())
	}
}

func (r *resourcePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	identity, err := parseResourcePolicyImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", err.Error())
		return
	}

	helpers.MarkImportState(ctx, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourcePolicyID(identity.ApplicationID, identity.ResourceID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), identity.ProjectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), identity.ApplicationID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), identity.ResourceID)...)
}

func parseResourcePolicyImportID(importID string) (resourcePolicyImportIdentity, error) {
	parts := strings.Split(importID, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return resourcePolicyImportIdentity{}, fmt.Errorf("import ID must be in the format 'project_id/application_id/resource_id'")
	}
	return resourcePolicyImportIdentity{ProjectID: parts[0], ApplicationID: parts[1], ResourceID: parts[2]}, nil
}

func resourcePolicyID(applicationID, resourceID string) string {
	return applicationID + "/" + resourceID
}
