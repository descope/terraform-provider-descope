package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_session_settings is the project-level session management settings singleton (id = project_id).
// The JWT templates are referenced directly by id (e.g. descope_jwt_template.example.id).

var SessionSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level session management settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          SessionSettingsAttributes,
}

var SessionSettingsAttributes = map[string]schema.Attribute{
	"id":                                  stringattr.Identifier(),
	"project_id":                          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"user_jwt_template":                   stringattr.Default(""),
	"access_key_jwt_template":             stringattr.Default(""),
	"refresh_token_expiration":            durationattr.Default("4 weeks", durationattr.MinimumValue("3 minutes")),
	"refresh_token_rotation":              boolattr.Default(false),
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
}

type SessionSettingsModel struct {
	ID                              stringattr.Type `tfsdk:"id"`
	ProjectID                       stringattr.Type `tfsdk:"project_id"`
	UserJWTTemplate                 stringattr.Type `tfsdk:"user_jwt_template"`
	AccessKeyJWTTemplate            stringattr.Type `tfsdk:"access_key_jwt_template"`
	RefreshTokenExpiration          stringattr.Type `tfsdk:"refresh_token_expiration"`
	RefreshTokenRotation            boolattr.Type   `tfsdk:"refresh_token_rotation"`
	RefreshTokenResponseMethod      stringattr.Type `tfsdk:"refresh_token_response_method"`
	RefreshTokenCookiePolicy        stringattr.Type `tfsdk:"refresh_token_cookie_policy"`
	RefreshTokenCookieDomain        stringattr.Type `tfsdk:"refresh_token_cookie_domain"`
	SessionTokenExpiration          stringattr.Type `tfsdk:"session_token_expiration"`
	SessionTokenResponseMethod      stringattr.Type `tfsdk:"session_token_response_method"`
	SessionTokenCookiePolicy        stringattr.Type `tfsdk:"session_token_cookie_policy"`
	SessionTokenCookieDomain        stringattr.Type `tfsdk:"session_token_cookie_domain"`
	StepUpTokenExpiration           stringattr.Type `tfsdk:"step_up_token_expiration"`
	TrustedDeviceTokenExpiration    stringattr.Type `tfsdk:"trusted_device_token_expiration"`
	AccessKeySessionTokenExpiration stringattr.Type `tfsdk:"access_key_session_token_expiration"`
	EnableInactivity                boolattr.Type   `tfsdk:"enable_inactivity"`
	InactivityTime                  stringattr.Type `tfsdk:"inactivity_time"`
}

func (m *SessionSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.UserJWTTemplate, data, "userTemplateId")
	stringattr.Get(m.AccessKeyJWTTemplate, data, "keyTemplateId")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	boolattr.Get(m.RefreshTokenRotation, data, "rotateJwt")
	getResponseMethod(m.RefreshTokenResponseMethod, data, "refreshTokenResponseMethod")
	stringattr.Get(m.RefreshTokenCookiePolicy, data, "refreshTokenCookiePolicy")
	stringattr.Get(m.RefreshTokenCookieDomain, data, "refreshTokenCookieDomain")
	durationattr.Get(m.SessionTokenExpiration, data, "sessionTokenExpiration")
	getResponseMethod(m.SessionTokenResponseMethod, data, "sessionTokenResponseMethod")
	stringattr.Get(m.SessionTokenCookiePolicy, data, "sessionTokenCookiePolicy")
	stringattr.Get(m.SessionTokenCookieDomain, data, "sessionTokenCookieDomain")
	durationattr.Get(m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Get(m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Get(m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Get(m.EnableInactivity, data, "enableInactivity")
	durationattr.Get(m.InactivityTime, data, "inactivityTime")
	return data
}

func (m *SessionSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.UserJWTTemplate, data, "userTemplateId")
	stringattr.Set(&m.AccessKeyJWTTemplate, data, "keyTemplateId")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	boolattr.Set(&m.RefreshTokenRotation, data, "rotateJwt")
	setResponseMethod(&m.RefreshTokenResponseMethod, data, "refreshTokenResponseMethod", h)
	setCookiePolicy(&m.RefreshTokenCookiePolicy, data, "refreshTokenCookiePolicy")
	stringattr.Set(&m.RefreshTokenCookieDomain, data, "refreshTokenCookieDomain")
	durationattr.Set(&m.SessionTokenExpiration, data, "sessionTokenExpiration")
	setResponseMethod(&m.SessionTokenResponseMethod, data, "sessionTokenResponseMethod", h)
	setCookiePolicy(&m.SessionTokenCookiePolicy, data, "sessionTokenCookiePolicy")
	stringattr.Set(&m.SessionTokenCookieDomain, data, "sessionTokenCookieDomain")
	durationattr.Set(&m.StepUpTokenExpiration, data, "stepupTokenExpiration")
	durationattr.Set(&m.TrustedDeviceTokenExpiration, data, "trustedDeviceTokenExpiration")
	durationattr.Set(&m.AccessKeySessionTokenExpiration, data, "keySessionTokenExpiration")
	boolattr.Set(&m.EnableInactivity, data, "enableInactivity")
	durationattr.Set(&m.InactivityTime, data, "inactivityTime")
}

func (m *SessionSettingsModel) ModifyPlan(h *helpers.Handler, _, _ *SessionSettingsModel) {
	session, ok := durationattr.GetSeconds(m.SessionTokenExpiration)
	if !ok {
		return
	}
	if refresh, ok := durationattr.GetSeconds(m.RefreshTokenExpiration); ok && session > refresh {
		h.Invalid("The session_token_expiration value %s cannot be longer than the refresh_token_expiration value %s", m.SessionTokenExpiration.ValueString(), m.RefreshTokenExpiration.ValueString())
	}
}

func (m *SessionSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *SessionSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *SessionSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

func getResponseMethod(field stringattr.Type, data map[string]any, key string) {
	switch s := field.ValueString(); s {
	case "cookies":
		data[key] = "cookie"
	case "response_body":
		data[key] = "onBody"
	}
}

func setResponseMethod(field *stringattr.Type, data map[string]any, key string, h *helpers.Handler) {
	switch s := data[key]; s {
	case "cookie":
		*field = stringattr.Value("cookies")
	case "onBody", "", nil:
		*field = stringattr.Value("response_body")
	default:
		h.Error("Unexpected token response method", "Expected value to be either 'cookie' or 'onBody', found: '%v'", s)
	}
}

func setCookiePolicy(field *stringattr.Type, data map[string]any, key string) {
	if s, _ := data[key].(string); s != "" {
		*field = stringattr.Value(s)
	} else {
		*field = stringattr.Value("none")
	}
}
