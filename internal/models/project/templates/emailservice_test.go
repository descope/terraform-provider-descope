package templates_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEmail(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(emailService(`
					`)),
			ExpectError: regexp.MustCompile(`attribute "connector" is required`),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
						connector = ""
					`)),
			ExpectError: regexp.MustCompile(`must not be empty`),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
						connector = "Foo"
					`)),
			ExpectError: regexp.MustCompile(`Unknown connector reference`),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
						connector = "Descope"
					`)),
			Check: p.Check(map[string]any{
				"authentication.magic_link.email_service.connector": "Descope",
			}),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
						connector = "Descope"
						templates = [
							{
								active = true
								name = "foo"
								html_body = "a"
								subject = "x"
							}
						]
					`)),
			ExpectError: regexp.MustCompile(`Invalid email service connector`),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
						connector = "Descope"
						templates = [
							{
								name = "foo"
								html_body = "a"
								subject = "x"
							},
							{
								name = "foo"
								html_body = "b"
								subject = "y"
							}
						]
					`)),
			ExpectError: regexp.MustCompile(`names must be unique`),
		},
		resource.TestStep{
			Config: p.Config(emailService(`
				connector = "Descope"
				templates = [
					{
						name = "foo"
						html_body = "a"
						subject = "x"
					},
					{
						name = "bar"
						html_body = "b"
						subject = "y"
					}
				]
			`)),
			Check: p.Check(map[string]any{
				"authentication.magic_link.email_service.connector":        "Descope",
				"authentication.magic_link.email_service.templates.#":      2,
				"authentication.magic_link.email_service.templates.0.name": "foo",
				"authentication.magic_link.email_service.templates.1.name": "bar",
			}),
		},
	)
}

func emailService(s string) string {
	return `authentication = {
				magic_link = {
					email_service = {
					` + s + `
					}
				}
			}`
}
