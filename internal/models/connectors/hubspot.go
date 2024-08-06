package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var HubSpotAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"access_token": stringattr.SecretRequired(),
	"base_url":     stringattr.Default(""),
}

// Model

type HubSpotModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessToken types.String `tfsdk:"access_token"`
	BaseURL     types.String `tfsdk:"base_url"`
}

func (m *HubSpotModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "hubspot"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *HubSpotModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *HubSpotModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessToken, c, "accessToken")
	stringattr.Get(m.BaseURL, c, "baseUrl")
	return c
}

// Matching

func (m *HubSpotModel) GetName() types.String {
	return m.Name
}

func (m *HubSpotModel) GetID() types.String {
	return m.ID
}

func (m *HubSpotModel) SetID(id types.String) {
	m.ID = id
}
