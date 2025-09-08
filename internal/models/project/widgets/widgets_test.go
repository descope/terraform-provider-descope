package widgets_test

import (
	_ "embed"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	//go:embed tests/testwidget.json
	testWidget string
)

func TestWidgets(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		// Widgets
		resource.TestStep{
			Config: p.Config(`
				widgets = {
					"test-widget" = {
						data = jsonencode(` + testWidget + `)
					}
				}
			`),
			Check: p.Check(map[string]any{
				"widgets.%":                1,
				"widgets.test-widget.data": testacc.AttributeMatchesJSON(testWidget),
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				widgets = {
					"invalidid!@#$" = {
						data = jsonencode(` + testWidget + `)
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"widgets.%": 1,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				widgets = {}
			`),
			Check: p.Check(map[string]any{
				"widgets.%": 0,
			}),
		},
	)
}
