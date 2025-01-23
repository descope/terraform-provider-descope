package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var HubSpotAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_token":   stringattr.SecretRequired(),
	"base_url":       stringattr.Default(""),
	"use_static_ips": boolattr.Default(false),
}

// Model

type HubSpotModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessToken  types.String `tfsdk:"access_token"`
	BaseURL      types.String `tfsdk:"base_url"`
	UseStaticIPs types.Bool   `tfsdk:"use_static_ips"`
}

func (m *HubSpotModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "hubspot"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *HubSpotModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.AccessToken, c, "accessToken")
		stringattr.Set(&m.BaseURL, c, "baseUrl")
		boolattr.Set(&m.UseStaticIPs, c, "useStaticIps")
	}
}

// Configuration

func (m *HubSpotModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessToken, c, "accessToken")
	stringattr.Get(m.BaseURL, c, "baseUrl")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
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
