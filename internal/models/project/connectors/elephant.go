// Code generated by terragen. DO NOT EDIT.

package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ElephantAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key": stringattr.SecretRequired(),
}

// Model

type ElephantModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccessKey stringattr.Type `tfsdk:"access_key"`
}

func (m *ElephantModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "elephant"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *ElephantModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *ElephantModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKey, c, "accessKey")
	return c
}

func (m *ElephantModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Nil(&m.AccessKey)
}

// Matching

func (m *ElephantModel) GetName() stringattr.Type {
	return m.Name
}

func (m *ElephantModel) GetID() stringattr.Type {
	return m.ID
}

func (m *ElephantModel) SetID(id stringattr.Type) {
	m.ID = id
}
