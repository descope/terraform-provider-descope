package tenant_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestTenant(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test requires TF_ACC=1")
	}

	projectID := os.Getenv("DESCOPE_PROJECT_ID")
	require.NotEmpty(t, projectID, "DESCOPE_PROJECT_ID must be set for tenant acceptance tests")
	tenant := testacc.Tenant(t)

	testacc.Run(t,
		resource.TestStep{
			Config: tenant.Config(`
				project_id = "` + projectID + `"
				role_inheritance = "invalid"
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: tenant.Config(`
				project_id = "` + projectID + `"
				self_provisioning_domains = ["example.com"]
			`),
			Check: tenant.Check(map[string]any{
				"id":                          testacc.AttributeIsSet,
				"project_id":                  projectID,
				"name":                        tenant.Name,
				"self_provisioning_domains.#": "1",
				"self_provisioning_domains.0": "example.com",
				"disabled":                    false,
				"enforce_sso":                 false,
				"enforce_sso_exclusions.#":    "0",
				"federated_application_ids.#": "0",
			}),
		},
		resource.TestStep{
			ResourceName:        tenant.Path(),
			ImportState:         true,
			ImportStateIdPrefix: projectID + "/",
			ImportStateVerify:   true,
		},
		resource.TestStep{
			PreConfig: func() {
				tenant.Name += " updated"
			},
			Config: tenant.Config(`
				project_id = "` + projectID + `"
				disabled = true
				enforce_sso = true
			`),
			Check: tenant.Check(map[string]any{
				"name":        tenant.Name,
				"disabled":    true,
				"enforce_sso": true,
			}),
		},
	)
}
