package provider

import (
	"context"
	"os"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &descopeProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &descopeProvider{
			version: version,
		}
	}
}

type descopeProvider struct {
	version string
}

type descopeProviderConfig struct {
	ProjectID     types.String `tfsdk:"project_id"`
	ManagementKey types.String `tfsdk:"management_key"`
	BaseURL       types.String `tfsdk:"base_url"`
}

func (p *descopeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "descope"
	resp.Version = p.version
}

func (p *descopeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Optional: true,
			},
			"management_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"base_url": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *descopeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Descope provider")

	var config descopeProviderConfig
	diags := req.Config.Get(ctx, &config)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if config.ProjectID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("project_id"), "Unknown Descope Project ID", "The provider cannot create the Descope client as there is an unknown configuration value for the Descope project ID. Either target apply the source of the value first, set the value statically in the configuration, or use the DESCOPE_PROJECT_ID environment variable.")
	}
	if config.ManagementKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("management_key"), "Unknown Descope Management Key", "The provider cannot create the Descope client as there is an unknown configuration value for the Descope management key. Either target apply the source of the value first, set the value statically in the configuration, or use the DESCOPE_MANAGEMENT_KEY environment variable.")
	}
	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("base_url"), "Unknown Descope Base URL", "The provider cannot create the Descope client as there is an unknown configuration value for the Descope base URL. Either target apply the source of the value first, set the value statically in the configuration, or use the DESCOPE_BASE_URL environment variable.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := os.Getenv("DESCOPE_PROJECT_ID")
	if !config.ProjectID.IsNull() {
		projectID = config.ProjectID.ValueString()
	}

	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	if !config.ManagementKey.IsNull() {
		managementKey = config.ManagementKey.ValueString()
	}

	baseURL := os.Getenv("DESCOPE_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	if projectID == "" {
		resp.Diagnostics.AddAttributeError(path.Root("project_id"), "Missing Descope Project ID", "The provider cannot create the Descope client as there is a missing or empty value for the Descope project ID. Set the project_id value in the configuration or use the DESCOPE_PROJECT_ID environment variable. If either is already set, ensure the value is not empty.")
	}
	if managementKey == "" {
		resp.Diagnostics.AddAttributeError(path.Root("management_key"), "Missing Descope Management Key", "The provider cannot create the Descope client as there is a missing or empty value for the Descope management key. Set the management_key value in the configuration or use the DESCOPE_MANAGEMENT_KEY environment variable. If either is already set, ensure the value is not empty.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := infra.NewClient(projectID, managementKey, baseURL)
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Descope provider")
}

func (p *descopeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *descopeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
	}
}
