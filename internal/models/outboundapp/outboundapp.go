package outboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_outbound_app is a custom outbound application with a server-assigned id. Outbound apps hold the
// third party OAuth credentials Descope uses to fetch and refresh tokens on behalf of a user or a tenant.
//
// The use_dcr and dcr_url fields of the backend model are deliberately not exposed: the management API's
// update endpoint does not carry them, so managing them here would reset them on every apply. See the
// resource documentation for the consequence when importing an app that was created with dynamic
// client registration.

var OutboundAppSchema = schema.Schema{
	MarkdownDescription: "Manages a custom outbound application in a Descope project. Outbound applications let Descope hold and refresh a third party OAuth token on behalf of a user or a tenant, so an application can call that provider's API without implementing the OAuth flow itself.",
	Attributes:          OutboundAppAttributes,
}

var OutboundAppAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),

	// The tenant is absent from the backend's update statement, so it can only ever be set at creation.
	"tenant_id": stringattr.Default("", stringplanmodifier.RequiresReplace()),

	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default("", stringattr.StandardLenValidator),
	"logo":        stringattr.Default(""),
	"app_type":    stringattr.Default("", stringvalidator.OneOf("", "oauth", "apikey")),

	"client_id":     stringattr.Default("", stringattr.StandardLenValidator),
	"client_secret": stringattr.SecretOptional(),

	"discovery_url":            stringattr.Default("", stringattr.StandardLenValidator),
	"authorization_url":        stringattr.Default("", stringattr.StandardLenValidator),
	"authorization_url_params": listattr.Default[URLParamModel](URLParamAttributes),
	"token_url":                stringattr.Default("", stringattr.StandardLenValidator),
	"token_url_params":         listattr.Default[URLParamModel](URLParamAttributes),
	"revocation_url":           stringattr.Default("", stringattr.StandardLenValidator),

	"default_scopes":       strsetattr.Default(stringattr.StandardLenValidator),
	"default_redirect_url": stringattr.Default("", stringattr.URLValidator),
	"callback_domain":      stringattr.Default("", stringattr.StandardLenValidator),

	"pkce":        boolattr.Default(false),
	"access_type": stringattr.Default("", stringvalidator.OneOf("", "offline", "online")),
	"prompt": strsetattr.Default(setvalidator.ValueStringsAre(
		stringvalidator.OneOf("none", "login", "consent", "select_account"),
	)),
}

type OutboundAppModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	TenantID  stringattr.Type `tfsdk:"tenant_id"`

	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Logo        stringattr.Type `tfsdk:"logo"`
	AppType     stringattr.Type `tfsdk:"app_type"`

	ClientID     stringattr.Type `tfsdk:"client_id"`
	ClientSecret stringattr.Type `tfsdk:"client_secret"`

	DiscoveryURL           stringattr.Type              `tfsdk:"discovery_url"`
	AuthorizationURL       stringattr.Type              `tfsdk:"authorization_url"`
	AuthorizationURLParams listattr.Type[URLParamModel] `tfsdk:"authorization_url_params"`
	TokenURL               stringattr.Type              `tfsdk:"token_url"`
	TokenURLParams         listattr.Type[URLParamModel] `tfsdk:"token_url_params"`
	RevocationURL          stringattr.Type              `tfsdk:"revocation_url"`

	DefaultScopes      strsetattr.Type `tfsdk:"default_scopes"`
	DefaultRedirectURL stringattr.Type `tfsdk:"default_redirect_url"`
	CallbackDomain     stringattr.Type `tfsdk:"callback_domain"`

	Pkce       boolattr.Type   `tfsdk:"pkce"`
	AccessType stringattr.Type `tfsdk:"access_type"`
	Prompt     strsetattr.Type `tfsdk:"prompt"`
}

func (m *OutboundAppModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.TenantID, data, "tenantId")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Logo, data, "logo")
	stringattr.Get(m.AppType, data, "appType")
	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	stringattr.Get(m.DiscoveryURL, data, "discoveryUrl")
	stringattr.Get(m.AuthorizationURL, data, "authorizationUrl")
	listattr.Get(m.AuthorizationURLParams, data, "authorizationUrlParams", h)
	stringattr.Get(m.TokenURL, data, "tokenUrl")
	listattr.Get(m.TokenURLParams, data, "tokenUrlParams", h)
	stringattr.Get(m.RevocationURL, data, "revocationUrl")
	strsetattr.Get(m.DefaultScopes, data, "defaultScopes", h)
	stringattr.Get(m.DefaultRedirectURL, data, "defaultRedirectUrl")
	stringattr.Get(m.CallbackDomain, data, "callbackDomain")
	boolattr.Get(m.Pkce, data, "pkce")
	stringattr.Get(m.AccessType, data, "accessType")
	strsetattr.Get(m.Prompt, data, "prompt", h)
	return data
}

func (m *OutboundAppModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.TenantID, data, "tenantId")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Logo, data, "logo")
	stringattr.Set(&m.AppType, data, "appType")
	stringattr.Set(&m.ClientID, data, "clientId")
	// The backend never returns the client secret, so the value given at create time is kept in state.
	stringattr.Set(&m.ClientSecret, data, "clientSecret", stringattr.SkipIfAlreadySet)
	stringattr.Set(&m.DiscoveryURL, data, "discoveryUrl")
	stringattr.Set(&m.AuthorizationURL, data, "authorizationUrl")
	listattr.Set(&m.AuthorizationURLParams, data, "authorizationUrlParams", h)
	stringattr.Set(&m.TokenURL, data, "tokenUrl")
	listattr.Set(&m.TokenURLParams, data, "tokenUrlParams", h)
	stringattr.Set(&m.RevocationURL, data, "revocationUrl")
	strsetattr.Set(&m.DefaultScopes, data, "defaultScopes", h)
	stringattr.Set(&m.DefaultRedirectURL, data, "defaultRedirectUrl")
	stringattr.Set(&m.CallbackDomain, data, "callbackDomain")
	boolattr.Set(&m.Pkce, data, "pkce")
	stringattr.Set(&m.AccessType, data, "accessType")
	strsetattr.Set(&m.Prompt, data, "prompt", h)
}

func (m *OutboundAppModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OutboundAppModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *OutboundAppModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
