package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TwilioVerifyAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"account_sid":    stringattr.Required(),
	"service_sid":    stringattr.Required(),
	"sender":         stringattr.Default(""),
	"authentication": objattr.Required[TwilioAuthFieldModel](TwilioAuthFieldAttributes, TwilioAuthFieldValidator),
}

// Model

type TwilioVerifyModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccountSID stringattr.Type                    `tfsdk:"account_sid"`
	ServiceSID stringattr.Type                    `tfsdk:"service_sid"`
	Sender     stringattr.Type                    `tfsdk:"sender"`
	Auth       objattr.Type[TwilioAuthFieldModel] `tfsdk:"authentication"`
}

func (m *TwilioVerifyModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "twilio-verify"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *TwilioVerifyModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *TwilioVerifyModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccountSID, c, "accountSid")
	stringattr.Get(m.ServiceSID, c, "verifyServiceSid")
	stringattr.Get(m.Sender, c, "from")
	objattr.Get(m.Auth, c, helpers.RootKey, h)
	return c
}

func (m *TwilioVerifyModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.AccountSID, c, "accountSid")
	stringattr.Set(&m.ServiceSID, c, "verifyServiceSid")
	stringattr.Set(&m.Sender, c, "from")
	objattr.Set(&m.Auth, c, helpers.RootKey, h)
}

// Matching

func (m *TwilioVerifyModel) GetName() stringattr.Type {
	return m.Name
}

func (m *TwilioVerifyModel) GetID() stringattr.Type {
	return m.ID
}

func (m *TwilioVerifyModel) SetID(id stringattr.Type) {
	m.ID = id
}
