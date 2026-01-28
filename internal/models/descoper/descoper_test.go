package descoper_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDescoper(t *testing.T) {
	email := testacc.GenerateAlias(t) + "@kljafshjlkafhjkafs.com"
	d := testacc.Descoper(t)
	p := testacc.Project(t)
	testacc.Run(t,
		// Test empty rbac fails validation
		resource.TestStep{
			Config: d.Config(`
				email = "` + email + `"
				rbac = {}
			`),
			ExpectError: regexp.MustCompile(`must have is_company_admin`),
		},
		// Test both is_company_admin and project_roles fails validation
		resource.TestStep{
			Config: d.Config(`
				email = "` + email + `"
				rbac = {
					is_company_admin = true
					project_roles = [
						{
							project_ids = ["P123"]
							role = "admin"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`cannot have both`),
		},
		// Test both is_company_admin and tag_roles fails validation
		resource.TestStep{
			Config: d.Config(`
				email = "` + email + `"
				rbac = {
					is_company_admin = true
					tag_roles = [
						{
							tags = ["production"]
							role = "developer"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`cannot have both`),
		},
		// Test basic creation with is_company_admin
		resource.TestStep{
			Config: d.Config(`
				email = "` + email + `"
				rbac = {
					is_company_admin = true
				}
			`),
			Check: d.Check(map[string]any{
				"id":                    testacc.AttributeIsSet,
				"email":                 email,
				"name":                  d.Name,
				"rbac.is_company_admin": "true",
				"rbac.tag_roles.#":      "0",
				"rbac.project_roles.#":  "0",
			}),
		},
		// Test import
		resource.TestStep{
			ResourceName: d.Path(),
			ImportState:  true,
		},
		// Test name update
		resource.TestStep{
			PreConfig: func() {
				d.Name += "bar"
			},
			Config: d.Config(`
				email = "` + email + `"
				rbac = {
					is_company_admin = true
				}
			`),
			Check: d.Check(map[string]any{
				"id":                    testacc.AttributeIsSet,
				"email":                 email,
				"name":                  d.Name,
				"rbac.is_company_admin": "true",
			}),
		},
		// Destroy resource
		resource.TestStep{
			Config: d.Config(`
				email = "` + email + `"
				rbac = {
					is_company_admin = true
				}
			`),
			Destroy: true,
		},
		// Test with tag_roles AAAAAAA
		resource.TestStep{
			Config: p.Config(`
				tags = ["production", "staging"]
			`) + d.Config(`
				email = "`+email+`"
				rbac = {
					tag_roles = [
						{
							tags = ["production", "staging"]
							role = "admin"
						}
					]
				}
			`),
			Check: d.Check(map[string]any{
				"id":                      testacc.AttributeIsSet,
				"email":                   email,
				"name":                    d.Name,
				"rbac.is_company_admin":   "false",
				"rbac.tag_roles.#":        "1",
				"rbac.tag_roles.0.tags.#": "2",
				"rbac.tag_roles.0.role":   "admin",
			}),
		},
		// Destroy resource
		resource.TestStep{
			Config: p.Config(`
				tags = ["production", "staging"]
			`) + d.Config(`
				email = "`+email+`"
				rbac = {
					tag_roles = [
						{
							tags = ["production"]
							role = "admin"
						}
					]
				}
			`),
			Destroy: true,
		},
		// Test with project_roles
		resource.TestStep{
			Config: p.Config(`
				tags = ["production", "staging"]
			`) + d.Config(`
				email = "`+email+`"
				rbac = {
					project_roles = [
						{
							project_ids = [`+p.Path()+`.id]
							role = "developer"
						}
					]
				}
			`),
			Check: d.Check(map[string]any{
				"id":                                 testacc.AttributeIsSet,
				"email":                              email,
				"name":                               d.Name,
				"rbac.is_company_admin":              "false",
				"rbac.project_roles.#":               "1",
				"rbac.project_roles.0.project_ids.#": "1",
				"rbac.project_roles.0.role":          "developer",
			}),
		},
		// Test with multiple project_roles
		resource.TestStep{
			Config: p.Config(`
				tags = ["production", "staging"]
			`) + d.Config(`
				email = "`+email+`"
				rbac = {
					project_roles = [
						{
							project_ids = [`+p.Path()+`.id]
							role = "developer"
						},
						{
							project_ids = [`+p.Path()+`.id]
							role = "support"
						}
					]
				}
			`),
			Check: d.Check(map[string]any{
				"rbac.project_roles.#": "2",
			}),
		},
	)
}
