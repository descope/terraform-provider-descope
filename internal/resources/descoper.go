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
	descoperEntity = "descoper"
)

var (
	_ resource.Resource                = &descoperResource{}
	_ resource.ResourceWithConfigure   = &descoperResource{}
	_ resource.ResourceWithImportState = &descoperResource{}
)

func NewDescoperResource() resource.Resource {
	return &descoperResource{}
}

type descoperResource struct {
	client *infra.Client
}

func (r *descoperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *descoperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + descoperEntity
}

func (r *descoperResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.DescoperSchema
}

func (r *descoperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating descoper resource")

	entity := entities.NewDescoperEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, infra.NoProjectID, descoperEntity, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid descoper configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating descoper", err.Error())
		return
	}

	entity.SetID(ctx, res.ID)
	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Descoper resource created")
}

func (r *descoperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading descoper resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	entity := entities.NewDescoperEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, infra.NoProjectID, descoperEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading descoper", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Descoper resource read")
}

func (r *descoperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating descoper resource")

	entity := entities.NewDescoperEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, infra.NoProjectID, descoperEntity, id, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid descoper configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating descoper", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Descoper resource updated")
}

func (r *descoperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting descoper resource")

	entity := entities.NewDescoperEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, infra.NoProjectID, descoperEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting descoper", err.Error())
		return
	}

	tflog.Info(ctx, "Descoper resource deleted")
}

func (r *descoperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing descoper resource")
	helpers.MarkImportState(ctx, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
