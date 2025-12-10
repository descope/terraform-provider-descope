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
				lists = {
					list = [
						{
							name = "JSON List"
							description = "A JSON list"
							type = "json"
							data = ` + jsonList + `
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#":             1,
				"lists.list.0.name":        "JSON List",
				"lists.list.0.description": "A JSON list",
				"lists.list.0.type":        "json",
				"lists.list.0.data":        testacc.AttributeIsSet,
			}),
		},
		// texts
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "Texts List"
							description = "A texts list"
							type = "texts"
							data = ` + textsList + `
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#":             1,
				"lists.list.0.name":        "Texts List",
				"lists.list.0.description": "A texts list",
				"lists.list.0.type":        "texts",
				"lists.list.0.data":        testacc.AttributeIsSet,
			}),
		},
		// ips
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "IPs List"
							description = "An IPs list"
							type = "ips"
							data = ` + ipsList + `
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#":             1,
				"lists.list.0.name":        "IPs List",
				"lists.list.0.description": "An IPs list",
				"lists.list.0.type":        "ips",
				"lists.list.0.data":        testacc.AttributeIsSet,
			}),
		},
		// Multiple lists
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "JSON List"
							type = "json"
							data = ` + jsonList + `
						},
						{
							name = "Texts List"
							type = "texts"
							data = ` + textsList + `
						},
						{
							name = "IPs List"
							type = "ips"
							data = ` + ipsList + `
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#": 3,
			}),
		},
		// invalid type
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "Invalid List"
							type = "invalid"
							data = ` + jsonList + `
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		// empty
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = []
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#": 0,
			}),
		},
		// list without type (should work, type is optional)
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "List Without Type"
							data = ` + jsonList + `
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"lists.list.#":      1,
				"lists.list.0.name": "List Without Type",
			}),
		},
		// list with invalid name (too long)
		resource.TestStep{
			Config: p.Config(`
				lists = {
					list = [
						{
							name = "` + strings.Repeat("a", 101) + `"
							type = "json"
							data = ` + jsonList + `
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
	)
}
