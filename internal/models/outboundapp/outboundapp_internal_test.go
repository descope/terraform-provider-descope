package outboundapp

import (
	"context"
	"slices"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestOutboundAppClientSecretPresence(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &OutboundAppModel{
		Name:                   stringattr.Value("app"),
		AuthorizationURLParams: strmapattr.Empty(),
		TokenURLParams:         strmapattr.Empty(),
		DefaultScopes:          strsetattr.Empty(),
		Prompt:                 strsetattr.Empty(),
	}

	if v, ok := m.Values(h)["clientSecret"]; ok {
		t.Errorf("an unset client_secret must be omitted so the stored secret is kept, got %v", v)
	}

	m.ClientSecret = stringattr.Value("real-secret")
	if m.Values(h)["clientSecret"] != "real-secret" {
		t.Errorf("a configured client_secret must be sent as is, got %v", m.Values(h)["clientSecret"])
	}
}

func TestOutboundAppURLParamOrder(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &OutboundAppModel{
		Name:                   stringattr.Value("app"),
		AuthorizationURLParams: strmapattr.Value(map[string]string{"c": "3", "a": "1", "b": "2"}),
		TokenURLParams:         strmapattr.Empty(),
		DefaultScopes:          strsetattr.Empty(),
		Prompt:                 strsetattr.Empty(),
	}

	for range 10 {
		entries, ok := m.Values(h)["authorizationUrlParams"].([]any)
		if !ok {
			t.Fatalf("expected an array of url params, got %T", m.Values(h)["authorizationUrlParams"])
		}
		keys := []string{}
		for _, entry := range entries {
			pair, ok := entry.(map[string]any)
			if !ok {
				t.Fatalf("expected a key and value pair, got %T", entry)
			}
			key, ok := pair["key"].(string)
			if !ok {
				t.Fatalf("expected a string key, got %T", pair["key"])
			}
			keys = append(keys, key)
		}
		if !slices.IsSorted(keys) {
			t.Fatalf("url params must serialize in a stable order or every plan shows a phantom diff, got %v", keys)
		}
	}
}
