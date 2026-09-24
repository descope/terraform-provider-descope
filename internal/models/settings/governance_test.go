package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGovernance(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.Governance(t)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":             testacc.AttributeIsSet,
				"project_id":     testacc.AttributeIsSet,
				"configured":     false,
				"auto_approval":  false,
				"suite_disabled": false,
				"logo":           "",
			}),
		},
		// set every field, including the suite kill switch
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				configured = true
				auto_approval = true
				suite_disabled = true
				logo = "data:image/png;base64,iVBORw0KGgo="
			`),
			Check: o.Check(map[string]any{
				"configured":     true,
				"auto_approval":  true,
				"suite_disabled": true,
				"logo":           "data:image/png;base64,iVBORw0KGgo=",
			}),
		},
		// a false has to stick: one that fails to round-trip reads as never set,
		// leaving the previous true in place and re-enabling a disabled suite
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				configured = true
				auto_approval = false
				suite_disabled = false
				logo = ""
			`),
			Check: o.Check(map[string]any{
				"configured":     true,
				"auto_approval":  false,
				"suite_disabled": false,
				"logo":           "",
			}),
		},
		// import using <project_id>/<id> (id == project_id for this singleton)
		resource.TestStep{
			ResourceName:      o.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id"),
		},
	)
}
