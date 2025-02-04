package templates_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestText(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(textService(`
					`)),
			ExpectError: regexp.MustCompile(`attribute "connector" is required`),
		},
		resource.TestStep{
			Config: p.Config(textService(`
						connector = ""
					`)),
			ExpectError: regexp.MustCompile(`must not be empty`),
		},
		resource.TestStep{
			Config: p.Config(textService(`
						connector = "Foo"
					`)),
			ExpectError: regexp.MustCompile(`Unknown connector reference`),
		},
		resource.TestStep{
			Config: p.Config(textService(`
						connector = "Descope"
					`)),
			Check: p.Check(map[string]any{
				"authentication.otp.text_service.connector": "Descope",
			}),
		},
		resource.TestStep{
			Config: p.Config(textService(`
						connector = "Descope"
						templates = [
							{
								active = true
								name = "foo"
								body = "a"
							}
						]
					`)),
			ExpectError: regexp.MustCompile(`Invalid text service connector`),
		},
		resource.TestStep{
			Config: p.Config(textService(`
						connector = "Descope"
						templates = [
							{
								name = "foo"
								body = "a"
							},
							{
								name = "foo"
								body = "b"
							}
						]
					`)),
			ExpectError: regexp.MustCompile(`names must be unique`),
		},
		resource.TestStep{
			Config: p.Config(textService(`
				connector = "Descope"
				templates = [
					{
						name = "foo"
						body = "a"
					},
					{
						name = "bar"
						body = "b"
					}
				]
			`)),
			Check: p.Check(map[string]any{
				"authentication.otp.text_service.connector":        "Descope",
				"authentication.otp.text_service.templates.#":      2,
				"authentication.otp.text_service.templates.0.name": "foo",
				"authentication.otp.text_service.templates.1.name": "bar",
			}),
		},
		resource.TestStep{
			Config: p.Config(textService(`
				connector = "Descope"
				templates = [
					{
						active = true
						name = "foo"
						body = "a"
					}
				]
			`)),
			ExpectError: regexp.MustCompile(`must not be set to Descope`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"generic_sms_gateway": [
						{
							name = "Generic SMS Gateway Connector"
							post_url = "https://example.com"
							sender = "test"
						}
					]
					"twilio_core": [
						{
							name = "Twilio Core Connector"
							account_sid = "foo"
							authentication = {
								auth_token = "bar"
							}
							senders = {
								sms = {
									phone_number = "1234"
								}
							}
						}
					]
				}
				authentication = {
					magic_link = {
						text_service = {
							connector = "Generic SMS Gateway Connector"
						}
					}
					otp = {
						text_service = {
							connector = "Twilio Core Connector"
							templates = [
								{
									active = true
									name = "foo"
									body = "bar"
								}
							]
						}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.generic_sms_gateway.#":                   1,
				"connectors.twilio_core.#":                           1,
				"authentication.magic_link.text_service.connector":   "Generic SMS Gateway Connector",
				"authentication.otp.text_service.connector":          "Twilio Core Connector",
				"authentication.otp.text_service.templates.#":        1,
				"authentication.otp.text_service.templates.0.active": true,
				"authentication.otp.text_service.templates.0.name":   "foo",
				"authentication.otp.text_service.templates.0.body":   "bar",
			}),
		},
	)
}

func textService(s string) string {
	return `authentication = {
				otp = {
					text_service = {
					` + s + `
					}
				}
			}`
}
