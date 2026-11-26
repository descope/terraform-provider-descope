package settings

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

func TestProjectSettingsAppsCustomDomain(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &ProjectSettingsModel{AppsCustomDomain: stringattr.Value("apps.acme.com"), ApprovedDomains: strsetattr.Value([]string{})}
	assert.Equal(t, "apps.acme.com", m.Values(h)["appsCustomDomain"])

	m.AppsCustomDomain = stringattr.Value("")
	data := m.Values(h)
	assert.Contains(t, data, "appsCustomDomain")
	assert.Equal(t, "", data["appsCustomDomain"])

	m = &ProjectSettingsModel{}
	m.SetValues(h, map[string]any{"appsCustomDomain": "apps.acme.com"})
	assert.Equal(t, "apps.acme.com", m.AppsCustomDomain.ValueString())

	m = &ProjectSettingsModel{}
	m.SetValues(h, map[string]any{})
	assert.False(t, m.AppsCustomDomain.IsNull())
	assert.Equal(t, "", m.AppsCustomDomain.ValueString())
	assert.False(t, diags.HasError())

	m = &ProjectSettingsModel{AppsCustomDomain: stringattr.Value("https://apps.acme.com")}
	m.Validate(h)
	assert.True(t, diags.HasError())
}
