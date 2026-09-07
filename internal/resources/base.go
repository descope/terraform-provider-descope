package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Creates a new resource driven through the CRUD APIs implemented by ops.
func newResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema, ops operations) resource.Resource {
	return &baseResource[T, M]{name: name, schema: sc, singleton: false, ops: ops}
}

// Creates a singleton resource: it always exists for a project, has no server-assigned id (the id is the project id), and imports
// with just that id. Singleton ops leave Create nil - creating applies Update with the project id as the entity id.
func newSingletonResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema, ops operations) resource.Resource {
	return &baseResource[T, M]{name: name, schema: sc, singleton: true, ops: ops}
}

// Use a random model to ensure interface conformance
var (
	_ resource.Resource                   = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithConfigure      = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithImportState    = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithModifyPlan     = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
	_ resource.ResourceWithValidateConfig = &baseResource[accesskey.AccessKeyModel, *accesskey.AccessKeyModel]{}
)

// DestroyChecker lets acceptance tests verify a destroyed entity is gone via the resource's own Read; meaningless for singletons.
type DestroyChecker interface {
	CheckDestroyed(ctx context.Context, client *infra.Client, projectID, appID, id string) error
}

func (r *baseResource[T, M]) CheckDestroyed(ctx context.Context, client *infra.Client, projectID, appID, id string) error {
	var err error
	if r.ops.ScopedRead != nil {
		_, err = r.ops.ScopedRead(ctx, client, projectID, appID, id)
	} else {
		_, err = r.ops.Read(ctx, client, projectID, id)
	}
	if err == nil {
		return fmt.Errorf("%s entity with id %s still exists on the backend after destroy", r.name, id)
	}
	if infra.AsNotFoundError(err) {
		return nil
	}
	return err
}

// validatableModel is the optional hook for plan-time cross-field validation; models without it are only validated per attribute.
type validatableModel interface {
	Validate(*helpers.Handler)
}

// planModifiableModel is the optional hook for cross-attribute plan logic, called after attribute plan modifiers on create and update
// plans (never destroy) with the planned model as the receiver to mutate in place; on create state is a zero-value model (all nulls).
type planModifiableModel[T any] interface {
	ModifyPlan(h *helpers.Handler, config *T, state *T)
}

type baseResource[T any, M helpers.ResourceModel[T]] struct {
	name      string
	schema    schema.Schema
	singleton bool
	ops       operations
	client    *infra.Client
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

func (r *baseResource[T, M]) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	model := M(new(T))
	validatable, ok := any(model).(validatableModel)
	if !ok {
		return
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	validatable.Validate(handler)
}

// Blocks destroy and replace plans for resources with the deletion protection attribute, and rejects changes to scoping attributes
// in all resources. Note: projectResource (project.go) duplicates this wiring by hand - keep it in sync.
func (r *baseResource[T, M]) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !req.Plan.Raw.IsNull() {
		r.modelModifyPlan(ctx, req, resp) // runs on create and update, not on destroy
	}
	if req.State.Raw.IsNull() {
		return // nothing to protect when the resource is being created
	}
	if !req.Plan.Raw.IsNull() {
		checkImmutableAttributes(ctx, r.schema, r.name, req, resp)
	}
	if _, ok := r.schema.Attributes[deletionProtectionAttribute]; !ok {
		return
	}
	if req.Plan.Raw.IsNull() {
		checkDestroyProtection(ctx, req.State, M(new(T)), r.name, &resp.Diagnostics)
		return
	}
	if isPlannedReplace(ctx, r.schema, req) {
		checkReplaceProtection(ctx, req.State, M(new(T)), r.name, &resp.Diagnostics)
	}
}

func (r *baseResource[T, M]) modelModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	model := M(new(T))
	modifiable, ok := any(model).(planModifiableModel[T])
	if !ok {
		return
	}
	config, state := new(T), new(T)
	resp.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	modifiable.ModifyPlan(helpers.NewHandler(ctx, &resp.Diagnostics), config, state)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.Plan.Set(ctx, model)...)
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

	projectID := model.GetProjectID().ValueString()
	var id string
	var data map[string]any
	var err error
	if r.singleton {
		id = projectID
		data, err = r.ops.Update(ctx, r.client, projectID, projectID, values)
	} else {
		id, data, err = r.ops.Create(ctx, r.client, projectID, values)
	}
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid "+r.name+" configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating "+r.name, err.Error())
		return
	}

	model.SetID(types.StringValue(id))
	model.SetValues(handler, data)
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

	var data map[string]any
	var err error
	if r.ops.ScopedRead != nil {
		data, err = r.ops.ScopedRead(ctx, r.client, model.GetProjectID().ValueString(), modelScope(model), model.GetID().ValueString())
	} else {
		data, err = r.ops.Read(ctx, r.client, model.GetProjectID().ValueString(), model.GetID().ValueString())
	}
	if err != nil {
		// an entity deleted out-of-band is dropped from state so the next plan re-creates it, unless it's protected
		if infra.AsNotFoundError(err) && !helpers.IsImportState(ctx) {
			if _, ok := r.schema.Attributes[deletionProtectionAttribute]; ok {
				checkRemovalProtection(ctx, req.State, model, r.name, &resp.Diagnostics)
				if resp.Diagnostics.HasError() {
					return
				}
			}
			tflog.Info(ctx, "Removing "+r.name+" resource from state: not found")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading "+r.name, err.Error())
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	model.SetValues(handler, data)
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

	data, err := r.ops.Update(ctx, r.client, model.GetProjectID().ValueString(), model.GetID().ValueString(), values)
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid "+r.name+" configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating "+r.name, err.Error())
		return
	}

	model.SetValues(handler, data)
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

	// safeguard for deletions that aren't visible at plan time, e.g. forced replacements
	if _, ok := r.schema.Attributes[deletionProtectionAttribute]; ok {
		checkDestroyProtection(ctx, req.State, model, r.name, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var err error
	if r.ops.ScopedDelete != nil {
		err = r.ops.ScopedDelete(ctx, r.client, model.GetProjectID().ValueString(), modelScope(model), model.GetID().ValueString())
	} else {
		err = r.ops.Delete(ctx, r.client, model.GetProjectID().ValueString(), model.GetID().ValueString())
	}
	// an entity that was already deleted out of band still counts as destroyed
	if err != nil && !infra.AsNotFoundError(err) {
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

	if r.singleton {
		if strings.Contains(req.ID, "/") {
			resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be the ID of the project that the "+r.name+" belong to")
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), req.ID)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
		return
	}

	if _, ok := r.schema.Attributes["method"]; ok {
		parts := strings.SplitN(req.ID, "/", 3)
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Import ID must be in the format 'project_id/method/%s_id'", r.name))
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("method"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
		return
	}

	if _, ok := r.schema.Attributes["app_id"]; ok {
		parts := strings.SplitN(req.ID, "/", 3)
		if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Import ID must be in the format 'project_id/app_id/%s_id'", r.name))
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
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

type appScopedModel interface {
	GetAppID() stringattr.Type
}

// methodScopedModel is implemented by models of entities scoped to an auth method, e.g. the messaging template resources.
type methodScopedModel interface {
	GetMethod() stringattr.Type
}

func modelScope(model any) string {
	if m, ok := model.(appScopedModel); ok {
		return m.GetAppID().ValueString()
	}
	if m, ok := model.(methodScopedModel); ok {
		return m.GetMethod().ValueString()
	}
	return ""
}
