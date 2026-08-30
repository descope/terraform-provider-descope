package engine_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEngine(t *testing.T) {
	// Engine creation also needs the enterprise KMS key for the engine secret, and without these the create hangs until the request times out.
	t.Skip("Temporarily skipping engine test: engineservice is not deployed in the acceptance-test environment")

	p := testacc.Project(t)
	e := testacc.Engine(t)
	renamed := testacc.Engine(t)

	testacc.Run(t,
		resource.TestStep{
			Config: p.Config() + e.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: e.Check(map[string]any{
				"id":           testacc.AttributeIsSet,
				"project_id":   testacc.AttributeIsSet,
				"name":         e.Name,
				"created_time": testacc.AttributeIsSet,
				"secret":       testacc.AttributeIsSet,
			}),
		},
		resource.TestStep{
			Config: p.Config() + renamed.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: renamed.Check(map[string]any{
				"name":   renamed.Name,
				"secret": testacc.AttributeIsSet,
			}),
		},
		resource.TestStep{
			ResourceName:            e.Path(),
			ImportState:             true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(e.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"secret"}, // secret is not returned by the API on read, so it can't be verified on import.
		},
		resource.TestStep{
			Config:  p.Config() + e.Config(`project_id = `+p.Path()+`.id`),
			Destroy: true,
		},
	)
}
