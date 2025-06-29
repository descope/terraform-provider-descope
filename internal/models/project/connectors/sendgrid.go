package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SendGridAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"sender":         objattr.Required[SenderFieldModel](SenderFieldAttributes),
	"authentication": objattr.Required[SendGridAuthFieldModel](SendGridAuthFieldAttributes),
}

// Model

type SendGridModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	Sender objattr.Type[SenderFieldModel]       `tfsdk:"sender"`
	Auth   objattr.Type[SendGridAuthFieldModel] `tfsdk:"authentication"`
}

func (m *SendGridModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "sendgrid"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SendGridModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *SendGridModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	objattr.Get(m.Sender, c, helpers.RootKey, h)
	objattr.Get(m.Auth, c, helpers.RootKey, h)
	return c
}

func (m *SendGridModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	objattr.Set(&m.Sender, c, helpers.RootKey, h)
	objattr.Set(&m.Auth, c, helpers.RootKey, h)
}

// Matching

func (m *SendGridModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SendGridModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SendGridModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Auth

var SendGridAuthFieldAttributes = map[string]schema.Attribute{
	"api_key": stringattr.SecretRequired(),
}

type SendGridAuthFieldModel struct {
	ApiKey stringattr.Type `tfsdk:"api_key"`
}

func (m *SendGridAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ApiKey, data, "apiKey")
	return data
}

func (m *SendGridAuthFieldModel) SetValues(h *helpers.Handler, _ map[string]any) {
	stringattr.Nil(&m.ApiKey)
}
