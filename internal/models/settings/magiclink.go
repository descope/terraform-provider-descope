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

// descope_magiclink_settings is the project-level magic link settings singleton (id = project_id).
// The message templates are managed by the descope_email_template and descope_text_template resources and selected here by id.

var MagicLinkSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level magic link authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          MagicLinkSettingsAttributes,
}

var MagicLinkSettingsAttributes = map[string]schema.Attribute{
	"id":                stringattr.Identifier(),
	"project_id":        stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":          boolattr.Default(false),
	"expiration_time":   durationattr.Default("3 minutes", durationattr.MinimumValue("1 minute")),
	"redirect_url":      stringattr.Default("", stringattr.URLValidator),
	"email_service":     objattr.Default[EmailServiceRefModel](nil, EmailServiceRefAttributes),
	"text_service":      objattr.Default[TextServiceRefModel](nil, TextServiceRefAttributes),
	"email_template_id": stringattr.Default(""),
	"text_template_id":  stringattr.Default(""),
}

type MagicLinkSettingsModel struct {
	ID              stringattr.Type                    `tfsdk:"id"`
	ProjectID       stringattr.Type                    `tfsdk:"project_id"`
	Disabled        boolattr.Type                      `tfsdk:"disabled"`
	ExpirationTime  stringattr.Type                    `tfsdk:"expiration_time"`
	RedirectURL     stringattr.Type                    `tfsdk:"redirect_url"`
	EmailService    objattr.Type[EmailServiceRefModel] `tfsdk:"email_service"`
	TextService     objattr.Type[TextServiceRefModel]  `tfsdk:"text_service"`
	EmailTemplateID stringattr.Type                    `tfsdk:"email_template_id"`
	TextTemplateID  stringattr.Type                    `tfsdk:"text_template_id"`
}

func (m *MagicLinkSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	objattr.Get(m.TextService, data, helpers.RootKey, h)
	stringattr.Get(m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Get(m.TextTemplateID, data, "textTemplateId")

	useDescopeService(m.EmailService, data, "emailServiceProvider")
	useDescopeService(m.TextService, data, "textServiceProvider")

	return data
}

func (m *MagicLinkSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	objattr.Set(&m.TextService, data, helpers.RootKey, h)
	stringattr.Set(&m.EmailTemplateID, data, "emailTemplateId")
	stringattr.Set(&m.TextTemplateID, data, "textTemplateId")
}

func (m *MagicLinkSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *MagicLinkSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *MagicLinkSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
