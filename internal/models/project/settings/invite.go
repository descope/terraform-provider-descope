package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var InviteSettingsAttributes = map[string]schema.Attribute{
	"require_invitation":  boolattr.Default(false),
	"invite_url":          stringattr.Default(""),
	"add_magiclink_token": boolattr.Default(false),
	"send_email":          boolattr.Default(true),
	"send_text":           boolattr.Default(false),
}

type InviteSettingsModel struct {
	RequireInvitation boolattr.Type   `tfsdk:"require_invitation"`
	InviteURL         stringattr.Type `tfsdk:"invite_url"`
	AddMagicLinkToken boolattr.Type   `tfsdk:"add_magiclink_token"`
	SendEmail         boolattr.Type   `tfsdk:"send_email"`
	SendText          boolattr.Type   `tfsdk:"send_text"`
}

var InviteSettingsDefault = &InviteSettingsModel{
	RequireInvitation: boolattr.Value(false),
	InviteURL:         stringattr.Value(""),
	AddMagicLinkToken: boolattr.Value(false),
	SendEmail:         boolattr.Value(true),
	SendText:          boolattr.Value(false),
}

func (m *InviteSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Get(m.InviteURL, data, "inviteUrl")
	boolattr.Get(m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Get(m.SendEmail, data, "inviteSendEmail")
	boolattr.Get(m.SendText, data, "inviteSendSms")
	return data
}

func (m *InviteSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Set(&m.InviteURL, data, "inviteUrl")
	boolattr.Set(&m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Set(&m.SendEmail, data, "inviteSendEmail")
	boolattr.Set(&m.SendText, data, "inviteSendSms")
}
