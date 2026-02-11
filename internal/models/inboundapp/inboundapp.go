package inboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var InboundAppAttributes = map[string]schema.Attribute{
	"id":                               stringattr.Identifier(),
	"project_id":                       stringattr.Required(stringplanmodifier.RequiresReplace()),
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
	"default_audience":                 stringattr.Default("", stringvalidator.OneOf("", "projectId", "clientId")), // XXX maybe switch to set
	"non_confidential_client":          boolattr.Default(false, boolplanmodifier.RequiresReplace()),
	"client_id":                        stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"client_secret":                    stringattr.SecretGenerated(true),
}

type InboundAppModel struct {
	ID                           stringattr.Type                      `tfsdk:"id"`
	ProjectID                    stringattr.Type                      `tfsdk:"project_id"`
	Name                         stringattr.Type                      `tfsdk:"name"`
	Description                  stringattr.Type                      `tfsdk:"description"`
	LogoUrl                      stringattr.Type                      `tfsdk:"logo_url"`
	LoginPageUrl                 stringattr.Type                      `tfsdk:"login_page_url"`
	ApprovedCallbackUrls         strsetattr.Type                      `tfsdk:"approved_callback_urls"`
	PermissionsScopes            listattr.Type[ApplicationScopeModel] `tfsdk:"permissions_scopes"`
	AttributesScopes             listattr.Type[ApplicationScopeModel] `tfsdk:"attributes_scopes"`
	ConnectionsScopes            listattr.Type[ApplicationScopeModel] `tfsdk:"connections_scopes"`
	SessionSettings              objattr.Type[SessionSettingsModel]   `tfsdk:"session_settings"`
	AudienceWhitelist            strsetattr.Type                      `tfsdk:"audience_whitelist"`
	ForceAddAllAuthorizationInfo boolattr.Type                        `tfsdk:"force_add_all_authorization_info"`
	DefaultAudience              stringattr.Type                      `tfsdk:"default_audience"`
	NonConfidentialClient        boolattr.Type                        `tfsdk:"non_confidential_client"`
	ClientId                     stringattr.Type                      `tfsdk:"client_id"`
	ClientSecret                 stringattr.Type                      `tfsdk:"client_secret"`
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
	stringattr.Get(m.DefaultAudience, data, "defaultAudience")
	boolattr.Get(m.NonConfidentialClient, data, "nonConfidentialClient")
	stringattr.Get(m.ClientId, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
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
	stringattr.Set(&m.DefaultAudience, data, "defaultAudience")
	boolattr.Set(&m.NonConfidentialClient, data, "nonConfidentialClient")
	stringattr.Set(&m.ClientId, data, "clientId")
	stringattr.Set(&m.ClientSecret, data, "clientSecret")
}
