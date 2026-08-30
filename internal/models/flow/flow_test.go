package flow_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFlow(t *testing.T) {
	projectID := testacc.ProjectID(t)
	f := testacc.Flow(t)
	testacc.Run(t,
		resource.TestStep{
			Config: f.Block(`
				project_id = "` + projectID + `"
				flow_id = "testacc-step-up"
				data = ` + testacc.FixtureJSON(t, "testdata/flow.json"),
			),
			Check: f.Check(map[string]any{
				"id":         "testacc-step-up",
				"project_id": testacc.AttributeIsSet,
				"flow_id":    "testacc-step-up",
				"data":       testacc.AttributeIsSet,
			}),
		},
		resource.TestStep{
			Config: f.Block(`
				project_id = "` + projectID + `"
				flow_id = "testacc-step-up"
				data = ` + testacc.FixtureJSON(t, "testdata/flow.json", "used for extra user validation", "used for extra user validation in tests"),
			),
			Check: f.Check(map[string]any{
				"id":      "testacc-step-up",
				"flow_id": "testacc-step-up",
			}),
		},
		resource.TestStep{
			ResourceName:            f.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(f.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"data"},
		},
	)
}
