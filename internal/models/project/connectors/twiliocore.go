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
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.AccountSID, c, "accountSid")
	}
	objectattr.Set(&m.Senders, data, "configuration", h)
	objectattr.Set(&m.Auth, data, "configuration", h)
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
	} `tfsdk:"sms"`
	Voice *struct {
		PhoneNumber types.String `tfsdk:"phone_number"`
	} `tfsdk:"voice"`
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
	m.SMS = helpers.ZVL(m.SMS)
	stringattr.Set(&m.SMS.PhoneNumber, data, "fromPhone")
	stringattr.Set(&m.SMS.MessagingServiceSID, data, "messagingServiceSid")
	m.Voice = helpers.ZVL(m.Voice)
	stringattr.Set(&m.Voice.PhoneNumber, data, "fromPhoneVoice")
}

func (m *TwilioCoreSendersFieldModel) Validate(h *helpers.Handler) {
	if v := m.SMS; v != nil {
		if helpers.HasUnknownValues(v.PhoneNumber, v.MessagingServiceSID) {
			return // skip validation if there are unknown values
		}
		if m.SMS.PhoneNumber.ValueString() == "" && m.SMS.MessagingServiceSID.ValueString() == "" {
			h.Missing("The Twilio Core connector SMS sender requires either the phone_number or messaging_service_sid attribute to be set")
		}
		if m.SMS.PhoneNumber.ValueString() != "" && m.SMS.MessagingServiceSID.ValueString() != "" {
			h.Invalid("The Twilio Core connector SMS sender must only have one of its attributes set")
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
	stringattr.Set(&m.AuthToken, data, "authToken")
	stringattr.Set(&m.APIKey, data, "apiKey")
	stringattr.Set(&m.APISecret, data, "apiSecret")
}

func (m *TwilioAuthFieldModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.AuthToken, m.APIKey, m.APISecret) {
		return // skip validation if there are unknown values
	}
	if m.AuthToken.ValueString() == "" && (m.APIKey.ValueString() == "" && m.APISecret.ValueString() == "") {
		h.Missing("The Twilio Core connector requires an authentication method to be set")
	}
	if m.AuthToken.ValueString() == "" && (m.APIKey.ValueString() == "" || m.APISecret.ValueString() == "") {
		h.Missing("The Twilio Core connector authentication attribute requires both api_key and api_secret to be specified together")
	}
	if m.AuthToken.ValueString() != "" && (m.APIKey.ValueString() != "" || m.APISecret.ValueString() != "") {
		h.Invalid("The Twilio Core connector must only have one authentication method set")
	}
}
