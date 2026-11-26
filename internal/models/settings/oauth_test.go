package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOAuthSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.OAuthSettings(t)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"disabled":   false,
			}),
		},
		// flip the disabled toggle
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				disabled = true
			`),
			Check: o.Check(map[string]any{
				"disabled": true,
			}),
		},
		// import using the project id (id == project_id for this singleton)
		resource.TestStep{
			ResourceName:      o.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id"),
		},
	)
}
