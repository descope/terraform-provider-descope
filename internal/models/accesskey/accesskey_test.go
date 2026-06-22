package accesskey_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccessKey(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.AccessKey(t)
	testacc.Run(t,
		// Test basic creation with required fields only
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: a.Check(map[string]any{
				"id":              testacc.AttributeIsSet,
				"project_id":      testacc.AttributeIsSet,
				"name":            a.Name,
				"description":     "",
				"status":          "active",
				"expire_time":     "0",
				"roles.#":         "0",
				"tenants.#":       "0",
				"permitted_ips.#": "0",
				"client_id":       testacc.AttributeIsSet,
				"cleartext":       testacc.AttributeIsSet,
			}),
		},
		// Test update of mutable fields
		resource.TestStep{
			Config: p.Config(`
				authorization = {
					roles = [
						{ name = "Viewer" }
					]
				}
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				description = "Updated description"
				permitted_ips = ["10.0.0.0/8"]
				roles = ["Viewer"]
			`),
			Check: a.Check(map[string]any{
				"description":     "Updated description",
				"permitted_ips.#": "1",
				"permitted_ips.0": "10.0.0.0/8",
				"roles.#":         "1",
				"roles.0":         "Viewer",
			}),
		},
		// Test tenants attribute with non-existent tenant
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				tenants = [
					{
						tenant_id = "T2foo"
						roles = ["Viewer", "Editor"]
					},
				]
			`),
			ExpectError: regexp.MustCompile(`Tenant .* does not exist in project`),
		},
		// Test roles attribute with non-existent role
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				roles = ["Quux"]
			`),
			ExpectError: regexp.MustCompile(`Role .* does not exist in project`),
		},
		// Test custom claims as JSON string
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				custom_claims = jsonencode({ env = "staging" })
			`),
			Check: a.Check(map[string]any{
				"custom_claims": testacc.AttributeIsSet,
			}),
		},
		// Test deactivate via status update
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				status = "inactive"
			`),
			Check: a.Check(map[string]any{
				"status": "inactive",
			}),
		},
		// Test reactivate
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				status = "active"
			`),
			Check: a.Check(map[string]any{
				"status": "active",
			}),
		},
		// Test expire_time triggers replacement (RequiresReplace)
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				expire_time = 1924991999
			`),
			Check: a.Check(map[string]any{
				"expire_time": "1924991999",
				"cleartext":   testacc.AttributeIsSet,
			}),
		},
		// Test import with composite ID
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		},
		// Destroy resource
		resource.TestStep{
			Config:  p.Config() + a.Config(`project_id = `+p.Path()+`.id`),
			Destroy: true,
		},
	)
}
