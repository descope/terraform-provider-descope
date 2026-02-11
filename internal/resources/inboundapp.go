package resources

import (
	"context"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/entities"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	inboundAppEntity = "inbound_app"
)

var (
	_ resource.Resource                = &inboundAppResource{}
	_ resource.ResourceWithConfigure   = &inboundAppResource{}
	_ resource.ResourceWithImportState = &inboundAppResource{}
)

func NewInboundAppResource() resource.Resource {
	return &inboundAppResource{}
}

type inboundAppResource struct {
	client *infra.Client
}

func (r *inboundAppResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *inboundAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + inboundAppEntity
}

func (r *inboundAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.InboundAppSchema
}

func (r *inboundAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating inbound app resource")

	entity := entities.NewInboundAppEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, entity.ProjectID(ctx), inboundAppEntity, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid inbound app configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating inbound app", err.Error())
		return
	}

	entity.SetID(ctx, res.ID)
	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Inbound app resource created")
}

func (r *inboundAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading inbound app resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	entity := entities.NewInboundAppEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, entity.ProjectID(ctx), inboundAppEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading inbound app", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Inbound app resource read")
}

func (r *inboundAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating inbound app resource")

	entity := entities.NewInboundAppEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, entity.ProjectID(ctx), inboundAppEntity, id, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid inbound app configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating inbound app", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Inbound app resource updated")
}

func (r *inboundAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting inbound app resource")

	entity := entities.NewInboundAppEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	id := entity.ID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, entity.ProjectID(ctx), inboundAppEntity, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting inbound app", err.Error())
		return
	}

	tflog.Info(ctx, "Inbound app resource deleted")
}

func (r *inboundAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing inbound app resource")
	helpers.MarkImportState(ctx, resp)

	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be in the format 'project_id/inboundapp_id'")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
