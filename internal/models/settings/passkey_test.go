package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPasskeySettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.PasskeySettings(t)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":           testacc.AttributeIsSet,
				"project_id":   testacc.AttributeIsSet,
				"disabled":     false,
				"display_name": "",
			}),
		},
		// set every field
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				disabled = true
				display_name = "My App"
				top_level_domain = "example.com"
				android_fingerprints = ["00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF"]
			`),
			Check: o.Check(map[string]any{
				"disabled":         true,
				"display_name":     "My App",
				"top_level_domain": "example.com",
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
