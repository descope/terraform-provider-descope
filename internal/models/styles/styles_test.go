package styles_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
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

// An apply must not delete styles this resource does not carry - styles created in the console are
// invisible to Terraform, and used to be wiped by every apply.
// See https://github.com/descope/etc/issues/12809
func TestStylesKeepsUnmanagedStyles(t *testing.T) {
	projectID := testacc.ProjectID(t)
	s := testacc.Styles(t)

	unmanaged := map[string]any{
		"unmanaged-light": map[string]any{"name": "Unmanaged", "type": "flows"},
		"unmanaged-dark":  map[string]any{"name": "Unmanaged", "type": "flows"},
	}

	// leave the project holding only the styles this test manages, whatever the outcome
	t.Cleanup(func() {
		if projectID == "" {
			return
		}
		// /v2/mgmt/theme/import upserts and never deletes, so it cannot undo what this test
		// created. The deprecated /v1 endpoint still replaces the whole theme, which is what a
		// reset needs - it takes the theme in its cssTemplate shape rather than as styles.
		testacc.OutOfBandPost(t, projectID, "/v1/mgmt/theme/import", map[string]any{
			"theme": map[string]any{
				"cssTemplate": map[string]any{"light": map[string]any{}, "dark": map[string]any{}},
			},
		})
	})

	testacc.Run(t,
		// terraform manages the default style only
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
				"id.==": "project_id",
				"data":  testacc.AttributeIsSet,
			}),
		},
		// a style appears outside terraform, the way the console creates one, and then the managed
		// data changes so an import actually runs
		resource.TestStep{
			PreConfig: func() {
				testacc.OutOfBandPost(t, projectID, "/v2/mgmt/theme/import", map[string]any{
					"theme": map[string]any{"styles": unmanaged},
				})
			},
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
			Check: func(*terraform.State) error {
				theme := testacc.OutOfBandPostData(t, projectID, "/v2/mgmt/theme/export", map[string]any{})
				inner, ok := theme["theme"].(map[string]any)
				require.True(t, ok, "exported theme has no theme object")
				styles, ok := inner["styles"].(map[string]any)
				require.True(t, ok, "exported theme has no styles object")
				for key := range unmanaged {
					require.Contains(t, styles, key, "an apply deleted a style Terraform does not manage")
				}
				require.Contains(t, styles, "light", "an apply deleted the managed style")
				return nil
			},
		},
		// the surviving unmanaged style does not turn into a permanent diff
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
			PlanOnly:           true,
			ExpectNonEmptyPlan: false,
		},
	)
}
