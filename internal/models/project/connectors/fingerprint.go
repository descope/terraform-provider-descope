// Code generated by terragen. DO NOT EDIT.

package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var FingerprintValidator = objattr.NewValidator[FingerprintModel]("must have a valid configuration")

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
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	PublicAPIKey             stringattr.Type `tfsdk:"public_api_key"`
	SecretAPIKey             stringattr.Type `tfsdk:"secret_api_key"`
	UseCloudflareIntegration boolattr.Type   `tfsdk:"use_cloudflare_integration"`
	CloudflareScriptURL      stringattr.Type `tfsdk:"cloudflare_script_url"`
	CloudflareEndpointURL    stringattr.Type `tfsdk:"cloudflare_endpoint_url"`
}

func (m *FingerprintModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "fingerprint"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *FingerprintModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

func (m *FingerprintModel) Validate(h *helpers.Handler) {
	if !m.CloudflareScriptURL.IsNull() && !m.UseCloudflareIntegration.ValueBool() {
		h.Conflict("The cloudflare_script_url field cannot be used unless use_cloudflare_integration is set to true")
	}
	if !m.CloudflareEndpointURL.IsNull() && !m.UseCloudflareIntegration.ValueBool() {
		h.Conflict("The cloudflare_endpoint_url field cannot be used unless use_cloudflare_integration is set to true")
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

func (m *FingerprintModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.PublicAPIKey, c, "publicApiKey")
	stringattr.Nil(&m.SecretAPIKey)
	boolattr.Set(&m.UseCloudflareIntegration, c, "useCloudflareIntegration")
	stringattr.Set(&m.CloudflareScriptURL, c, "cloudflareScriptUrl")
	stringattr.Set(&m.CloudflareEndpointURL, c, "cloudflareEndpointUrl")
}

// Matching

func (m *FingerprintModel) GetName() stringattr.Type {
	return m.Name
}

func (m *FingerprintModel) GetID() stringattr.Type {
	return m.ID
}

func (m *FingerprintModel) SetID(id stringattr.Type) {
	m.ID = id
}
