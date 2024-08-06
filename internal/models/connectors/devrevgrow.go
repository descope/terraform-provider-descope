package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DevRevGrowAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"api_key": stringattr.SecretRequired(),
}

// Model

type DevRevGrowModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey types.String `tfsdk:"api_key"`
}

func (m *DevRevGrowModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "devrev-grow"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *DevRevGrowModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *DevRevGrowModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	return c
}

// Matching

func (m *DevRevGrowModel) GetName() types.String {
	return m.Name
}

func (m *DevRevGrowModel) GetID() types.String {
	return m.ID
}

func (m *DevRevGrowModel) SetID(id types.String) {
	m.ID = id
}
