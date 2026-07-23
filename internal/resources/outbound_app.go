package resources

import (
	"context"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/outboundapp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &outboundAppResource{}
	_ resource.ResourceWithConfigure   = &outboundAppResource{}
	_ resource.ResourceWithImportState = &outboundAppResource{}
)

type outboundAppResource struct {
	client *infra.Client
}

func NewOutboundAppResource() resource.Resource {
	return &outboundAppResource{}
}

func (r *outboundAppResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *outboundAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_app"
}

func (r *outboundAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = outboundapp.Schema
}

func (r *outboundAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model outboundapp.OutboundAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	secret := configuredOutboundAppSecret(ctx, req.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &descope.CreateOutboundAppRequest{OutboundApp: *model.OutboundApp(ctx, &resp.Diagnostics)}
	if secret != nil {
		request.ClientSecret = *secret
	}
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateOutboundApp(ctx, model.ProjectID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating outbound app", err.Error())
		return
	}
	model.SetOutboundApp(ctx, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *outboundAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = helpers.ContextWithImportState(ctx, req, resp)
	var model outboundapp.OutboundAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.LoadOutboundApp(ctx, model.ProjectID.ValueString(), model.ID.ValueString())
	if isOutboundAppNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound app", err.Error())
		return
	}
	model.SetOutboundApp(ctx, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *outboundAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model outboundapp.OutboundAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	secret := configuredOutboundAppSecret(ctx, req.Config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateOutboundApp(ctx, model.ProjectID.ValueString(), infra.OutboundAppUpdateRequest{
		App:          model.OutboundApp(ctx, &resp.Diagnostics),
		ClientSecret: secret,
	})
	if resp.Diagnostics.HasError() {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating outbound app", err.Error())
		return
	}
	model.SetOutboundApp(ctx, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *outboundAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model outboundapp.OutboundAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOutboundApp(ctx, model.ProjectID.ValueString(), model.ID.ValueString())
	if err != nil && !isOutboundAppNotFound(err) {
		resp.Diagnostics.AddError("Error deleting outbound app", err.Error())
	}
}

func (r *outboundAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.MarkImportState(ctx, resp)
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be in the format 'project_id/app_id'.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func configuredOutboundAppSecret(ctx context.Context, config tfsdk.Config, diagnostics *diag.Diagnostics) *string {
	var secret types.String
	diagnostics.Append(config.GetAttribute(ctx, path.Root("client_secret"), &secret)...)
	if secret.IsUnknown() {
		diagnostics.AddAttributeError(path.Root("client_secret"), "Unknown outbound app client secret", "The client_secret value must be known before the outbound app can be created or updated.")
		return nil
	}
	if secret.IsNull() {
		return nil
	}
	value := secret.ValueString()
	return &value
}

func isOutboundAppNotFound(err error) bool {
	return descope.IsNotFoundError(err)
}
