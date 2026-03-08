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
	accessKeyEntity = "access_key"
)

var (
	_ resource.Resource                = &accessKeyResource{}
	_ resource.ResourceWithConfigure   = &accessKeyResource{}
	_ resource.ResourceWithImportState = &accessKeyResource{}
)

func NewAccessKeyResource() resource.Resource {
	return &accessKeyResource{}
}

type accessKeyResource struct {
	client *infra.Client
}

func (r *accessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *accessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + accessKeyEntity
}

func (r *accessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.AccessKeySchema
}

func (r *accessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating access key resource")

	entity := entities.NewAccessKeyEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, infra.NoProjectID, accessKeyEntity, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid access key configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating access key", err.Error())
		return
	}

	entity.SetID(ctx, res.ID)
	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Access key resource created")
}

func (r *accessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading access key resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	entity := entities.NewAccessKeyEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, infra.NoProjectID, accessKeyEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access key", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Access key resource read")
}

func (r *accessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating access key resource")

	entity := entities.NewAccessKeyEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, infra.NoProjectID, accessKeyEntity, id, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid access key configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating access key", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Access key resource updated")
}

func (r *accessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting access key resource")

	entity := entities.NewAccessKeyEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, infra.NoProjectID, accessKeyEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting access key", err.Error())
		return
	}

	tflog.Info(ctx, "Access key resource deleted")
}

func (r *accessKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing access key resource")
	helpers.MarkImportState(ctx, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
