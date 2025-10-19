package authentication

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Attributes

var OAuthProviderCustomAttributes = map[string]schema.Attribute{
	// shared in all provider attributes
	"disabled":               boolattr.Default(false),
	"client_id":              stringattr.Default(""),
	"client_secret":          stringattr.SecretOptional(),
	"manage_provider_tokens": boolattr.Default(false),
	"callback_domain":        stringattr.Optional(),
	"redirect_url":           stringattr.Optional(),
	"prompts":                strlistattr.Optional(stringvalidator.OneOf("none", "login", "consent", "select_account")),
	"allowed_grant_types":    strlistattr.Optional(stringvalidator.OneOf("authorization_code", "implicit")),
	"scopes":                 strlistattr.Optional(),
	"merge_user_accounts":    boolattr.Default(true),

	// custom only
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

var OAuthProviderSystemAttributes = map[string]schema.Attribute{
	// shared in all provider attributes
	"disabled":               boolattr.Default(false),
	"client_id":              stringattr.Default(""),
	"client_secret":          stringattr.SecretOptional(),
	"manage_provider_tokens": boolattr.Default(false),
	"callback_domain":        stringattr.Optional(),
	"redirect_url":           stringattr.Optional(),
	"prompts":                strlistattr.Optional(stringvalidator.OneOf("none", "login", "consent", "select_account")),
	"allowed_grant_types":    strlistattr.Optional(stringvalidator.OneOf("authorization_code", "implicit")),
	"scopes":                 strlistattr.Optional(),
	"merge_user_accounts":    boolattr.Default(true),
}

var OAuthProviderAppleAttributes = map[string]schema.Attribute{
	// shared in all provider attributes
	"disabled":               boolattr.Default(false),
	"client_id":              stringattr.Default(""),
	"client_secret":          stringattr.SecretOptional(),
	"manage_provider_tokens": boolattr.Default(false),
	"callback_domain":        stringattr.Optional(),
	"redirect_url":           stringattr.Optional(),
	"prompts":                strlistattr.Optional(stringvalidator.OneOf("none", "login", "consent", "select_account")),
	"allowed_grant_types":    strlistattr.Optional(stringvalidator.OneOf("authorization_code", "implicit")),
	"scopes":                 strlistattr.Optional(),
	"merge_user_accounts":    boolattr.Default(true),

	// apple only
	"native_client_id":           stringattr.Default(""),
	"native_client_secret":       stringattr.SecretOptional(),
	"apple_key_generator":        objattr.Default[OAuthProviderAppleKeyGeneratorModel](nil, OAuthProviderAppleKeyGeneratorAttributes),
	"native_apple_key_generator": objattr.Default[OAuthProviderAppleKeyGeneratorModel](nil, OAuthProviderAppleKeyGeneratorAttributes),
}

// Models

type OAuthProviderBaseModel struct {
	Disabled             boolattr.Type    `tfsdk:"disabled"`
	ClientID             stringattr.Type  `tfsdk:"client_id"`
	ClientSecret         stringattr.Type  `tfsdk:"client_secret"`
	ManageProviderTokens boolattr.Type    `tfsdk:"manage_provider_tokens"`
	CallbackDomain       stringattr.Type  `tfsdk:"callback_domain"`
	RedirectURL          stringattr.Type  `tfsdk:"redirect_url"`
	Prompts              strlistattr.Type `tfsdk:"prompts"`
	AllowedGrantTypes    strlistattr.Type `tfsdk:"allowed_grant_types"`
	Scopes               strlistattr.Type `tfsdk:"scopes"`
	MergeUserAccounts    boolattr.Type    `tfsdk:"merge_user_accounts"`
}

type OAuthProviderCustomModel struct {
	OAuthProviderBaseModel

	Description           stringattr.Type `tfsdk:"description"`
	Logo                  stringattr.Type `tfsdk:"logo"`
	Issuer                stringattr.Type `tfsdk:"issuer"`
	AuthorizationEndpoint stringattr.Type `tfsdk:"authorization_endpoint"`
	TokenEndpoint         stringattr.Type `tfsdk:"token_endpoint"`
	UserInfoEndpoint      stringattr.Type `tfsdk:"user_info_endpoint"`
	JWKsEndpoint          stringattr.Type `tfsdk:"jwks_endpoint"`
	UseClientAssertion    boolattr.Type   `tfsdk:"use_client_assertion"`
	ClaimMapping          strmapattr.Type `tfsdk:"claim_mapping"`
}

type OAuthProviderSystemModel struct {
	OAuthProviderBaseModel
}

type OAuthProviderAppleModel struct {
	OAuthProviderBaseModel

	NativeClientID          stringattr.Type                                   `tfsdk:"native_client_id"`
	NativeClientSecret      stringattr.Type                                   `tfsdk:"native_client_secret"`
	AppleKeyGenerator       objattr.Type[OAuthProviderAppleKeyGeneratorModel] `tfsdk:"apple_key_generator"`
	NativeAppleKeyGenerator objattr.Type[OAuthProviderAppleKeyGeneratorModel] `tfsdk:"native_apple_key_generator"`
}

// Base

var OAuthProviderBaseValidator = objattr.NewValidator[OAuthProviderBaseModel]("must have a valid OAuth provider configuration")

func (m *OAuthProviderBaseModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
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
	return data
}

func (m *OAuthProviderBaseModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.ClientID, data, "clientId")
	stringattr.Nil(&m.ClientSecret)
	boolattr.Set(&m.ManageProviderTokens, data, "manageProviderTokens")
	stringattr.Set(&m.CallbackDomain, data, "callbackDomain")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	strlistattr.Set(&m.Prompts, data, "prompts", h)
	strlistattr.Set(&m.AllowedGrantTypes, data, "allowedGrantTypes", h)
	strlistattr.Set(&m.Scopes, data, "scopes", h)
	boolattr.Set(&m.MergeUserAccounts, data, "trustProvidedEmails")
}

func (m *OAuthProviderBaseModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.ClientID) {
		return // skip validation if there are unknown values
	}

	if m.ClientID.ValueString() == "" {
		if !m.Scopes.IsEmpty() {
			h.Invalid("The scopes attribute can only be set in an OAuth provider when a client_id is set")
		}
		if m.ManageProviderTokens.ValueBool() {
			h.Invalid("The manage_provider_tokens attribute can only be set in an OAuth provider when a client_id is set")
		}
		if m.CallbackDomain.ValueString() != "" {
			h.Invalid("The callback_domain attribute can only be set in an OAuth provider when a client_id is set")
		}
	}
}

// Custom

var OAuthProviderCustomValidator = objattr.NewValidator[OAuthProviderCustomModel]("must have a valid custom OAuth provider configuration")

func (m *OAuthProviderCustomModel) Values(h *helpers.Handler) map[string]any {
	data := m.OAuthProviderBaseModel.Values(h)
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Logo, data, "logo")
	stringattr.Get(m.Issuer, data, "issuer")
	stringattr.Get(m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Get(m.TokenEndpoint, data, "tokenUrl")
	stringattr.Get(m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Get(m.JWKsEndpoint, data, "jwksUrl")
	boolattr.Get(m.UseClientAssertion, data, "useClientAssertion")
	data["userDataClaimsMapping"] = customClaimValues(m.ClaimMapping, h)

	data["useSelfAccount"] = true
	data["custom"] = true
	data["useNonce"] = true

	return data
}

func (m *OAuthProviderCustomModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.OAuthProviderBaseModel.SetValues(h, data)
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Logo, data, "logo")
	stringattr.Set(&m.Issuer, data, "issuer")
	stringattr.Set(&m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Set(&m.TokenEndpoint, data, "tokenUrl")
	stringattr.Set(&m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Set(&m.JWKsEndpoint, data, "jwksUrl")
	boolattr.Set(&m.UseClientAssertion, data, "useClientAssertion")
	strmapattr.Nil(&m.ClaimMapping, h) // XXX empty defaults are added by the backend, add parsing for refresh
}

func (m *OAuthProviderCustomModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.ClientID, m.ClientSecret, m.UseClientAssertion, m.AllowedGrantTypes, m.AuthorizationEndpoint, m.TokenEndpoint, m.UserInfoEndpoint) {
		return // skip validation if there are unknown values
	}

	if m.ClientID.ValueString() == "" {
		h.Missing("The client_id attribute is required in a custom OAuth provider")
	}
	if !m.UseClientAssertion.ValueBool() && m.ClientSecret.ValueString() == "" {
		h.Missing("The client_secret attribute is required in a custom OAuth provider unless use_client_assertion is set to true")
	}
	if m.AllowedGrantTypes.IsEmpty() {
		h.Missing("The allowed_grant_types attribute must not be empty in a custom OAuth provider")
	}
	if m.AuthorizationEndpoint.ValueString() == "" {
		h.Missing("The authorization_endpoint attribute is required in a custom OAuth provider")
	}
	if m.TokenEndpoint.ValueString() == "" {
		h.Missing("The token_endpoint attribute is required in a custom OAuth provider")
	}
	if m.UserInfoEndpoint.ValueString() == "" {
		h.Missing("The user_info_endpoint attribute is required in a custom OAuth provider")
	}
	if v, _ := m.ClaimMapping.ToMap(h.Ctx); len(v) > 0 && v["loginId"].ValueString() == "" {
		h.Invalid("Claim mapping set for custom provider but 'loginId' was not mapped")
	}

	m.OAuthProviderBaseModel.Validate(h)
}

// System

var OAuthProviderSystemValidator = objattr.NewValidator[OAuthProviderSystemModel]("must have a valid system OAuth provider configuration")

func (m *OAuthProviderSystemModel) Values(h *helpers.Handler) map[string]any {
	data := m.OAuthProviderBaseModel.Values(h)

	data["useSelfAccount"] = m.ClientID.ValueString() != ""
	data["custom"] = false

	return data
}

func (m *OAuthProviderSystemModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.OAuthProviderBaseModel.SetValues(h, data)
}

func (m *OAuthProviderSystemModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.ClientID, m.ClientSecret) {
		return // skip validation if there are unknown values
	}

	if m.ClientID.ValueString() != "" && m.ClientSecret.ValueString() == "" {
		h.Missing("The client_secret attribute is required when client_id is set in a system OAuth provider")
	}

	m.OAuthProviderBaseModel.Validate(h)
}

// Apple

var OAuthProviderAppleValidator = objattr.NewValidator[OAuthProviderAppleModel]("must have a valid Apple system OAuth provider configuration")

func (m *OAuthProviderAppleModel) Values(h *helpers.Handler) map[string]any {
	data := m.OAuthProviderBaseModel.Values(h)
	stringattr.Get(m.NativeClientID, data, "nativeClientId")
	stringattr.Get(m.NativeClientSecret, data, "nativeClientSecret")
	objattr.Get(m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	objattr.Get(m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)

	data["useSelfAccount"] = m.ClientID.ValueString() != ""
	data["custom"] = false

	return data
}

func (m *OAuthProviderAppleModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.OAuthProviderBaseModel.SetValues(h, data)
	stringattr.Set(&m.NativeClientID, data, "nativeClientId")
	stringattr.Nil(&m.NativeClientSecret)
	objattr.Set(&m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	objattr.Set(&m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)
}

func (m *OAuthProviderAppleModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.ClientID, m.NativeClientID, m.ClientSecret, m.NativeClientSecret, m.AppleKeyGenerator, m.NativeAppleKeyGenerator) {
		return // skip validation if there are unknown values
	}

	if m.ClientID.ValueString() != "" {
		if m.ClientSecret.ValueString() == "" && !m.AppleKeyGenerator.IsSet() {
			h.Missing("The client_secret attribute or the apple_key_generator attribute must be set when client_id is set in an Apple OAuth provider")
		}
		if m.ClientSecret.ValueString() != "" && m.AppleKeyGenerator.IsSet() {
			h.Invalid("The client_secret and apple_key_generator attributes cannot both be set in an Apple OAuth provider")
		}
	}

	if m.NativeClientID.ValueString() != "" {
		if m.NativeClientSecret.ValueString() == "" && !m.NativeAppleKeyGenerator.IsSet() {
			h.Missing("The native_client_secret attribute or the native_apple_key_generator attribute must be set when native_client_id is set in an Apple OAuth provider")
		}
		if m.NativeClientSecret.ValueString() != "" && m.NativeAppleKeyGenerator.IsSet() {
			h.Invalid("The native_client_secret and native_apple_key_generator attributes cannot both be set in an Apple OAuth provider")
		}
	}

	m.OAuthProviderBaseModel.Validate(h)
}

// Apple Key Generator

var OAuthProviderAppleKeyGeneratorAttributes = map[string]schema.Attribute{
	"key_id":      stringattr.Required(),
	"team_id":     stringattr.Required(),
	"private_key": stringattr.SecretRequired(),
}

type OAuthProviderAppleKeyGeneratorModel struct {
	KeyID      stringattr.Type `tfsdk:"key_id"`
	TeamID     stringattr.Type `tfsdk:"team_id"`
	PrivateKey stringattr.Type `tfsdk:"private_key"`
}

func (m *OAuthProviderAppleKeyGeneratorModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.KeyID, data, "keyId")
	stringattr.Get(m.TeamID, data, "teamId")
	stringattr.Get(m.PrivateKey, data, "privateKey")
	return data
}

func (m *OAuthProviderAppleKeyGeneratorModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.KeyID, data, "keyId")
	stringattr.Set(&m.TeamID, data, "teamId")
	stringattr.Nil(&m.PrivateKey)
}

// Custom Claims

var descopeClaimMapping = []string{"loginId", "username", "name", "email", "phoneNumber", "verifiedEmail", "verifiedPhone", "picture", "givenName", "middleName", "familyName"}

func customClaimValues(m strmapattr.Type, h *helpers.Handler) map[string]any {
	claimMapping := map[string]any{}
	customAttributes := map[string]string{}
	for k, v := range strmapattr.Iterator(m, h) {
		if slices.Contains(descopeClaimMapping, k) {
			claimMapping[k] = v
		} else {
			customAttributes[k] = v
		}
	}
	claimMapping["customAttributes"] = customAttributes
	return claimMapping
}
