package settings

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// The token is write-only, so the request has to tell apart leaving a token that was configured out of
// band alone, which omits the field, from clearing it, which sends an empty one.
func TestSessionMigrationAPITokenPresence(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &SessionMigrationSettingsModel{
		Vendor:                   stringattr.Value("okta"),
		ClientID:                 stringattr.Value("cid"),
		Issuer:                   stringattr.Value("https://dev.okta.com"),
		LoginIDMatchedAttributes: strsetattr.Value([]string{"email"}),
	}

	migration, _ := m.Values(h)["sessionMigration"].(map[string]any)
	if v, ok := migration["apiToken"]; ok {
		t.Errorf("an unset api_token must be omitted so the stored token is kept, got %v", v)
	}

	m.APIToken = stringattr.Value("real-token")
	migration, _ = m.Values(h)["sessionMigration"].(map[string]any)
	if migration["apiToken"] != "real-token" {
		t.Errorf("a configured api_token must be sent as is, got %v", migration["apiToken"])
	}

	m.APIToken = stringattr.Value("")
	migration, _ = m.Values(h)["sessionMigration"].(map[string]any)
	if v, ok := migration["apiToken"]; !ok || v != "" {
		t.Errorf("an empty api_token must be sent so the stored token is cleared, got %v", v)
	}
}
