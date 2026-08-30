package styles_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestStyles(t *testing.T) {
	projectID := testacc.ProjectID(t)
	s := testacc.Styles(t)
	testacc.Run(t,
		// import the styles from the theme representation
		resource.TestStep{
			Config: s.Block(`
				project_id = "` + projectID + `"
				data = jsonencode({
					styles = {
						light = {}
						dark = {}
					}
				})
			`),
			Check: s.Check(map[string]any{
				"id.==":      "project_id",
				"project_id": testacc.AttributeIsSet,
				"data":       testacc.AttributeIsSet,
			}),
		},
		// changing the styles data updates the theme in place
		resource.TestStep{
			Config: s.Block(`
				project_id = "` + projectID + `"
				data = jsonencode({
					styles = {
						light = {
							designTokens = {}
						}
						dark = {}
					}
				})
			`),
			Check: s.Check(map[string]any{
				"id.==": "project_id",
				"data":  testacc.AttributeIsSet,
			}),
		},
		// data is re-serialized server-side on read, so verify the identity attributes only
		resource.TestStep{
			ResourceName:            s.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(s.Path(), "project_id"),
			ImportStateVerifyIgnore: []string{"data"},
		},
		// removing the resource is a no-op that leaves the theme in place
	)
}
