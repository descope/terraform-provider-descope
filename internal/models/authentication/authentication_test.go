package authentication_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthentication(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"authentication.%": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						enabled = true
						redirect_url = "1"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`The redirectUrl field must be a valid URL`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						enabled = true
						redirect_url = "https://example.com"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.magic_link.enabled":      true,
				"authentication.magic_link.redirect_url": "https://example.com",
			}),
		},
	)
}
