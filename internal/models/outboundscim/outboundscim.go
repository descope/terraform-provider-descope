package outboundscim

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var appIDPattern = regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,30}$`)

var Schema = schema.Schema{
	MarkdownDescription: "Manages an outbound SCIM configuration for a federated application.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Resource identifier, equal to the federated application ID.",
		},
		"project_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Descope project ID that owns the outbound SCIM configuration.",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"app_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Federated SAML or OIDC application ID to provision from.",
			Validators: []validator.String{
				stringvalidator.RegexMatches(appIDPattern, "must contain 1 to 30 letters, numbers, hyphens, or underscores"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"configuration": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("{}"),
			MarkdownDescription: "JSON object containing non-secret connector configuration. Put token, password, and secret fields in `secrets`.",
			Validators:          []validator.String{jsonObjectValidator{rejectSecrets: true}},
		},
		"secrets": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			WriteOnly:           true,
			MarkdownDescription: "Write-only JSON object deep-merged into `configuration` for create and update requests. Secret and token values are never stored in Terraform state.",
			Validators:          []validator.String{jsonObjectValidator{}},
		},
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
			MarkdownDescription: "Whether outbound SCIM provisioning is enabled.",
		},
		"last_export_time": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Unix timestamp of the last export.",
		},
		"last_processing_time": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Unix timestamp of the last processing run.",
		},
		"failures": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Current consecutive failure count.",
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "Configuration version used for optimistic concurrency.",
		},
	},
}

type OutboundSCIMConfigurationModel struct {
	ID                 types.String `tfsdk:"id"`
	ProjectID          types.String `tfsdk:"project_id"`
	AppID              types.String `tfsdk:"app_id"`
	Configuration      types.String `tfsdk:"configuration"`
	Secrets            types.String `tfsdk:"secrets"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	LastExportTime     types.Int64  `tfsdk:"last_export_time"`
	LastProcessingTime types.Int64  `tfsdk:"last_processing_time"`
	Failures           types.Int64  `tfsdk:"failures"`
	Version            types.Int64  `tfsdk:"version"`
}
