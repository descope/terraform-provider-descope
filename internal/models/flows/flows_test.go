package flows_test

import (
	_ "embed"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	//go:embed flows_test_emptystyle.json
	emptyStyle string

	//go:embed flows_test_basicflow.json
	basicFlow string

	//go:embed flows_test_referencesflow.json
	referencesFlow string
)

func TestFlows(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
			`),
			Check: p.Check(map[string]any{
				"styles.data": testacc.AttributeIsNotSet,
				"flows.%":     0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				styles = {
					data = jsonencode(` + emptyStyle + `)
				}
			`),
			Check: p.Check(map[string]any{
				"styles.data": testacc.AttributeMatchesJSON(emptyStyle),
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"basic-flow" = {
						data = jsonencode(` + basicFlow + `)
					}
				}
			`),
			Check: p.Check(map[string]any{
				"flows.basic-flow.data": testacc.AttributeMatchesJSON(basicFlow),
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + referencesFlow + `)
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown connector reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + referencesFlow + `)
					}
				}
				connectors = {
					"http": [
						{
							name = "My HTTP Connector"
							base_url = "https://example.com"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"flows.references-flow.data": testacc.AttributeMatchesJSON(referencesFlow),
				"connectors.http.#":          1,
				"connectors.http.0.id":       testacc.AttributeHasPrefix("CI"),
				"connectors.http.0.name":     "My HTTP Connector",
			}),
		},
	)
}
