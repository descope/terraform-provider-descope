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
	"cookie_policy":            stringattr.Optional(stringvalidator.OneOf("strict", "lax", "none")),
	"domain":                   stringattr.Optional(),
	"trusted_domains":          strlistattr.Optional(strlistattr.CommaSeparatedListValidator),
	"enable_inactivity":        boolattr.Optional(),
	"inactivity_time":          durationattr.Optional(durationattr.MinimumValue("10 minutes")),
	"refresh_token_expiration": durationattr.Optional(durationattr.MinimumValue("2 minutes")),
	"user_jwt_template":        stringattr.Optional(),
	"access_key_jwt_template":  stringattr.Optional(),
}

type SettingsModel struct {
	CookiePolicy           types.String `tfsdk:"cookie_policy"`
	Domain                 types.String `tfsdk:"domain"`
	TrustedDomains         []string     `tfsdk:"trusted_domains"`
	EnableInactivity       types.Bool   `tfsdk:"enable_inactivity"`
	InactivityTime         types.String `tfsdk:"inactivity_time"`
	RefreshTokenExpiration types.String `tfsdk:"refresh_token_expiration"`
	UserJWTTemplate        types.String `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate   types.String `tfsdk:"access_key_jwt_template"`
}

func (m *SettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.CookiePolicy, data, "cookiePolicy")
	stringattr.Get(m.Domain, data, "domain")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	getJWTTemplate(m.UserJWTTemplate, data, "userTemplateId", "user", h)
	getJWTTemplate(m.AccessKeyJWTTemplate, data, "keyTemplateId", "key", h)
	strlistattr.GetCommaSeparated(m.TrustedDomains, data, "trustedDomains")
	return data
}

func (m *SettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.CookiePolicy, data, "cookiePolicy")
	stringattr.Set(&m.Domain, data, "domain")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	stringattr.EnsureKnown(&m.UserJWTTemplate)
	stringattr.EnsureKnown(&m.AccessKeyJWTTemplate)
	strlistattr.SetCommaSeparated(&m.TrustedDomains, data, "trustedDomains")
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
