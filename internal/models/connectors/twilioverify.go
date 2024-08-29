package connectors

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TwilioVerifyAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"account_sid":    stringattr.Required(),
	"service_sid":    stringattr.Required(),
	"sender":         stringattr.Default(""),
	"authentication": objectattr.Required(TwilioAuthFieldAttributes, TwilioAuthFieldValidator),
}

// Model

type TwilioVerifyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccountSID types.String          `tfsdk:"account_sid"`
	ServiceSID types.String          `tfsdk:"service_sid"`
	Sender     types.String          `tfsdk:"sender"`
	Auth       *TwilioAuthFieldModel `tfsdk:"authentication"`
}

func (m *TwilioVerifyModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "twilio-verify"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *TwilioVerifyModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the configuration
}

// Configuration

func (m *TwilioVerifyModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccountSID, c, "accountSid")
	stringattr.Get(m.ServiceSID, c, "verifyServiceSid")
	stringattr.Get(m.Sender, c, "from")
	maps.Copy(c, m.Auth.Values(h))
	return c
}

// Matching

func (m *TwilioVerifyModel) GetName() types.String {
	return m.Name
}

func (m *TwilioVerifyModel) GetID() types.String {
	return m.ID
}

func (m *TwilioVerifyModel) SetID(id types.String) {
	m.ID = id
}
