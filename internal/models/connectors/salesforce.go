package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SalesforceAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"base_url":      stringattr.Required(),
	"client_id":     stringattr.Required(),
	"client_secret": stringattr.SecretRequired(),
	"version":       stringattr.Required(),
}

// Model

type SalesforceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	BaseURL      types.String `tfsdk:"base_url"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Version      types.String `tfsdk:"version"`
}

func (m *SalesforceModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "salesforce"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SalesforceModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SalesforceModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.BaseURL, c, "baseUrl")
	stringattr.Get(m.ClientID, c, "clientId")
	stringattr.Get(m.ClientSecret, c, "clientSecret")
	stringattr.Get(m.Version, c, "version")
	return c
}

// Matching

func (m *SalesforceModel) GetName() types.String {
	return m.Name
}

func (m *SalesforceModel) GetID() types.String {
	return m.ID
}

func (m *SalesforceModel) SetID(id types.String) {
	m.ID = id
}
