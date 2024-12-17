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
	strlistattr.Get(m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks")
	data["saml"] = settings
	return data
}

func (m *SAMLModel) SetValues(h *Handler, data map[string]any) {
	// all saml application values are specified in the configuration
}

// Attribute Mapping

var AttributeMappingAttributes = map[string]schema.Attribute{
	"name":  stringattr.Required(),
	"value": stringattr.Required(),
}

type AttributeMappingModel struct {
	Name  string `tfsdk:"name"`
	Value string `tfsdk:"value"`
}

func (m *AttributeMappingModel) Values(h *Handler) map[string]any {
	return map[string]any{
		"name":  m.Name,
		"value": m.Value,
	}
}

func (m *AttributeMappingModel) SetValues(h *Handler, data map[string]any) {
	// all attribute mapping values are specified in the configuration
}

// Dynamic Configuration

var DynamicConfigurationAttributes = map[string]schema.Attribute{
	"metadata_url": stringattr.Required(),
}

type DynamicConfigurationModel struct {
	MetadataURL string `tfsdk:"metadata_url"`
}

func (m *DynamicConfigurationModel) Values(h *Handler) map[string]any {
	return map[string]any{
		"metadataUrl": m.MetadataURL,
	}
}

func (m *DynamicConfigurationModel) SetValues(h *Handler, data map[string]any) {
	// all dynamic configuration mapping values are specified in the configuration
}

// Manual Configuration

var ManualConfigurationAttributes = map[string]schema.Attribute{
	"acs_url":     stringattr.Required(),
	"entity_id":   stringattr.Required(),
	"certificate": stringattr.Required(),
}

type ManualConfigurationModel struct {
	ACSURL      string `tfsdk:"acs_url"`
	EntityID    string `tfsdk:"entity_id"`
	Certificate string `tfsdk:"certificate"`
}

func (m *ManualConfigurationModel) Values(h *Handler) map[string]any {
	return map[string]any{
		"acsUrl":      m.ACSURL,
		"entityId":    m.EntityID,
		"certificate": m.Certificate,
	}
}

func (m *ManualConfigurationModel) SetValues(h *Handler, data map[string]any) {
	// all manual configuration mapping values are specified in the configuration
}
