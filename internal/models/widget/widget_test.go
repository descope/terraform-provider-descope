package widget_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		// the widget's flows are embedded in its data when read back
		resource.TestStep{
			ResourceName:      w.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(w.Path(), "project_id", "id"),
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 {
					return fmt.Errorf("expected a single imported state, got %d", len(states))
				}
				var data struct {
					Metadata struct {
						FlowsIDs []string `json:"flowsIds"`
					} `json:"metadata"`
					Flows []struct {
						FlowID string `json:"flowId"`
					} `json:"flows"`
				}
				if err := json.Unmarshal([]byte(states[0].Attributes["data"]), &data); err != nil {
					return fmt.Errorf("failed to parse the imported widget data: %w", err)
				}
				if !slices.Contains(data.Metadata.FlowsIDs, "testacc-portal-delete-name") {
					return fmt.Errorf("expected the widget flow in flowsIds, got %v", data.Metadata.FlowsIDs)
				}
				var flowIDs []string
				for _, flow := range data.Flows {
					flowIDs = append(flowIDs, flow.FlowID)
				}
				if !slices.Contains(flowIDs, "testacc-portal-delete-name") {
					return fmt.Errorf("expected the widget flow embedded in the imported data, got %v", flowIDs)
				}
				return nil
			},
		},
	)
}
