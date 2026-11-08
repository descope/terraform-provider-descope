package apps

import (
	"context"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
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
	"logout_redirect_url":         stringattr.Default(""),
	"error_redirect_url":          stringattr.Default(""),
	"idp_logout_url":              stringattr.Generated(),
	"single_logout_enabled":       boolattr.Default(false),
	"sp_single_logout_url":        stringattr.Default(""),
	"sp_single_logout_binding":    stringattr.Default("", stringvalidator.OneOf("", "redirect", "post", "soap")),
	"single_logout_propagation":   stringattr.Default("", stringvalidator.OneOf("", "allApps", "none")),
	"single_logout_scope":         stringattr.Default("", stringvalidator.OneOf("", "session", "allSessions")),
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
	LogoutRedirectURL         stringattr.Type                         `tfsdk:"logout_redirect_url"`
	ErrorRedirectURL          stringattr.Type                         `tfsdk:"error_redirect_url"`
	IDPLogoutURL              stringattr.Type                         `tfsdk:"idp_logout_url"`
	SingleLogoutEnabled       boolattr.Type                           `tfsdk:"single_logout_enabled"`
	SPSingleLogoutURL         stringattr.Type                         `tfsdk:"sp_single_logout_url"`
	SPSingleLogoutBinding     stringattr.Type                         `tfsdk:"sp_single_logout_binding"`
	SingleLogoutPropagation   stringattr.Type                         `tfsdk:"single_logout_propagation"`
	SingleLogoutScope         stringattr.Type                         `tfsdk:"single_logout_scope"`
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
	stringattr.Get(m.LogoutRedirectURL, data, "logoutRedirectUrl")
	stringattr.Get(m.ErrorRedirectURL, data, "errorRedirectUrl")
	boolattr.Get(m.SingleLogoutEnabled, data, "singleLogoutEnabled")
	stringattr.Get(m.SPSingleLogoutURL, data, "spSingleLogoutUrl")
	stringattr.Get(m.SPSingleLogoutBinding, data, "spSingleLogoutBinding")
	stringattr.Get(m.SingleLogoutPropagation, data, "singleLogoutPropagation")
	stringattr.Get(m.SingleLogoutScope, data, "singleLogoutScope")
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
		stringattr.Set(&m.LogoutRedirectURL, settings, "logoutRedirectUrl")
		stringattr.Set(&m.ErrorRedirectURL, settings, "errorRedirectUrl")
		stringattr.Set(&m.IDPLogoutURL, settings, "idpLogoutUrl")
		boolattr.Set(&m.SingleLogoutEnabled, settings, "singleLogoutEnabled")
		stringattr.Set(&m.SPSingleLogoutURL, settings, "spSingleLogoutUrl")
		stringattr.Set(&m.SPSingleLogoutBinding, settings, "spSingleLogoutBinding")
		stringattr.Set(&m.SingleLogoutPropagation, settings, "singleLogoutPropagation")
		stringattr.Set(&m.SingleLogoutScope, settings, "singleLogoutScope")
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

	if helpers.HasUnknownValues(m.SingleLogoutEnabled, m.SPSingleLogoutURL, m.SPSingleLogoutBinding) || !m.SingleLogoutEnabled.ValueBool() {
		return
	}
	if url := m.SPSingleLogoutURL.ValueString(); url == "" {
		// with dynamic_configuration the URL can come from the service provider's metadata instead
		if m.ManualConfiguration.IsSet() {
			h.Missing("The sp_single_logout_url attribute must be set when single_logout_enabled is true and manual_configuration is used")
		}
	} else if m.SPSingleLogoutBinding.ValueString() == "soap" && !strings.HasPrefix(strings.ToLower(url), "https://") {
		h.Invalid("The sp_single_logout_url attribute must use https when sp_single_logout_binding is soap")
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
