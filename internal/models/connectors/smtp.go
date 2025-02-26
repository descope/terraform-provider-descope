package connectors

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SMTPAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"sender":         objectattr.Required(SenderFieldAttributes),
	"server":         objectattr.Required(ServerFieldAttributes),
	"authentication": objectattr.Required(SMTPAuthFieldAttributes),
}

// Model

type SMTPModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Sender *SenderFieldModel   `tfsdk:"sender"`
	Server *ServerFieldModel   `tfsdk:"server"`
	Auth   *SMTPAuthFieldModel `tfsdk:"authentication"`
}

func (m *SMTPModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "smtp"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SMTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	objectattr.Set(&m.Sender, data, "configuration", h)
	objectattr.Set(&m.Auth, data, "configuration", h)
	objectattr.Set(&m.Server, data, "configuration", h)
}

// Configuration

func (m *SMTPModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	maps.Copy(c, m.Sender.Values(h))
	maps.Copy(c, m.Server.Values(h))
	maps.Copy(c, m.Auth.Values(h))
	return c
}

// Matching

func (m *SMTPModel) GetName() types.String {
	return m.Name
}

func (m *SMTPModel) GetID() types.String {
	return m.ID
}

func (m *SMTPModel) SetID(id types.String) {
	m.ID = id
}

// Auth

var SMTPAuthFieldAttributes = map[string]schema.Attribute{
	"username": stringattr.Required(),
	"password": stringattr.SecretRequired(),
	"method":   stringattr.Default("plain", stringvalidator.OneOf("plain", "login")),
}

type SMTPAuthFieldModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Method   types.String `tfsdk:"method"`
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
	stringattr.Set(&m.Password, data, "password")
	stringattr.Set(&m.Method, data, "authMethod")
}
