package settings

import (
	"net/url"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SettingsValidator = objattr.NewValidator[SettingsModel]("must have a valid configuration")

var SettingsAttributes = map[string]schema.Attribute{
	"app_url":                             stringattr.Default(""),
	"custom_domain":                       stringattr.Default(""),
	"approved_domains":                    strsetattr.Default(strsetattr.CommaSeparatedValidator),
	"default_no_sso_apps":                 boolattr.Default(false),
	"refresh_token_rotation":              boolattr.Default(false),
	"refresh_token_expiration":            durationattr.Default("4 weeks", durationattr.MinimumValue("3 minutes")),
	"refresh_token_response_method":       stringattr.Default("response_body", stringvalidator.OneOf("cookies", "response_body")),
	"refresh_token_cookie_policy":         stringattr.Default("none", stringvalidator.OneOf("strict", "lax", "none")),
	"refresh_token_cookie_domain":         stringattr.Default(""),
	"session_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"session_token_response_method":       stringattr.Default("response_body", stringvalidator.OneOf("cookies", "response_body")),
	"session_token_cookie_policy":         stringattr.Default("none", stringvalidator.OneOf("strict", "lax", "none")),
	"session_token_cookie_domain":         stringattr.Default(""),
	"step_up_token_expiration":            durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"trusted_device_token_expiration":     durationattr.Default("365 days", durationattr.MinimumValue("3 minutes")),
	"access_key_session_token_expiration": durationattr.Default("10 minutes", durationattr.MinimumValue("3 minutes")),
	"enable_inactivity":                   boolattr.Default(false),
	"inactivity_time":                     durationattr.Default("12 minutes", durationattr.MinimumValue("10 minutes")),
	"test_users_loginid_regexp":           stringattr.Default(""),
	"test_users_verifier_regexp":          stringattr.Default(""),
	"test_users_static_otp":               stringattr.Default("", stringattr.OTPValidator),
	"user_jwt_template":                   stringattr.Default(""),
	"access_key_jwt_template":             stringattr.Default(""),
	"session_migration":                   objattr.Default[SessionMigrationModel](nil, SessionMigrationAttributes, SessionMigrationValidator),
}

type SettingsModel struct {
	AppURL                          stringattr.Type                     `tfsdk:"app_url"`
	CustomDomain                    stringattr.Type                     `tfsdk:"custom_domain"`
	ApprovedDomain                  strsetattr.Type                     `tfsdk:"approved_domains"`
	DefaultNoSSOApps                boolattr.Type                       `tfsdk:"default_no_sso_apps"`
	RefreshTokenRotation            boolattr.Type                       `tfsdk:"refresh_token_rotation"`
	RefreshTokenExpiration          stringattr.Type                     `tfsdk:"refresh_token_expiration"`
	RefreshTokenResponseMethod      stringattr.Type                     `tfsdk:"refresh_token_response_method"`
	RefreshTokenCookiePolicy        stringattr.Type                     `tfsdk:"refresh_token_cookie_policy"`
	RefreshTokenCookieDomain        stringattr.Type                     `tfsdk:"refresh_token_cookie_domain"`
	SessionTokenExpiration          stringattr.Type                     `tfsdk:"session_token_expiration"`
	SessionTokenResponseMethod      stringattr.Type                     `tfsdk:"session_token_response_method"`
	SessionTokenCookiePolicy        stringattr.Type                     `tfsdk:"session_token_cookie_policy"`
	SessionTokenCookieDomain        stringattr.Type                     `tfsdk:"session_token_cookie_domain"`
	StepUpTokenExpiration           stringattr.Type                     `tfsdk:"step_up_token_expiration"`
	TrustedDeviceTokenExpiration    stringattr.Type                     `tfsdk:"trusted_device_token_expiration"`
	AccessKeySessionTokenExpiration stringattr.Type                     `tfsdk:"access_key_session_token_expiration"`
	EnableInactivity                boolattr.Type                       `tfsdk:"enable_inactivity"`
	InactivityTime                  stringattr.Type                     `tfsdk:"inactivity_time"`
	TestUsersLoginIDRegExp          stringattr.Type                     `tfsdk:"test_users_loginid_regexp"`
	TestUsersVerifierRegExp         stringattr.Type                     `tfsdk:"test_users_verifier_regexp"`
	TestUsersStaticOTP              stringattr.Type                     `tfsdk:"test_users_static_otp"`
	UserJWTTemplate                 stringattr.Type                     `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate            stringattr.Type                     `tfsdk:"access_key_jwt_template"`
	SessionMigration                objattr.Type[SessionMigrationModel] `tfsdk:"session_migration"`
}

func (m *SettingsModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	data := map[string]any{}
	stringattr.Get(m.AppURL, data, "appUrl")
	stringattr.Get(m.CustomDomain, data, "customDomain")
	strsetattr.GetCommaSeparated(m.ApprovedDomain, data, "trustedDomains", h)
	boolattr.Get(m.RefreshTokenRotation, data, "rotateJwt")
	boolattr.Get(m.DefaultNoSSOApps, data, "defaultNoSSOApps")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	if s := m.RefreshTokenResponseMethod.ValueString(); s == "cookies" {
		data["tokenResponseMethod"] = "cookie"
	} else if s == "response_body" {
		data["tokenResponseMethod"] = "onBody"
	} else if s != "" {
		panic("unexpected refresh_token_response_method value: " + s)
	}
	stringattr.Get(m.RefreshTokenCookiePolicy, data, "cookiePolicy")
	stringattr.Get(m.RefreshTokenCookieDomain, data, "domain")
	durationattr.Get(m.SessionTokenExpiration, data, "sessionTokenExpiration")
	if s := m.SessionTokenResponseMethod.ValueString(); s == "cookies" {
		data["sessionTokenResponseMethod"] = "cookie"
	} else if s == "response_body" {
		data["sessionTokenResponseMethod"] = "onBody"
	} else if s != "" {
		panic("unexpected session_token_response_method value: " + s)
	}
	stringattr.Get(m.SessionTokenCookiePolicy, data, "sessionTokenCookiePolicy")
	stringattr.Get(m.SessionTokenCookieDomain, data, "sessionTokenCookieDomain")
	durationattr.Get(m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Get(m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Get(m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	stringattr.Get(m.TestUsersLoginIDRegExp, data, "testUserRegex")
	stringattr.Get(m.TestUsersVerifierRegExp, data, "testUserFixedAuthVerifierRegex")
	stringattr.Get(m.TestUsersStaticOTP, data, "testUserFixedAuthToken")
	data["testUserAllowFixedAuth"] = m.TestUsersStaticOTP.ValueString() != ""
	getJWTTemplate(m.UserJWTTemplate, data, "userTemplateId", "user", h)
	getJWTTemplate(m.AccessKeyJWTTemplate, data, "keyTemplateId", "key", h)
	objattr.Get(m.SessionMigration, data, "externalAuthConfig", h)
	return data
}

func (m *SettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.AppURL, data, "appUrl")
	stringattr.Set(&m.CustomDomain, data, "customDomain")
	strsetattr.SetCommaSeparated(&m.ApprovedDomain, data, "trustedDomains", h)
	boolattr.Set(&m.RefreshTokenRotation, data, "rotateJwt")
	boolattr.Set(&m.DefaultNoSSOApps, data, "defaultNoSSOApps")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	if s := data["tokenResponseMethod"]; s == "cookie" {
		m.RefreshTokenResponseMethod = stringattr.Value("cookies")
	} else if s == "onBody" || s == nil {
		m.RefreshTokenResponseMethod = stringattr.Value("response_body")
	} else {
		h.Error("Unexpected refresh token response method", "Expected value to be either 'cookie' or 'onBody', found: '%v'", s)
	}
	stringattr.Set(&m.RefreshTokenCookiePolicy, data, "cookiePolicy")
	stringattr.Set(&m.RefreshTokenCookieDomain, data, "domain")
	durationattr.Set(&m.SessionTokenExpiration, data, "sessionTokenExpiration")
	if s := data["sessionTokenResponseMethod"]; s == "cookie" {
		m.SessionTokenResponseMethod = stringattr.Value("cookies")
	} else if s == "onBody" || s == nil {
		m.SessionTokenResponseMethod = stringattr.Value("response_body")
	} else {
		h.Error("Unexpected session token response method", "Expected value to be either 'cookie' or 'onBody', found: '%v'", s)
	}
	stringattr.Set(&m.SessionTokenCookiePolicy, data, "sessionTokenCookiePolicy")
	stringattr.Set(&m.SessionTokenCookieDomain, data, "sessionTokenCookieDomain")
	durationattr.Set(&m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Set(&m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Set(&m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
	stringattr.Set(&m.TestUsersLoginIDRegExp, data, "testUserRegex")
	stringattr.Set(&m.TestUsersVerifierRegExp, data, "testUserFixedAuthVerifierRegex")
	if data["testUserAllowFixedAuth"] == true {
		stringattr.Set(&m.TestUsersStaticOTP, data, "testUserFixedAuthToken")
	} else {
		m.TestUsersStaticOTP = stringattr.Value("")
	}
	stringattr.Set(&m.UserJWTTemplate, data, "userTemplateId")     // replaced by template name by UpdateReferences later
	stringattr.Set(&m.AccessKeyJWTTemplate, data, "keyTemplateId") // replaced by template name by UpdateReferences later
	objattr.Set(&m.SessionMigration, data, "externalAuthConfig", h)
}

func (m *SettingsModel) Check(h *helpers.Handler) {
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

	validateCookieDomain := func(value, key string) {
		if value != "" && !strings.HasSuffix(value, ".descope.com") && !strings.HasSuffix(value, ".descope.org") && !strings.HasSuffix(value, ".descope.app") {
			if customDomain == "" {
				h.Missing("The %s attribute requires the custom_domain attribute to be set", key)
			} else if value != customDomain && !strings.HasSuffix(customDomain, "."+value) {
				h.Invalid("The %s attribute must be set to the same domain as the custom_domain attribute or one of its top level domains", key)
			}
		}
	}
	validateCookieDomain(m.RefreshTokenCookieDomain.ValueString(), "refresh_token_cookie_domain")
	validateCookieDomain(m.SessionTokenCookieDomain.ValueString(), "session_token_cookie_domain")

	if (m.TestUsersStaticOTP.ValueString() == "") != (m.TestUsersVerifierRegExp.ValueString() == "") {
		h.Invalid("The test_users_static_otp and test_users_verifier_regexp attributes must be set together")
	}
}

func (m *SettingsModel) Validate(h *helpers.Handler) {
	// XXX move Check here eventually
}

func getJWTTemplate(field stringattr.Type, data map[string]any, key string, typ string, h *helpers.Handler) {
	if v := field; !v.IsNull() && !v.IsUnknown() {
		jwtTemplateName := v.ValueString()
		if jwtTemplateName == "" {
			data[key] = ""
		} else if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, jwtTemplateName); ref == nil {
			h.Error("Unknown JWT template reference", "No %s JWT template named '%s' was defined in the project", typ, jwtTemplateName)
		} else if ref.Type != typ {
			h.Error("Invalid JWT template reference", "The JWT template named '%s' is not a %s template", jwtTemplateName, typ)
		} else {
			h.Log("Setting %s reference to JWT template '%s'", key, jwtTemplateName)
			data[key] = ref.ReferenceValue()
		}
	}
}

func (m *SettingsModel) UpdateReferences(h *helpers.Handler) {
	if m.AccessKeyJWTTemplate.ValueString() != "" {
		replaceJWTTemplateIDWithReference(&m.AccessKeyJWTTemplate, h)
	}
	if m.UserJWTTemplate.ValueString() != "" {
		replaceJWTTemplateIDWithReference(&m.UserJWTTemplate, h)
	}
}

func replaceJWTTemplateIDWithReference(s *stringattr.Type, h *helpers.Handler) {
	if id := s.ValueString(); id != "" {
		ref := h.Refs.Name(id)
		if ref != "" {
			*s = stringattr.Value(ref)
		}
	}
}
