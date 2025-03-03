package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IntercomAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"token":  stringattr.SecretRequired(),
	"region": stringattr.Default("US"),
}

// Model

type IntercomModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Token  types.String `tfsdk:"token"`
	Region types.String `tfsdk:"region"`
}

func (m *IntercomModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "intercom"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *IntercomModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.Token, c, "token")
		stringattr.Set(&m.Region, c, "region")
	}
}

// Configuration

func (m *IntercomModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.Token, c, "token")
	stringattr.Get(m.Region, c, "region")
	return c
}

// Matching

func (m *IntercomModel) GetName() types.String {
	return m.Name
}

func (m *IntercomModel) GetID() types.String {
	return m.ID
}

func (m *IntercomModel) SetID(id types.String) {
	m.ID = id
}
