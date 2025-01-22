package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SlackAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"token": stringattr.SecretRequired(),
}

// Model

type SlackModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Token types.String `tfsdk:"token"`
}

func (m *SlackModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "slack"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SlackModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SlackModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.Token, c, "token")
	return c
}

// Matching

func (m *SlackModel) GetName() types.String {
	return m.Name
}

func (m *SlackModel) GetID() types.String {
	return m.ID
}

func (m *SlackModel) SetID(id types.String) {
	m.ID = id
}
