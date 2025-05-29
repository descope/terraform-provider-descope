package attributes_test

import (
	"regexp"
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
				"attributes.user": map[string]any{
					"#":      2,
					"0.name": "foo",
					"0.type": "string",
					"1.name": "bar",
					"1.type": "number",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				attributes = {
					user = [
						{
							name = var.quxname
							type = "string"
						}
					]
				}
			`) + p.Variables(`
				variable "quxname" {
					type    = string
					default = "qux"
				}
			`),
			Check: p.Check(map[string]any{
				"attributes.user": map[string]any{
					"#":      1,
					"0.name": "qux",
					"0.type": "string",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				attributes = {
					user = var.user_attributes
				}
			`) + p.Variables(`
				variable "user_attributes" {
					type = list(object({
						name = string
						type = string
					}))
					default = [
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
				"attributes.user": map[string]any{
					"#":      2,
					"0.name": "foo",
					"0.type": "string",
					"1.name": "bar",
					"1.type": "number",
				},
			}),
		},
		resource.TestStep{
			// Not entirely sure why this simple substitution fails when it's not wrapped in a conditional
			Config: p.Config(`
				attributes = {
					user = [ var.bar ]
				}
			`) + p.Variables(`
				variable "bar" {
					type = object({
						name = string
						type = string
					})
					default = {
						name = "bar"
						type = "number"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Missing Configuration for Required Attribute`),
		},
		resource.TestStep{
			// Same as above, but it works when wrapped in a conditional
			Config: p.Config(`
				attributes = {
					user = var.bar != null ? [var.bar] : []
				}
			`) + p.Variables(`
				variable "bar" {
					type = object({
						name = string
						type = string
					})
					default = {
						name = "bar"
						type = "number"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"attributes.user": map[string]any{
					"#":      1,
					"0.name": "bar",
					"0.type": "number",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				attributes = {}
			`),
			Check: p.Check(map[string]any{
				"attributes.user.#": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				attributes = {
					tenant = var.foo != null ? [var.foo] : []
				}
			`) + p.Variables(`
				variable "foo" {
					default = {
						name = "bar"
						type = "number"
						select_options = ["x", "y", "z"]
						authorization = {
         					view_permissions = ["foo", "bar"]
        				}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"attributes.tenant": map[string]any{
					"#":                1,
					"0.name":           "bar",
					"0.type":           "number",
					"0.select_options": []string{"x", "y", "z"},
					"0.authorization": map[string]any{
						"view_permissions": []string{"foo", "bar"},
					},
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"attributes.user.#": 0,
			}),
		},
	)
}
