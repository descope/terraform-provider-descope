package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The testdata fixture is a default applications portal widget export with the widgetId key removed, as the widget_id attribute is authoritative.

func TestAdminPortal(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.AdminPortal(t)
	w := testacc.Widget(t)
	widget := w.Block(`
		project_id = "` + projectID + `"
		widget_id = "testacc-adminportal"
		data = ` + testacc.FixtureJSON(t, "testdata/widget.json"),
	)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"enabled":    false,
				"style_id":   "",
				"widgets.#":  0,
			}),
		},
		// enabling the portal requires at least one widget
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				enabled = true
			`),
			ExpectError: regexp.MustCompile(`admin_portal must have at least one widget when enabled`),
		},
		// enable with a widget managed in the same configuration
		resource.TestStep{
			Config: widget + o.Block(`
				project_id = "`+projectID+`"
				enabled = true
				style_id = "theme"
				widgets = [{ widget_id = `+w.Path()+`.widget_id, type = "applications-portal-widget" }]
			`),
			Check: o.Check(map[string]any{
				"enabled":             true,
				"style_id":            "theme",
				"widgets.#":           1,
				"widgets.0.widget_id": "testacc-adminportal",
				"widgets.0.type":      "applications-portal-widget",
			}),
		},
		// update the widget list
		resource.TestStep{
			Config: widget + o.Block(`
				project_id = "`+projectID+`"
				enabled = true
				widgets = [
					{ widget_id = `+w.Path()+`.widget_id, type = "applications-portal-widget" },
					{ widget_id = "user-management", type = "user-management" },
				]
			`),
			Check: o.Check(map[string]any{
				"enabled":             true,
				"style_id":            "",
				"widgets.#":           2,
				"widgets.1.widget_id": "user-management",
			}),
		},
		// import using <project_id>/<id> (id == project_id for this singleton)
		resource.TestStep{
			ResourceName:      o.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id"),
		},
		// restore the shared project to its default disabled state
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"enabled":   false,
				"widgets.#": 0,
			}),
		},
	)
}
