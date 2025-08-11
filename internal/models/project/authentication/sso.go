package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SSOAttributes = map[string]schema.Attribute{
	"disabled":                      boolattr.Default(false),
	"merge_users":                   boolattr.Default(false),
	"redirect_url":                  stringattr.Default(""),
	"sso_suite_style_id":            stringattr.Default(""),
	"hide_sso_suite_scim":           boolattr.Default(false),
	"hide_sso_suite_groups_mapping": boolattr.Default(false),
	"hide_sso_suite_domains":        boolattr.Default(false),
	"hide_sso_suite_saml":           boolattr.Default(false),
	"hide_sso_suite_oidc":           boolattr.Default(false),
}

type SSOModel struct {
	Disabled                  boolattr.Type   `tfsdk:"disabled"`
	MergeUsers                boolattr.Type   `tfsdk:"merge_users"`
	RedirectURL               stringattr.Type `tfsdk:"redirect_url"`
	SsoSuiteStyleId           stringattr.Type `tfsdk:"sso_suite_style_id"`
	HideSsoSuiteScim          boolattr.Type   `tfsdk:"hide_sso_suite_scim"`
	HideSsoSuiteGroupsMapping boolattr.Type   `tfsdk:"hide_sso_suite_groups_mapping"`
	HideSsoSuiteDomains       boolattr.Type   `tfsdk:"hide_sso_suite_domains"`
	HideSsoSuiteSaml          boolattr.Type   `tfsdk:"hide_sso_suite_saml"`
	HideSsoSuiteOidc          boolattr.Type   `tfsdk:"hide_sso_suite_oidc"`
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	stringattr.Get(m.SsoSuiteStyleId, data, "ssoSuiteStyleId")
	boolattr.Get(m.HideSsoSuiteScim, data, "hideSsoSuiteScim")
	boolattr.Get(m.HideSsoSuiteGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Get(m.HideSsoSuiteDomains, data, "hideSsoSuiteDomains")
	boolattr.Get(m.HideSsoSuiteSaml, data, "hideSsoSuiteSaml")
	boolattr.Get(m.HideSsoSuiteOidc, data, "hideSsoSuiteOidc")
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	stringattr.Set(&m.SsoSuiteStyleId, data, "ssoSuiteStyleId")
	boolattr.Set(&m.HideSsoSuiteScim, data, "hideSsoSuiteScim")
	boolattr.Set(&m.HideSsoSuiteGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Set(&m.HideSsoSuiteDomains, data, "hideSsoSuiteDomains")
	boolattr.Set(&m.HideSsoSuiteSaml, data, "hideSsoSuiteSaml")
	boolattr.Set(&m.HideSsoSuiteOidc, data, "hideSsoSuiteOidc")
}
