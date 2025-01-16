package settings

import (
	"net/url"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SettingsValidator = objectattr.NewValidator[SettingsModel]("must have a valid configuration")

var SettingsAttributes = map[string]schema.Attribute{
	"app_url":                             stringattr.Optional(),
	"custom_domain":                       stringattr.Optional(),
	"approved_domains":                    strlistattr.Optional(strlistattr.CommaSeparatedListValidator),
	"token_response_method":               stringattr.Default("response_body", stringvalidator.OneOf("cookies", "response_body")),
	"cookie_policy":                       stringattr.Optional(stringvalidator.OneOf("strict", "lax", "none")),
	"cookie_domain":                       stringattr.Default(""),
	"refresh_token_rotation":              boolattr.Default(false),
	"refresh_token_expiration":            durationattr.Default("4 weeks", durationattr.MinimumValue("3 minutes")),
	"session_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"step_up_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"trusted_device_token_expiration":     durationattr.Default("365 days", durationattr.MinimumValue("3 minutes")),
	"access_key_session_token_expiration": durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"enable_inactivity":                   boolattr.Default(false),
	"inactivity_time":                     durationattr.Default("12 minutes", durationattr.MinimumValue("10 minutes")),
	"test_users_loginid_regexp":           stringattr.Default(""),
	"user_jwt_template":                   stringattr.Optional(),
	"access_key_jwt_template":             stringattr.Optional(),

	// Deprecated
	"domain": stringattr.Renamed("domain", "cookie_domain"),
}

type SettingsModel struct {
	AppURL                          types.String `tfsdk:"app_url"`
	CustomDomain                    types.String `tfsdk:"custom_domain"`
	ApprovedDomain                  []string     `tfsdk:"approved_domains"`
	TokenResponseMethod             types.String `tfsdk:"token_response_method"`
	CookiePolicy                    types.String `tfsdk:"cookie_policy"`
	CookieDomain                    types.String `tfsdk:"cookie_domain"`
	RefreshTokenRotation            types.Bool   `tfsdk:"refresh_token_rotation"`
	RefreshTokenExpiration          types.String `tfsdk:"refresh_token_expiration"`
	SessionTokenExpiration          types.String `tfsdk:"session_token_expiration"`
	StepUpTokenExpiration           types.String `tfsdk:"step_up_token_expiration"`
	TrustedDeviceTokenExpiration    types.String `tfsdk:"trusted_device_token_expiration"`
	AccessKeySessionTokenExpiration types.String `tfsdk:"access_key_session_token_expiration"`
	EnableInactivity                types.Bool   `tfsdk:"enable_inactivity"`
	InactivityTime                  types.String `tfsdk:"inactivity_time"`
	TestUsersLoginIDRegExp          types.String `tfsdk:"test_users_loginid_regexp"`
	UserJWTTemplate                 types.String `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate            types.String `tfsdk:"access_key_jwt_template"`

	// Deprecated
	Domain types.String `tfsdk:"domain"`
}

func (m *SettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AppURL, data, "appUrl")
	stringattr.Get(m.CustomDomain, data, "customDomain")
	strlistattr.GetCommaSeparated(m.ApprovedDomain, data, "trustedDomains")
	if s := m.TokenResponseMethod.ValueString(); s == "cookies" {
		data["tokenResponseMethod"] = "cookie"
	} else if s == "response_body" {
		data["tokenResponseMethod"] = "onBody"
	} else if s != "" {
		panic("unexpected token_response_method value: " + s)
	}
	stringattr.Get(m.CookiePolicy, data, "cookiePolicy")
	stringattr.Get(m.CookieDomain, data, "domain")
	boolattr.Get(m.RefreshTokenRotation, data, "rotateJwt")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Get(m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Get(m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Get(m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Get(m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	stringattr.Get(m.TestUsersLoginIDRegExp, data, "testUserRegex")
	getJWTTemplate(m.UserJWTTemplate, data, "userTemplateId", "user", h)
	getJWTTemplate(m.AccessKeyJWTTemplate, data, "keyTemplateId", "key", h)
	stringattr.Get(m.Domain, data, "domain") // deprecated, replaced by cookie_domain
	return data
}

func (m *SettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Get(m.AppURL, data, "appUrl")
	stringattr.Get(m.CustomDomain, data, "customDomain")
	strlistattr.SetCommaSeparated(&m.ApprovedDomain, data, "trustedDomains")
	if data["tokenResponseMethod"] == "cookie" {
		m.TokenResponseMethod = types.StringValue("cookies")
	} else if data["tokenResponseMethod"] == "onBody" {
		m.TokenResponseMethod = types.StringValue("response_body")
	} else {
		h.Error("Unexpected token response method", "Expected value to be either 'cookie' or 'onBody', found: '%v'", data["tokenResponseMethod"])
	}
	stringattr.Set(&m.CookiePolicy, data, "cookiePolicy")
	// stringattr.Set(&m.CookieDomain, data, "domain") temporarily ignored until domain is removed to prevent inconsistent values
	boolattr.Set(&m.RefreshTokenRotation, data, "rotateJwt")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Set(&m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Set(&m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Set(&m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Set(&m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
	stringattr.Set(&m.TestUsersLoginIDRegExp, data, "testUserRegex")
	stringattr.EnsureKnown(&m.UserJWTTemplate)
	stringattr.EnsureKnown(&m.AccessKeyJWTTemplate)
	stringattr.EnsureKnown(&m.Domain) // deprecated, replaced by cookie_domain
}

func (m *SettingsModel) Validate(h *helpers.Handler) {
	if m.Domain.ValueString() != "" && m.CookieDomain.ValueString() != "" {
		h.Error("Conflicting Attribute Value", "The deprecated domain attribute should not be used together with the cookie_domain attribute")
	}
}

func (m *SettingsModel) Check(h *helpers.Handler) {
	if m.CookieDomain.ValueString() != "" && m.AppURL.ValueString() == "" { // temporary warning instead of error
		h.Warn("Missing Attribute Value", "The cookie_domain attribute should be used together with app_url and custom_domain")
		return
	}

	appDomain := ""
	if v := m.AppURL.ValueString(); v != "" {
		if appURL, err := url.Parse(v); err == nil {
			appDomain = appURL.Hostname()
		}
		if appDomain == "" {
			h.Invalid("The app_url attribute must be a valid URL")
		}
	}

	customDomain := ""
	if v := m.CustomDomain.ValueString(); v != "" {
		if appDomain == "" {
			h.Missing("The custom_domain attribute requires the app_url attribute to be set")
		} else if strings.Contains(v, "://") {
			h.Missing("The custom_domain attribute must be a domain name and not a full URL")
		} else if !strings.HasSuffix(v, "."+appDomain) {
			h.Invalid("The custom_domain attribute must be a subdomain of the app_url domain")
		} else if strings.HasSuffix(v, ".localhost") {
			h.Invalid("The custom_domain attribute cannot be used with the reserved domain 'localhost'")
		}
		for _, domain := range []string{"test", "example", "invalid"} {
			for _, tld := range []string{"com", "net", "org"} {
				if strings.HasSuffix(v, "."+domain+"."+tld) {
					h.Invalid("The custom_domain attribute cannot be used with the reserved domain '%s'", domain+"."+tld)
				}
			}
		}
		customDomain = v
	}

	if v := m.CookieDomain.ValueString(); v != "" && !strings.HasSuffix(v, ".descope.com") && !strings.HasSuffix(v, ".descope.org") {
		if customDomain == "" {
			h.Missing("The cookie_domain attribute requires the custom_domain attribute to be set")
		} else if v != customDomain && !strings.HasSuffix(customDomain, "."+v) {
			h.Invalid("The cookie_domain attribute must be set to the same domain as the custom_domain attribute or one of its top level domains")
		}
	}
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
