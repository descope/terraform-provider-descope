package apps

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var SAMLAppSchema = schema.Schema{
	MarkdownDescription: "Manages a single SAML federated application where Descope acts as the identity provider.",
	Attributes:          SAMLAppAttributes,
}

var SAMLAppAttributes = map[string]schema.Attribute{
	"id":                  stringattr.Optional(stringplanmodifier.RequiresReplace(), appIDValidator),
	"project_id":          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"deletion_protection": boolattr.Tristate(),
	"name":                stringattr.Required(stringattr.StandardLenValidator),
	"description":         stringattr.Default("", stringvalidator.LengthAtMost(100000)),
	"logo":                stringattr.Default(""),
	"disabled":            boolattr.Default(false),

	"login_page_url":              stringattr.Optional(),
	"dynamic_configuration":       objattr.Default[DynamicConfigurationModel](nil, DynamicConfigurationAttributes),
	"manual_configuration":        objattr.Default[ManualConfigurationModel](nil, ManualConfigurationAttributes),
	"acs_allowed_callback_urls":   strsetattr.Default(),
	"subject_name_id_type":        stringattr.Default("", stringvalidator.OneOf("", "email", "phone")),
	"subject_name_id_format":      stringattr.Default("", stringvalidator.OneOf("", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent", "urn:oasis:names:tc:SAML:2.0:nameid-format:transient")),
	"default_relay_state":         stringattr.Default(""),
	"default_signature_algorithm": stringattr.Default("", stringvalidator.OneOf("", "sha256")),
	"attribute_mapping":           listattr.Default[AttributeMappingModel](AttributeMappingAttributes),
	"force_authentication":        boolattr.Default(false),
}

// Model

type SAMLAppModel struct {
	ID                 stringattr.Type `tfsdk:"id"`
	ProjectID          stringattr.Type `tfsdk:"project_id"`
	DeletionProtection boolattr.Type   `tfsdk:"deletion_protection"`
	Name               stringattr.Type `tfsdk:"name"`
	Description        stringattr.Type `tfsdk:"description"`
	Logo               stringattr.Type `tfsdk:"logo"`
	Disabled           boolattr.Type   `tfsdk:"disabled"`

	LoginPageURL              stringattr.Type                         `tfsdk:"login_page_url"`
	DynamicConfiguration      objattr.Type[DynamicConfigurationModel] `tfsdk:"dynamic_configuration"`
	ManualConfiguration       objattr.Type[ManualConfigurationModel]  `tfsdk:"manual_configuration"`
	ACSAllowedCallbackURLs    strsetattr.Type                         `tfsdk:"acs_allowed_callback_urls"`
	SubjectNameIDType         stringattr.Type                         `tfsdk:"subject_name_id_type"`
	SubjectNameIDFormat       stringattr.Type                         `tfsdk:"subject_name_id_format"`
	DefaultRelayState         stringattr.Type                         `tfsdk:"default_relay_state"`
	DefaultSignatureAlgorithm stringattr.Type                         `tfsdk:"default_signature_algorithm"`
	AttributeMapping          listattr.Type[AttributeMappingModel]    `tfsdk:"attribute_mapping"`
	ForceAuthentication       boolattr.Type                           `tfsdk:"force_authentication"`
}

func (m *SAMLAppModel) Values(h *helpers.Handler) map[string]any {
	data := sharedAppValues(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	stringattr.Get(m.LoginPageURL, data, "loginPageUrl")
	if m.DynamicConfiguration.IsSet() {
		data["useMetadataInfo"] = true
		objattr.Get(m.DynamicConfiguration, data, helpers.RootKey, h)
	} else if m.ManualConfiguration.IsSet() {
		data["useMetadataInfo"] = false
		objattr.Get(m.ManualConfiguration, data, helpers.RootKey, h)
	}
	strsetattr.Get(m.ACSAllowedCallbackURLs, data, "acsAllowedCallbacks", h)
	stringattr.Get(m.SubjectNameIDType, data, "subjectNameIdType")
	stringattr.Get(m.SubjectNameIDFormat, data, "subjectNameIdFormat")
	stringattr.Get(m.DefaultRelayState, data, "defaultRelayState")
	stringattr.Get(m.DefaultSignatureAlgorithm, data, "defaultSignatureAlgorithm")
	listattr.Get(m.AttributeMapping, data, "attributeMapping", h)
	boolattr.Get(m.ForceAuthentication, data, "forceAuthentication")
	return data
}

func (m *SAMLAppModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedAppValues(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["samlSettings"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		if useMetadataInfo, ok := settings["useMetadataInfo"].(bool); ok && useMetadataInfo {
			objattr.Set(&m.DynamicConfiguration, settings, helpers.RootKey, h)
		} else {
			objattr.Set(&m.ManualConfiguration, settings, helpers.RootKey, h)
		}
		strsetattr.Set(&m.ACSAllowedCallbackURLs, settings, "acsAllowedCallbacks", h)
		stringattr.Set(&m.SubjectNameIDType, settings, "subjectNameIdType")
		stringattr.Set(&m.SubjectNameIDFormat, settings, "subjectNameIdFormat")
		stringattr.Set(&m.DefaultRelayState, settings, "defaultRelayState")
		stringattr.Set(&m.DefaultSignatureAlgorithm, settings, "defaultSignatureAlgorithm")
		listattr.Set(&m.AttributeMapping, settings, "attributeMapping", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
	}
}

func (m *SAMLAppModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.DynamicConfiguration, m.ManualConfiguration) {
		return
	}
	if !m.DynamicConfiguration.IsSet() && !m.ManualConfiguration.IsSet() {
		h.Missing("Either the dynamic_configuration or manual_configuration attribute must be set")
	} else if m.DynamicConfiguration.IsSet() && m.ManualConfiguration.IsSet() {
		h.Conflict("Only one of the dynamic_configuration and manual_configuration attributes can be set")
	}
}

func (m *SAMLAppModel) DeletionProtectionDefault(_ context.Context) bool {
	return true
}

func (m *SAMLAppModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SAMLAppModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *SAMLAppModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

// Dynamic Configuration

var DynamicConfigurationAttributes = map[string]schema.Attribute{
	"metadata_url": stringattr.Required(),
}

type DynamicConfigurationModel struct {
	MetadataURL stringattr.Type `tfsdk:"metadata_url"`
}

func (m *DynamicConfigurationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.MetadataURL, data, "metadataUrl")
	return data
}

func (m *DynamicConfigurationModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.MetadataURL, data, "metadataUrl")
}

// Manual Configuration

var ManualConfigurationAttributes = map[string]schema.Attribute{
	"acs_url":     stringattr.Required(),
	"entity_id":   stringattr.Required(),
	"certificate": stringattr.Default(""),
}

type ManualConfigurationModel struct {
	ACSURL      stringattr.Type `tfsdk:"acs_url"`
	EntityID    stringattr.Type `tfsdk:"entity_id"`
	Certificate stringattr.Type `tfsdk:"certificate"`
}

func (m *ManualConfigurationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ACSURL, data, "acsUrl")
	stringattr.Get(m.EntityID, data, "entityId")
	stringattr.Get(m.Certificate, data, "certificate")
	return data
}

func (m *ManualConfigurationModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ACSURL, data, "acsUrl")
	stringattr.Set(&m.EntityID, data, "entityId")
	stringattr.Set(&m.Certificate, data, "certificate")
}
