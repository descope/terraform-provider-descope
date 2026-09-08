package outboundapp

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestOutboundAppClientSecretPresence(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &OutboundAppModel{
		Name:                   stringattr.Value("app"),
		AuthorizationURLParams: listattr.Empty[URLParamModel](),
		TokenURLParams:         listattr.Empty[URLParamModel](),
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
