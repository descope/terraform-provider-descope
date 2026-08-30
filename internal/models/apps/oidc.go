package apps

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OIDCAppSchema = schema.Schema{
	MarkdownDescription: "Manages a single OIDC federated application where Descope acts as the identity provider.",
	Attributes:          OIDCAppAttributes,
}

var OIDCAppAttributes = map[string]schema.Attribute{
	"id":                  stringattr.Optional(stringplanmodifier.RequiresReplace(), appIDValidator),
	"project_id":          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"deletion_protection": boolattr.Tristate(),
	"name":                stringattr.Required(stringattr.StandardLenValidator),
	"description":         stringattr.Default("", stringvalidator.LengthAtMost(100000)),
	"logo":                stringattr.Default(""),
	"disabled":            boolattr.Default(false),

	"login_page_url":                      stringattr.Optional(),
	"claims":                              strlistattr.Default(),
	"force_authentication":                boolattr.Default(false),
	"backchannel_logout_url":              stringattr.Default(""),
	"custom_idp_initiated_login_page_url": stringattr.Default(""),
	"client_id":                           stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"client_secret":                       stringattr.SecretGenerated(true, stringplanmodifier.RequiresReplace()),
	"client_type":                         stringattr.Default("", stringvalidator.OneOf("", "confidential", "public")),
	"approved_redirect_urls":              strsetattr.Default(),
	"authorization_code_disabled":         boolattr.Default(false),
	"client_credentials_disabled":         boolattr.Default(false),
	"refresh_token_disabled":              boolattr.Default(false),
	"jwt_bearer_disabled":                 boolattr.Default(false),
	"device_code_disabled":                boolattr.Default(false),
	"force_pkce":                          boolattr.Default(false),
	"default_audience":                    stringattr.Default("", stringvalidator.OneOf("", "projectId", "clientId", "appId", "empty")),
	"trusted_apps_audience":               stringattr.Default("", stringvalidator.OneOf("", "projectId", "clientId", "appId", "empty")),
}

// Model

type OIDCAppModel struct {
	ID                 stringattr.Type `tfsdk:"id"`
	ProjectID          stringattr.Type `tfsdk:"project_id"`
	DeletionProtection boolattr.Type   `tfsdk:"deletion_protection"`
	Name               stringattr.Type `tfsdk:"name"`
	Description        stringattr.Type `tfsdk:"description"`
	Logo               stringattr.Type `tfsdk:"logo"`
	Disabled           boolattr.Type   `tfsdk:"disabled"`

	LoginPageURL                   stringattr.Type  `tfsdk:"login_page_url"`
	Claims                         strlistattr.Type `tfsdk:"claims"`
	ForceAuthentication            boolattr.Type    `tfsdk:"force_authentication"`
	BackchannelLogoutURL           stringattr.Type  `tfsdk:"backchannel_logout_url"`
	CustomIDPInitiatedLoginPageURL stringattr.Type  `tfsdk:"custom_idp_initiated_login_page_url"`
	ClientID                       stringattr.Type  `tfsdk:"client_id"`
	ClientSecret                   stringattr.Type  `tfsdk:"client_secret"`
	ClientType                     stringattr.Type  `tfsdk:"client_type"`
	ApprovedRedirectURLs           strsetattr.Type  `tfsdk:"approved_redirect_urls"`
	AuthorizationCodeDisabled      boolattr.Type    `tfsdk:"authorization_code_disabled"`
	ClientCredentialsDisabled      boolattr.Type    `tfsdk:"client_credentials_disabled"`
	RefreshTokenDisabled           boolattr.Type    `tfsdk:"refresh_token_disabled"`
	JWTBearerDisabled              boolattr.Type    `tfsdk:"jwt_bearer_disabled"`
	DeviceCodeDisabled             boolattr.Type    `tfsdk:"device_code_disabled"`
	ForcePkce                      boolattr.Type    `tfsdk:"force_pkce"`
	DefaultAudience                stringattr.Type  `tfsdk:"default_audience"`
	TrustedAppsAudience            stringattr.Type  `tfsdk:"trusted_apps_audience"`
}

func (m *OIDCAppModel) Values(h *helpers.Handler) map[string]any {
	data := sharedAppValues(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	stringattr.Get(m.LoginPageURL, data, "loginPageUrl")
	strlistattr.Get(m.Claims, data, "claims", h)
	boolattr.Get(m.ForceAuthentication, data, "forceAuthentication")
	stringattr.Get(m.BackchannelLogoutURL, data, "backChannelLogoutUrl")
	stringattr.Get(m.CustomIDPInitiatedLoginPageURL, data, "customIdpInitiatedLoginPageUrl")

	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	stringattr.Get(m.ClientType, data, "clientType")
	strsetattr.Get(m.ApprovedRedirectURLs, data, "approvedRedirectUrls", h)
	boolattr.Get(m.AuthorizationCodeDisabled, data, "authorizationCodeDisabled")
	boolattr.Get(m.ClientCredentialsDisabled, data, "clientCredentialsDisabled")
	boolattr.Get(m.RefreshTokenDisabled, data, "refreshTokenDisabled")
	boolattr.Get(m.JWTBearerDisabled, data, "jwtBearerDisabled")
	boolattr.Get(m.DeviceCodeDisabled, data, "deviceCodeDisabled")
	boolattr.Get(m.ForcePkce, data, "forcePkce")
	stringattr.Get(m.DefaultAudience, data, "defaultAudience")
	stringattr.Get(m.TrustedAppsAudience, data, "trustedAppsAudience")
	return data
}

func (m *OIDCAppModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedAppValues(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["oidcSettings"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		strlistattr.Set(&m.Claims, settings, "claims", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
		stringattr.Set(&m.BackchannelLogoutURL, settings, "backChannelLogoutUrl")
		stringattr.Set(&m.CustomIDPInitiatedLoginPageURL, settings, "customIdpInitiatedLoginPageUrl")

		// The read returns an empty clientId for server-generated clients until a client_type is set, so it's only adopted when the state has no value.
		if m.ClientID.IsNull() || m.ClientID.IsUnknown() {
			stringattr.Set(&m.ClientID, settings, "clientId")
		}
		// The read always blanks the clientSecret, so the create-time value in the state is kept.
		stringattr.Set(&m.ClientSecret, settings, "clientSecret", stringattr.SkipIfAlreadySet)
		stringattr.Set(&m.ClientType, settings, "clientType")
		strsetattr.Set(&m.ApprovedRedirectURLs, settings, "approvedRedirectUrls", h)
		boolattr.Set(&m.AuthorizationCodeDisabled, settings, "authorizationCodeDisabled")
		boolattr.Set(&m.ClientCredentialsDisabled, settings, "clientCredentialsDisabled")
		boolattr.Set(&m.RefreshTokenDisabled, settings, "refreshTokenDisabled")
		boolattr.Set(&m.JWTBearerDisabled, settings, "jwtBearerDisabled")
		boolattr.Set(&m.DeviceCodeDisabled, settings, "deviceCodeDisabled")
		boolattr.Set(&m.ForcePkce, settings, "forcePkce")
		stringattr.Set(&m.DefaultAudience, settings, "defaultAudience")
		stringattr.Set(&m.TrustedAppsAudience, settings, "trustedAppsAudience")
	}
}

// The backend materializes a client_secret when an app without one is switched to the confidential client type.
func (m *OIDCAppModel) ModifyPlan(_ *helpers.Handler, config, state *OIDCAppModel) {
	if !config.ClientSecret.IsNull() || state.ClientSecret.ValueString() != "" {
		return // existing secrets never change on update
	}
	if m.ClientType.IsUnknown() || m.ClientType.ValueString() == "confidential" {
		m.ClientSecret = types.StringUnknown()
	}
}

func (m *OIDCAppModel) DeletionProtectionDefault(_ context.Context) bool {
	return true
}

func (m *OIDCAppModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OIDCAppModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *OIDCAppModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
