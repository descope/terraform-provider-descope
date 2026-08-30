package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEmbeddedLinkSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.EmbeddedLinkSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":              testacc.AttributeIsSet,
				"project_id":      testacc.AttributeIsSet,
				"disabled":        false,
				"expiration_time": "3 minutes",
			}),
		},
		// update the settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				disabled = true
				expiration_time = "10 minutes"
			`),
			Check: m.Check(map[string]any{
				"disabled":        true,
				"expiration_time": "10 minutes",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"disabled":        false,
				"expiration_time": "3 minutes",
			}),
		},
		resource.TestStep{
			ResourceName:      m.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(m.Path(), "project_id"),
		},
	)
}
