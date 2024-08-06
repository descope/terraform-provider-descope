package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var VeriffAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"api_key":    stringattr.Required(),
	"secret_key": stringattr.SecretRequired(),
	"base_url":   stringattr.Default(""),
}

// Model

type VeriffModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey    types.String `tfsdk:"api_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	BaseURL   types.String `tfsdk:"base_url"`
}

func (m *VeriffModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "veriff"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *VeriffModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *VeriffModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.SecretKey, c, "secretKey")
	stringattr.Get(m.BaseURL, c, "baseUrl")
	return c
}

// Matching

func (m *VeriffModel) GetName() types.String {
	return m.Name
}

func (m *VeriffModel) GetID() types.String {
	return m.ID
}

func (m *VeriffModel) SetID(id types.String) {
	m.ID = id
}
