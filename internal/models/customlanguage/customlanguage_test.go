package customlanguage_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCustomLanguage(t *testing.T) {
	p := testacc.Project(t)
	c := testacc.CustomLanguage(t)
	renamed := testacc.CustomLanguage(t) // same address (descope_custom_language.test), different name → exercises update
	testacc.Run(t,
		// Create with a language code only (no region)
		resource.TestStep{
			Config: p.Config() + c.Config(`
				project_id = `+p.Path()+`.id
				language   = "phl"
			`),
			Check: c.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"language":   "phl",
				"name":       c.Name,
			}),
		},
		// Add a region: the code is immutable, so this replaces the resource
		resource.TestStep{
			Config: p.Config() + c.Config(`
				project_id = `+p.Path()+`.id
				language   = "phl"
				region     = "PH"
			`),
			Check: c.Check(map[string]any{
				"language": "phl",
				"region":   "PH",
				"name":     c.Name,
			}),
		},
		// Update the name in place (language/region are immutable and unchanged)
		resource.TestStep{
			Config: p.Config() + renamed.Config(`
				project_id = `+p.Path()+`.id
				language   = "phl"
				region     = "PH"
			`),
			Check: renamed.Check(map[string]any{
				"language": "phl",
				"region":   "PH",
				"name":     renamed.Name,
			}),
		},
		// Import with composite ID
		resource.TestStep{
			ResourceName:      c.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(c.Path(), "project_id", "id"),
		},
		// Destroy resource (config must carry the required language + match the final state)
		resource.TestStep{
			Config: p.Config() + c.Config(`
				project_id = `+p.Path()+`.id
				language   = "phl"
				region     = "PH"
			`),
			Destroy: true,
		},
	)
}
