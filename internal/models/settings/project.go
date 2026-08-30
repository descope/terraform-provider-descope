package settings

import (
	"net/url"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_project_settings is the project-level general settings singleton (id = project_id), matching the console's project page.

var ProjectSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level domain, security, and test user settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          ProjectSettingsAttributes,
}

var ProjectSettingsAttributes = map[string]schema.Attribute{
	"id":                                  stringattr.Identifier(),
	"project_id":                          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"app_url":                             stringattr.Default("", stringattr.URLValidator),
	"custom_domain":                       stringattr.Default(""),
	"approved_domains":                    strsetattr.Default(strsetattr.CommaSeparatedValidator),
	"default_no_sso_apps":                 boolattr.Default(false),
	"tenant_user_isolation":               boolattr.Default(false),
	"allow_auth_hosting_iframe_embedding": boolattr.Default(false),
	"test_users_loginid_regexp":           stringattr.Default(""),
	"test_users_static_otp":               stringattr.Default("", stringattr.OTPValidator),
	"test_users_verifier_regexp":          stringattr.Default(""),
}

type ProjectSettingsModel struct {
	ID                              stringattr.Type `tfsdk:"id"`
	ProjectID                       stringattr.Type `tfsdk:"project_id"`
	AppURL                          stringattr.Type `tfsdk:"app_url"`
	CustomDomain                    stringattr.Type `tfsdk:"custom_domain"`
	ApprovedDomains                 strsetattr.Type `tfsdk:"approved_domains"`
	DefaultNoSSOApps                boolattr.Type   `tfsdk:"default_no_sso_apps"`
	TenantUserIsolation             boolattr.Type   `tfsdk:"tenant_user_isolation"`
	AllowAuthHostingIframeEmbedding boolattr.Type   `tfsdk:"allow_auth_hosting_iframe_embedding"`
	TestUsersLoginIDRegExp          stringattr.Type `tfsdk:"test_users_loginid_regexp"`
	TestUsersStaticOTP              stringattr.Type `tfsdk:"test_users_static_otp"`
	TestUsersVerifierRegExp         stringattr.Type `tfsdk:"test_users_verifier_regexp"`
}

func (m *ProjectSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AppURL, data, "appUrl")
	stringattr.Get(m.CustomDomain, data, "customDomain")
	strsetattr.GetCommaSeparated(m.ApprovedDomains, data, "trustedDomains", h)
	boolattr.Get(m.DefaultNoSSOApps, data, "defaultNoSSOApps")
	boolattr.Get(m.TenantUserIsolation, data, "tenantUserIsolation")
	boolattr.Get(m.AllowAuthHostingIframeEmbedding, data, "allowAuthHostingIframeEmbedding")
	stringattr.Get(m.TestUsersLoginIDRegExp, data, "testUserRegex")
	stringattr.Get(m.TestUsersStaticOTP, data, "testUserFixedAuthToken")
	stringattr.Get(m.TestUsersVerifierRegExp, data, "testUserFixedAuthVerifierRegex")
	data["testUserAllowFixedAuth"] = m.TestUsersStaticOTP.ValueString() != ""
	return data
}

func (m *ProjectSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.AppURL, data, "appUrl")
	stringattr.Set(&m.CustomDomain, data, "customDomain")
	strsetattr.SetCommaSeparated(&m.ApprovedDomains, data, "trustedDomains", h)
	boolattr.Set(&m.DefaultNoSSOApps, data, "defaultNoSSOApps")
	boolattr.Set(&m.TenantUserIsolation, data, "tenantUserIsolation")
	boolattr.Set(&m.AllowAuthHostingIframeEmbedding, data, "allowAuthHostingIframeEmbedding")
	stringattr.Set(&m.TestUsersLoginIDRegExp, data, "testUserRegex")
	if data["testUserAllowFixedAuth"] == true {
		stringattr.Set(&m.TestUsersStaticOTP, data, "testUserFixedAuthToken")
	} else {
		m.TestUsersStaticOTP = stringattr.Value("")
	}
	stringattr.Set(&m.TestUsersVerifierRegExp, data, "testUserFixedAuthVerifierRegex")
}

func (m *ProjectSettingsModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.AppURL, m.CustomDomain, m.TestUsersStaticOTP, m.TestUsersVerifierRegExp) {
		return
	}

	appDomain := ""
	if v := m.AppURL.ValueString(); v != "" {
		// app_url format is enforced by URLValidator; here we only derive the domain for custom_domain
		if appURL, err := url.Parse(v); err == nil {
			appDomain = appURL.Hostname()
		}
	}

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
	}

	if (m.TestUsersStaticOTP.ValueString() == "") != (m.TestUsersVerifierRegExp.ValueString() == "") {
		h.Invalid("The test_users_static_otp and test_users_verifier_regexp attributes must be set together")
	}
}

func (m *ProjectSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *ProjectSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *ProjectSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }
