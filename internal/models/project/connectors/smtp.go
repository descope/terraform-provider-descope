package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SMTPAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"sender":         objattr.Required[SenderFieldModel](SenderFieldAttributes),
	"server":         objattr.Required[ServerFieldModel](ServerFieldAttributes),
	"authentication": objattr.Required[SMTPAuthFieldModel](SMTPAuthFieldAttributes),
	"use_static_ips": boolattr.Default(false),
}

// Model

type SMTPModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	Sender       objattr.Type[SenderFieldModel]   `tfsdk:"sender"`
	Server       objattr.Type[ServerFieldModel]   `tfsdk:"server"`
	Auth         objattr.Type[SMTPAuthFieldModel] `tfsdk:"authentication"`
	UseStaticIPs boolattr.Type                    `tfsdk:"use_static_ips"`
}

func (m *SMTPModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "smtp"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SMTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *SMTPModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	objattr.Get(m.Sender, c, helpers.RootKey, h)
	objattr.Get(m.Server, c, helpers.RootKey, h)
	objattr.Get(m.Auth, c, helpers.RootKey, h)
	if m.UseStaticIPs.ValueBool() { // don't send field if false in old MP connectors otherwise we'll get an unrecognized key error
		boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	}
	return c
}

func (m *SMTPModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	objattr.Set(&m.Sender, c, helpers.RootKey, h)
	objattr.Set(&m.Auth, c, helpers.RootKey, h)
	objattr.Set(&m.Server, c, helpers.RootKey, h)
	boolattr.Set(&m.UseStaticIPs, c, "useStaticIps")
}

// Matching

func (m *SMTPModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SMTPModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SMTPModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Auth

var SMTPAuthFieldAttributes = map[string]schema.Attribute{
	"username": stringattr.Required(),
	"password": stringattr.SecretRequired(),
	"method":   stringattr.Default("plain", stringvalidator.OneOf("plain", "login")),
}

type SMTPAuthFieldModel struct {
	Username stringattr.Type `tfsdk:"username"`
	Password stringattr.Type `tfsdk:"password"`
	Method   stringattr.Type `tfsdk:"method"`
}

func (m *SMTPAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Username, data, "username")
	stringattr.Get(m.Password, data, "password")
	stringattr.Get(m.Method, data, "authMethod")
	return data
}

func (m *SMTPAuthFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Username, data, "username")
	stringattr.Nil(&m.Password)
	stringattr.Set(&m.Method, data, "authMethod")
}
