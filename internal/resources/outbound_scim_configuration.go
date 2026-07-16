package resources

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/outboundscim"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &outboundSCIMConfigurationResource{}
	_ resource.ResourceWithConfigure   = &outboundSCIMConfigurationResource{}
	_ resource.ResourceWithImportState = &outboundSCIMConfigurationResource{}
)

type outboundSCIMConfigurationResource struct {
	client *infra.Client
}

func NewOutboundSCIMConfigurationResource() resource.Resource {
	return &outboundSCIMConfigurationResource{}
}

func (r *outboundSCIMConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		r.client = client
	}
}

func (r *outboundSCIMConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_scim_configuration"
}

func (r *outboundSCIMConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = outboundscim.Schema
}

func (r *outboundSCIMConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model outboundscim.OutboundSCIMConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := model.ConfigurationForWrite()
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("configuration"), "Invalid outbound SCIM configuration", err.Error())
		return
	}
	created, err := r.client.CreateOutboundSCIMConfiguration(ctx, model.ProjectID.ValueString(), infra.OutboundSCIMWriteRequest{
		AppID:         model.AppID.ValueString(),
		Configuration: configuration,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating outbound SCIM configuration", err.Error())
		return
	}
	created.Enabled = model.Enabled.ValueBool()
	if !created.Enabled {
		created, err = r.client.SetOutboundSCIMEnabled(ctx, model.ProjectID.ValueString(), infra.OutboundSCIMEnabledRequest{AppID: model.AppID.ValueString(), Enabled: false})
		if err != nil {
			resp.Diagnostics.AddError("Error disabling outbound SCIM configuration", err.Error())
			return
		}
	}
	if err := setOutboundSCIMState(&model, created); err != nil {
		resp.Diagnostics.AddError("Error normalizing outbound SCIM configuration", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *outboundSCIMConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = helpers.ContextWithImportState(ctx, req, resp)
	var model outboundscim.OutboundSCIMConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := r.client.LoadOutboundSCIMConfiguration(ctx, model.ProjectID.ValueString(), model.AppID.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound SCIM configuration", err.Error())
		return
	}
	if err := setOutboundSCIMState(&model, configuration); err != nil {
		resp.Diagnostics.AddError("Error normalizing outbound SCIM configuration", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *outboundSCIMConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan outboundscim.OutboundSCIMConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state outboundscim.OutboundSCIMConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	configuration, err := plan.ConfigurationForWrite()
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("configuration"), "Invalid outbound SCIM configuration", err.Error())
		return
	}
	version := plan.Version.ValueInt64()
	if plan.Version.IsNull() || plan.Version.IsUnknown() {
		version = state.Version.ValueInt64()
	}
	updated, err := r.client.UpdateOutboundSCIMConfiguration(ctx, plan.ProjectID.ValueString(), infra.OutboundSCIMWriteRequest{
		AppID:         plan.AppID.ValueString(),
		Configuration: configuration,
		Version:       version,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating outbound SCIM configuration", err.Error())
		return
	}
	if updated.Enabled != plan.Enabled.ValueBool() {
		updated, err = r.client.SetOutboundSCIMEnabled(ctx, plan.ProjectID.ValueString(), infra.OutboundSCIMEnabledRequest{AppID: plan.AppID.ValueString(), Enabled: plan.Enabled.ValueBool()})
		if err != nil {
			resp.Diagnostics.AddError("Error updating outbound SCIM enabled state", err.Error())
			return
		}
	}
	if err := setOutboundSCIMState(&plan, updated); err != nil {
		resp.Diagnostics.AddError("Error normalizing outbound SCIM configuration", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *outboundSCIMConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model outboundscim.OutboundSCIMConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteOutboundSCIMConfiguration(ctx, model.ProjectID.ValueString(), model.AppID.ValueString())
	if err != nil && !isNotFound(err) {
		resp.Diagnostics.AddError("Error deleting outbound SCIM configuration", err.Error())
	}
}

func (r *outboundSCIMConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.MarkImportState(ctx, resp)
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid Import ID", "Import ID must be in the format 'project_id/app_id'.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func setOutboundSCIMState(model *outboundscim.OutboundSCIMConfigurationModel, configuration *infra.OutboundSCIMConfiguration) error {
	normalized, err := outboundscim.NormalizeConfigurationForState(configuration.Configuration)
	if err != nil {
		return err
	}
	model.ID = types.StringValue(configuration.AppID)
	model.AppID = types.StringValue(configuration.AppID)
	model.Configuration = types.StringValue(normalized)
	model.Enabled = types.BoolValue(configuration.Enabled)
	model.LastExportTime = types.Int64Value(configuration.LastExportTime)
	model.LastProcessingTime = types.Int64Value(configuration.LastProcessingTime)
	model.Failures = types.Int64Value(configuration.Failures)
	model.Version = types.Int64Value(configuration.Version)
	return nil
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	descopeError := descope.AsError(err)
	if descopeError == nil {
		return false
	}
	status, ok := descopeError.Info[descope.ErrorInfoKeys.HTTPResponseStatusCode]
	return ok && fmt.Sprint(status) == fmt.Sprint(http.StatusNotFound)
}
