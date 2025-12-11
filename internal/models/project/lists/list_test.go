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
		// JSON type list
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
				"lists.#":             1,
				"lists.0.name":        "JSON List",
				"lists.0.description": "A JSON list",
				"lists.0.type":        "json",
				"lists.0.data":        testacc.AttributeIsSet,
			}),
		},
		// texts
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
				"lists.#":             1,
				"lists.0.name":        "Texts List",
				"lists.0.description": "A texts list",
				"lists.0.type":        "texts",
				"lists.0.data":        testacc.AttributeIsSet,
			}),
		},
		// ips
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
				"lists.#":             1,
				"lists.0.name":        "IPs List",
				"lists.0.description": "An IPs list",
				"lists.0.type":        "ips",
				"lists.0.data":        testacc.AttributeIsSet,
			}),
		},
		// Multiple lists
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
		// invalid type
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
		// empty
		resource.TestStep{
			Config: p.Config(`
				lists = [
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#": 0,
			}),
		},
		// list without type (should work, type is optional)
		resource.TestStep{
			Config: p.Config(`
				lists = [
					{
						name = "List Without Type"
						data = jsonencode(` + jsonList + `)
					}
				]
			`),
			Check: p.Check(map[string]any{
				"lists.#":      1,
				"lists.0.name": "List Without Type",
			}),
		},
		// list with invalid name (too long)
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
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
	)
}
