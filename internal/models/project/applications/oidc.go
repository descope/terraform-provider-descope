package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var OIDCAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"login_page_url":       stringattr.Default(""),
	"claims":               strlistattr.Default(),
	"force_authentication": boolattr.Default(false),

	// Dedicated client credentials and per-app policy (config-driven; defaults preserve legacy behavior).
	// client_id/client_secret may import an existing OIDC client; set on create only (RequiresReplace),
	// computed/generated server-side when omitted. Mirrors the inbound third-party app attributes.
	"client_id":     stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"client_secret": stringattr.SecretGenerated(true),
	"client_type":   stringattr.Default("", stringvalidator.OneOf("", "confidential", "public")), // "", "confidential", or "public"
	// Set (not list): the backend may reorder the URLs, and a list would show perpetual plan diffs -
	// matching approved_callback_urls (inbound) and acs_allowed_callback_urls (SAML).
	"approved_redirect_urls": strsetattr.Default(),
	// Per-app modular grant types (disabled-polarity; all false = all grant types enabled = legacy).
	"authorization_code_disabled": boolattr.Default(false),
	"client_credentials_disabled": boolattr.Default(false),
	"refresh_token_disabled":      boolattr.Default(false),
	"jwt_bearer_disabled":         boolattr.Default(false),
	"device_code_disabled":        boolattr.Default(false),
	"force_pkce":                  boolattr.Default(false),
	// Default audience policy for issued tokens (modern apps only): "projectId", "clientId", or "" (both).
	// Legacy apps (empty client_type) always use the project ID; the empty default preserves that.
	"default_audience": stringattr.Default("", stringvalidator.OneOf("", "projectId", "clientId")),

	"permissions": listattr.Default[SSOAppPermissionModel](SSOAppPermissionAttributes),
	"roles":       listattr.Default[SSOAppRoleModel](SSOAppRoleAttributes),
}

// Model

type OIDCModel struct {
	ID                  stringattr.Type  `tfsdk:"id"`
	Name                stringattr.Type  `tfsdk:"name"`
	Description         stringattr.Type  `tfsdk:"description"`
	Logo                stringattr.Type  `tfsdk:"logo"`
	Disabled            boolattr.Type    `tfsdk:"disabled"`
	LoginPageURL        stringattr.Type  `tfsdk:"login_page_url"`
	Claims              strlistattr.Type `tfsdk:"claims"`
	ForceAuthentication boolattr.Type    `tfsdk:"force_authentication"`

	ClientID                  stringattr.Type `tfsdk:"client_id"`
	ClientSecret              stringattr.Type `tfsdk:"client_secret"`
	ClientType                stringattr.Type `tfsdk:"client_type"`
	ApprovedRedirectURLs      strsetattr.Type `tfsdk:"approved_redirect_urls"`
	AuthorizationCodeDisabled boolattr.Type   `tfsdk:"authorization_code_disabled"`
	ClientCredentialsDisabled boolattr.Type   `tfsdk:"client_credentials_disabled"`
	RefreshTokenDisabled      boolattr.Type   `tfsdk:"refresh_token_disabled"`
	JWTBearerDisabled         boolattr.Type   `tfsdk:"jwt_bearer_disabled"`
	DeviceCodeDisabled        boolattr.Type   `tfsdk:"device_code_disabled"`
	ForcePkce                 boolattr.Type   `tfsdk:"force_pkce"`
	DefaultAudience           stringattr.Type `tfsdk:"default_audience"`

	Permissions listattr.Type[SSOAppPermissionModel] `tfsdk:"permissions"`
	Roles       listattr.Type[SSOAppRoleModel]       `tfsdk:"roles"`
}

func (m *OIDCModel) Values(h *helpers.Handler) map[string]any {
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	strlistattr.Get(m.Claims, settings, "claims", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")

	stringattr.Get(m.ClientID, settings, "clientId")
	stringattr.Get(m.ClientSecret, settings, "clientSecret")
	stringattr.Get(m.ClientType, settings, "clientType")
	strsetattr.Get(m.ApprovedRedirectURLs, settings, "approvedRedirectUrls", h)
	boolattr.Get(m.AuthorizationCodeDisabled, settings, "authorizationCodeDisabled")
	boolattr.Get(m.ClientCredentialsDisabled, settings, "clientCredentialsDisabled")
	boolattr.Get(m.RefreshTokenDisabled, settings, "refreshTokenDisabled")
	boolattr.Get(m.JWTBearerDisabled, settings, "jwtBearerDisabled")
	boolattr.Get(m.DeviceCodeDisabled, settings, "deviceCodeDisabled")
	boolattr.Get(m.ForcePkce, settings, "forcePkce")
	stringattr.Get(m.DefaultAudience, settings, "defaultAudience")

	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	data["oidc"] = settings
	emitSSOAppRoles(h, data, m.Permissions, m.Roles)
	return data
}

func (m *OIDCModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["oidc"].(map[string]any); ok {
		stringattr.Nil(&m.LoginPageURL) // XXX reset by the backend on response for now
		strlistattr.Set(&m.Claims, settings, "claims", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")

		stringattr.Set(&m.ClientID, settings, "clientId")
		stringattr.Set(&m.ClientSecret, settings, "clientSecret")
		stringattr.Set(&m.ClientType, settings, "clientType")
		strsetattr.Set(&m.ApprovedRedirectURLs, settings, "approvedRedirectUrls", h)
		boolattr.Set(&m.AuthorizationCodeDisabled, settings, "authorizationCodeDisabled")
		boolattr.Set(&m.ClientCredentialsDisabled, settings, "clientCredentialsDisabled")
		boolattr.Set(&m.RefreshTokenDisabled, settings, "refreshTokenDisabled")
		boolattr.Set(&m.JWTBearerDisabled, settings, "jwtBearerDisabled")
		boolattr.Set(&m.DeviceCodeDisabled, settings, "deviceCodeDisabled")
		boolattr.Set(&m.ForcePkce, settings, "forcePkce")
		stringattr.Set(&m.DefaultAudience, settings, "defaultAudience")
	}
	listattr.SetMatchingNames(&m.Permissions, data, "permissions", "name", h)
	listattr.SetMatchingNames(&m.Roles, data, "roles", "name", h)
}

// Matching

func (m *OIDCModel) GetName() stringattr.Type {
	return m.Name
}

func (m *OIDCModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OIDCModel) SetID(id stringattr.Type) {
	m.ID = id
}
