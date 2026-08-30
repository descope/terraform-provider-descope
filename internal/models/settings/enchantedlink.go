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

// descope_enchantedlink_settings is the project-level enchanted link settings singleton (id = project_id).
// The message templates are managed by the descope_email_template resource and selected here by id.

var EnchantedLinkSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level enchanted link authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          EnchantedLinkSettingsAttributes,
}

var EnchantedLinkSettingsAttributes = map[string]schema.Attribute{
	"id":                stringattr.Identifier(),
	"project_id":        stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":          boolattr.Default(false),
	"expiration_time":   durationattr.Default("3 minutes", durationattr.MinimumValue("1 minute")),
	"redirect_url":      stringattr.Default("", stringattr.URLValidator),
	"email_service":     objattr.Default[EmailServiceRefModel](nil, EmailServiceRefAttributes),
	"email_template_id": stringattr.Default(""),
}

type EnchantedLinkSettingsModel struct {
	ID              stringattr.Type                    `tfsdk:"id"`
	ProjectID       stringattr.Type                    `tfsdk:"project_id"`
	Disabled        boolattr.Type                      `tfsdk:"disabled"`
	ExpirationTime  stringattr.Type                    `tfsdk:"expiration_time"`
	RedirectURL     stringattr.Type                    `tfsdk:"redirect_url"`
	EmailService    objattr.Type[EmailServiceRefModel] `tfsdk:"email_service"`
	EmailTemplateID stringattr.Type                    `tfsdk:"email_template_id"`
}

func (m *EnchantedLinkSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	stringattr.Get(m.EmailTemplateID, data, "emailTemplateId")

	useDescopeService(m.EmailService, data, "emailServiceProvider")

	return data
}

func (m *EnchantedLinkSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	stringattr.Set(&m.EmailTemplateID, data, "emailTemplateId")
}

func (m *EnchantedLinkSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *EnchantedLinkSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *EnchantedLinkSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
