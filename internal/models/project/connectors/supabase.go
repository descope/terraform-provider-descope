// Code generated by terragen. DO NOT EDIT.

package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/floatattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SupabaseAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"signing_secret":  stringattr.SecretRequired(),
	"expiration_time": floatattr.Default(60),
}

// Model

type SupabaseModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	SigningSecret  stringattr.Type `tfsdk:"signing_secret"`
	ExpirationTime floatattr.Type  `tfsdk:"expiration_time"`
}

func (m *SupabaseModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "supabase"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SupabaseModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *SupabaseModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.SigningSecret, c, "signingSecret")
	floatattr.Get(m.ExpirationTime, c, "expirationTimeMinutes")
	return c
}

func (m *SupabaseModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Nil(&m.SigningSecret)
	floatattr.Set(&m.ExpirationTime, c, "expirationTimeMinutes")
}

// Matching

func (m *SupabaseModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SupabaseModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SupabaseModel) SetID(id stringattr.Type) {
	m.ID = id
}
