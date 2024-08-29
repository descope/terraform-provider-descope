package connectors

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TwilioCoreAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"account_sid":    stringattr.Required(),
	"senders":        objectattr.Required(TwilioCoreSendersFieldAttributes, TwilioCoreSendersFieldValidator),
	"authentication": objectattr.Required(TwilioAuthFieldAttributes, TwilioAuthFieldValidator),
}

// Model

type TwilioCoreModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccountSID types.String                 `tfsdk:"account_sid"`
	Senders    *TwilioCoreSendersFieldModel `tfsdk:"senders"`
	Auth       *TwilioAuthFieldModel        `tfsdk:"authentication"`
}

func (m *TwilioCoreModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "twilio-core"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *TwilioCoreModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the configuration
}

// Configuration

func (m *TwilioCoreModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccountSID, c, "accountSid")
	maps.Copy(c, m.Senders.Values(h))
	maps.Copy(c, m.Auth.Values(h))
	return c
}

// Matching

func (m *TwilioCoreModel) GetName() types.String {
	return m.Name
}

func (m *TwilioCoreModel) GetID() types.String {
	return m.ID
}

func (m *TwilioCoreModel) SetID(id types.String) {
	m.ID = id
}

// Senders

var TwilioCoreSendersFieldValidator = objectattr.NewValidator[TwilioCoreSendersFieldModel]("must have valid senders configured")

var TwilioCoreSendersFieldAttributes = map[string]schema.Attribute{
	"sms":   objectattr.Required(TwilioCoreSendersSMSFieldAttributes),
	"voice": objectattr.Optional(TwilioCoreSendersVoiceFieldAttributes),
}

var TwilioCoreSendersSMSFieldAttributes = map[string]schema.Attribute{
	"phone_number":          stringattr.Default(""),
	"messaging_service_sid": stringattr.Default(""),
}

var TwilioCoreSendersVoiceFieldAttributes = map[string]schema.Attribute{
	"phone_number": stringattr.Required(),
}

type TwilioCoreSendersFieldModel struct {
	SMS *struct {
		PhoneNumber         types.String `tfsdk:"phone_number"`
		MessagingServiceSID types.String `tfsdk:"messaging_service_sid"`
	} `tfksd:"sms"`
	Voice *struct {
		PhoneNumber types.String `tfsdk:"phone_number"`
	} `tfksd:"voice"`
}

func (m *TwilioCoreSendersFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	if v := m.SMS; v != nil {
		stringattr.Get(v.PhoneNumber, data, "fromPhone")
		stringattr.Get(v.MessagingServiceSID, data, "messagingServiceSid")
		if v.PhoneNumber.ValueString() != "" {
			data["selectedProp"] = "fromPhone"
		} else {
			data["selectedProp"] = "messagingServiceSid"
		}
	}
	if v := m.Voice; v != nil {
		stringattr.Get(v.PhoneNumber, data, "fromPhoneVoice")
	}
	return data
}

func (m *TwilioCoreSendersFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the configuration
}

func (m *TwilioCoreSendersFieldModel) Validate(h *helpers.Handler) {
	if v := m.SMS; v != nil {
		if m.SMS.PhoneNumber.ValueString() == "" && m.SMS.MessagingServiceSID.ValueString() == "" {
			h.Error("Invalid Twilio Core senders field", "The connector requires an SMS sender field to be set")
		}
		if m.SMS.PhoneNumber.ValueString() != "" && m.SMS.MessagingServiceSID.ValueString() != "" {
			h.Error("Invalid Twilio Core senders field", "The connector requires only one one SMS sender field to be set")
		}
	}
}

// Auth

var TwilioAuthFieldValidator = objectattr.NewValidator[TwilioAuthFieldModel]("must have valid senders configured")

var TwilioAuthFieldAttributes = map[string]schema.Attribute{
	"auth_token": stringattr.SecretOptional(),
	"api_key":    stringattr.SecretOptional(),
	"api_secret": stringattr.SecretOptional(),
}

type TwilioAuthFieldModel struct {
	AuthToken types.String `tfsdk:"auth_token"`
	APIKey    types.String `tfsdk:"api_key"`
	APISecret types.String `tfsdk:"api_secret"`
}

func (m *TwilioAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AuthToken, data, "authToken")
	stringattr.Get(m.APIKey, data, "apiKey")
	stringattr.Get(m.APISecret, data, "apiSecret")
	if m.AuthToken.ValueString() != "" {
		data["selectedAuthProp"] = "methodAuthToken"
	} else {
		data["selectedAuthProp"] = "methodApiSecret"
	}
	return data
}

func (m *TwilioAuthFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the configuration
}

func (m *TwilioAuthFieldModel) Validate(h *helpers.Handler) {
	if m.AuthToken.ValueString() == "" && (m.APIKey.ValueString() == "" && m.APISecret.ValueString() == "") {
		h.Error("Invalid Twilio authentication field", "The connector requires an authentication method to be set")
	}
	if m.AuthToken.ValueString() == "" && (m.APIKey.ValueString() == "" || m.APISecret.ValueString() == "") {
		h.Error("Invalid Twilio authentication field", "The api_key and api_secret fields must be specified together")
	}
	if m.AuthToken.ValueString() != "" && (m.APIKey.ValueString() != "" || m.APISecret.ValueString() != "") {
		h.Error("Invalid Twilio authentication field", "The connector requires only one one authentication method to be set")
	}
}
