package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SSOAttributes = map[string]schema.Attribute{
	"disabled":                boolattr.Default(false),
	"merge_users":             boolattr.Default(false),
	"redirect_url":            stringattr.Default(""),
	"sso_suite_settings":      objattr.Default(SSOSuiteDefault, SSOSuiteAttributes, SSOSuiteValidator),
	"allow_duplicate_domains": boolattr.Default(false),
	"allow_override_roles":    boolattr.Default(false),
	"groups_priority":         boolattr.Default(false),
}

type SSOModel struct {
	Disabled                               boolattr.Type               `tfsdk:"disabled"`
	MergeUsers                             boolattr.Type               `tfsdk:"merge_users"`
	RedirectURL                            stringattr.Type             `tfsdk:"redirect_url"`
	SSOSuiteSettings                       objattr.Type[SSOSuiteModel] `tfsdk:"sso_suite_settings"`
	AllowDuplicateSSODomainsInOtherTenants boolattr.Type               `tfsdk:"allow_duplicate_domains"`
	AllowOverrideRoles                     boolattr.Type               `tfsdk:"allow_override_roles"`
	GroupsPriority                         boolattr.Type               `tfsdk:"groups_priority"`
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	objattr.Get(m.SSOSuiteSettings, data, helpers.RootKey, h)
	boolattr.Get(m.AllowDuplicateSSODomainsInOtherTenants, data, "allowDuplicateSSODomainsInOtherTenants")
	boolattr.Get(m.AllowOverrideRoles, data, "allowOverrideRoles")
	boolattr.Get(m.GroupsPriority, data, "groupPriorityEnabled")
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	objattr.Set(&m.SSOSuiteSettings, data, helpers.RootKey, h)
	boolattr.Set(&m.AllowDuplicateSSODomainsInOtherTenants, data, "allowDuplicateSSODomainsInOtherTenants")
	boolattr.Set(&m.AllowOverrideRoles, data, "allowOverrideRoles")
	boolattr.Set(&m.GroupsPriority, data, "groupPriorityEnabled")
}

// SSO Suite Settings

var SSOSuiteValidator = objattr.NewValidator[SSOSuiteModel]("must have a valid configuration")

var SSOSuiteAttributes = map[string]schema.Attribute{
	"style_id":                  stringattr.Default(""),
	"hide_scim":                 boolattr.Default(false),
	"hide_groups_mapping":       boolattr.Default(false),
	"hide_domains":              boolattr.Default(false),
	"hide_saml":                 boolattr.Default(false),
	"hide_oidc":                 boolattr.Default(false),
	"force_domain_verification": boolattr.Default(false),
}

type SSOSuiteModel struct {
	StyleID                 stringattr.Type `tfsdk:"style_id"`
	HideSCIM                boolattr.Type   `tfsdk:"hide_scim"`
	HideGroupsMapping       boolattr.Type   `tfsdk:"hide_groups_mapping"`
	HideDomains             boolattr.Type   `tfsdk:"hide_domains"`
	HideSAML                boolattr.Type   `tfsdk:"hide_saml"`
	HideOIDC                boolattr.Type   `tfsdk:"hide_oidc"`
	ForceDomainVerification boolattr.Type   `tfsdk:"force_domain_verification"`
}

var SSOSuiteDefault = &SSOSuiteModel{
	StyleID:                 stringattr.Value(""),
	HideSCIM:                boolattr.Value(false),
	HideGroupsMapping:       boolattr.Value(false),
	HideDomains:             boolattr.Value(false),
	HideSAML:                boolattr.Value(false),
	HideOIDC:                boolattr.Value(false),
	ForceDomainVerification: boolattr.Value(false),
}

func (m *SSOSuiteModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.StyleID, data, "ssoSuiteStyleId")
	boolattr.Get(m.HideSCIM, data, "hideSsoSuiteScim")
	boolattr.Get(m.HideGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Get(m.HideDomains, data, "hideSsoSuiteDomains")
	boolattr.Get(m.HideSAML, data, "hideSsoSuiteSaml")
	boolattr.Get(m.HideOIDC, data, "hideSsoSuiteOidc")
	boolattr.Get(m.ForceDomainVerification, data, "ssoSuiteForceDomainVerification")
	return data
}

func (m *SSOSuiteModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.StyleID, data, "ssoSuiteStyleId")
	boolattr.Set(&m.HideSCIM, data, "hideSsoSuiteScim")
	boolattr.Set(&m.HideGroupsMapping, data, "hideSsoSuiteGroupsMapping")
	boolattr.Set(&m.HideDomains, data, "hideSsoSuiteDomains")
	boolattr.Set(&m.HideSAML, data, "hideSsoSuiteSaml")
	boolattr.Set(&m.HideOIDC, data, "hideSsoSuiteOidc")
	boolattr.Set(&m.ForceDomainVerification, data, "ssoSuiteForceDomainVerification")
}

func (m *SSOSuiteModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.HideSAML, m.HideOIDC) {
		return
	} else if m.HideSAML.ValueBool() && m.HideOIDC.ValueBool() {
		h.Invalid("The attributes hide_oidc and hide_saml cannot both be true")
	}

	if helpers.HasUnknownValues(m.HideDomains, m.ForceDomainVerification) {
		return
	} else if m.HideDomains.ValueBool() && m.ForceDomainVerification.ValueBool() {
		h.Invalid("The attributes force_domain_verification and hide_domains cannot both be true")
	}
}
