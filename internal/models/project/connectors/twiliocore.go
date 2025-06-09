package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TwilioCoreAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"account_sid":    stringattr.Required(),
	"senders":        objattr.Required[TwilioCoreSendersFieldModel](TwilioCoreSendersFieldAttributes),
	"authentication": objattr.Required[TwilioAuthFieldModel](TwilioAuthFieldAttributes, TwilioAuthFieldValidator),
}

// Model

type TwilioCoreModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccountSID stringattr.Type                           `tfsdk:"account_sid"`
	Senders    objattr.Type[TwilioCoreSendersFieldModel] `tfsdk:"senders"`
	Auth       objattr.Type[TwilioAuthFieldModel]        `tfsdk:"authentication"`
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
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *TwilioCoreModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccountSID, c, "accountSid")
	objattr.Get(m.Senders, c, helpers.RootKey, h)
	objattr.Get(m.Auth, c, helpers.RootKey, h)
	return c
}

func (m *TwilioCoreModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.AccountSID, c, "accountSid")
	objattr.Set(&m.Senders, c, helpers.RootKey, h)
	objattr.Set(&m.Auth, c, helpers.RootKey, h)
}

// Matching

func (m *TwilioCoreModel) GetName() stringattr.Type {
	return m.Name
}

func (m *TwilioCoreModel) GetID() stringattr.Type {
	return m.ID
}

func (m *TwilioCoreModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Senders

var TwilioCoreSendersFieldAttributes = map[string]schema.Attribute{
	"sms":   objattr.Required[TwilioCoreSendersSMSFieldModel](TwilioCoreSendersSMSFieldAttributes, TwilioCoreSendersSMSFieldValidator),
	"voice": objattr.Optional[TwilioCoreSendersVoiceFieldModel](TwilioCoreSendersVoiceFieldAttributes),
}

type TwilioCoreSendersFieldModel struct {
	SMS   objattr.Type[TwilioCoreSendersSMSFieldModel]   `tfsdk:"sms"`
	Voice objattr.Type[TwilioCoreSendersVoiceFieldModel] `tfsdk:"voice"`
}

func (m *TwilioCoreSendersFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.SMS, data, helpers.RootKey, h)
	objattr.Get(m.Voice, data, helpers.RootKey, h)
	return data
}

func (m *TwilioCoreSendersFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.SMS, data, helpers.RootKey, h)
	objattr.Set(&m.Voice, data, helpers.RootKey, h)
}

// TwilioCoreSendersSMSField

var TwilioCoreSendersSMSFieldValidator = objattr.NewValidator[TwilioCoreSendersSMSFieldModel]("must have valid senders configured")

var TwilioCoreSendersSMSFieldAttributes = map[string]schema.Attribute{
	"phone_number":          stringattr.Default(""),
	"messaging_service_sid": stringattr.Default(""),
}

type TwilioCoreSendersSMSFieldModel struct {
	PhoneNumber         stringattr.Type `tfsdk:"phone_number"`
	MessagingServiceSID stringattr.Type `tfsdk:"messaging_service_sid"`
}

func (m *TwilioCoreSendersSMSFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.PhoneNumber, data, "fromPhone")
	stringattr.Get(m.MessagingServiceSID, data, "messagingServiceSid")
	if m.PhoneNumber.ValueString() != "" {
		data["selectedProp"] = "fromPhone"
	} else {
		data["selectedProp"] = "messagingServiceSid"
	}
	return data
}

func (m *TwilioCoreSendersSMSFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.PhoneNumber, data, "fromPhone")
	stringattr.Set(&m.MessagingServiceSID, data, "messagingServiceSid")
}

func (m *TwilioCoreSendersSMSFieldModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.PhoneNumber, m.MessagingServiceSID) {
		return // skip validation if there are unknown values
	}
	if m.PhoneNumber.ValueString() == "" && m.MessagingServiceSID.ValueString() == "" {
		h.Missing("The Twilio Core connector SMS sender requires either the phone_number or messaging_service_sid attribute to be set")
	}
	if m.PhoneNumber.ValueString() != "" && m.MessagingServiceSID.ValueString() != "" {
		h.Invalid("The Twilio Core connector SMS sender must only have one of its attributes set")
	}
}

// TwilioCoreSendersVoiceField

var TwilioCoreSendersVoiceFieldAttributes = map[string]schema.Attribute{
	"phone_number": stringattr.Required(),
}

type TwilioCoreSendersVoiceFieldModel struct {
	PhoneNumber stringattr.Type `tfsdk:"phone_number"`
}

func (m *TwilioCoreSendersVoiceFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.PhoneNumber, data, "fromPhoneVoice")
	return data
}

func (m *TwilioCoreSendersVoiceFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.PhoneNumber, data, "fromPhoneVoice")
}

// Auth

var TwilioAuthFieldValidator = objattr.NewValidator[TwilioAuthFieldModel]("must have valid senders configured")

var TwilioAuthFieldAttributes = map[string]schema.Attribute{
	"auth_token": stringattr.SecretOptional(),
	"api_key":    stringattr.SecretOptional(),
	"api_secret": stringattr.SecretOptional(),
}

type TwilioAuthFieldModel struct {
	AuthToken stringattr.Type `tfsdk:"auth_token"`
	APIKey    stringattr.Type `tfsdk:"api_key"`
	APISecret stringattr.Type `tfsdk:"api_secret"`
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

func (m *TwilioAuthFieldModel) SetValues(h *helpers.Handler, _ map[string]any) {
	stringattr.Nil(&m.AuthToken)
	stringattr.Nil(&m.APIKey)
	stringattr.Nil(&m.APISecret)
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
