package authentication

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var OAuthValidator = objattr.NewValidator[OAuthModel]("must have a valid OAuth configuration")

var OAuthAttributes = map[string]schema.Attribute{
	"disabled": boolattr.Default(false),
	"system":   objattr.Default[OAuthSystemProvidersModel](nil, OAuthSystemProviderAttributes),
	"custom":   mapattr.Default[OAuthProviderCustomModel](nil, OAuthProviderCustomAttributes, OAuthProviderCustomValidator),
}

type OAuthModel struct {
	Disabled boolattr.Type                           `tfsdk:"disabled"`
	System   objattr.Type[OAuthSystemProvidersModel] `tfsdk:"system"`
	Custom   mapattr.Type[OAuthProviderCustomModel]  `tfsdk:"custom"`
}

func (m *OAuthModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")

	// collect all system and custom provider objects into a single map
	providers := map[string]any{}
	objattr.Get(m.System, providers, helpers.RootKey, h)
	for name, provider := range mapattr.Iterator(m.Custom, h) {
		providers[name] = provider.Values(h)
	}

	// add the name explicitly to each provider object as expected by the backend
	for name, v := range providers {
		if obj, ok := v.(map[string]any); ok {
			obj["name"] = name
		}
	}

	data["providerSettings"] = providers
	return data
}

func (m *OAuthModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")

	system := map[string]any{}
	custom := map[string]any{}

	// split the single list of providers from the backend into system and custom
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

func (m *OAuthModel) Validate(h *helpers.Handler) {
	for name := range mapattr.Iterator(m.Custom, h) {
		if slices.Contains(systemProviderNames, name) {
			h.Error("Reserved OAuth Provider Name", "The %s name is reserved for system providers and cannot be used for a custom provider", name)
		}
	}
}
