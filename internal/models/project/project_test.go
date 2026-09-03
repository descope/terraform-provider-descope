package project_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectDeletionProtection(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				environment = "production"
			`),
			Check: p.Check(map[string]any{
				"deletion_protection": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = "production"
			`),
			Destroy:     true,
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = ""
				deletion_protection = true
			`),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = ""
				deletion_protection = true
			`),
			Destroy:     true,
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = "production"
				deletion_protection = false
			`),
			Check: p.Check(map[string]any{
				"environment":         "production",
				"deletion_protection": "false",
			}),
		},
	)
}

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
	)
}
