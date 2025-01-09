package authentication

import (
	"fmt"
	"maps"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OAuthValidator = objectattr.NewValidator[OAuthModel]("must have a valid OAuth configuration")

var OAuthAttributes = map[string]schema.Attribute{
	"disabled": boolattr.Default(false),
	"system":   objectattr.Optional(OAuthSystemProviderAttributes),
	"custom":   mapattr.Optional(OAuthProviderAttributes),
}

type OAuthModel struct {
	Disabled types.Bool                     `tfsdk:"disabled"`
	System   *OAuthSystemProvidersModel     `tfsdk:"system"`
	Custom   map[string]*OAuthProviderModel `tfsdk:"custom"`
}

func (m *OAuthModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")

	providers := map[string]any{}

	if v := m.System; v != nil {
		ensureSystemProvider(h, v.Apple, "apple")
		ensureSystemProvider(h, v.Discord, "discord")
		ensureSystemProvider(h, v.Facebook, "facebook")
		ensureSystemProvider(h, v.Github, "github")
		ensureSystemProvider(h, v.Gitlab, "gitlab")
		ensureSystemProvider(h, v.Google, "google")
		ensureSystemProvider(h, v.Linkedin, "linkedin")
		ensureSystemProvider(h, v.Microsoft, "microsoft")
		ensureSystemProvider(h, v.Slack, "slack")

		maps.Copy(providers, v.Values(h))
	}

	for name, provider := range m.Custom {
		if _, ok := provider.ClaimMapping["loginId"]; !ok && len(provider.ClaimMapping) > 0 {
			h.Error("Invalid Claim Mapping", "Claim mapping set for custom provider %s but 'loginId' was not mapped", name)
		}

		ensureRequiredCustomProviderField(h, provider.ClientID, "client_id", name)
		ensureRequiredCustomProviderField(h, provider.ClientSecret, "client_secret", name)
		ensureRequiredCustomProviderField(h, provider.AllowedGrantTypes, "allowed_grant_types", name)
		ensureRequiredCustomProviderField(h, provider.AuthorizationEndpoint, "authorization_endpoint", name)
		ensureRequiredCustomProviderField(h, provider.TokenEndpoint, "token_endpoint", name)
		ensureRequiredCustomProviderField(h, provider.UserInfoEndpoint, "user_info_endpoint", name)

		data := provider.Values(h)
		data["useSelfAccount"] = true
		data["custom"] = true
		data["useNonce"] = true
		data["name"] = name
		providers[name] = data
	}

	data["providerSettings"] = providers
	return data
}

func (m *OAuthModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	if providers, ok := data["providerSettings"].(map[string]any); ok {
		if v := m.System; v != nil {
			v.SetValues(h, providers)
		}
		for k, v := range m.Custom {
			if data, ok := providers[k].(map[string]any); ok {
				v.SetValues(h, data)
			}
		}
	}
}

var systemProviderNames = []string{"apple", "discord", "facebook", "github", "gitlab", "google", "linkedin", "microsoft", "slack"}

func (m *OAuthModel) Validate(h *helpers.Handler) {
	if v := m.System; v != nil {
		validateSystemProvider(h, v.Apple, "apple")
		validateSystemProvider(h, v.Discord, "discord")
		validateSystemProvider(h, v.Facebook, "facebook")
		validateSystemProvider(h, v.Github, "github")
		validateSystemProvider(h, v.Gitlab, "gitlab")
		validateSystemProvider(h, v.Google, "google")
		validateSystemProvider(h, v.Linkedin, "linkedin")
		validateSystemProvider(h, v.Microsoft, "microsoft")
		validateSystemProvider(h, v.Slack, "slack")
	}
	for name := range m.Custom {
		if slices.Contains(systemProviderNames, name) {
			h.Error("Reserved OAuth Provider Name", "The %s name is reserved for system providers and cannot be used for a custom provider", name)
			continue
		}
	}
}

func ensureRequiredCustomProviderField(h *helpers.Handler, field any, fieldKey, name string) {
	invalid := false

	switch v := field.(type) {
	case types.String:
		invalid = v.ValueString() == ""
	case []string:
		invalid = len(v) == 0
	default:
		panic(fmt.Sprintf("unexpected type %T for attribute %s in custom provider %s", field, fieldKey, name))
	}

	if invalid {
		h.Error("Invalid Custom OAuth Provider", "Custom provider %s must set a non-empty value for the %s attribute", name, fieldKey)
	}
}

func ensureSystemProvider(h *helpers.Handler, m *OAuthProviderModel, name string) {
	if m == nil {
		return
	}
	// own account specific validations
	ownAccount := m.ClientID.ValueString() != ""
	if ownAccount {
		if m.ClientSecret.ValueString() == "" {
			h.Error("Missing Client Secret", "The client_id attribute was set for the %s system provider but the client_secret attribute was not", name)
		}
	} else {
		if len(m.Scopes) > 0 {
			h.Error("Invalid Attribute Value", "Set a client_id and client_secret for the %s system provider in order to set the scopes attribute", name)
		}
		if m.ProviderTokenManagement != nil {
			h.Error("Invalid Attribute Value", "Set a client_id and client_secret for the %s system provider in order to set the provider_token_management attribute", name)
		}
	}
}

func validateSystemProvider(h *helpers.Handler, m *OAuthProviderModel, name string) {
	if m == nil {
		return
	}
	ensureNoCustomProviderFields(h, m.Description, "description", name)
	ensureNoCustomProviderFields(h, m.Logo, "logo", name)
	ensureNoCustomProviderFields(h, m.Issuer, "issuer", name)
	ensureNoCustomProviderFields(h, m.AuthorizationEndpoint, "authorization_endpoint", name)
	ensureNoCustomProviderFields(h, m.TokenEndpoint, "token_endpoint", name)
	ensureNoCustomProviderFields(h, m.UserInfoEndpoint, "user_info_endpoint", name)
	ensureNoCustomProviderFields(h, m.JWKsEndpoint, "jwks_endpoint", name)
	if len(m.ClaimMapping) > 0 {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its claim_mapping attribute is reserved", name)
	}
}

func ensureNoCustomProviderFields(h *helpers.Handler, field types.String, fieldKey, name string) {
	if !field.IsUnknown() && !field.IsNull() {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its %s attribute is reserved", name, fieldKey)
	}
}

// System OAuth Providers

var OAuthSystemProviderAttributes = map[string]schema.Attribute{
	"apple":     objectattr.Optional(OAuthProviderAttributes),
	"discord":   objectattr.Optional(OAuthProviderAttributes),
	"facebook":  objectattr.Optional(OAuthProviderAttributes),
	"github":    objectattr.Optional(OAuthProviderAttributes),
	"gitlab":    objectattr.Optional(OAuthProviderAttributes),
	"google":    objectattr.Optional(OAuthProviderAttributes),
	"linkedin":  objectattr.Optional(OAuthProviderAttributes),
	"microsoft": objectattr.Optional(OAuthProviderAttributes),
	"slack":     objectattr.Optional(OAuthProviderAttributes),
}

type OAuthSystemProvidersModel struct {
	Apple     *OAuthProviderModel `tfsdk:"apple"`
	Discord   *OAuthProviderModel `tfsdk:"discord"`
	Facebook  *OAuthProviderModel `tfsdk:"facebook"`
	Github    *OAuthProviderModel `tfsdk:"github"`
	Gitlab    *OAuthProviderModel `tfsdk:"gitlab"`
	Google    *OAuthProviderModel `tfsdk:"google"`
	Linkedin  *OAuthProviderModel `tfsdk:"linkedin"`
	Microsoft *OAuthProviderModel `tfsdk:"microsoft"`
	Slack     *OAuthProviderModel `tfsdk:"slack"`
}

func (m *OAuthSystemProvidersModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	getSystemProvider(h, data, m.Apple, "apple")
	getSystemProvider(h, data, m.Discord, "discord")
	getSystemProvider(h, data, m.Facebook, "facebook")
	getSystemProvider(h, data, m.Github, "github")
	getSystemProvider(h, data, m.Gitlab, "gitlab")
	getSystemProvider(h, data, m.Google, "google")
	getSystemProvider(h, data, m.Linkedin, "linkedin")
	getSystemProvider(h, data, m.Microsoft, "microsoft")
	getSystemProvider(h, data, m.Slack, "slack")
	return data
}

func getSystemProvider(h *helpers.Handler, providers map[string]any, provider *OAuthProviderModel, name string) {
	if provider == nil {
		return
	}
	data := provider.Values(h)
	data["useSelfAccount"] = provider.ClientID.ValueString() != ""
	data["custom"] = false
	providers[name] = data
}

func (m *OAuthSystemProvidersModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSystemProvider(h, data, m.Apple, "apple")
	setSystemProvider(h, data, m.Discord, "discord")
	setSystemProvider(h, data, m.Facebook, "facebook")
	setSystemProvider(h, data, m.Github, "github")
	setSystemProvider(h, data, m.Gitlab, "gitlab")
	setSystemProvider(h, data, m.Google, "google")
	setSystemProvider(h, data, m.Linkedin, "linkedin")
	setSystemProvider(h, data, m.Microsoft, "microsoft")
	setSystemProvider(h, data, m.Slack, "slack")
}

func setSystemProvider(h *helpers.Handler, providers map[string]any, provider *OAuthProviderModel, name string) {
	if provider == nil {
		return
	}
	if data, ok := providers[name].(map[string]any); ok {
		provider.SetValues(h, data)
	}
}

// OAuth Provider

var systemClaimMapping = []string{"loginId", "username", "name", "email", "phoneNumber", "verifiedEmail", "verifiedPhone", "picture", "givenName", "middleName", "familyName"}

var OAuthProviderAttributes = map[string]schema.Attribute{
	"disabled":                  boolattr.Default(false),
	"client_id":                 stringattr.Optional(),
	"client_secret":             stringattr.SecretOptional(),
	"provider_token_management": objectattr.Optional(OAuthProviderTokenManagementAttribute),
	"prompts":                   strlistattr.Optional(listvalidator.ValueStringsAre(stringvalidator.OneOf("none", "login", "consent", "select_account"))),
	"allowed_grant_types":       strlistattr.Optional(listvalidator.ValueStringsAre(stringvalidator.OneOf("authorization_code", "implicit"))),
	"scopes":                    strlistattr.Optional(),
	"merge_user_accounts":       boolattr.Default(true),
	// editable for custom only
	"description":            stringattr.Optional(),
	"logo":                   stringattr.Optional(),
	"issuer":                 stringattr.Optional(),
	"authorization_endpoint": stringattr.Optional(),
	"token_endpoint":         stringattr.Optional(),
	"user_info_endpoint":     stringattr.Optional(),
	"jwks_endpoint":          stringattr.Optional(),
	"claim_mapping":          mapattr.StringOptional(),
}

type OAuthProviderModel struct {
	Disabled                types.Bool                         `tfsdk:"disabled"`
	ClientID                types.String                       `tfsdk:"client_id"`
	ClientSecret            types.String                       `tfsdk:"client_secret"`
	ProviderTokenManagement *OAuthProviderTokenManagementModel `tfsdk:"provider_token_management"`
	Prompts                 []string                           `tfsdk:"prompts"`
	Scopes                  []string                           `tfsdk:"scopes"`
	MergeUserAccounts       types.Bool                         `tfsdk:"merge_user_accounts"`
	Description             types.String                       `tfsdk:"description"`
	Logo                    types.String                       `tfsdk:"logo"`
	AllowedGrantTypes       []string                           `tfsdk:"allowed_grant_types"`
	Issuer                  types.String                       `tfsdk:"issuer"`
	AuthorizationEndpoint   types.String                       `tfsdk:"authorization_endpoint"`
	TokenEndpoint           types.String                       `tfsdk:"token_endpoint"`
	UserInfoEndpoint        types.String                       `tfsdk:"user_info_endpoint"`
	JWKsEndpoint            types.String                       `tfsdk:"jwks_endpoint"`
	ClaimMapping            map[string]string                  `tfsdk:"claim_mapping"`
}

func (m *OAuthProviderModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{
		"enabled": !m.Disabled.ValueBool(),
	}
	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	if ptk := m.ProviderTokenManagement; ptk != nil {
		data["manageProviderTokens"] = true
		stringattr.Get(ptk.CallbackDomain, data, "callbackDomain")
		stringattr.Get(ptk.RedirectURL, data, "redirectUrl")
	} else {
		data["manageProviderTokens"] = false
	}
	if len(m.Prompts) > 0 {
		strlistattr.Get(m.Prompts, data, "prompts")
	}
	if len(m.Scopes) > 0 {
		strlistattr.Get(m.Scopes, data, "scopes")
	}
	boolattr.Get(m.MergeUserAccounts, data, "trustProvidedEmails")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Logo, data, "logo")
	if len(m.AllowedGrantTypes) > 0 {
		strlistattr.Get(m.AllowedGrantTypes, data, "allowedGrantTypes")
	}
	stringattr.Get(m.Issuer, data, "issuer")
	stringattr.Get(m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Get(m.TokenEndpoint, data, "tokenUrl")
	stringattr.Get(m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Get(m.JWKsEndpoint, data, "jwksUrl")
	claimMapping := map[string]any{}
	customAttributes := map[string]string{}
	for k, v := range m.ClaimMapping {
		if slices.Contains(systemClaimMapping, k) {
			claimMapping[k] = v
		} else {
			customAttributes[k] = v
		}
	}
	claimMapping["customAttributes"] = customAttributes
	data["userDataClaimsMapping"] = claimMapping
	return data
}

func (m *OAuthProviderModel) SetValues(h *helpers.Handler, data map[string]any) {
	if b, ok := data["enabled"].(bool); ok {
		m.Disabled = types.BoolValue(!b)
	}
	stringattr.Set(&m.ClientID, data, "clientId")
	// m.ClientSecret - Not setting the secret as it is returned obfuscated
	if data["manageProviderTokens"] == true {
		m.ProviderTokenManagement = &OAuthProviderTokenManagementModel{}
		stringattr.Set(&m.ProviderTokenManagement.CallbackDomain, data, "callbackDomain")
		stringattr.Set(&m.ProviderTokenManagement.RedirectURL, data, "redirectUrl")
	}
	// Skipped for now: m.Prompts = helpers.AnySliceToStringSlice(data, "prompts")
	m.Scopes = helpers.AnySliceToStringSlice(data, "scopes")
	boolattr.Set(&m.MergeUserAccounts, data, "trustProvidedEmails")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Logo, data, "logo")
	// Skipped for now: m.AllowedGrantTypes = helpers.AnySliceToStringSlice(data, "allowedGrantTypes")
	stringattr.Set(&m.Issuer, data, "issuer")
	stringattr.Set(&m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Set(&m.TokenEndpoint, data, "tokenUrl")
	stringattr.Set(&m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Set(&m.JWKsEndpoint, data, "jwksUrl")
	// m.ClaimMapping - Not setting the claims, as empty defaults are added by the BE
}

// Provider Token Management

var OAuthProviderTokenManagementAttribute = map[string]schema.Attribute{
	"callback_domain": stringattr.Optional(),
	"redirect_url":    stringattr.Optional(),
}

type OAuthProviderTokenManagementModel struct {
	CallbackDomain types.String `tfsdk:"callback_domain"`
	RedirectURL    types.String `tfsdk:"redirect_url"`
}

func (m *OAuthProviderTokenManagementModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.CallbackDomain, data, "callbackDomain")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	return data
}

func (m *OAuthProviderTokenManagementModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.CallbackDomain, data, "callbackDomain")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
}
