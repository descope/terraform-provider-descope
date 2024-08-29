package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FingerprintValidator = objectattr.NewValidator[FingerprintModel]("must have a valid configuration")

var FingerprintAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"public_api_key":             stringattr.Required(),
	"secret_api_key":             stringattr.SecretRequired(),
	"use_cloudflare_integration": boolattr.Default(false),
	"cloudflare_script_url":      stringattr.Default(""),
	"cloudflare_endpoint_url":    stringattr.Default(""),
}

// Model

type FingerprintModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	PublicAPIKey             types.String `tfsdk:"public_api_key"`
	SecretAPIKey             types.String `tfsdk:"secret_api_key"`
	UseCloudflareIntegration types.Bool   `tfsdk:"use_cloudflare_integration"`
	CloudflareScriptURL      types.String `tfsdk:"cloudflare_script_url"`
	CloudflareEndpointURL    types.String `tfsdk:"cloudflare_endpoint_url"`
}

func (m *FingerprintModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "fingerprint"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *FingerprintModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

func (m *FingerprintModel) Validate(h *helpers.Handler) {
	if !m.CloudflareScriptURL.IsNull() && !m.UseCloudflareIntegration.ValueBool() {
		h.Error("Invalid connector configuration", "The cloudflare_script_url field cannot be used unless use_cloudflare_integration is set to true")
	}
	if !m.CloudflareEndpointURL.IsNull() && !m.UseCloudflareIntegration.ValueBool() {
		h.Error("Invalid connector configuration", "The cloudflare_endpoint_url field cannot be used unless use_cloudflare_integration is set to true")
	}
}

// Configuration

func (m *FingerprintModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.PublicAPIKey, c, "publicApiKey")
	stringattr.Get(m.SecretAPIKey, c, "secretApiKey")
	boolattr.Get(m.UseCloudflareIntegration, c, "useCloudflareIntegration")
	stringattr.Get(m.CloudflareScriptURL, c, "cloudflareScriptUrl")
	stringattr.Get(m.CloudflareEndpointURL, c, "cloudflareEndpointUrl")
	return c
}

// Matching

func (m *FingerprintModel) GetName() types.String {
	return m.Name
}

func (m *FingerprintModel) GetID() types.String {
	return m.ID
}

func (m *FingerprintModel) SetID(id types.String) {
	m.ID = id
}
