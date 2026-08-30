package widget_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWidget(t *testing.T) {
	projectID := testacc.ProjectID(t)
	w := testacc.Widget(t)
	testacc.Run(t,
		// create the widget from its exported representation
		resource.TestStep{
			Config: w.Block(`
				project_id = "` + projectID + `"
				widget_id = "testacc-portal"
				data = ` + testacc.FixtureJSON(t, "testdata/widget.json"),
			),
			Check: w.Check(map[string]any{
				"id":         "testacc-portal",
				"project_id": testacc.AttributeIsSet,
				"widget_id":  "testacc-portal",
				"data":       testacc.AttributeIsSet,
			}),
		},
		// changing the widget data updates the widget in place
		resource.TestStep{
			Config: w.Block(`
				project_id = "` + projectID + `"
				widget_id = "testacc-portal"
				data = ` + testacc.FixtureJSON(t, "testdata/widget.json", "users view all their associated applications", "users view all their applications"),
			),
			Check: w.Check(map[string]any{
				"id":        "testacc-portal",
				"widget_id": "testacc-portal",
			}),
		},
		// data is re-serialized server-side on read, so verify the identity attributes only
		resource.TestStep{
			ResourceName:            w.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(w.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"data"},
		},
	)
}
