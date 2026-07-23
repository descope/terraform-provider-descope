package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/tenant"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tenantResourceName = "tenant"

var (
	_ resource.Resource                = &tenantResource{}
	_ resource.ResourceWithConfigure   = &tenantResource{}
	_ resource.ResourceWithImportState = &tenantResource{}
)

func NewTenantResource() resource.Resource {
	return &tenantResource{}
}

type tenantResource struct {
	client *infra.Client
}

func (r *tenantResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *tenantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tenantResourceName
}

func (r *tenantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: tenant.TenantAttributes}
}

func (r *tenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating tenant resource")

	var plan tenant.Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	values := plan.Values(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateTenant(ctx, values.ProjectID, infra.TenantCreateRequest{
		ID:                      values.ID,
		Name:                    values.Name,
		SelfProvisioningDomains: values.SelfProvisioningDomains,
		Disabled:                values.Disabled,
		EnforceSSO:              values.EnforceSSO,
		EnforceSSOExclusions:    values.EnforceSSOExclusions,
		FederatedApplicationIDs: values.FederatedApplicationIDs,
		Parent:                  values.Parent,
		RoleInheritance:         values.RoleInheritance,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating tenant", err.Error())
		return
	}

	remote, err := r.client.ReadTenant(ctx, values.ProjectID, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created tenant", err.Error())
		return
	}
	plan.SetTenant(ctx, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading tenant resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	var state tenant.Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	values := state.Values(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	remote, err := r.client.ReadTenant(ctx, values.ProjectID, values.ID)
	if infra.IsTenantNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant", err.Error())
		return
	}

	state.SetTenant(ctx, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *tenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating tenant resource")

	var plan tenant.Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	values := plan.Values(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	current, err := r.client.ReadTenant(ctx, values.ProjectID, values.ID)
	if infra.IsTenantNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant before update", err.Error())
		return
	}
	err = r.client.UpdateTenant(ctx, values.ProjectID, infra.TenantUpdateRequest{
		ID:                      values.ID,
		Name:                    values.Name,
		SelfProvisioningDomains: values.SelfProvisioningDomains,
		CustomAttributes:        current.CustomAttributes,
		AuthType:                current.AuthType,
		Disabled:                values.Disabled,
		EnforceSSO:              values.EnforceSSO,
		EnforceSSOExclusions:    values.EnforceSSOExclusions,
		FederatedApplicationIDs: values.FederatedApplicationIDs,
		RoleInheritance:         values.RoleInheritance,
		IDJagSettings:           current.IDJagSettings,
		IDJagEnabled:            current.IDJagEnabled,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant", err.Error())
		return
	}

	remote, err := r.client.ReadTenant(ctx, values.ProjectID, values.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated tenant", err.Error())
		return
	}
	plan.SetTenant(ctx, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *tenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting tenant resource")

	var state tenant.Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	values := state.Values(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteTenant(ctx, values.ProjectID, values.ID); err != nil && !infra.IsTenantNotFound(err) {
		resp.Diagnostics.AddError("Error deleting tenant", err.Error())
	}
}

func (r *tenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	projectID, tenantID, err := parseTenantImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid tenant import ID", err.Error())
		return
	}

	helpers.MarkImportState(ctx, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), tenantID)...)
}

func parseTenantImportID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("tenant import ID must use the format <project_id>/<tenant_id>")
	}
	return parts[0], parts[1], nil
}
