package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/project/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var InviteSettingsAttributes = map[string]schema.Attribute{
	"require_invitation":   boolattr.Default(false),
	"invite_url":           stringattr.Default(""),
	"add_magiclink_token":  boolattr.Default(false),
	"expire_invited_users": boolattr.Default(false),
	"invite_expiration":    durationattr.Default("1 week", durationattr.MinimumValue("1 hour")),
	"send_email":           boolattr.Default(true),
	"send_text":            boolattr.Default(false),
	"email_service":        objattr.Optional[templates.EmailServiceModel](templates.EmailServiceAttributes, templates.EmailServiceValidator),
}

type InviteSettingsModel struct {
	RequireInvitation  boolattr.Type                             `tfsdk:"require_invitation"`
	InviteURL          stringattr.Type                           `tfsdk:"invite_url"`
	AddMagicLinkToken  boolattr.Type                             `tfsdk:"add_magiclink_token"`
	ExpireInvitedUsers boolattr.Type                             `tfsdk:"expire_invited_users"`
	InviteExpiration   durationattr.Type                         `tfsdk:"invite_expiration"`
	SendEmail          boolattr.Type                             `tfsdk:"send_email"`
	SendText           boolattr.Type                             `tfsdk:"send_text"`
	EmailService       objattr.Type[templates.EmailServiceModel] `tfsdk:"email_service"`
}

var InviteSettingsDefault = &InviteSettingsModel{
	RequireInvitation:  boolattr.Value(false),
	InviteURL:          stringattr.Value(""),
	AddMagicLinkToken:  boolattr.Value(false),
	ExpireInvitedUsers: boolattr.Value(false),
	InviteExpiration:   durationattr.Value("1 week"),
	SendEmail:          boolattr.Value(true),
	SendText:           boolattr.Value(false),
	EmailService:       objattr.Value[templates.EmailServiceModel](nil),
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
	convertKeysFromService(data)
	return data
}

func (m *InviteSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	convertKeysToService(data)
	boolattr.SetNot(&m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Set(&m.InviteURL, data, "inviteUrl")
	boolattr.Set(&m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Set(&m.ExpireInvitedUsers, data, "inviteExpireUser")
	durationattr.Set(&m.InviteExpiration, data, "inviteExpirationTime")
	boolattr.Set(&m.SendEmail, data, "inviteSendEmail")
	boolattr.Set(&m.SendText, data, "inviteSendSms")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
}

func (m *InviteSettingsModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.EmailService, h)
}

func convertKeysFromService(data map[string]any) {
	if v, ok := data["emailServiceProvider"]; ok {
		data["inviteEmailProviderId"] = v
		delete(data, "emailServiceProvider")
	}
	if v, ok := data["emailTemplates"]; ok {
		data["inviteEmailTemplates"] = v
		delete(data, "emailTemplates")
	}
}

func convertKeysToService(data map[string]any) {
	if v, ok := data["inviteEmailProviderId"]; ok {
		data["emailServiceProvider"] = v
		delete(data, "inviteEmailProviderId")
	}
	if v, ok := data["inviteEmailTemplates"]; ok {
		data["emailTemplates"] = v
		delete(data, "inviteEmailTemplates")
	}
}
