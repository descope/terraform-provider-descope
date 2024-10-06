package models_test

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
			}),
		},
		resource.TestStep{
			ResourceName:      p.Path(),
			ImportState:       true,
			ImportStateVerify: true,
		},
		resource.TestStep{
			PreConfig: func() {
				p.Name += "bar"
			},
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"name": p.Name,
			}),
		},
	)
}
