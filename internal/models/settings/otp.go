package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_otp_settings is the project-level one-time passcode settings singleton (id = project_id).
// The message templates are managed by the descope_email_template, descope_text_template and descope_voice_template resources, selected here by id.

var OTPSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level one-time passcode authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          OTPSettingsAttributes,
}

var OTPSettingsAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"project_id":         stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":           boolattr.Default(false),
	"domain":             stringattr.Default("", stringattr.StandardLenValidator),
	"expiration_time":    durationattr.Default("3 minutes", durationattr.MinimumValue("1 minute")),
	"email_connector_id": stringattr.Default(""),
	"text_connector_id":  stringattr.Default(""),
	"voice_connector_id": stringattr.Default(""),
	"email_template_id":  stringattr.Default(""),
	"text_template_id":   stringattr.Default(""),
	"voice_template_id":  stringattr.Default(""),
}

type OTPSettingsModel struct {
	ID               stringattr.Type `tfsdk:"id"`
	ProjectID        stringattr.Type `tfsdk:"project_id"`
	Disabled         boolattr.Type   `tfsdk:"disabled"`
	Domain           stringattr.Type `tfsdk:"domain"`
	ExpirationTime   stringattr.Type `tfsdk:"expiration_time"`
	EmailConnectorID stringattr.Type `tfsdk:"email_connector_id"`
	TextConnectorID  stringattr.Type `tfsdk:"text_connector_id"`
	VoiceConnectorID stringattr.Type `tfsdk:"voice_connector_id"`
	EmailTemplateID  stringattr.Type `tfsdk:"email_template_id"`
	TextTemplateID   stringattr.Type `tfsdk:"text_template_id"`
	VoiceTemplateID  stringattr.Type `tfsdk:"voice_template_id"`
}

func (m *OTPSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.Domain, data, "domain")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.EmailConnectorID, data, "emailConnectorId")
	stringattr.Get(m.TextConnectorID, data, "textConnectorId")
	stringattr.Get(m.VoiceConnectorID, data, "voiceConnectorId")
	stringattr.Get(m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Get(m.TextTemplateID, data, "textTemplateId")
	stringattr.Get(m.VoiceTemplateID, data, "voiceTemplateId")
	return data
}

func (m *OTPSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.Domain, data, "domain")
	durationattr.SetDefault(&m.ExpirationTime, data, "expirationTime", "3 minutes")
	stringattr.Set(&m.EmailConnectorID, data, "emailConnectorId")
	stringattr.Set(&m.TextConnectorID, data, "textConnectorId")
	stringattr.Set(&m.VoiceConnectorID, data, "voiceConnectorId")
	stringattr.Set(&m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Set(&m.TextTemplateID, data, "textTemplateId")
	stringattr.Set(&m.VoiceTemplateID, data, "voiceTemplateId")
}

func (m *OTPSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *OTPSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *OTPSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
