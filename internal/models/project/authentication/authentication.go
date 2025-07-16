package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AuthenticationAttributes = map[string]schema.Attribute{
	"otp":            objattr.Default[OTPModel](nil, OTPAttributes),
	"magic_link":     objattr.Default[MagicLinkModel](nil, MagicLinkAttributes),
	"enchanted_link": objattr.Default[EnchantedLinkModel](nil, EnchantedLinkAttributes),
	"embedded_link":  objattr.Default[EmbeddedLinkModel](nil, EmbeddedLinkAttributes),
	"password":       objattr.Default[PasswordModel](nil, PasswordAttributes),
	"oauth":          objattr.Default[OAuthModel](nil, OAuthAttributes, OAuthValidator),
	"sso":            objattr.Default[SSOModel](nil, SSOAttributes),
	"totp":           objattr.Default[TOTPModel](nil, TOTPAttributes),
	"passkeys":       objattr.Default[PasskeysModel](nil, PasskeysAttributes),
}

type AuthenticationModel struct {
	OTP           objattr.Type[OTPModel]           `tfsdk:"otp"`
	MagicLink     objattr.Type[MagicLinkModel]     `tfsdk:"magic_link"`
	EnchantedLink objattr.Type[EnchantedLinkModel] `tfsdk:"enchanted_link"`
	EmbeddedLink  objattr.Type[EmbeddedLinkModel]  `tfsdk:"embedded_link"`
	Password      objattr.Type[PasswordModel]      `tfsdk:"password"`
	OAuth         objattr.Type[OAuthModel]         `tfsdk:"oauth"`
	SSO           objattr.Type[SSOModel]           `tfsdk:"sso"`
	TOTP          objattr.Type[TOTPModel]          `tfsdk:"totp"`
	Passkeys      objattr.Type[PasskeysModel]      `tfsdk:"passkeys"`
}

func (m *AuthenticationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.OTP, data, "otp", h)
	objattr.Get(m.MagicLink, data, "magiclink", h)
	objattr.Get(m.EnchantedLink, data, "enchantedlink", h)
	objattr.Get(m.EmbeddedLink, data, "embeddedlink", h)
	objattr.Get(m.Password, data, "password", h)
	objattr.Get(m.OAuth, data, "oauth", h)
	objattr.Get(m.SSO, data, "sso", h)
	objattr.Get(m.TOTP, data, "totp", h)
	objattr.Get(m.Passkeys, data, "webauthn", h)
	return data
}

func (m *AuthenticationModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.OTP, data, "otp", h)
	objattr.Set(&m.MagicLink, data, "magiclink", h)
	objattr.Set(&m.EnchantedLink, data, "enchantedlink", h)
	objattr.Set(&m.EmbeddedLink, data, "embeddedlink", h)
	objattr.Set(&m.Password, data, "password", h)
	objattr.Set(&m.OAuth, data, "oauth", h)
	objattr.Set(&m.SSO, data, "sso", h)
	objattr.Set(&m.TOTP, data, "totp", h)
	objattr.Set(&m.Passkeys, data, "webauthn", h)
}

func (m *AuthenticationModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.OTP, h)
	objattr.UpdateReferences(&m.MagicLink, h)
	objattr.UpdateReferences(&m.EnchantedLink, h)
	objattr.UpdateReferences(&m.Password, h)
}
