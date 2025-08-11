package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SSOAttributes = map[string]schema.Attribute{
	"disabled":           boolattr.Default(false),
	"merge_users":        boolattr.Default(false),
	"redirect_url":       stringattr.Default(""),
	"sso_suite_settings": objattr.Default(SSOSuiteDefault, SSOSuiteAttributes, SSOSuiteValidator),
}

type SSOModel struct {
	Disabled         boolattr.Type               `tfsdk:"disabled"`
	MergeUsers       boolattr.Type               `tfsdk:"merge_users"`
	RedirectURL      stringattr.Type             `tfsdk:"redirect_url"`
	SSOSuiteSettings objattr.Type[SSOSuiteModel] `tfsdk:"sso_suite_settings"`
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	objattr.Get(m.SSOSuiteSettings, data, helpers.RootKey, h)
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	objattr.Set(&m.SSOSuiteSettings, data, helpers.RootKey, h)
}

// SSO Suite Settings

var SSOSuiteValidator = objattr.NewValidator[SSOSuiteModel]("must have a valid configuration")

var SSOSuiteAttributes = map[string]schema.Attribute{
	"style_id":            stringattr.Default(""),
	"hide_scim":           boolattr.Default(false),
	"hide_groups_mapping": boolattr.Default(false),
	"hide_domains":        boolattr.Default(false),
	"hide_saml":           boolattr.Default(false),
	"hide_oidc":           boolattr.Default(false),
}

type SSOSuiteModel struct {
	StyleId           stringattr.Type `tfsdk:"style_id"`
	HideScim          boolattr.Type   `tfsdk:"hide_scim"`
	HideGroupsMapping boolattr.Type   `tfsdk:"hide_groups_mapping"`
	HideDomains       boolattr.Type   `tfsdk:"hide_domains"`
	HideSaml          boolattr.Type   `tfsdk:"hide_saml"`
	HideOidc          boolattr.Type   `tfsdk:"hide_oidc"`
}

var SSOSuiteDefault = &SSOSuiteModel{
	StyleId:           stringattr.Value(""),
	HideScim:          boolattr.Value(false),
	HideGroupsMapping: boolattr.Value(false),
	HideDomains:       boolattr.Value(false),
	HideSaml:          boolattr.Value(false),
	HideOidc:          boolattr.Value(false),
}

func (m *SSOSuiteModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.StyleId, data, "ssoSuiteStyleId")
	boolattr.Get(m.HideScim, data, "hideSsoSuiteScim")
	boolattr.Get(m.HideGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Get(m.HideDomains, data, "hideSsoSuiteDomains")
	boolattr.Get(m.HideSaml, data, "hideSsoSuiteSaml")
	boolattr.Get(m.HideOidc, data, "hideSsoSuiteOidc")
	return data
}

func (m *SSOSuiteModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.StyleId, data, "ssoSuiteStyleId")
	boolattr.Set(&m.HideScim, data, "hideSsoSuiteScim")
	boolattr.Set(&m.HideGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Set(&m.HideDomains, data, "hideSsoSuiteDomains")
	boolattr.Set(&m.HideSaml, data, "hideSsoSuiteSaml")
	boolattr.Set(&m.HideOidc, data, "hideSsoSuiteOidc")
}

func (m *SSOSuiteModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.HideSaml, m.HideOidc) {
		return
	}

	if m.HideSaml.ValueBool() && m.HideOidc.ValueBool() {
		h.Invalid("The attributes hide_oidc and hide_saml cannot both be true")
	}
}
