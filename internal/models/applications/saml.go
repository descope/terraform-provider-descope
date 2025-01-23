package applications

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SAMLAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"login_page_url":            stringattr.Default(""),
	"dynamic_configuration":     objectattr.Optional(DynamicConfigurationAttributes),
	"manual_configuration":      objectattr.Optional(ManualConfigurationAttributes),
	"acs_allowed_callback_urls": strlistattr.Optional(),
	"subject_name_id_type":      stringattr.Default("", stringvalidator.OneOf("", "email", "phone")),
	"subject_name_id_format":    stringattr.Default("", stringvalidator.OneOf("", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent", "urn:oasis:names:tc:SAML:2.0:nameid-format:transient")),
	"default_relay_state":       stringattr.Default(""),
	"attribute_mapping":         listattr.Optional(AttributeMappingAttributes),
	"force_authentication":      boolattr.Default(false),
}

// Model

type SAMLModel struct {
	ID                     types.String               `tfsdk:"id"`
	Name                   types.String               `tfsdk:"name"`
	Description            types.String               `tfsdk:"description"`
	Logo                   types.String               `tfsdk:"logo"`
	Disabled               types.Bool                 `tfsdk:"disabled"`
	LoginPageURL           types.String               `tfsdk:"login_page_url"`
	DynamicConfiguration   *DynamicConfigurationModel `tfsdk:"dynamic_configuration"`
	ManualConfiguration    *ManualConfigurationModel  `tfsdk:"manual_configuration"`
	ACSAllowedCallbackURLs []string                   `tfsdk:"acs_allowed_callback_urls"`
	SubjectNameIDType      types.String               `tfsdk:"subject_name_id_type"`
	SubjectNameIDFormat    types.String               `tfsdk:"subject_name_id_format"`
	DefaultRelayState      types.String               `tfsdk:"default_relay_state"`
	AttributeMapping       []*AttributeMappingModel   `tfsdk:"attribute_mapping"`
	ForceAuthentication    types.Bool                 `tfsdk:"force_authentication"`
}

func (m *SAMLModel) Values(h *Handler) map[string]any {
	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	if m.DynamicConfiguration != nil {
		settings["useMetadataInfo"] = true
		maps.Copy(settings, m.DynamicConfiguration.Values(h))
	} else if m.ManualConfiguration != nil {
		settings["useMetadataInfo"] = false
		maps.Copy(settings, m.ManualConfiguration.Values(h))
	}
	stringattr.Get(m.SubjectNameIDType, settings, "subjectNameIdType")
	stringattr.Get(m.SubjectNameIDFormat, settings, "subjectNameIdFormat")
	stringattr.Get(m.DefaultRelayState, settings, "defaultRelayState")
	listattr.Get(m.AttributeMapping, settings, "attributeMapping", h)
	strlistattr.Get(m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")
	data["saml"] = settings
	return data
}

func (m *SAMLModel) SetValues(h *Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["saml"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		if useMetadataInfo, ok := settings["useMetadataInfo"].(bool); ok && useMetadataInfo {
			m.DynamicConfiguration = &DynamicConfigurationModel{}
			m.DynamicConfiguration.SetValues(h, settings)
		} else {
			m.ManualConfiguration = &ManualConfigurationModel{}
			m.ManualConfiguration.SetValues(h, settings)
		}
		stringattr.Set(&m.SubjectNameIDType, settings, "subjectNameIdType")
		stringattr.Set(&m.SubjectNameIDFormat, settings, "subjectNameIdFormat")
		stringattr.Set(&m.DefaultRelayState, settings, "defaultRelayState")
		m.AttributeMapping = []*AttributeMappingModel{}
		listattr.Set(&m.AttributeMapping, settings, "attributeMapping", h)
		m.ACSAllowedCallbackURLs = []string{}
		strlistattr.Set(&m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks", h)
	}
}

// Attribute Mapping

var AttributeMappingAttributes = map[string]schema.Attribute{
	"name":  stringattr.Required(),
	"value": stringattr.Required(),
}

type AttributeMappingModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (m *AttributeMappingModel) Values(h *Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Value, data, "value")
	return data
}

func (m *AttributeMappingModel) SetValues(h *Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Value, data, "value")
}

// Dynamic Configuration

var DynamicConfigurationAttributes = map[string]schema.Attribute{
	"metadata_url": stringattr.Required(),
}

type DynamicConfigurationModel struct {
	MetadataURL types.String `tfsdk:"metadata_url"`
}

func (m *DynamicConfigurationModel) Values(h *Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.MetadataURL, data, "metadataUrl")
	return data
}

func (m *DynamicConfigurationModel) SetValues(h *Handler, data map[string]any) {
	stringattr.Set(&m.MetadataURL, data, "metadataUrl")
}

// Manual Configuration

var ManualConfigurationAttributes = map[string]schema.Attribute{
	"acs_url":     stringattr.Required(),
	"entity_id":   stringattr.Required(),
	"certificate": stringattr.Required(),
}

type ManualConfigurationModel struct {
	ACSURL      types.String `tfsdk:"acs_url"`
	EntityID    types.String `tfsdk:"entity_id"`
	Certificate types.String `tfsdk:"certificate"`
}

func (m *ManualConfigurationModel) Values(h *Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ACSURL, data, "acsUrl")
	stringattr.Get(m.EntityID, data, "entityId")
	stringattr.Get(m.Certificate, data, "certificate")
	return data
}

func (m *ManualConfigurationModel) SetValues(h *Handler, data map[string]any) {
	stringattr.Set(&m.ACSURL, data, "acsUrl")
	stringattr.Set(&m.EntityID, data, "entityId")
	stringattr.Set(&m.Certificate, data, "certificate")
}
