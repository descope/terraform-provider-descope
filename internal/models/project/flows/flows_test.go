package flows_test

import (
	_ "embed"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	//go:embed tests/emptystyle.json
	emptyStyle string

	//go:embed tests/basicstyle.json
	basicStyle string

	//go:embed tests/basicflow.json
	basicFlow string

	//go:embed tests/connectorflow.json
	connectorFlow string

	//go:embed tests/roleflow.json
	roleFlow string
)

func TestFlows(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		// Styles
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"styles.data": testacc.AttributeIsNotSet,
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
				styles = {
					data = jsonencode(` + basicStyle + `)
				}
			`),
			Check: p.Check(map[string]any{
				"styles.data": testacc.AttributeMatchesJSON(basicStyle),
			}),
		},
		// Flows
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"basic-flow" = {
						data = jsonencode(` + basicFlow + `)
					}
				}
			`),
			Check: p.Check(map[string]any{
				"flows.%":               1,
				"flows.basic-flow.data": testacc.AttributeMatchesJSON(basicFlow),
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"invalidid!@#$" = {
						data = jsonencode(` + basicFlow + `)
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + roleFlow + `)
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown role reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + connectorFlow + `)
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown connector reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + connectorFlow + `)
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
				"flows.references-flow.data": testacc.AttributeMatchesJSON(connectorFlow),
				"connectors.http.#":          1,
				"connectors.http.0.id":       testacc.AttributeHasPrefix("CI"),
				"connectors.http.0.name":     "My HTTP Connector",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"references-flow" = {
						data = jsonencode(` + connectorFlow + `)
					}
				}
				connectors = {
					"http": [
						{
							name = "Renamed Connector"
							base_url = "https://example.com"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown connector reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {
					"basic-flow" = {
						data = jsonencode(` + basicFlow + `)
					}
					"references-flow" = {
						data = jsonencode(` + connectorFlow + `)
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
				"flows.%":                    2,
				"flows.basic-flow.data":      testacc.AttributeMatchesJSON(basicFlow),
				"flows.references-flow.data": testacc.AttributeMatchesJSON(connectorFlow),
				"connectors.http.#":          1,
				"connectors.http.0.id":       testacc.AttributeHasPrefix("CI"),
				"connectors.http.0.name":     "My HTTP Connector",
			}),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"flows.%": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				flows = {}
			`),
			Check: p.Check(map[string]any{
				"flows.%": 0,
			}),
		},
	)
}
