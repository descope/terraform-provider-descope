package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_invite_settings is the project-level user invitation settings singleton (id = project_id).

var InviteSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level user invitation settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          InviteSettingsAttributes,
}

var InviteSettingsAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Identifier(),
	"project_id":           stringattr.Required(stringplanmodifier.RequiresReplace()),
	"require_invitation":   boolattr.Default(false),
	"invite_url":           stringattr.Default("", stringattr.URLValidator),
	"add_magiclink_token":  boolattr.Default(false),
	"expire_invited_users": boolattr.Default(false),
	"invite_expiration":    durationattr.Default("1 week", durationattr.MinimumValue("1 hour"), durationattr.MaximumValue("1000 weeks")),
	"send_email":           boolattr.Default(true),
	"send_text":            boolattr.Default(false),
	"email_service":        objattr.Default[templates.EmailServiceIDModel](nil, templates.EmailServiceIDAttributes, templates.EmailServiceIDValidator),
}

type InviteSettingsModel struct {
	ID                 stringattr.Type                             `tfsdk:"id"`
	ProjectID          stringattr.Type                             `tfsdk:"project_id"`
	RequireInvitation  boolattr.Type                               `tfsdk:"require_invitation"`
	InviteURL          stringattr.Type                             `tfsdk:"invite_url"`
	AddMagicLinkToken  boolattr.Type                               `tfsdk:"add_magiclink_token"`
	ExpireInvitedUsers boolattr.Type                               `tfsdk:"expire_invited_users"`
	InviteExpiration   stringattr.Type                             `tfsdk:"invite_expiration"`
	SendEmail          boolattr.Type                               `tfsdk:"send_email"`
	SendText           boolattr.Type                               `tfsdk:"send_text"`
	EmailService       objattr.Type[templates.EmailServiceIDModel] `tfsdk:"email_service"`
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
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	useDescopeService(m.EmailService, data, "emailServiceProvider")
	nestInviteEmailService(data)
	return data
}

func (m *InviteSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	flattenInviteEmailService(data)
	boolattr.SetNot(&m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Set(&m.InviteURL, data, "inviteUrl")
	boolattr.Set(&m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Set(&m.ExpireInvitedUsers, data, "inviteExpireUser")
	durationattr.Set(&m.InviteExpiration, data, "inviteExpirationTime")
	boolattr.Set(&m.SendEmail, data, "inviteSendEmail")
	boolattr.Set(&m.SendText, data, "inviteSendSms")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
}

func (m *InviteSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *InviteSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *InviteSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

// nestInviteEmailService moves the flat email service keys into the inviteEmailService object the API expects.
func nestInviteEmailService(data map[string]any) {
	service := map[string]any{"providerId": data["emailServiceProvider"]}
	delete(data, "emailServiceProvider")
	if v, ok := data["emailTemplates"]; ok {
		service["templates"] = v
		delete(data, "emailTemplates")
	} else {
		service["templates"] = []any{}
	}
	data["inviteEmailService"] = service
}

// flattenInviteEmailService is the inverse, exposing the API's inviteEmailService object as flat keys.
func flattenInviteEmailService(data map[string]any) {
	service, _ := data["inviteEmailService"].(map[string]any)
	if service == nil {
		return
	}
	data["emailServiceProvider"] = service["providerId"]
	data["emailTemplates"] = service["templates"]
	delete(data, "inviteEmailService")
}
