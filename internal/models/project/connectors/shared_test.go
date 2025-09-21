package connectors_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectorsShared(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"smtp": [
						{
							name = "My SMTP Connector"
							description = ""
							server = {
								host = "example.com"
								port = 587
							}
							sender = {
								email = "foo@bar.com"
								name = "Foo Bar"
							}
							authentication = {
								username = "foo"
								password = "bar"
							}
						}
					]
					"sns" = [
						{
							name = "My SNS Connector"
							description = "Foo Bar"
							access_key_id = "Foo"
							secret = "Bar"
							region = "us-west-2"
							organization_number = "123456789012"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.smtp.#": 1,
				"connectors.smtp.0": map[string]any{
					"id":             testacc.AttributeMatchesPattern(`^(CI|MP)`),
					"name":           "My SMTP Connector",
					"description":    "",
					"use_static_ips": false,
				},
				"connectors.sns.#": 1,
				"connectors.sns.0": map[string]any{
					"id":                  testacc.AttributeMatchesPattern(`^(CI|MP)`),
					"name":                "My SNS Connector",
					"description":         "Foo Bar",
					"access_key_id":       "Foo",
					"secret":              "Bar",
					"region":              "us-west-2",
					"organization_number": "123456789012",
				},
			}),
		},
	)
}
