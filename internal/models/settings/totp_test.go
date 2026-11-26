package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTOTPSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.TOTPSettings(t)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":            testacc.AttributeIsSet,
				"project_id":    testacc.AttributeIsSet,
				"disabled":      false,
				"service_label": "",
				"account_label": "",
			}),
		},
		// set every field
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				disabled = true
				service_label = "MyService"
				account_label = "{{user.email}}"
			`),
			Check: o.Check(map[string]any{
				"disabled":      true,
				"service_label": "MyService",
				"account_label": "{{user.email}}",
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
