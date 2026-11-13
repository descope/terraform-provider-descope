package resources

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewOutboundAppDataSource() datasource.DataSource {
	return &outboundAppDataSource{}
}

var (
	_ datasource.DataSource              = &outboundAppDataSource{}
	_ datasource.DataSourceWithConfigure = &outboundAppDataSource{}
)

type outboundAppDataSource struct {
	client *infra.Client
}

type outboundAppDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	AppType     types.String `tfsdk:"app_type"`
}

func (d *outboundAppDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if client, ok := req.ProviderData.(*infra.Client); ok {
		d.client = client
	}
}

func (d *outboundAppDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_app"
}

func (d *outboundAppDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an existing outbound application in a Descope project, and fails when the project has no outbound application with the given ID. Use it to depend on an application that is managed outside of Terraform, such as one whose client secret is entered in the Descope console.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the outbound application."},
			"project_id":  schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the Descope project this outbound application belongs to."},
			"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the outbound application."},
			"description": schema.StringAttribute{Computed: true, MarkdownDescription: "The description of the outbound application."},
			"app_type":    schema.StringAttribute{Computed: true, MarkdownDescription: "The kind of credential the application holds, either `oauth` or `apikey`."},
		},
	}
}

func (d *outboundAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model outboundAppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID, id := model.ProjectID.ValueString(), model.ID.ValueString()
	body, err := d.client.Get(ctx, projectID, "/v1/mgmt/outbound/app/"+id, nil)
	var app map[string]any
	if err == nil {
		app, err = unwrapEntity(body, "app")
	}
	if infra.AsNotFoundError(err) {
		resp.Diagnostics.AddError("Outbound app not found", fmt.Sprintf("No outbound app with id %q exists in project %q.", id, projectID))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading outbound_app", err.Error())
		return
	}

	name, _ := app["name"].(string)
	description, _ := app["description"].(string)
	appType, _ := app["appType"].(string)
	model.Name = types.StringValue(name)
	model.Description = types.StringValue(description)
	model.AppType = types.StringValue(appType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
