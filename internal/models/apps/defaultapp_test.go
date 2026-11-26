package apps_test

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/apps"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestDefaultOIDCAppChecksWhenImportedWithoutID(t *testing.T) {
	for _, tc := range []struct {
		name    string
		stateID string
		appName string
		disable bool
		invalid string
	}{
		{name: "RenamedDefaultApp", stateID: apps.DefaultOIDCAppID, appName: "Renamed", invalid: "The name of the built-in default OIDC application cannot be changed"},
		{name: "DisabledDefaultApp", stateID: apps.DefaultOIDCAppID, appName: "OIDC default application", disable: true, invalid: "The built-in default OIDC application cannot be disabled"},
		{name: "UnchangedDefaultApp", stateID: apps.DefaultOIDCAppID, appName: "OIDC default application"},
		{name: "OtherApp", stateID: "SA123", appName: "Renamed", disable: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := &apps.OIDCAppModel{
				ID:          stringattr.Value(tc.stateID),
				Name:        stringattr.Value(tc.appName),
				Description: stringattr.Value("Default OIDC APP"),
				Disabled:    boolattr.Value(tc.disable),
			}
			config := &apps.OIDCAppModel{ID: types.StringNull(), Name: plan.Name, Description: plan.Description, Disabled: plan.Disabled}
			state := &apps.OIDCAppModel{ID: stringattr.Value(tc.stateID)}

			var diags diag.Diagnostics
			plan.ModifyPlan(helpers.NewHandler(context.Background(), &diags), config, state)

			if tc.invalid == "" {
				require.False(t, diags.HasError(), "unexpected errors: %v", diags)
				return
			}
			require.True(t, diags.HasError(), "expected an error containing %q", tc.invalid)
			require.Contains(t, diags.Errors()[0].Detail(), tc.invalid)
		})
	}
}
