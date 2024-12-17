package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var MParticleAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_key":    stringattr.SecretRequired(),
	"api_secret": stringattr.SecretRequired(),
	"base_url":   stringattr.Default(""),
}

// Model

type MParticleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
	BaseURL   types.String `tfsdk:"base_url"`
}

func (m *MParticleModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "mparticle"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *MParticleModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *MParticleModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.APISecret, c, "apiSecret")
	stringattr.Get(m.BaseURL, c, "baseUrl")
	return c
}

// Matching

func (m *MParticleModel) GetName() types.String {
	return m.Name
}

func (m *MParticleModel) GetID() types.String {
	return m.ID
}

func (m *MParticleModel) SetID(id types.String) {
	m.ID = id
}
