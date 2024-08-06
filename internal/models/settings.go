package models

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SettingsAttributes = map[string]schema.Attribute{
	"cookie_policy":            intattr.Optional(),
	"domain":                   stringattr.Optional(),
	"enable_inactivity":        boolattr.Optional(),
	"inactivity_time":          durationattr.Optional(durationattr.MinimumValue("10 minutes")),
	"refresh_token_expiration": durationattr.Optional(durationattr.MinimumValue("2 minutes")),
	"user_jwt_template":        stringattr.Optional(),
	"access_key_jwt_template":  stringattr.Optional(),
}

type SettingsModel struct {
	CookiePolicy           types.Int64  `tfsdk:"cookie_policy"`
	Domain                 types.String `tfsdk:"domain"`
	EnableInactivity       types.Bool   `tfsdk:"enable_inactivity"`
	InactivityTime         types.String `tfsdk:"inactivity_time"`
	RefreshTokenExpiration types.String `tfsdk:"refresh_token_expiration"`
	UserJWTTemplate        types.String `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate   types.String `tfsdk:"access_key_jwt_template"`
}

func (m *SettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	intattr.Get(m.CookiePolicy, data, "cookiePolicy")
	stringattr.Get(m.Domain, data, "domain")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	getJWTTemplate(m.UserJWTTemplate, data, "userTemplateId", h)
	getJWTTemplate(m.AccessKeyJWTTemplate, data, "keyTemplateId", h)
	return data
}

func (m *SettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	intattr.Set(&m.CookiePolicy, data, "cookiePolicy")
	stringattr.Set(&m.Domain, data, "domain")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	stringattr.EnsureKnown(&m.UserJWTTemplate)
	stringattr.EnsureKnown(&m.AccessKeyJWTTemplate)
}

func getJWTTemplate(field types.String, data map[string]any, key string, h *helpers.Handler) {
	if v := field; !v.IsNull() && !v.IsUnknown() {
		jwtTemplateName := v.ValueString()
		if jwtTemplateName == "" {
			data[key] = ""
		} else if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, jwtTemplateName); ref != nil {
			h.Log("Setting %s reference to JWT template '%s'", key, jwtTemplateName)
			data[key] = ref.ReferenceValue()
		} else {
			h.Error("Unknown JWT template reference", "No JWT template named '"+jwtTemplateName+"' for project settings was defined")
		}
	}
}
