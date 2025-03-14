package attributes_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAttributes(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				attributes = {
					user = [
						{
							name = "foo"
							type = "string"
						},
						{
							name = "bar"
							type = "number"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"attributes.user.#":      2,
				"attributes.user.0.name": "foo",
				"attributes.user.0.type": "string",
				"attributes.user.1.name": "bar",
				"attributes.user.1.type": "number",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				attributes = {
					user = [
						{
							name = "bar"
							type = "string"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"attributes.user.#":      1,
				"attributes.user.0.name": "bar",
				"attributes.user.0.type": "number",
			}),
		},
	)
}
