package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AuthenticationAttributes = map[string]schema.Attribute{
	"otp":            objectattr.Optional(OTPAttributes),
	"magic_link":     objectattr.Optional(MagicLinkAttributes),
	"enchanted_link": objectattr.Optional(EnchantedLinkAttributes),
	"embedded_link":  objectattr.Optional(EmbeddedLinkAttributes),
	"password":       objectattr.Optional(PasswordAttributes),
	"oauth":          objectattr.Optional(OAuthAttributes, OAuthValidator),
	"sso":            objectattr.Optional(SSOAttributes),
	"totp":           objectattr.Optional(TOTPAttributes),
	"passkeys":       objectattr.Optional(PasskeysAttributes),
}

type AuthenticationModel struct {
	OTP           *OTPModel           `tfsdk:"otp"`
	MagicLink     *MagicLinkModel     `tfsdk:"magic_link"`
	EnchantedLink *EnchantedLinkModel `tfsdk:"enchanted_link"`
	EmbeddedLink  *EmbeddedLinkModel  `tfsdk:"embedded_link"`
	Password      *PasswordModel      `tfsdk:"password"`
	OAuth         *OAuthModel         `tfsdk:"oauth"`
	SSO           *SSOModel           `tfsdk:"sso"`
	TOTP          *TOTPModel          `tfsdk:"totp"`
	Passkeys      *PasskeysModel      `tfsdk:"passkeys"`
}

func (m *AuthenticationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objectattr.Get(m.OTP, data, "otp", h)
	objectattr.Get(m.MagicLink, data, "magiclink", h)
	objectattr.Get(m.EnchantedLink, data, "enchantedlink", h)
	objectattr.Get(m.EmbeddedLink, data, "embeddedlink", h)
	objectattr.Get(m.Password, data, "password", h)
	objectattr.Get(m.OAuth, data, "oauth", h)
	objectattr.Get(m.SSO, data, "sso", h)
	objectattr.Get(m.TOTP, data, "totp", h)
	objectattr.Get(m.Passkeys, data, "webauthn", h)
	return data
}

func (m *AuthenticationModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.OTP = utils.ZVL(m.OTP)
	m.MagicLink = utils.ZVL(m.MagicLink)
	m.EnchantedLink = utils.ZVL(m.EnchantedLink)
	m.EmbeddedLink = utils.ZVL(m.EmbeddedLink)
	m.Password = utils.ZVL(m.Password)
	m.OAuth = utils.ZVL(m.OAuth)
	m.SSO = utils.ZVL(m.SSO)
	m.TOTP = utils.ZVL(m.TOTP)
	m.Passkeys = utils.ZVL(m.Passkeys)

	objectattr.Set(&m.OTP, data, "otp", h)
	objectattr.Set(&m.MagicLink, data, "magiclink", h)
	objectattr.Set(&m.EnchantedLink, data, "enchantedlink", h)
	objectattr.Set(&m.EmbeddedLink, data, "embeddedlink", h)
	objectattr.Set(&m.Password, data, "password", h)
	objectattr.Set(&m.OAuth, data, "oauth", h)
	objectattr.Set(&m.SSO, data, "sso", h)
	objectattr.Set(&m.TOTP, data, "totp", h)
	objectattr.Set(&m.Passkeys, data, "webauthn", h)
}
