package connectors_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Generated tests pick each boolean's value by hashing the field name, so these pin the toggle-off case for the three shapes that misfire.
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
