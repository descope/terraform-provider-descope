package adminportal_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAdminPortal(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"admin_portal": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				admin_portal = {
					enabled = true
				}
			`),
			ExpectError: regexp.MustCompile(`admin_portal must have at least one widget when enabled`),
		},
		resource.TestStep{
			Config: p.Config(`
				admin_portal = {
					enabled  = true
					widgets = [
						{
							widget_id = "w1"
							type      = "users"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"admin_portal.enabled":   true,
				"admin_portal.widgets.#": 1,
				"admin_portal.widgets.0": map[string]any{"widget_id": "w1", "type": "users"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				admin_portal = {
					enabled  = true
					style_id = "style-1"
					widgets = [
						{
							widget_id = "w1"
							type      = "users"
						},
						{
							widget_id = "w2"
							type      = "roles"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"admin_portal.enabled":   true,
				"admin_portal.style_id":  "style-1",
				"admin_portal.widgets.#": 2,
				"admin_portal.widgets.0": map[string]any{"widget_id": "w1", "type": "users"},
				"admin_portal.widgets.1": map[string]any{"widget_id": "w2", "type": "roles"},
			}),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"admin_portal": testacc.AttributeIsNotSet,
			}),
		},
	)
}
