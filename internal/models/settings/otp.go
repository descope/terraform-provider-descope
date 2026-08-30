package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
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
	"id":                stringattr.Identifier(),
	"project_id":        stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":          boolattr.Default(false),
	"domain":            stringattr.Default("", stringattr.StandardLenValidator),
	"expiration_time":   durationattr.Default("3 minutes", durationattr.MinimumValue("1 minute")),
	"email_service":     objattr.Default[EmailServiceRefModel](nil, EmailServiceRefAttributes),
	"text_service":      objattr.Default[TextServiceRefModel](nil, TextServiceRefAttributes),
	"voice_service":     objattr.Default[VoiceServiceRefModel](nil, VoiceServiceRefAttributes),
	"email_template_id": stringattr.Default(""),
	"text_template_id":  stringattr.Default(""),
	"voice_template_id": stringattr.Default(""),
}

type OTPSettingsModel struct {
	ID              stringattr.Type                    `tfsdk:"id"`
	ProjectID       stringattr.Type                    `tfsdk:"project_id"`
	Disabled        boolattr.Type                      `tfsdk:"disabled"`
	Domain          stringattr.Type                    `tfsdk:"domain"`
	ExpirationTime  stringattr.Type                    `tfsdk:"expiration_time"`
	EmailService    objattr.Type[EmailServiceRefModel] `tfsdk:"email_service"`
	TextService     objattr.Type[TextServiceRefModel]  `tfsdk:"text_service"`
	VoiceService    objattr.Type[VoiceServiceRefModel] `tfsdk:"voice_service"`
	EmailTemplateID stringattr.Type                    `tfsdk:"email_template_id"`
	TextTemplateID  stringattr.Type                    `tfsdk:"text_template_id"`
	VoiceTemplateID stringattr.Type                    `tfsdk:"voice_template_id"`
}

func (m *OTPSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.Domain, data, "domain")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	objattr.Get(m.TextService, data, helpers.RootKey, h)
	objattr.Get(m.VoiceService, data, helpers.RootKey, h)
	stringattr.Get(m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Get(m.TextTemplateID, data, "textTemplateId")
	stringattr.Get(m.VoiceTemplateID, data, "voiceTemplateId")

	useDescopeService(m.EmailService, data, "emailServiceProvider")
	useDescopeService(m.TextService, data, "textServiceProvider")
	useDescopeService(m.VoiceService, data, "voiceServiceProvider")

	return data
}

func (m *OTPSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.Domain, data, "domain")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	objattr.Set(&m.TextService, data, helpers.RootKey, h)
	objattr.Set(&m.VoiceService, data, helpers.RootKey, h)
	stringattr.Set(&m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Set(&m.TextTemplateID, data, "textTemplateId")
	stringattr.Set(&m.VoiceTemplateID, data, "voiceTemplateId")
}

func (m *OTPSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *OTPSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *OTPSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
