package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEmailVerificationSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.EmailVerificationSettings(t)
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
			}),
		},
		// enable
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				enabled = true
			`),
			Check: o.Check(map[string]any{
				"enabled": true,
			}),
		},
		// import using <project_id>
		resource.TestStep{
			ResourceName:      o.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id"),
		},
		// disable
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"enabled": false,
			}),
		},
	)
}
