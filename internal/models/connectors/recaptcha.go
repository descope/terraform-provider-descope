package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RecaptchaAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"site_key":   stringattr.Required(),
	"secret_key": stringattr.SecretRequired(),
}

// Model

type RecaptchaModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	SiteKey   types.String `tfsdk:"site_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (m *RecaptchaModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "recaptcha"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *RecaptchaModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *RecaptchaModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.SiteKey, c, "siteKey")
	stringattr.Get(m.SecretKey, c, "secretKey")
	return c
}

// Matching

func (m *RecaptchaModel) GetName() types.String {
	return m.Name
}

func (m *RecaptchaModel) GetID() types.String {
	return m.ID
}

func (m *RecaptchaModel) SetID(id types.String) {
	m.ID = id
}
