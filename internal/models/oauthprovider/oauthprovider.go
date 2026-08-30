package oauthprovider

import (
	"fmt"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single OAuth provider in a Descope project, either the custom configuration of a built-in system provider such as Google or GitHub, or a fully custom OAuth provider. " +
		"Applying this resource to an OAuth provider that was previously configured in the Descope Console takes over its configuration declaratively: console-only fields that aren't expressible in this schema are normalized to their schema defaults on the first apply. " +
		"Client secrets are write-only, so they are never read back into state, and a configured secret is preserved on the backend when the attribute is omitted from the plan. " +
		"The connection endpoints of built-in system providers (such as `authorization_endpoint` and `token_endpoint`) are managed by Descope and cannot be configured, so those attributes should only be set for custom providers.",
	Attributes: ProviderAttributes,
}

var ProviderAttributes = map[string]schema.Attribute{
	"id":         stringattr.Required(systemCasingValidator{}, stringplanmodifier.RequiresReplace()),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),

	"disabled":               boolattr.Default(false),
	"client_id":              stringattr.Optional(),
	"client_secret":          stringattr.SecretOptional(),
	"manage_provider_tokens": boolattr.Default(false),
	"callback_domain":        stringattr.Optional(),
	"redirect_url":           stringattr.Optional(stringattr.URLValidator),
	"prompts":                strlistattr.Optional(stringvalidator.OneOf("none", "login", "consent", "select_account")),
	"allowed_grant_types":    strlistattr.Optional(stringvalidator.OneOf("authorization_code", "implicit")),
	"scopes":                 strlistattr.Optional(),
	"merge_user_accounts":    boolattr.Default(true),
	"disable_jit_updates":    boolattr.Default(false),
	"native_client_id":       stringattr.Optional(),
	"native_client_secret":   stringattr.SecretOptional(),

	// apple system provider only
	"apple_key_generator":        objattr.Default[AppleKeyGeneratorModel](nil, AppleKeyGeneratorModelAttributes),
	"native_apple_key_generator": objattr.Default[AppleKeyGeneratorModel](nil, AppleKeyGeneratorModelAttributes),

	// custom providers only
	"description":            stringattr.Optional(),
	"logo":                   stringattr.Optional(),
	"issuer":                 stringattr.Optional(stringattr.URLValidator),
	"authorization_endpoint": stringattr.Optional(stringattr.URLValidator),
	"token_endpoint":         stringattr.Optional(stringattr.URLValidator),
	"user_info_endpoint":     stringattr.Optional(stringattr.URLValidator),
	"jwks_endpoint":          stringattr.Optional(stringattr.URLValidator),
	"use_client_assertion":   boolattr.Default(false),
	"claim_mapping":          strmapattr.Optional(),
}

type OAuthProviderModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`

	Disabled                boolattr.Type                        `tfsdk:"disabled"`
	ClientID                stringattr.Type                      `tfsdk:"client_id"`
	ClientSecret            stringattr.Type                      `tfsdk:"client_secret"`
	ManageProviderTokens    boolattr.Type                        `tfsdk:"manage_provider_tokens"`
	CallbackDomain          stringattr.Type                      `tfsdk:"callback_domain"`
	RedirectURL             stringattr.Type                      `tfsdk:"redirect_url"`
	Prompts                 strlistattr.Type                     `tfsdk:"prompts"`
	Scopes                  strlistattr.Type                     `tfsdk:"scopes"`
	MergeUserAccounts       boolattr.Type                        `tfsdk:"merge_user_accounts"`
	DisableJITUpdates       boolattr.Type                        `tfsdk:"disable_jit_updates"`
	NativeClientID          stringattr.Type                      `tfsdk:"native_client_id"`
	NativeClientSecret      stringattr.Type                      `tfsdk:"native_client_secret"`
	AppleKeyGenerator       objattr.Type[AppleKeyGeneratorModel] `tfsdk:"apple_key_generator"`
	NativeAppleKeyGenerator objattr.Type[AppleKeyGeneratorModel] `tfsdk:"native_apple_key_generator"`
	Description             stringattr.Type                      `tfsdk:"description"`
	Logo                    stringattr.Type                      `tfsdk:"logo"`
	Issuer                  stringattr.Type                      `tfsdk:"issuer"`
	AllowedGrantTypes       strlistattr.Type                     `tfsdk:"allowed_grant_types"`
	AuthorizationEndpoint   stringattr.Type                      `tfsdk:"authorization_endpoint"`
	TokenEndpoint           stringattr.Type                      `tfsdk:"token_endpoint"`
	UserInfoEndpoint        stringattr.Type                      `tfsdk:"user_info_endpoint"`
	JWKsEndpoint            stringattr.Type                      `tfsdk:"jwks_endpoint"`
	UseClientAssertion      boolattr.Type                        `tfsdk:"use_client_assertion"`
	ClaimMapping            strmapattr.Type                      `tfsdk:"claim_mapping"`
}

func (m *OAuthProviderModel) Values(h *helpers.Handler) map[string]any {
	// The plan-time run of these checks skips comparisons involving unknown values, so they run again here where everything is resolved.
	m.Validate(h)

	id := m.ID.ValueString()
	isSystem := slices.Contains(systemProviderNames, id)

	data := map[string]any{
		"id":      id,
		"enabled": !m.Disabled.ValueBool(),
	}
	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	boolattr.Get(m.ManageProviderTokens, data, "manageProviderTokens")
	stringattr.Get(m.CallbackDomain, data, "callbackDomain")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	// the backend patches writes onto the stored entry, so the lists are always sent: an empty list clears
	strlistattr.Get(m.Prompts, data, "prompts", h)
	strlistattr.Get(m.AllowedGrantTypes, data, "allowedGrantTypes", h)
	strlistattr.Get(m.Scopes, data, "scopes", h)
	boolattr.Get(m.MergeUserAccounts, data, "trustProvidedEmails")
	boolattr.Get(m.DisableJITUpdates, data, "jitUpdatesDisabled")
	boolattr.Get(m.UseClientAssertion, data, "useClientAssertion")

	// the configuration fields are Descope-managed on system providers and the backend rejects them as reserved
	if !isSystem {
		stringattr.Get(m.Description, data, "description")
		stringattr.Get(m.Logo, data, "logo")
		stringattr.Get(m.Issuer, data, "issuer")
		stringattr.Get(m.AuthorizationEndpoint, data, "authUrl")
		stringattr.Get(m.TokenEndpoint, data, "tokenUrl")
		stringattr.Get(m.UserInfoEndpoint, data, "userDataUrl")
		stringattr.Get(m.JWKsEndpoint, data, "jwksUrl")

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
	}

	stringattr.Get(m.NativeClientID, data, "nativeClientId")
	stringattr.Get(m.NativeClientSecret, data, "nativeClientSecret")
	// patch semantics on the backend: the null sent for an unset generator clears the stored one, while an absent field would keep it
	objattr.Get(m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	objattr.Get(m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)

	if isSystem {
		// system providers use Descope-managed credentials unless the user supplies their own client_id
		data["useSelfAccount"] = m.ClientID.ValueString() != ""
	} else {
		data["useSelfAccount"] = true
		data["useNonce"] = true
	}
	return data
}

func (m *OAuthProviderModel) SetValues(h *helpers.Handler, data map[string]any) {
	if b, ok := data["enabled"].(bool); ok {
		m.Disabled = boolattr.Value(!b)
	}
	stringattr.Set(&m.ClientID, data, "clientId")
	// client_secret is write-only: never read back from the server, the config value is authoritative.
	boolattr.Set(&m.ManageProviderTokens, data, "manageProviderTokens")
	stringattr.Set(&m.CallbackDomain, data, "callbackDomain")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	strlistattr.Set(&m.Prompts, data, "prompts", h)
	strlistattr.Set(&m.AllowedGrantTypes, data, "allowedGrantTypes", h)
	strlistattr.Set(&m.Scopes, data, "scopes", h)
	boolattr.Set(&m.MergeUserAccounts, data, "trustProvidedEmails")
	boolattr.Set(&m.DisableJITUpdates, data, "jitUpdatesDisabled")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Logo, data, "logo")
	stringattr.Set(&m.Issuer, data, "issuer")
	stringattr.Set(&m.AuthorizationEndpoint, data, "authUrl")
	stringattr.Set(&m.TokenEndpoint, data, "tokenUrl")
	stringattr.Set(&m.UserInfoEndpoint, data, "userDataUrl")
	stringattr.Set(&m.JWKsEndpoint, data, "jwksUrl")
	boolattr.Set(&m.UseClientAssertion, data, "useClientAssertion")
	stringattr.Set(&m.NativeClientID, data, "nativeClientId")
	stringattr.Nil(&m.NativeClientSecret)
	objattr.Set(&m.AppleKeyGenerator, data, "appleKeyGenerator", h)
	objattr.Set(&m.NativeAppleKeyGenerator, data, "nativeAppleKeyGenerator", h)
	strmapattr.Nil(&m.ClaimMapping, h) // empty defaults are added by the backend

	// the backend omits the custom-only fields on system provider reads, so nulls are coerced to empty values to match the state an apply produces
	if id, ok := data["id"].(string); ok && slices.Contains(systemProviderNames, id) {
		for _, s := range []*stringattr.Type{&m.Description, &m.Logo, &m.Issuer, &m.AuthorizationEndpoint, &m.TokenEndpoint, &m.UserInfoEndpoint, &m.JWKsEndpoint} {
			if s.IsNull() {
				*s = stringattr.Value("")
			}
		}
	}

	// the id is never read from the response data, the framework populates it from the request or import id via SetID
}

func (m *OAuthProviderModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OAuthProviderModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *OAuthProviderModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

// Validation

// Runs at plan time so errors are reported before an apply, and so tooling such as tfexport can tell which conditional attributes are missing.
func (m *OAuthProviderModel) Validate(h *helpers.Handler) {
	id := m.ID.ValueString()
	m.validate(h, id, slices.Contains(systemProviderNames, id))
}

func (m *OAuthProviderModel) validate(h *helpers.Handler, name string, isSystem bool) {
	if name == "" {
		return // the Required validator already reports the missing provider
	}
	if isSystem {
		m.validateSystem(h, name)
	} else {
		m.validateCustom(h, name)
	}
}

func (m *OAuthProviderModel) validateCustom(h *helpers.Handler, name string) {
	ensureRequiredField(h, m.ClientID, "client_id", name)
	if !m.UseClientAssertion.ValueBool() {
		ensureRequiredField(h, m.ClientSecret, "client_secret", name)
	}
	ensureRequiredField(h, m.AllowedGrantTypes, "allowed_grant_types", name)
	ensureRequiredField(h, m.AuthorizationEndpoint, "authorization_endpoint", name)
	ensureRequiredField(h, m.TokenEndpoint, "token_endpoint", name)
	ensureRequiredField(h, m.UserInfoEndpoint, "user_info_endpoint", name)

	claimMapping, _ := m.ClaimMapping.ToMap(h.Ctx)
	if _, ok := claimMapping["loginId"]; !ok && len(claimMapping) > 0 {
		h.Error("Invalid Claim Mapping", "Claim mapping set for custom provider %s but 'loginId' was not mapped", name)
	}

	if m.AppleKeyGenerator.IsSet() || m.NativeAppleKeyGenerator.IsSet() {
		h.Error("Reserved Attribute", "The apple_key_generator and native_apple_key_generator attributes are only valid for the apple system provider")
	}
}

func (m *OAuthProviderModel) validateSystem(h *helpers.Handler, name string) {
	ensureReservedField(h, m.Description, "description", name)
	ensureReservedField(h, m.Logo, "logo", name)
	ensureReservedField(h, m.Issuer, "issuer", name)
	ensureReservedField(h, m.AuthorizationEndpoint, "authorization_endpoint", name)
	ensureReservedField(h, m.TokenEndpoint, "token_endpoint", name)
	ensureReservedField(h, m.UserInfoEndpoint, "user_info_endpoint", name)
	ensureReservedField(h, m.JWKsEndpoint, "jwks_endpoint", name)
	if m.UseClientAssertion.ValueBool() {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its use_client_assertion attribute is reserved", name)
	}
	if !m.ClaimMapping.IsEmpty() {
		h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its claim_mapping attribute is reserved", name)
	}
	if name != "apple" && (m.AppleKeyGenerator.IsSet() || m.NativeAppleKeyGenerator.IsSet()) {
		h.Error("Reserved Attribute", "The apple_key_generator and native_apple_key_generator attributes are only valid for the apple system provider")
	}

	if helpers.HasUnknownValues(m.ClientID, m.ClientSecret) {
		return // skip pairing checks while values are still unknown
	}

	ownAccount := m.ClientID.ValueString() != ""
	if ownAccount {
		if name != "apple" {
			if m.ClientSecret.ValueString() == "" {
				h.Missing("The client_id attribute was set for the %s system provider but the client_secret attribute was not set", name)
			}
		} else {
			if m.ClientSecret.ValueString() == "" && !m.AppleKeyGenerator.IsSet() {
				h.Missing("The client_id attribute was set for the apple system provider but the client_secret or apple_key_generator attribute was not set")
			}
			if m.ClientSecret.ValueString() != "" && m.AppleKeyGenerator.IsSet() {
				h.Invalid("The client_secret and the apple_key_generator attributes cannot both be set for the apple system provider")
			}
			if m.NativeClientID.ValueString() != "" {
				if m.NativeClientSecret.ValueString() == "" && !m.NativeAppleKeyGenerator.IsSet() {
					h.Missing("The native_client_id attribute was set for the apple system provider but the native_client_secret or native_apple_key_generator attribute was not set")
				}
				if m.NativeClientSecret.ValueString() != "" && m.NativeAppleKeyGenerator.IsSet() {
					h.Invalid("The native_client_secret and the native_apple_key_generator attributes cannot both be set for the apple system provider")
				}
			}
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

func ensureRequiredField(h *helpers.Handler, field attr.Value, fieldKey, name string) {
	if field.IsUnknown() {
		return // a value that isn't resolved yet can't be judged
	}
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

func ensureReservedField(h *helpers.Handler, field attr.Value, fieldKey, name string) {
	if field.IsUnknown() || field.IsNull() {
		return
	}
	// unset values round-trip as empty strings rather than nulls once the state has been read, so only non-empty values count as set
	if v, ok := field.(stringattr.Type); ok && v.ValueString() == "" {
		return
	}
	h.Error("Reserved Attribute", "The %s OAuth provider is a system provider and its %s attribute is reserved", name, fieldKey)
}

// Apple Key Generator

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
	stringattr.Get(m.KeyID, data, "keyId")
	stringattr.Get(m.TeamID, data, "teamId")
	stringattr.Get(m.PrivateKey, data, "privateKey")
	return data
}

func (m *AppleKeyGeneratorModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.KeyID, data, "keyId")
	stringattr.Set(&m.TeamID, data, "teamId")
	stringattr.Nil(&m.PrivateKey)
}
