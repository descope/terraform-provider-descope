package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AbuseIPDBAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"api_key": stringattr.SecretRequired(),
}

// Model

type AbuseIPDBModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey types.String `tfsdk:"api_key"`
}

func (m *AbuseIPDBModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "abuseipdb"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AbuseIPDBModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *AbuseIPDBModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	return c
}

// Matching

func (m *AbuseIPDBModel) GetName() types.String {
	return m.Name
}

func (m *AbuseIPDBModel) GetID() types.String {
	return m.ID
}

func (m *AbuseIPDBModel) SetID(id types.String) {
	m.ID = id
}
