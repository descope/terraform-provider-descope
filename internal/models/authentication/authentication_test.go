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
						disabled = true
						redirect_url = "https://example.com"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.magic_link.disabled":        true,
				"authentication.magic_link.redirect_url":    "https://example.com",
				"authentication.magic_link.expiration_time": "3 minutes",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "2000 seconds"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`space and one of the valid time units`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "1 second"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "5 minutes"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.magic_link.disabled":        false,
				"authentication.magic_link.redirect_url":    "https://example.com",
				"authentication.magic_link.expiration_time": "5 minutes",
			}),
		},
	)
}
