package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSCIMConnector(t *testing.T) {
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

func TestHTTPConnectorOAuth2ClientCredentials(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"http": [
						{
							name     = "My HTTP Connector"
							base_url = "https://example.com"
							authentication = {
								oauth2_client_credentials = {
									client_id     = "test-client-id"
									client_secret = "test-client-secret"
									auth_url      = "https://example.com/oauth/token"
									auth_style    = "header"
									scopes        = "read write"
									token_request_headers = {
										"X-Custom" = "value"
									}
								}
							}
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.http.#": 1,
				"connectors.http.0": map[string]any{
					"id":       testacc.AttributeMatchesPattern(`^(CI|MP)`),
					"name":     "My HTTP Connector",
					"base_url": "https://example.com",
					"authentication.oauth2_client_credentials.client_id":                      "test-client-id",
					"authentication.oauth2_client_credentials.client_secret":                  "test-client-secret",
					"authentication.oauth2_client_credentials.auth_url":                       "https://example.com/oauth/token",
					"authentication.oauth2_client_credentials.auth_style":                     "header",
					"authentication.oauth2_client_credentials.scopes":                         "read write",
					"authentication.oauth2_client_credentials.token_request_headers.X-Custom": "value",
				},
			}),
		},
	)
}

func TestHTTPConnectorOAuth2ClientCredentialsValidation(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		// A single connector cannot declare more than one authentication method.
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"http": [
						{
							name     = "invalid-two-auth"
							base_url = "https://example.com"
							authentication = {
								bearer_token = "some-bearer-token"
								oauth2_client_credentials = {
									client_id     = "test"
									client_secret = "test"
									auth_url      = "https://example.com/oauth/token"
								}
							}
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Cannot specify more than one connector authentication method`),
		},
		// auth_style only accepts "header" or "body".
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"http": [
						{
							name     = "invalid-auth-style"
							base_url = "https://example.com"
							authentication = {
								oauth2_client_credentials = {
									client_id     = "test"
									client_secret = "test"
									auth_url      = "https://example.com/oauth/token"
									auth_style    = "invalid"
								}
							}
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`value must be one of`),
		},
	)
}

func TestHTTPConnectorEngine(t *testing.T) {
	// Assigning a connector to an engine requires a running engineservice to resolve the
	// engine, which is not deployed in the acceptance-test environment. The engine_id
	// round-trip is covered by the unit test in engine_internal_test.go.
	t.Skip("Temporarily skipping HTTP connector engine test: engineservice is not deployed in the acceptance-test environment")

	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"http": [
						{
							name      = "My HTTP Connector"
							base_url  = "https://example.com"
							engine_id = "CIEngineExample"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.http.#": 1,
				"connectors.http.0": map[string]any{
					"id":        testacc.AttributeHasPrefix("CI"),
					"name":      "My HTTP Connector",
					"base_url":  "https://example.com",
					"engine_id": "CIEngineExample",
				},
			}),
		},
	)
}
