package project_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProject(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				environment = "foo"
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = "production"
			`),
			Check: p.Check(map[string]any{
				"id":          testacc.AttributeIsSet,
				"name":        p.Name,
				"environment": "production",
				"tags":        []string{},
			}),
		},
		resource.TestStep{
			ResourceName: p.Path(),
			ImportState:  true,
		},
		resource.TestStep{
			PreConfig: func() {
				p.Name += "bar"
			},
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"name":        p.Name,
				"environment": "production",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = ""
				tags = ["foo", "bar"]
			`),
			Check: p.Check(map[string]any{
				"name":        p.Name,
				"tags":        []string{"foo", "bar"},
				"environment": "",
			}),
		},
		// Project-wide scope claim mapping (no use_project_mapping / mandatory at this level)
		resource.TestStep{
			Config: p.Config(`
				scope_claim_mapping = [
					{
						scope = "profile"
						description = "Profile info"
						claims = { name = "{{user.name}}" }
					},
					{
						scope = "address"
						description = "Address info"
					},
				]
			`),
			Check: p.Check(map[string]any{
				"scope_claim_mapping.#": 2,
				"scope_claim_mapping.0": map[string]any{
					"scope":       "profile",
					"description": "Profile info",
					"claims":      map[string]any{"name": "{{user.name}}"},
				},
				"scope_claim_mapping.1": map[string]any{
					"scope":       "address",
					"description": "Address info",
				},
			}),
		},
	)
}
