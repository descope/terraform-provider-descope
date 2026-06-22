package connectors_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSCIMConnector(t *testing.T) {
	t.Skip("Temporarily skipping SCIM test because of backend problems")
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"scim": [
						{
							name             = "My SCIM Connector"
							description      = "A SCIM connector for provisioning"
							federated_app_id = "fake-app-id"
							base_url         = "https://example.com/scim"
							authentication   = {
								bearer_token = "test-bearer-token"
							}
							headers = {
								"X-Custom-Header" = "header-value"
							}
							hmac_secret = "test-hmac-secret"
							insecure    = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.scim.#": 1,
				"connectors.scim.0": map[string]any{
					"id":                          testacc.AttributeHasPrefix("CI"),
					"name":                        "My SCIM Connector",
					"description":                 "A SCIM connector for provisioning",
					"federated_app_id":            "fake-app-id",
					"base_url":                    "https://example.com/scim",
					"authentication.bearer_token": "test-bearer-token",
					"headers.X-Custom-Header":     "header-value",
					"insecure":                    true,
					"disabled":                    false,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"scim": [
						{
							name             = "My SCIM Connector"
							description      = "Updated description"
							federated_app_id = "fake-app-id"
							base_url         = "https://updated.example.com/scim/v2"
							insecure         = false
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.scim.#": 1,
				"connectors.scim.0": map[string]any{
					"id":                          testacc.AttributeHasPrefix("CI"),
					"name":                        "My SCIM Connector",
					"description":                 "Updated description",
					"federated_app_id":            "fake-app-id",
					"base_url":                    "https://updated.example.com/scim/v2",
					"authentication.bearer_token": "",
					"headers.%":                   0,
					"insecure":                    false,
					"disabled":                    false,
				},
			}),
		},
	)
}

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
