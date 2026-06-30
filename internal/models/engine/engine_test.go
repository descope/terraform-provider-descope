package engine_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEngine(t *testing.T) {
	// Skipped until the acceptance-test environment supports engines: engine creation requires
	// engineservice (not deployed in sandbox) and the enterprise KMS key for the engine secret.
	// Without them the create hangs and the request times out. Re-enable once those are available.
	t.Skip("Temporarily skipping engine test: engineservice is not deployed in the acceptance-test environment")

	p := testacc.Project(t)
	e := testacc.Engine(t)
	renamed := testacc.Engine(t) // same address (descope_engine.test), different name → exercises update

	testacc.Run(t,
		// Create with required fields only. The secret and created_time are computed.
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
		// Update the name; the create-time secret is preserved in state.
		resource.TestStep{
			Config: p.Config() + renamed.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: renamed.Check(map[string]any{
				"name":   renamed.Name,
				"secret": testacc.AttributeIsSet,
			}),
		},
		// Import using the composite "project_id/id" identifier.
		resource.TestStep{
			ResourceName:      e.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(e.Path(), "project_id", "id"),
			// secret is not returned by the API on read, so it can't be verified on import.
			ImportStateVerifyIgnore: []string{"secret"},
		},
		// Destroy.
		resource.TestStep{
			Config:  p.Config() + e.Config(`project_id = `+p.Path()+`.id`),
			Destroy: true,
		},
	)
}
