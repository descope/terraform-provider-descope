package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/entities"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const projectEntity = "project"

var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *infra.Client
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + projectEntity
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = entities.ProjectSchema
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating project resource")

	entity := entities.NewProjectEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, infra.PrincipalProjectID, projectEntity, values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	entity.SetProjectID(ctx, res.ID)
	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Project resource created")
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading project resource")

	entity := entities.NewProjectEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	projectID := entity.ProjectID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, projectID, projectEntity, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Project resource read")
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating project resource")

	entity := entities.NewProjectEntity(ctx, req.Plan, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	values := entity.Values(ctx)
	projectID := entity.ProjectID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, projectID, projectEntity, projectID, values)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	entity.SetValues(ctx, res.Data)
	entity.Save(ctx, &resp.State)

	tflog.Info(ctx, "Project resource updated")
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting project resource")

	entity := entities.NewProjectEntity(ctx, req.State, &resp.Diagnostics)
	if entity.Diagnostics.HasError() {
		return
	}

	projectID := entity.ProjectID(ctx)
	if entity.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, projectID, projectEntity, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}

	tflog.Info(ctx, "Project resource deleted")
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
