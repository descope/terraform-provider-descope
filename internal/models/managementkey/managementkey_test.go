package managementkey_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestManagementKey(t *testing.T) {
	m := testacc.ManagementKey(t)
	p := testacc.Project(t)
	testacc.Run(t,
		// Test empty rebac fails validation
		resource.TestStep{
			Config: m.Config(`
				rebac = {}
			`),
			ExpectError: regexp.MustCompile(`must have at least one role`),
		},
		// Test both company_roles and project_roles fails validation
		resource.TestStep{
			Config: m.Config(`
				rebac = {
					company_roles = ["company-full-access"]
					project_roles = [
						{
							project_ids = ["P123"]
							roles = ["project-full-access"]
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`cannot have both`),
		},
		// Test creating with status = "inactive" fails
		resource.TestStep{
			Config: m.Config(`
				status = "inactive"
				rebac = {
					company_roles = ["company-full-access"]
				}
			`),
			ExpectError: regexp.MustCompile(`Cannot set status`),
		},
		// Test basic creation with company_roles
		resource.TestStep{
			Config: m.Config(`
				rebac = {
					company_roles = ["company-full-access"]
				}
			`),
			Check: m.Check(map[string]any{
				"id":                    testacc.AttributeIsSet,
				"name":                  m.Name,
				"status":                "active",
				"cleartext":             testacc.AttributeIsSet,
				"rebac.company_roles.#": "1",
				"rebac.company_roles.0": "company-full-access",
			}),
		},
		// Test status update (should succeed)
		resource.TestStep{
			Config: m.Config(`
				status = "inactive"
				rebac = {
					company_roles = ["company-full-access"]
				}
			`),
			Check: m.Check(map[string]any{
				"id":                    testacc.AttributeIsSet,
				"name":                  m.Name,
				"status":                "inactive",
				"cleartext":             testacc.AttributeIsSet,
				"rebac.company_roles.#": "1",
				"rebac.company_roles.0": "company-full-access",
			}),
		},
		// Test import
		resource.TestStep{
			ResourceName: m.Path(),
			ImportState:  true,
		},
		// Test basic creation with company_roles
		resource.TestStep{
			PreConfig: func() {
				m.Name += "bar"
			},
			Config: m.Config(`
				rebac = {
					company_roles = ["company-full-access"]
				}
				status = "inactive"
			`),
			Check: m.Check(map[string]any{
				"id":                    testacc.AttributeIsSet,
				"name":                  m.Name,
				"status":                "inactive",
				"cleartext":             testacc.AttributeIsSet,
				"rebac.company_roles.#": "1",
				"rebac.company_roles.0": "company-full-access",
			}),
		},
		// Test with permitted_ips
		resource.TestStep{
			Config: m.Config(`
				description = "With permitted IPs"
				permitted_ips = ["192.168.1.0/24", "10.0.0.1"]
				rebac = {
					company_roles = ["company-full-access"]
				}
			`),
			Check: m.Check(map[string]any{
				"description":           "With permitted IPs",
				"permitted_ips.#":       "2",
				"permitted_ips.0":       "192.168.1.0/24",
				"permitted_ips.1":       "10.0.0.1",
				"rebac.company_roles.#": "1",
				"rebac.company_roles.0": "company-full-access",
			}),
		},
		// Destroy resource
		resource.TestStep{
			Config: m.Config(`
				rebac = {
					company_roles = ["company-full-access"]
				}
			`),
			Destroy: true,
		},
		// Test with project_roles and expire_time
		resource.TestStep{
			Config: p.Config() + m.Config(`
				description = "With project roles"
				status = "active"
				expire_time = 1893456000
				rebac = {
					project_roles = [
						{
							project_ids = [`+p.Path()+`.id]
							roles = ["project-full-access"]
						}
					]
				}
			`),
			Check: m.Check(map[string]any{
				"description": "With project roles",
				"status":      "active",
				"expire_time": "1893456000",
				"rebac.project_roles": map[string]any{
					"#":               "1",
					"0.project_ids.#": "1",
					"0.roles.#":       "1",
					"0.roles.0":       "project-full-access",
				},
			}),
		},
		// Destroy resource
		resource.TestStep{
			Config: p.Config() + m.Config(`
				expire_time = 1893456000
				rebac = {
					project_roles = [
						{
							project_ids = [`+p.Path()+`.id]
							roles = ["project-full-access"]
						}
					]
				}
			`),
			Destroy: true,
		},
		// Test with tag_roles
		resource.TestStep{
			Config: m.Config(`
				description = "With tag roles"
				rebac = {
					tag_roles = [
						{
							tags = ["production", "staging"]
							roles = ["tag-infra-read-write"]
						}
					]
				}
			`),
			Check: m.Check(map[string]any{
				"description": "With tag roles",
				"rebac.tag_roles": map[string]any{
					"#":         "1",
					"0.tags.#":  "2",
					"0.tags.0":  "production",
					"0.tags.1":  "staging",
					"0.roles.#": "1",
					"0.roles.0": "tag-infra-read-write",
				},
			}),
		},
	)
}
