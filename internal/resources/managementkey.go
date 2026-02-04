package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/entities"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	managementKeyEntity = "management_key"
)

var (
	_ resource.Resource                = &managementKeyResource{}
	_ resource.ResourceWithConfigure   = &managementKeyResource{}
	_ resource.ResourceWithImportState = &managementKeyResource{}
)

func NewManagementKeyResource() resource.Resource {
	return &managementKeyResource{}
}

type managementKeyResource struct {
	client *infra.Client
}

func (r *managementKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *managementKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + managementKeyEntity
}

func (r *managementKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ManagementKeySchema
}

func (r *managementKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating management key resource")

	entity := entities.NewManagementKeyEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, infra.NoProjectID, managementKeyEntity, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid management key configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating management key", err.Error())
		return
	}

	entity.SetID(ctx, res.ID)
	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Management key resource created")
}

func (r *managementKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading management key resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	entity := entities.NewManagementKeyEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, infra.NoProjectID, managementKeyEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading management key", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Management key resource read")
}

func (r *managementKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating management key resource")

	entity := entities.NewManagementKeyEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, infra.NoProjectID, managementKeyEntity, id, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid management key configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating management key", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Management key resource updated")
}

func (r *managementKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting management key resource")

	entity := entities.NewManagementKeyEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, infra.NoProjectID, managementKeyEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting management key", err.Error())
		return
	}

	tflog.Info(ctx, "Management key resource deleted")
}

func (r *managementKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing management key resource")
	helpers.MarkImportState(ctx, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
