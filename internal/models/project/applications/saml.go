package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SAMLAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"login_page_url":            stringattr.Default(""),
	"dynamic_configuration":     objattr.Default[DynamicConfigurationModel](nil, DynamicConfigurationAttributes),
	"manual_configuration":      objattr.Default[ManualConfigurationModel](nil, ManualConfigurationAttributes),
	"acs_allowed_callback_urls": strsetattr.Default(),
	"subject_name_id_type":      stringattr.Default("", stringvalidator.OneOf("", "email", "phone")),
	"subject_name_id_format":    stringattr.Default("", stringvalidator.OneOf("", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent", "urn:oasis:names:tc:SAML:2.0:nameid-format:transient")),
	"default_relay_state":       stringattr.Default(""),
	"attribute_mapping":         listattr.Default[AttributeMappingModel](AttributeMappingAttributes),
	"force_authentication":      boolattr.Default(false),
}

// Model

type SAMLModel struct {
	ID                     stringattr.Type                         `tfsdk:"id"`
	Name                   stringattr.Type                         `tfsdk:"name"`
	Description            stringattr.Type                         `tfsdk:"description"`
	Logo                   stringattr.Type                         `tfsdk:"logo"`
	Disabled               boolattr.Type                           `tfsdk:"disabled"`
	LoginPageURL           stringattr.Type                         `tfsdk:"login_page_url"`
	DynamicConfiguration   objattr.Type[DynamicConfigurationModel] `tfsdk:"dynamic_configuration"`
	ManualConfiguration    objattr.Type[ManualConfigurationModel]  `tfsdk:"manual_configuration"`
	ACSAllowedCallbackURLs strsetattr.Type                         `tfsdk:"acs_allowed_callback_urls"`
	SubjectNameIDType      stringattr.Type                         `tfsdk:"subject_name_id_type"`
	SubjectNameIDFormat    stringattr.Type                         `tfsdk:"subject_name_id_format"`
	DefaultRelayState      stringattr.Type                         `tfsdk:"default_relay_state"`
	AttributeMapping       listattr.Type[AttributeMappingModel]    `tfsdk:"attribute_mapping"`
	ForceAuthentication    boolattr.Type                           `tfsdk:"force_authentication"`
}

func (m *SAMLModel) Values(h *Handler) map[string]any {
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	if m.DynamicConfiguration.IsSet() {
		settings["useMetadataInfo"] = true
		objattr.Get(m.DynamicConfiguration, settings, helpers.RootKey, h)
	} else if m.ManualConfiguration.IsSet() {
		settings["useMetadataInfo"] = false
		objattr.Get(m.ManualConfiguration, settings, helpers.RootKey, h)
	}
	stringattr.Get(m.SubjectNameIDType, settings, "subjectNameIdType")
	stringattr.Get(m.SubjectNameIDFormat, settings, "subjectNameIdFormat")
	stringattr.Get(m.DefaultRelayState, settings, "defaultRelayState")
	listattr.Get(m.AttributeMapping, settings, "attributeMapping", h)
	strsetattr.Get(m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")

	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	data["saml"] = settings
	return data
}

func (m *SAMLModel) SetValues(h *Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["saml"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		if useMetadataInfo, ok := settings["useMetadataInfo"].(bool); ok && useMetadataInfo {
			objattr.Set(&m.DynamicConfiguration, settings, helpers.RootKey, h)
		} else {
			objattr.Set(&m.ManualConfiguration, settings, helpers.RootKey, h)
		}
		stringattr.Set(&m.SubjectNameIDType, settings, "subjectNameIdType")
		stringattr.Set(&m.SubjectNameIDFormat, settings, "subjectNameIdFormat")
		stringattr.Set(&m.DefaultRelayState, settings, "defaultRelayState")
		listattr.Set(&m.AttributeMapping, settings, "attributeMapping", h)
		strsetattr.Set(&m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
	}
}

// Attribute Mapping

var AttributeMappingAttributes = map[string]schema.Attribute{
	"name":  stringattr.Required(),
	"value": stringattr.Required(),
}

type AttributeMappingModel struct {
	Name  stringattr.Type `tfsdk:"name"`
	Value stringattr.Type `tfsdk:"value"`
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
	MetadataURL stringattr.Type `tfsdk:"metadata_url"`
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
	ACSURL      stringattr.Type `tfsdk:"acs_url"`
	EntityID    stringattr.Type `tfsdk:"entity_id"`
	Certificate stringattr.Type `tfsdk:"certificate"`
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
