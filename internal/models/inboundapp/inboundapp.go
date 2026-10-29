package inboundapp

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const clientTypeReplacementDescription = "Changing the client type replaces the app, except from the legacy empty value to confidential."

var clientTypeReplacement = stringplanmodifier.RequiresReplaceIf(isClientTypeReplacement, clientTypeReplacementDescription, clientTypeReplacementDescription)

var InboundAppAttributes = map[string]schema.Attribute{
	"id":                               stringattr.Identifier(),
	"project_id":                       stringattr.Required(stringplanmodifier.RequiresReplace()),
	"deletion_protection":              boolattr.Tristate(),
	"name":                             stringattr.Required(),
	"description":                      stringattr.Default(""),
	"logo_url":                         stringattr.Optional(),
	"login_page_url":                   stringattr.Optional(),
	"approved_callback_urls":           strsetattr.Default(),
	"permissions_scopes":               listattr.Default[ApplicationScopeModel](ApplicationScopeAttributes),
	"attributes_scopes":                listattr.Default[ApplicationScopeModel](ApplicationScopeAttributes),
	"connections_scopes":               listattr.Default[ApplicationScopeModel](ApplicationScopeAttributes),
	"session_settings":                 objattr.Optional[SessionSettingsModel](SessionSettingsAttributes),
	"audience_whitelist":               strsetattr.Default(),
	"force_add_all_authorization_info": boolattr.Default(false),
	"force_dpop":                       boolattr.Default(false),
	"default_audience":                 stringattr.Default("", stringvalidator.OneOf("", "projectId", "clientId")), // XXX maybe switch to set
	"client_type":                      stringattr.Optional(stringvalidator.OneOf("", "confidential", "public"), clientTypeReplacement),
	"client_id":                        stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"client_secret":                    stringattr.SecretGenerated(true),
	"force_pkce":                       boolattr.Default(false),
	"allowed_tenants":                  strsetattr.Default(),
	"scope_claim_mapping":              listattr.Default[ScopeClaimMappingModel](ScopeClaimMappingAttributes),
	"authorization_code_disabled":      boolattr.Default(false),
	"client_credentials_disabled":      boolattr.Default(false),
	"refresh_token_disabled":           boolattr.Default(false),
	"jwt_bearer_disabled":              boolattr.Default(false),
	"device_code_disabled":             boolattr.Default(false),
	"token_exchange_disabled":          boolattr.Default(false),
}

var Schema = schema.Schema{
	Attributes: InboundAppAttributes,
}

type InboundAppModel struct {
	ID                           stringattr.Type                       `tfsdk:"id"`
	ProjectID                    stringattr.Type                       `tfsdk:"project_id"`
	DeletionProtection           boolattr.Type                         `tfsdk:"deletion_protection"`
	Name                         stringattr.Type                       `tfsdk:"name"`
	Description                  stringattr.Type                       `tfsdk:"description"`
	LogoUrl                      stringattr.Type                       `tfsdk:"logo_url"`
	LoginPageUrl                 stringattr.Type                       `tfsdk:"login_page_url"`
	ApprovedCallbackUrls         strsetattr.Type                       `tfsdk:"approved_callback_urls"`
	PermissionsScopes            listattr.Type[ApplicationScopeModel]  `tfsdk:"permissions_scopes"`
	AttributesScopes             listattr.Type[ApplicationScopeModel]  `tfsdk:"attributes_scopes"`
	ConnectionsScopes            listattr.Type[ApplicationScopeModel]  `tfsdk:"connections_scopes"`
	SessionSettings              objattr.Type[SessionSettingsModel]    `tfsdk:"session_settings"`
	AudienceWhitelist            strsetattr.Type                       `tfsdk:"audience_whitelist"`
	ForceAddAllAuthorizationInfo boolattr.Type                         `tfsdk:"force_add_all_authorization_info"`
	ForceDpop                    boolattr.Type                         `tfsdk:"force_dpop"`
	DefaultAudience              stringattr.Type                       `tfsdk:"default_audience"`
	ClientType                   stringattr.Type                       `tfsdk:"client_type"`
	ClientId                     stringattr.Type                       `tfsdk:"client_id"`
	ClientSecret                 stringattr.Type                       `tfsdk:"client_secret"`
	ForcePkce                    boolattr.Type                         `tfsdk:"force_pkce"`
	AllowedTenants               strsetattr.Type                       `tfsdk:"allowed_tenants"`
	ScopeClaimMapping            listattr.Type[ScopeClaimMappingModel] `tfsdk:"scope_claim_mapping"`
	AuthorizationCodeDisabled    boolattr.Type                         `tfsdk:"authorization_code_disabled"`
	ClientCredentialsDisabled    boolattr.Type                         `tfsdk:"client_credentials_disabled"`
	RefreshTokenDisabled         boolattr.Type                         `tfsdk:"refresh_token_disabled"`
	JWTBearerDisabled            boolattr.Type                         `tfsdk:"jwt_bearer_disabled"`
	DeviceCodeDisabled           boolattr.Type                         `tfsdk:"device_code_disabled"`
	TokenExchangeDisabled        boolattr.Type                         `tfsdk:"token_exchange_disabled"`
}

func (m *InboundAppModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.LogoUrl, data, "logoUrl")
	stringattr.Get(m.LoginPageUrl, data, "loginPageUrl")
	strsetattr.Get(m.ApprovedCallbackUrls, data, "approvedCallbackUrls", h)
	listattr.Get(m.PermissionsScopes, data, "permissionsScopes", h)
	listattr.Get(m.AttributesScopes, data, "attributesScopes", h)
	listattr.Get(m.ConnectionsScopes, data, "connectionsScopes", h)
	objattr.Get(m.SessionSettings, data, "sessionSettings", h)
	strsetattr.Get(m.AudienceWhitelist, data, "audienceWhitelist", h)
	boolattr.Get(m.ForceAddAllAuthorizationInfo, data, "forceAddAllAuthorizationInfo")
	boolattr.Get(m.ForceDpop, data, "forceDpop")
	stringattr.Get(m.DefaultAudience, data, "defaultAudience")
	stringattr.Get(m.ClientType, data, "clientType")
	stringattr.Get(m.ClientId, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	boolattr.Get(m.ForcePkce, data, "forcePkce")
	strsetattr.Get(m.AllowedTenants, data, "allowedTenants", h)
	listattr.Get(m.ScopeClaimMapping, data, "scopeClaimMapping", h)
	boolattr.Get(m.AuthorizationCodeDisabled, data, "authorizationCodeDisabled")
	boolattr.Get(m.ClientCredentialsDisabled, data, "clientCredentialsDisabled")
	boolattr.Get(m.RefreshTokenDisabled, data, "refreshTokenDisabled")
	boolattr.Get(m.JWTBearerDisabled, data, "jwtBearerDisabled")
	boolattr.Get(m.DeviceCodeDisabled, data, "deviceCodeDisabled")
	boolattr.Get(m.TokenExchangeDisabled, data, "tokenExchangeDisabled")
	return data
}

func (m *InboundAppModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.LogoUrl, data, "logoUrl", stringattr.SkipIfAlreadySet)
	stringattr.Set(&m.LoginPageUrl, data, "loginPageUrl", stringattr.SkipIfAlreadySet)
	strsetattr.Set(&m.ApprovedCallbackUrls, data, "approvedCallbackUrls", h)
	listattr.Set(&m.PermissionsScopes, data, "permissionsScopes", h)
	listattr.Set(&m.AttributesScopes, data, "attributesScopes", h)
	listattr.Set(&m.ConnectionsScopes, data, "connectionsScopes", h)
	objattr.Set(&m.SessionSettings, data, "sessionSettings", h)
	strsetattr.Set(&m.AudienceWhitelist, data, "audienceWhitelist", h)
	boolattr.Set(&m.ForceAddAllAuthorizationInfo, data, "forceAddAllAuthorizationInfo")
	boolattr.Set(&m.ForceDpop, data, "forceDpop")
	stringattr.Set(&m.DefaultAudience, data, "defaultAudience")
	stringattr.SetDefault(&m.ClientType, data, "clientType", "") // omitted by the backend for the legacy type, and Set would leave a null on import
	stringattr.Set(&m.ClientId, data, "clientId")
	stringattr.Set(&m.ClientSecret, data, "clientSecret")
	boolattr.Set(&m.ForcePkce, data, "forcePkce")
	strsetattr.Set(&m.AllowedTenants, data, "allowedTenants", h)
	listattr.Set(&m.ScopeClaimMapping, data, "scopeClaimMapping", h)
	boolattr.Set(&m.AuthorizationCodeDisabled, data, "authorizationCodeDisabled")
	boolattr.Set(&m.ClientCredentialsDisabled, data, "clientCredentialsDisabled")
	boolattr.Set(&m.RefreshTokenDisabled, data, "refreshTokenDisabled")
	boolattr.Set(&m.JWTBearerDisabled, data, "jwtBearerDisabled")
	boolattr.Set(&m.DeviceCodeDisabled, data, "deviceCodeDisabled")
	boolattr.Set(&m.TokenExchangeDisabled, data, "tokenExchangeDisabled")
}

// The backend silently discards a client_secret given for a public client, so without this the apply would fail with an inconsistent-result error.
func (m *InboundAppModel) Validate(h *helpers.Handler) {
	if m.ClientSecret.ValueString() != "" && m.ClientType.ValueString() == "public" {
		h.Conflict("The client_secret field cannot be used when client_type is public, as public clients do not have a secret")
	}
	if m.ClientType.ValueString() == "public" {
		validatePublicGrantTypes(h, m)
	}
}

// The backend always disables the client_credentials and jwt-bearer grant types for a public client, so they default to true for one.
func (m *InboundAppModel) ModifyPlan(h *helpers.Handler, config, _ *InboundAppModel) {
	var value types.Bool
	switch {
	case m.ClientType.IsUnknown():
		value = types.BoolUnknown() // the client type is chosen by the server or only known after apply, and might be public
	case m.ClientType.ValueString() == "public":
		validatePublicGrantTypes(h, config) // Validate misses a public client type that is kept from the state
		value = types.BoolValue(true)
	default:
		return
	}
	if config.ClientCredentialsDisabled.IsNull() {
		m.ClientCredentialsDisabled = value
	}
	if config.JWTBearerDisabled.IsNull() {
		m.JWTBearerDisabled = value
	}
}

func validatePublicGrantTypes(h *helpers.Handler, config *InboundAppModel) {
	if isFalse(config.ClientCredentialsDisabled) {
		h.Conflict("The client_credentials_disabled field cannot be false when client_type is public, as public clients cannot use the client_credentials grant type")
	}
	if isFalse(config.JWTBearerDisabled) {
		h.Conflict("The jwt_bearer_disabled field cannot be false when client_type is public, as public clients cannot use the jwt-bearer grant type")
	}
}

func isFalse(b boolattr.Type) bool {
	return !b.IsNull() && !b.IsUnknown() && !b.ValueBool()
}

func isClientTypeReplacement(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
	state, plan := req.StateValue.ValueString(), req.PlanValue.ValueString()
	resp.RequiresReplace = req.PlanValue.IsUnknown() || state == "public" || plan == "public" || (state != "" && plan == "")
}

// Protected by default because replacing an inbound app issues a new client ID and secret.
func (m *InboundAppModel) DeletionProtectionDefault(_ context.Context) bool {
	return true
}

func (m *InboundAppModel) GetID() stringattr.Type {
	return m.ID
}

func (m *InboundAppModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *InboundAppModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
