package settings

import (
	"strconv"
	"strings"

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
	// require_invitation is the logical inverse of the server's "projectSelfProvisioning"
	// flag, which controls whether users may self sign-up. The value is always emitted: the
	// server treats an absent projectSelfProvisioning as false (self sign-up disabled), so
	// dropping it would silently require invitations even though the default is to allow
	// self sign-up. A null/unknown require_invitation therefore maps to self-provisioning
	// enabled, matching the server default.
	data["projectSelfProvisioning"] = !requireInvitationValue(m.RequireInvitation)
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
	// Read projectSelfProvisioning tolerantly (the server may encode it as a bool or as a
	// "true"/"false" string) and store its inverse. An absent value mirrors the server
	// default of self sign-up being enabled, i.e. require_invitation = false.
	m.RequireInvitation = boolattr.Value(!selfProvisioningValue(data))
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

// requireInvitationValue resolves require_invitation to a concrete boolean, treating a null
// or unknown value as false: by default invitations are not required and self sign-up is
// enabled.
func requireInvitationValue(b boolattr.Type) bool {
	if b.IsNull() || b.IsUnknown() {
		return false
	}
	return b.ValueBool()
}

// selfProvisioningValue extracts the server's projectSelfProvisioning flag, accepting either
// a native boolean or a string encoding ("true"/"false"). A missing or unrecognized value
// falls back to true, matching the server default of self sign-up being enabled, so the flag
// is never silently interpreted as disabled.
func selfProvisioningValue(data map[string]any) bool {
	switch v := data["projectSelfProvisioning"].(type) {
	case bool:
		return v
	case string:
		if parsed, err := strconv.ParseBool(strings.TrimSpace(v)); err == nil {
			return parsed
		}
	}
	return true
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
