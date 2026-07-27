package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Creates a new resource with the given name. If the schema contains a `project_id` attribute then
// the resource will be assumed to be a project-level resource (like a connector or flow), otherwise
// it'll be assumed to be a company-level resource.
func newResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema) resource.Resource {
	return &baseResource[T, M]{name: name, schema: sc}
}

// Use a random model to ensure interface conformance
var (
	_ resource.Resource                = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithConfigure   = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithImportState = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
)

type baseResource[T any, M helpers.ResourceModel[T]] struct {
	name   string
	schema schema.Schema
	client *infra.Client
}

func (r *baseResource[T, M]) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *baseResource[T, M]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.name
}

func (r *baseResource[T, M]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema
}

func (r *baseResource[T, M]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating "+r.name+" resource")

	model := M(new(T))
	resp.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	values := model.Values(handler)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Create(ctx, model.GetProjectID().ValueString(), r.name, values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid "+r.name+" configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating "+r.name, err.Error())
		return
	}

	model.SetID(types.StringValue(res.ID))
	model.SetValues(handler, res.Data)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	tflog.Info(ctx, "Created "+r.name+" resource")
}

func (r *baseResource[T, M]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading "+r.name+" resource")
	ctx = helpers.ContextWithImportState(ctx, req, resp)

	model := M(new(T))
	resp.Diagnostics.Append(req.State.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Read(ctx, model.GetProjectID().ValueString(), r.name, model.GetID().ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading "+r.name, err.Error())
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	model.SetValues(handler, res.Data)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	tflog.Info(ctx, "Read "+r.name+" resource")
}

func (r *baseResource[T, M]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating "+r.name+" resource")

	model := M(new(T))
	resp.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	values := model.Values(handler)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.Update(ctx, model.GetProjectID().ValueString(), r.name, model.GetID().ValueString(), values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid "+r.name+" configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating "+r.name, err.Error())
		return
	}

	model.SetValues(handler, res.Data)
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)

	tflog.Info(ctx, "Updated "+r.name+" resource")
}

func (r *baseResource[T, M]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting "+r.name+" resource")

	model := M(new(T))
	resp.Diagnostics.Append(req.State.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, model.GetProjectID().ValueString(), r.name, model.GetID().ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting "+r.name, err.Error())
		return
	}

	tflog.Info(ctx, "Deleted "+r.name+" resource")
}

func (r *baseResource[T, M]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing "+r.name+" resource")
	helpers.MarkImportState(ctx, resp)

	if _, ok := r.schema.Attributes["project_id"]; !ok {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
		return
	}

	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Import ID must be in the format 'project_id/%s_id'", r.name))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
