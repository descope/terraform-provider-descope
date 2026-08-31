package governance_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGovernance(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			// A project that never configured AGS carries no governance.json, so the
			// attribute reads back unset rather than as a zeroed object.
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"governance": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			// Everything set at once, including suite_disabled true. That one is the
			// suite kill switch and the only field here whose meaningful value is
			// false-y, so it is what a "drop empty fields" pass would silently lose.
			Config: p.Config(`
				governance = {
					configured     = true
					auto_approval  = true
					suite_disabled = true
					logo           = "data:image/png;base64,iVBORw0KGgo="
				}
			`),
			Check: p.Check(map[string]any{
				"governance.configured":     true,
				"governance.auto_approval":  true,
				"governance.suite_disabled": true,
				"governance.logo":           "data:image/png;base64,iVBORw0KGgo=",
			}),
		},
		resource.TestStep{
			// Flipping the booleans back off has to stick. The write path replaces the
			// file whole, so a false that fails to round-trip would read as "never set"
			// and leave the previous true in place.
			Config: p.Config(`
				governance = {
					configured     = true
					auto_approval  = false
					suite_disabled = false
					logo           = "data:image/png;base64,iVBORw0KGgo="
				}
			`),
			Check: p.Check(map[string]any{
				"governance.configured":     true,
				"governance.auto_approval":  false,
				"governance.suite_disabled": false,
			}),
		},
		resource.TestStep{
			// The logo is carried by the model rather than left out, because the infra
			// layer replaces governance.json instead of merging: a config that omits
			// the logo blanks whatever the console had set.
			Config: p.Config(`
				governance = {
					configured = true
					logo       = ""
				}
			`),
			Check: p.Check(map[string]any{
				"governance.configured": true,
				"governance.logo":       "",
			}),
		},
	)
}
