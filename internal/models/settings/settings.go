package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SettingsAttributes = map[string]schema.Attribute{
	"cookie_policy":                       stringattr.Optional(stringvalidator.OneOf("strict", "lax", "none")),
	"domain":                              stringattr.Optional(),
	"approved_domains":                    strlistattr.Optional(strlistattr.CommaSeparatedListValidator),
	"refresh_token_rotation":              boolattr.Default(false),
	"enable_inactivity":                   boolattr.Default(false),
	"inactivity_time":                     durationattr.Default("12 minutes", durationattr.MinimumValue("10 minutes")),
	"refresh_token_expiration":            durationattr.Default("4 weeks", durationattr.MinimumValue("3 minutes")),
	"session_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"step_up_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"trusted_device_token_expiration":     durationattr.Default("365 days", durationattr.MinimumValue("3 minutes")),
	"access_key_session_token_expiration": durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"user_jwt_template":                   stringattr.Optional(),
	"access_key_jwt_template":             stringattr.Optional(),
}

type SettingsModel struct {
	CookiePolicy                    types.String `tfsdk:"cookie_policy"`
	Domain                          types.String `tfsdk:"domain"`
	ApprovedDomain                  []string     `tfsdk:"approved_domains"`
	RefreshTokenRotation            types.Bool   `tfsdk:"refresh_token_rotation"`
	EnableInactivity                types.Bool   `tfsdk:"enable_inactivity"`
	InactivityTime                  types.String `tfsdk:"inactivity_time"`
	RefreshTokenExpiration          types.String `tfsdk:"refresh_token_expiration"`
	SessionTokenExpiration          types.String `tfsdk:"session_token_expiration"`
	StepUpTokenExpiration           types.String `tfsdk:"step_up_token_expiration"`
	TrustedDeviceTokenExpiration    types.String `tfsdk:"trusted_device_token_expiration"`
	AccessKeySessionTokenExpiration types.String `tfsdk:"access_key_session_token_expiration"`
	UserJWTTemplate                 types.String `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate            types.String `tfsdk:"access_key_jwt_template"`
}

func (m *SettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.CookiePolicy, data, "cookiePolicy")
	stringattr.Get(m.Domain, data, "domain")
	strlistattr.GetCommaSeparated(m.ApprovedDomain, data, "trustedDomains")
	boolattr.Get(m.RefreshTokenRotation, data, "rotateJwt")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Get(m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Get(m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Get(m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Get(m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	getJWTTemplate(m.UserJWTTemplate, data, "userTemplateId", "user", h)
	getJWTTemplate(m.AccessKeyJWTTemplate, data, "keyTemplateId", "key", h)
	return data
}

func (m *SettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.CookiePolicy, data, "cookiePolicy")
	stringattr.Set(&m.Domain, data, "domain")
	strlistattr.SetCommaSeparated(&m.ApprovedDomain, data, "trustedDomains")
	boolattr.Set(&m.RefreshTokenRotation, data, "rotateJwt")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Set(&m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Set(&m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Set(&m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Set(&m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	stringattr.EnsureKnown(&m.UserJWTTemplate)
	stringattr.EnsureKnown(&m.AccessKeyJWTTemplate)
}

func getJWTTemplate(field types.String, data map[string]any, key string, typ string, h *helpers.Handler) {
	if v := field; !v.IsNull() && !v.IsUnknown() {
		jwtTemplateName := v.ValueString()
		if jwtTemplateName == "" {
			data[key] = ""
		} else if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, jwtTemplateName); ref == nil {
			h.Error("Unknown JWT template reference", "No JWT template named '%s' for project settings was defined", jwtTemplateName)
		} else if ref.Type != typ {
			h.Error("Invalid JWT template reference", "The JWT template named '%s' is not a %s template", jwtTemplateName, typ)
		} else {
			h.Log("Setting %s reference to JWT template '%s'", key, jwtTemplateName)
			data[key] = ref.ReferenceValue()
		}
	}
}
