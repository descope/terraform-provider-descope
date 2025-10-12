package authentication

import (
	"fmt"
	"maps"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var OAuthValidator = objattr.NewValidator[OAuthModel]("must have a valid OAuth configuration")

var OAuthAttributes = map[string]schema.Attribute{
	"disabled": boolattr.Default(false),
	"system":   objattr.Default[OAuthSystemProvidersModel](nil, OAuthSystemProviderAttributes),
	"custom":   mapattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
}

type OAuthModel struct {
	Disabled boolattr.Type                           `tfsdk:"disabled"`
	System   objattr.Type[OAuthSystemProvidersModel] `tfsdk:"system"`
	Custom   mapattr.Type[OAuthProviderModel]        `tfsdk:"custom"`
}

func (m *OAuthModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")

	providers := map[string]any{}

	if v, _ := m.System.ToObject(h.Ctx); v != nil {
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

	for name, provider := range mapattr.Iterator(m.Custom, h) {
		claimMapping, _ := provider.ClaimMapping.ToMap(h.Ctx)
		if _, ok := claimMapping["loginId"]; !ok && len(claimMapping) > 0 {
			h.Error("Invalid Claim Mapping", "Claim mapping set for custom provider %s but 'loginId' was not mapped", name)
		}

		ensureRequiredCustomProviderField(h, provider.ClientID, "client_id", name)
		if !provider.UseClientAssertion.ValueBool() {
			ensureRequiredCustomProviderField(h, provider.ClientSecret, "client_secret", name)
		}
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

	system := map[string]any{}
	custom := map[string]any{}

	providers, _ := data["providerSettings"].(map[string]any)
	for name, provider := range providers {
		if slices.Contains(systemProviderNames, name) {
			system[name] = provider
		} else {
			custom[name] = provider
		}
	}

	objattr.Set(&m.System, system, helpers.RootKey, h)
	mapattr.Set(&m.Custom, custom, helpers.RootKey, h)
}

var systemProviderNames = []string{"apple", "discord", "facebook", "github", "gitlab", "google", "linkedin", "microsoft", "slack"}

func (m *OAuthModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.System, m.Custom) {
		return
	}

	if v, _ := m.System.ToObject(h.Ctx); v != nil {
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

	for name := range mapattr.Iterator(m.Custom, h) {
		if slices.Contains(systemProviderNames, name) {
			h.Error("Reserved OAuth Provider Name", "The %s name is reserved for system providers and cannot be used for a custom provider", name)
			continue
		}
	}
}

func ensureRequiredCustomProviderField(h *helpers.Handler, field attr.Value, fieldKey, name string) {
	var invalid bool

	switch v := field.(type) {
	case stringattr.Type:
		invalid = v.ValueString() == ""
	case strlistattr.Type:
		invalid = v.IsEmpty()
	default:
		panic(fmt.Sprintf("unexpected type %T for attribute %s in custom provider %s", field, fieldKey, name))
	}

	if invalid {
		h.Error("Invalid Custom OAuth Provider", "Custom provider %s needs a non-empty value for the %s attribute", name, fieldKey)
	}
}

func ensureSystemProvider(h *helpers.Handler, provider objattr.Type[OAuthProviderModel], name string) {
	m, _ := provider.ToObject(h.Ctx)
	if m == nil || helpers.HasUnknownValues(m.ClientID, m.ClientSecret) {
		return // skip validation if there are unknown values
	}

	ownAccount := m.ClientID.ValueString() != ""
	if ownAccount {
		if m.ClientSecret.ValueString() == "" {
			h.Missing("The client_id attribute was set for the %s system provider but the client_secret attribute was not", name)
		}
	} else {
		if !m.Scopes.IsEmpty() {
			h.Invalid("Set a client_id and client_secret for the %s system provider in order to set the scopes attribute", name)
		}
		if m.ManageProviderTokens.ValueBool() {
			h.Invalid("Set a client_id and client_secret for the %s system provider in order to set the manage_provider_tokens attribute", name)
		}
		if m.CallbackDomain.ValueString() != "" {
			h.Invalid("Set a client_id and client_secret for the %s system provider in order to set the callback_domain attribute", name)
		}
	}
}

func validateSystemProvider(h *helpers.Handler, provider objattr.Type[OAuthProviderModel], name string) {
	m, _ := provider.ToObject(h.Ctx)
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
	ensureNoCustomProviderFields(h, m.UseClientAssertion, "use_client_assertion", name)

	if !m.ClaimMapping.IsEmpty() {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its claim_mapping attribute is reserved", name)
	}
}

func ensureNoCustomProviderFields(h *helpers.Handler, field attr.Value, fieldKey, name string) {
	if !field.IsUnknown() && !field.IsNull() {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its %s attribute is reserved", name, fieldKey)
	}
}

// System OAuth Providers

var OAuthSystemProviderAttributes = map[string]schema.Attribute{
	"apple":     objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"discord":   objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"facebook":  objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"github":    objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"gitlab":    objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"google":    objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"linkedin":  objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"microsoft": objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
	"slack":     objattr.Default[OAuthProviderModel](nil, OAuthProviderAttributes, OAuthProviderValidator),
}

type OAuthSystemProvidersModel struct {
	Apple     objattr.Type[OAuthProviderModel] `tfsdk:"apple"`
	Discord   objattr.Type[OAuthProviderModel] `tfsdk:"discord"`
	Facebook  objattr.Type[OAuthProviderModel] `tfsdk:"facebook"`
	Github    objattr.Type[OAuthProviderModel] `tfsdk:"github"`
	Gitlab    objattr.Type[OAuthProviderModel] `tfsdk:"gitlab"`
	Google    objattr.Type[OAuthProviderModel] `tfsdk:"google"`
	Linkedin  objattr.Type[OAuthProviderModel] `tfsdk:"linkedin"`
	Microsoft objattr.Type[OAuthProviderModel] `tfsdk:"microsoft"`
	Slack     objattr.Type[OAuthProviderModel] `tfsdk:"slack"`
}

func (m *OAuthSystemProvidersModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	getProviderValue(h, data, m.Apple, "apple")
	getProviderValue(h, data, m.Discord, "discord")
	getProviderValue(h, data, m.Facebook, "facebook")
	getProviderValue(h, data, m.Github, "github")
	getProviderValue(h, data, m.Gitlab, "gitlab")
	getProviderValue(h, data, m.Google, "google")
	getProviderValue(h, data, m.Linkedin, "linkedin")
	getProviderValue(h, data, m.Microsoft, "microsoft")
	getProviderValue(h, data, m.Slack, "slack")
	return data
}

func (m *OAuthSystemProvidersModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.Apple, data, "apple", h)
	objattr.Set(&m.Discord, data, "discord", h)
	objattr.Set(&m.Facebook, data, "facebook", h)
	objattr.Set(&m.Github, data, "github", h)
	objattr.Set(&m.Gitlab, data, "gitlab", h)
	objattr.Set(&m.Google, data, "google", h)
	objattr.Set(&m.Linkedin, data, "linkedin", h)
	objattr.Set(&m.Microsoft, data, "microsoft", h)
	objattr.Set(&m.Slack, data, "slack", h)
}

func getProviderValue(h *helpers.Handler, providers map[string]any, obj objattr.Type[OAuthProviderModel], name string) {
	provider, _ := obj.ToObject(h.Ctx)
	if provider == nil {
		return
	}

	data := provider.Values(h)
	data["useSelfAccount"] = provider.ClientID.ValueString() != ""
	data["custom"] = false
	providers[name] = data
}

// OAuth Provider

var systemClaimMapping = []string{"loginId", "username", "name", "email", "phoneNumber", "verifiedEmail", "verifiedPhone", "picture", "givenName", "middleName", "familyName"}

var OAuthProviderValidator = objattr.NewValidator[OAuthProviderModel]("must have a valid OAuth provider configuration")

var OAuthProviderAttributes = map[string]schema.Attribute{
	"disabled":                   boolattr.Default(false),
	"client_id":                  stringattr.Optional(),
	"client_secret":              stringattr.SecretOptional(),
	"manage_provider_tokens":     boolattr.Default(false),
	"callback_domain":            stringattr.Optional(),
	"redirect_url":               stringattr.Optional(),
	"provider_token_management":  objattr.Default[OAuthProviderTokenManagementModel](nil, OAuthProviderTokenManagementAttributes),
	"prompts":                    strlistattr.Optional(stringvalidator.OneOf("none", "login", "consent", "select_account")),
	"allowed_grant_types":        strlistattr.Optional(stringvalidator.OneOf("authorization_code", "implicit")),
	"scopes":                     strlistattr.Optional(),
	"merge_user_accounts":        boolattr.Default(true),
	"default_sso_roles":          strlistattr.Optional(),
	"native_client_id":           stringattr.Optional(),
	"native_client_secret":       stringattr.SecretOptional(),
	"apple_key_generator":        objattr.Default[AppleKeyGeneratorModel](nil, AppleKeyGeneratorModelAttributes),
	"native_apple_key_generator": objattr.Default[AppleKeyGeneratorModel](nil, AppleKeyGeneratorModelAttributes),
	// editable for custom only
	"description":            stringattr.Optional(),
	"logo":                   stringattr.Optional(),
	"issuer":                 stringattr.Optional(),
	"authorization_endpoint": stringattr.Optional(),
	"token_endpoint":         stringattr.Optional(),
	"user_info_endpoint":     stringattr.Optional(),
	"jwks_endpoint":          stringattr.Optional(),
	"use_client_assertion":   boolattr.Default(false),
	"claim_mapping":          strmapattr.Optional(),
}

type AppleKeyGeneratorModel struct {
	KeyID      stringattr.Type `tfsdk:"key_id"`
	TeamID     stringattr.Type `tfsdk:"team_id"`
	PrivateKey stringattr.Type `tfsdk:"private_key"`
}

var AppleKeyGeneratorModelAttributes = map[string]schema.Attribute{
	"key_id":      stringattr.Required(),
	"team_id":     stringattr.Required(),
	"private_key": stringattr.SecretRequired(),
}

func (m *AppleKeyGeneratorModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.KeyID, data, "key_id")
	stringattr.Get(m.TeamID, data, "team_id")
	stringattr.Get(m.PrivateKey, data, "private_key")
	return data
}

func (m *AppleKeyGeneratorModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.KeyID, data, "key_id")
	stringattr.Set(&m.TeamID, data, "team_id")
	stringattr.Set(&m.PrivateKey, data, "private_key")
}

type OAuthProviderModel struct {
	Disabled                boolattr.Type                                   `tfsdk:"disabled"`
	ClientID                stringattr.Type                                 `tfsdk:"client_id"`
	ClientSecret            stringattr.Type                                 `tfsdk:"client_secret"`
	ManageProviderTokens    boolattr.Type                                   `tfsdk:"manage_provider_tokens"`
	CallbackDomain          stringattr.Type                                 `tfsdk:"callback_domain"`
	RedirectURL             stringattr.Type                                 `tfsdk:"redirect_url"`
	ProviderTokenManagement objattr.Type[OAuthProviderTokenManagementModel] `tfsdk:"provider_token_management"`
	Prompts                 strlistattr.Type                                `tfsdk:"prompts"`
	Scopes                  strlistattr.Type                                `tfsdk:"scopes"`
	MergeUserAccounts       boolattr.Type                                   `tfsdk:"merge_user_accounts"`
	Description             stringattr.Type                                 `tfsdk:"description"`
	Logo                    stringattr.Type                                 `tfsdk:"logo"`
	AllowedGrantTypes       strlistattr.Type                                `tfsdk:"allowed_grant_types"`
	Issuer                  stringattr.Type                                 `tfsdk:"issuer"`
	AuthorizationEndpoint   stringattr.Type                                 `tfsdk:"authorization_endpoint"`
	TokenEndpoint           stringattr.Type                                 `tfsdk:"token_endpoint"`
	UserInfoEndpoint        stringattr.Type                                 `tfsdk:"user_info_endpoint"`
	JWKsEndpoint            stringattr.Type                                 `tfsdk:"jwks_endpoint"`
	UseClientAssertion      boolattr.Type                                   `tfsdk:"use_client_assertion"`
	UseNativeClient         boolattr.Type                                   `tfsdk:"use_native_client"`
	ClaimMapping            strmapattr.Type                                 `tfsdk:"claim_mapping"`
	DefaultSSORoles         strlistattr.Type                                `tfsdk:"default_sso_roles"`
	NativeClientID          stringattr.Type                                 `tfsdk:"native_client_id"`
	NativeClientSecret      stringattr.Type                                 `tfsdk:"native_client_secret"`
	AppleKeyGenerator       objattr.Type[AppleKeyGeneratorModel]            `tfsdk:"apple_key_generator"`
	NativeAppleKeyGenerator objattr.Type[AppleKeyGeneratorModel]            `tfsdk:"native_apple_key_generator"`
}

func (m *OAuthProviderModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{
		"enabled": !m.Disabled.ValueBool(),
	}
	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	boolattr.Get(m.ManageProviderTokens, data, "manageProviderTokens")
	stringattr.Get(m.CallbackDomain, data, "callbackDomain")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	if !m.Prompts.IsEmpty() {
		strlistattr.Get(m.Prompts, data, "prompts", h)
	}
	if !m.AllowedGrantTypes.IsEmpty() {
		strlistattr.Get(m.AllowedGrantTypes, data, "allowedGrantTypes", h)
	}
	if !m.Scopes.IsEmpty() {
		strlistattr.Get(m.Scopes, data, "scopes", h)
	}
	boolattr.Get(m.MergeUserAccounts, data, "trustProvidedEmails")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Logo, data, "logo")
	stringattr.Get(m.Issuer, data, "issuer")
	stringattr.Get(m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Get(m.TokenEndpoint, data, "tokenUrl")
	stringattr.Get(m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Get(m.JWKsEndpoint, data, "jwksUrl")
	boolattr.Get(m.UseClientAssertion, data, "useClientAssertion")
	claimMapping := map[string]any{}
	customAttributes := map[string]string{}
	for k, v := range strmapattr.Iterator(m.ClaimMapping, h) {
		if slices.Contains(systemClaimMapping, k) {
			claimMapping[k] = v
		} else {
			customAttributes[k] = v
		}
	}
	claimMapping["customAttributes"] = customAttributes
	data["userDataClaimsMapping"] = claimMapping
	if !m.DefaultSSORoles.IsEmpty() {
		strlistattr.Get(m.DefaultSSORoles, data, "defaultSSORoles", h)
	}
	stringattr.Get(m.NativeClientID, data, "nativeClientId")
	stringattr.Get(m.NativeClientSecret, data, "nativeClientSecret")
	if !m.AppleKeyGenerator.IsNull() && !m.AppleKeyGenerator.IsUnknown() {
		objattr.Get(m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	}
	if !m.NativeAppleKeyGenerator.IsNull() && !m.NativeAppleKeyGenerator.IsUnknown() {
		objattr.Get(m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)
	}
	return data
}

func (m *OAuthProviderModel) SetValues(h *helpers.Handler, data map[string]any) {
	if b, ok := data["enabled"].(bool); ok {
		m.Disabled = boolattr.Value(!b)
	}
	stringattr.Set(&m.ClientID, data, "clientId")
	boolattr.Set(&m.ManageProviderTokens, data, "manageProviderTokens")
	stringattr.Set(&m.CallbackDomain, data, "callbackDomain")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	strlistattr.Set(&m.Prompts, data, "prompts", h)                     // XXX was skipped
	strlistattr.Set(&m.AllowedGrantTypes, data, "allowedGrantTypes", h) // XXX was skipped
	strlistattr.Set(&m.Scopes, data, "scopes", h)
	boolattr.Set(&m.MergeUserAccounts, data, "trustProvidedEmails")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Logo, data, "logo")
	stringattr.Set(&m.Issuer, data, "issuer")
	stringattr.Set(&m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Set(&m.TokenEndpoint, data, "tokenUrl")
	stringattr.Set(&m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Set(&m.JWKsEndpoint, data, "jwksUrl")
	boolattr.Set(&m.UseClientAssertion, data, "useClientAssertion")
	m.UseNativeClient = boolattr.Value(false)
	if _, ok := data["nativeClientId"]; ok {
		stringattr.Set(&m.NativeClientID, data, "nativeClientId")
		m.UseNativeClient = boolattr.Value(true)
	}
	if _, ok := data["nativeClientSecret"]; ok {
		stringattr.Set(&m.NativeClientSecret, data, "nativeClientSecret")
		m.UseNativeClient = boolattr.Value(true)
	}

	if _, ok := data["defaultSSORoles"]; ok {
		strlistattr.Set(&m.DefaultSSORoles, data, "defaultSSORoles", h)
	}
	if _, ok := data["appleKeyGenerator"]; ok {
		objattr.Set(&m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	}
	if _, ok := data["nativeAppleKeyGenerator"]; ok {
		objattr.Set(&m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)
	}
	strmapattr.Nil(&m.ClaimMapping, h) // XXX empty defaults are added by the backend, add parsing for refresh
}

func (m *OAuthProviderModel) Validate(h *helpers.Handler) {
	if m.ProviderTokenManagement.IsSet() {
		h.Error("Deprecated Field", "The provider_token_management field is deprecated. Use manage_provider_tokens, callback_domain, and redirect_url fields directly on the OAuth provider instead.")
	}
}

// Provider Token Management

var OAuthProviderTokenManagementAttributes = map[string]schema.Attribute{}

type OAuthProviderTokenManagementModel struct{}
