package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AmplitudeAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_key":     stringattr.SecretRequired(),
	"server_url":  stringattr.Default(""),
	"server_zone": stringattr.Default(""),
}

// Model

type AmplitudeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey     types.String `tfsdk:"api_key"`
	ServerURL  types.String `tfsdk:"server_url"`
	ServerZone types.String `tfsdk:"server_zone"`
}

func (m *AmplitudeModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "amplitude"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AmplitudeModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *AmplitudeModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.ServerURL, c, "serverUrl")
	stringattr.Get(m.ServerZone, c, "serverZone")
	return c
}

// Matching

func (m *AmplitudeModel) GetName() types.String {
	return m.Name
}

func (m *AmplitudeModel) GetID() types.String {
	return m.ID
}

func (m *AmplitudeModel) SetID(id types.String) {
	m.ID = id
}
