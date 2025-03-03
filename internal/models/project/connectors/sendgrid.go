package connectors

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SendGridAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"sender":         objectattr.Required(SenderFieldAttributes),
	"authentication": objectattr.Required(SendGridAuthFieldAttributes),
}

// Model

type SendGridModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Sender *SenderFieldModel       `tfsdk:"sender"`
	Auth   *SendGridAuthFieldModel `tfsdk:"authentication"`
}

func (m *SendGridModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "sendgrid"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SendGridModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	objectattr.Set(&m.Sender, data, "configuration", h)
	objectattr.Set(&m.Auth, data, "configuration", h)
}

// Configuration

func (m *SendGridModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	maps.Copy(c, m.Sender.Values(h))
	maps.Copy(c, m.Auth.Values(h))
	return c
}

// Matching

func (m *SendGridModel) GetName() types.String {
	return m.Name
}

func (m *SendGridModel) GetID() types.String {
	return m.ID
}

func (m *SendGridModel) SetID(id types.String) {
	m.ID = id
}

// Auth

var SendGridAuthFieldAttributes = map[string]schema.Attribute{
	"api_key": stringattr.SecretRequired(),
}

type SendGridAuthFieldModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (m *SendGridAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ApiKey, data, "apiKey")
	return data
}

func (m *SendGridAuthFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ApiKey, data, "apiKey")
}
