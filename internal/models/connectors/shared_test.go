package connectors_test

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/connectors"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestConnectorDependencyValidation(t *testing.T) {
	projectID := testacc.ProjectID(t)

	t.Run("non-zero default", func(t *testing.T) {
		c := testacc.NewResource(t, "hcaptcha_connector")
		testacc.Run(t,
			resource.TestStep{
				Config: c.Config(`
					project_id = "` + projectID + `"
					site_key = "ikzbbly"
					secret_key = "wi4bhwt7a"
				`),
				Check: c.Check(map[string]any{
					"override_assessment": false,
					"assessment_score":    "0.5",
				}),
			},
		)
	})

	t.Run("zero default", func(t *testing.T) {
		c := testacc.NewResource(t, "fingerprint_connector")
		testacc.Run(t,
			resource.TestStep{
				Config: c.Config(`
					project_id = "` + projectID + `"
					public_api_key = "htt624yz4z6i"
					secret_api_key = "qxt75gbg4234"
				`),
				Check: c.Check(map[string]any{
					"use_cloudflare_integration": false,
					"cloudflare_script_url":      "",
				}),
			},
		)
	})

	t.Run("empty list default", func(t *testing.T) {
		c := testacc.NewResource(t, "mixpanel_connector")
		testacc.Run(t,
			resource.TestStep{
				Config: c.Config(`
					project_id = "` + projectID + `"
					project_token = "inazr4ilpcxv"
					api_secret = "hgg666mus"
					config_project_id = "yhw7b6yel"
					service_account_username = "26hhhmmzsm"
					service_account_secret = "xqbxbtxf"
					audit_enabled = false
				`),
				Check: c.Check(map[string]any{
					"audit_enabled": false,
				}),
			},
		)
	})
}

func TestConnectorDependencyUnknown(t *testing.T) {
	validate := func(t *testing.T, m interface{ Validate(*helpers.Handler) }) diag.Diagnostics {
		t.Helper()
		var diags diag.Diagnostics
		m.Validate(helpers.NewHandler(context.Background(), &diags))
		return diags
	}

	t.Run("dependency wants non-default", func(t *testing.T) {
		m := &connectors.HTTPConnectorModel{
			ClientCertificate: stringattr.Value("cert"),
			ClientKey:         stringattr.Value("key"),
		}

		m.UseMTLS = boolattr.Value(false)
		require.True(t, validate(t, m).HasError())

		m.UseMTLS = types.BoolUnknown()
		require.Empty(t, validate(t, m))
	})

	t.Run("dependency wants default", func(t *testing.T) {
		m := &connectors.LDAPConnectorModel{
			ClientCertificate: stringattr.Value("cert"),
			ClientKey:         stringattr.Value("key"),
			BindDN:            stringattr.Value("cn=admin"),
			BindPassword:      stringattr.Value("secret"),
		}

		m.UseMTLS = boolattr.Value(true)
		require.True(t, validate(t, m).HasError())

		m.UseMTLS = types.BoolUnknown()
		require.Empty(t, validate(t, m))
	})
}
