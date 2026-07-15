package customlanguage_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCustomLanguage(t *testing.T) {
	p := testacc.Project(t)
	c := testacc.CustomLanguage(t)
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
		// Import with composite ID
		resource.TestStep{
			ResourceName:      c.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(c.Path(), "project_id", "id"),
		},
		// Destroy resource
		resource.TestStep{
			Config:  p.Config() + c.Config(`project_id = `+p.Path()+`.id`),
			Destroy: true,
		},
	)
}
