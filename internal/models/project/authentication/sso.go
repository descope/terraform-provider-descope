package authentication

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SSOAttributes = map[string]schema.Attribute{
	"disabled":                              boolattr.Default(false),
	"merge_users":                           boolattr.Default(false),
	"redirect_url":                          stringattr.Default(""),
	"sso_suite_settings":                    objattr.Default(SSOSuiteDefault, SSOSuiteAttributes, SSOSuiteValidator),
	"allow_duplicate_domains":               boolattr.Default(false),
	"allow_override_roles":                  boolattr.Default(false),
	"groups_priority":                       boolattr.Default(false),
	"mandatory_user_attributes":             listattr.Default[MandatoryUserAttributeModel](MandatoryUserAttributeAttributes),
	"limit_mapping_to_mandatory_attributes": boolattr.Default(false),
	"sso_domains_required":                  boolattr.Default(false),
	"groups_attribute_name_required":        boolattr.Default(false),
}

var customAttributePrefix = "customAttributes."
var groupsAttributeName = "group"
var ssoDomainsAttributeName = "ssoDomains"

type SSOModel struct {
	Disabled                               boolattr.Type                              `tfsdk:"disabled"`
	MergeUsers                             boolattr.Type                              `tfsdk:"merge_users"`
	RedirectURL                            stringattr.Type                            `tfsdk:"redirect_url"`
	SSOSuiteSettings                       objattr.Type[SSOSuiteModel]                `tfsdk:"sso_suite_settings"`
	AllowDuplicateSSODomainsInOtherTenants boolattr.Type                              `tfsdk:"allow_duplicate_domains"`
	AllowOverrideRoles                     boolattr.Type                              `tfsdk:"allow_override_roles"`
	GroupsPriority                         boolattr.Type                              `tfsdk:"groups_priority"`
	MandatoryUserAttributes                listattr.Type[MandatoryUserAttributeModel] `tfsdk:"mandatory_user_attributes"`
	LimitMappingToMandatoryAttributes      boolattr.Type                              `tfsdk:"limit_mapping_to_mandatory_attributes"`
	SSODomainsRequired                     boolattr.Type                              `tfsdk:"sso_domains_required"`
	GroupsAttributeNameRequired            boolattr.Type                              `tfsdk:"groups_attribute_name_required"`
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	boolattr.Get(m.AllowDuplicateSSODomainsInOtherTenants, data, "allowDuplicateSSODomainsInOtherTenants")
	boolattr.Get(m.GroupsPriority, data, "groupPriorityEnabled")
	boolattr.Get(m.AllowOverrideRoles, data, "allowOverrideRoles")
	boolattr.Get(m.LimitMappingToMandatoryAttributes, data, "limitMappingToMandatoryAttributes")

	getMandatoryUserAttributesValues(m.MandatoryUserAttributes, m.SSODomainsRequired, m.GroupsAttributeNameRequired, h, data)

	objattr.Get(m.SSOSuiteSettings, data, helpers.RootKey, h)
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	boolattr.Set(&m.AllowDuplicateSSODomainsInOtherTenants, data, "allowDuplicateSSODomainsInOtherTenants")
	boolattr.Set(&m.GroupsPriority, data, "groupPriorityEnabled")
	boolattr.Set(&m.AllowOverrideRoles, data, "allowOverrideRoles")

	boolattr.Set(&m.LimitMappingToMandatoryAttributes, data, "limitMappingToMandatoryAttributes")

	setMandatoryUserAttributesValues(m.MandatoryUserAttributes, m.SSODomainsRequired, m.GroupsAttributeNameRequired, h, data)

	objattr.Set(&m.SSOSuiteSettings, data, helpers.RootKey, h)
}

// User Attribute

type MandatoryUserAttributeModel struct {
	ID     stringattr.Type `tfsdk:"id"`
	Custom boolattr.Type   `tfsdk:"custom"`
}

var MandatoryUserAttributeAttributes = map[string]schema.Attribute{
	"id":     stringattr.Required(),
	"custom": boolattr.Default(false),
}

func (m *MandatoryUserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "value")
	boolattr.Get(m.Custom, data, "custom")
	return data
}

func (m *MandatoryUserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "value")
	boolattr.Set(&m.Custom, data, "custom")
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

// mandatoryUserAttributes field includes user attributes as strings, custom attributes are prefixed with "customAttributes." and "ssoDomains" and "group" are special attributes.
func setMandatoryUserAttributesValues(mandatoryUserAttributes listattr.Type[MandatoryUserAttributeModel], ssoDomainsRequired boolattr.Type, groupsAttributeNameRequired boolattr.Type, h *helpers.Handler, data map[string]any) {
	attributes, ok := data["mandatoryUserAttributes"]
	if !ok {
		return
	}
	domainsRequired := false
	groupsRequired := false
	mandatoryUserAttributesData := []map[string]any{}

	for _, attribute := range attributes.([]interface{}) {
		attributeStr, ok := attribute.(string)
		if !ok {
			continue
		}
		if attribute == ssoDomainsAttributeName {
			domainsRequired = true
		} else if attribute == groupsAttributeName {
			groupsRequired = true
		} else {
			if strings.HasPrefix(attributeStr, customAttributePrefix) {
				mandatoryUserAttributesData = append(mandatoryUserAttributesData, map[string]any{
					"value":  strings.TrimPrefix(attributeStr, customAttributePrefix),
					"custom": true,
				})
			} else {
				mandatoryUserAttributesData = append(mandatoryUserAttributesData, map[string]any{
					"value":  attribute,
					"custom": false,
				})
			}
		}
	}

	tempData := map[string]any{
		"domainsRequired":         domainsRequired,
		"groupsRequired":          groupsRequired,
		"mandatoryUserAttributes": mandatoryUserAttributesData,
	}
	listattr.Set(&mandatoryUserAttributes, tempData, "mandatoryUserAttributes", h)
	boolattr.Set(&ssoDomainsRequired, tempData, "domainsRequired")
	boolattr.Set(&groupsAttributeNameRequired, tempData, "groupsRequired")
}

func getMandatoryUserAttributesValues(mandatoryUserAttributes listattr.Type[MandatoryUserAttributeModel], ssoDomainsRequired boolattr.Type, groupsAttributeNameRequired boolattr.Type, h *helpers.Handler, data map[string]any) {
	tempData := map[string]any{}
	listattr.Get(mandatoryUserAttributes, tempData, "mandatoryUserAttributes", h)
	boolattr.Get(ssoDomainsRequired, tempData, "domainsRequired")
	boolattr.Get(groupsAttributeNameRequired, tempData, "groupsRequired")

	attributes := []string{}
	for _, attribute := range tempData["mandatoryUserAttributes"].([]any) {
		attributeData, ok := attribute.(map[string]any)
		if !ok {
			continue
		}
		if attributeData["custom"].(bool) {
			attributes = append(attributes, customAttributePrefix+attributeData["value"].(string))
		} else {
			attributes = append(attributes, attributeData["value"].(string))
		}
	}
	if tempData["domainsRequired"].(bool) {
		attributes = append(attributes, ssoDomainsAttributeName)
	}
	if tempData["groupsRequired"].(bool) {
		attributes = append(attributes, groupsAttributeName)
	}

	data["mandatoryUserAttributes"] = attributes
}
