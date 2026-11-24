package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_invite_settings is the project-level user invitation settings singleton (id = project_id).
// The invitation email templates are managed by the descope_email_template resource and selected here by id,
// for the invitation itself and for the automatic reminders.

var InviteSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level user invitation settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          InviteSettingsAttributes,
}

var InviteSettingsAttributes = map[string]schema.Attribute{
	"id":                         stringattr.Identifier(),
	"project_id":                 stringattr.Required(stringplanmodifier.RequiresReplace()),
	"require_invitation":         boolattr.Default(false),
	"invite_url":                 stringattr.Default("", stringattr.URLValidator),
	"add_magiclink_token":        boolattr.Default(false),
	"expire_invited_users":       boolattr.Default(false),
	"invite_expiration":          durationattr.Default("1 week", durationattr.MinimumValue("1 hour"), durationattr.MaximumValue("1000 weeks")),
	"send_email":                 boolattr.Default(true),
	"send_text":                  boolattr.Default(false),
	"email_connector_id":         stringattr.Default(""),
	"email_template_id":          stringattr.Default(""),
	"send_reminders":             boolattr.Default(false),
	"reminder_interval":          durationattr.Default("1 day", durationattr.MinimumValue("1 day")),
	"max_reminders":              intattr.Default(3, int64validator.Between(1, 5)),
	"reminder_email_template_id": stringattr.Default(""),
	"reminder_text_template_id":  stringattr.Default(""),
}

type InviteSettingsModel struct {
	ID                      stringattr.Type `tfsdk:"id"`
	ProjectID               stringattr.Type `tfsdk:"project_id"`
	RequireInvitation       boolattr.Type   `tfsdk:"require_invitation"`
	InviteURL               stringattr.Type `tfsdk:"invite_url"`
	AddMagicLinkToken       boolattr.Type   `tfsdk:"add_magiclink_token"`
	ExpireInvitedUsers      boolattr.Type   `tfsdk:"expire_invited_users"`
	InviteExpiration        stringattr.Type `tfsdk:"invite_expiration"`
	SendEmail               boolattr.Type   `tfsdk:"send_email"`
	SendText                boolattr.Type   `tfsdk:"send_text"`
	EmailConnectorID        stringattr.Type `tfsdk:"email_connector_id"`
	EmailTemplateID         stringattr.Type `tfsdk:"email_template_id"`
	SendReminders           boolattr.Type   `tfsdk:"send_reminders"`
	ReminderInterval        stringattr.Type `tfsdk:"reminder_interval"`
	MaxReminders            intattr.Type    `tfsdk:"max_reminders"`
	ReminderEmailTemplateID stringattr.Type `tfsdk:"reminder_email_template_id"`
	ReminderTextTemplateID  stringattr.Type `tfsdk:"reminder_text_template_id"`
}

func (m *InviteSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Get(m.InviteURL, data, "inviteUrl")
	boolattr.Get(m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Get(m.ExpireInvitedUsers, data, "inviteExpireUser")
	durationattr.Get(m.InviteExpiration, data, "inviteExpirationTime")
	boolattr.Get(m.SendEmail, data, "inviteSendEmail")
	boolattr.Get(m.SendText, data, "inviteSendSms")
	stringattr.Get(m.EmailConnectorID, data, "inviteEmailConnectorId")
	stringattr.Get(m.EmailTemplateID, data, "inviteEmailTemplateId")
	boolattr.Get(m.SendReminders, data, "inviteReminderEnabled")
	durationattr.Get(m.ReminderInterval, data, "inviteReminderInterval")
	intattr.Get(m.MaxReminders, data, "inviteReminderMaxCount")
	stringattr.Get(m.ReminderEmailTemplateID, data, "inviteReminderEmailTemplateId")
	stringattr.Get(m.ReminderTextTemplateID, data, "inviteReminderSmsTemplateId")
	return data
}

func (m *InviteSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Set(&m.InviteURL, data, "inviteUrl")
	boolattr.Set(&m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Set(&m.ExpireInvitedUsers, data, "inviteExpireUser")
	durationattr.SetDefault(&m.InviteExpiration, data, "inviteExpirationTime", "1 week")
	boolattr.Set(&m.SendEmail, data, "inviteSendEmail")
	boolattr.Set(&m.SendText, data, "inviteSendSms")
	stringattr.Set(&m.EmailConnectorID, data, "inviteEmailConnectorId")
	stringattr.Set(&m.EmailTemplateID, data, "inviteEmailTemplateId")
	boolattr.Set(&m.SendReminders, data, "inviteReminderEnabled")
	durationattr.SetDefault(&m.ReminderInterval, data, "inviteReminderInterval", "1 day")
	// a project that never turned reminders on stores 0, which is not a valid count, so it reads as the default
	if count, _ := data["inviteReminderMaxCount"].(float64); count > 0 {
		intattr.Set(&m.MaxReminders, data, "inviteReminderMaxCount")
	} else if m.MaxReminders.IsNull() || m.MaxReminders.IsUnknown() {
		m.MaxReminders = intattr.Value(3)
	}
	stringattr.Set(&m.ReminderEmailTemplateID, data, "inviteReminderEmailTemplateId")
	stringattr.Set(&m.ReminderTextTemplateID, data, "inviteReminderSmsTemplateId")
}

func (m *InviteSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *InviteSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *InviteSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
