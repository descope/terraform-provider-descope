package lists_test

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	//go:embed tests/jsonlist.json
	jsonList string

	//go:embed tests/textslist.json
	textsList string

	//go:embed tests/ipslist.json
	ipsList string
)

func TestLists(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "JSON List"
						description = "A JSON list"
						type = "json"
						data = jsonencode(` + jsonList + `)
					}
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#": 1,
				"lists.0": map[string]any{
					"name":        "JSON List",
					"description": "A JSON list",
					"type":        "json",
					"data":        testacc.AttributeIsSet,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "Texts List"
						description = "A texts list"
						type = "texts"
						data = jsonencode(` + textsList + `)
					}
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#": 1,
				"lists.0": map[string]any{
					"name":        "Texts List",
					"description": "A texts list",
					"type":        "texts",
					"data":        testacc.AttributeIsSet,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "IPs List"
						description = "An IPs list"
						type = "ips"
						data = jsonencode(` + ipsList + `)
					}
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#": 1,
				"lists.0": map[string]any{
					"name":        "IPs List",
					"description": "An IPs list",
					"type":        "ips",
					"data":        testacc.AttributeIsSet,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "JSON List"
						type = "json"
						data = jsonencode(` + jsonList + `)
					},
					{
						name = "Texts List"
						type = "texts"
						data = jsonencode(` + textsList + `)
					},
					{
						name = "IPs List"
						type = "ips"
						data = jsonencode(` + ipsList + `)
					}
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#": 3,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "Invalid List"
						type = "invalid"
						data = jsonencode(` + jsonList + `)
					}
				]
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "` + strings.Repeat("a", 101) + `"
						type = "json"
						data = jsonencode(` + jsonList + `)
					}
				]
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
		},
		resource.TestStep{
			Config: p.Config(`
				lists = []
			`),
			Check: p.Check(map[string]any{
				"lists.#": 0,
			}),
		},
	)
}
