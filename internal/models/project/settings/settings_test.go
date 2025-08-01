package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSettings(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "3 weeks"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "3 weeks",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "1 day"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "1 day",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "1 minute"
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_token_expiration = "1 hour"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_token_expiration": "1 hour",
			}),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "4 weeks",
				"project_settings.session_token_expiration": "1 hour",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "4 weeks",
				"project_settings.session_token_expiration": "10 minutes",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_rotation = true
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_rotation": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_rotation": false,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					default_no_sso_apps = true
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.default_no_sso_apps": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.default_no_sso_apps": false,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = ["example.com"]
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{"example.com"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = ["example.com"]
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{"example.com"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = []
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = ["example.com"]
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{"example.com"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = null
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.approved_domains": []string{},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					approved_domains = ["example.com",","]
				}
			`),
			ExpectError: regexp.MustCompile(`commas`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					user_jwt_template = "foo"
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown JWT template reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					step_up_token_expiration = "12 minutes"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.step_up_token_expiration": "12 minutes",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					trusted_device_token_expiration = "52 weeks"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.trusted_device_token_expiration": "52 weeks",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					access_key_session_token_expiration = "2 minutes"
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					test_users_loginid_regexp = "^foo-[0-9]+@acmecorp.com$"
					test_users_verifier_regexp = "^bar-[0-9]+@acmecorp.com$"
					test_users_static_otp = "123456"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.test_users_loginid_regexp":  "^foo-[0-9]+@acmecorp.com$",
				"project_settings.test_users_verifier_regexp": "^bar-[0-9]+@acmecorp.com$",
				"project_settings.test_users_static_otp":      "123456",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_token_response_method": "response_body",
				"project_settings.refresh_token_response_method": "response_body",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_token_response_method = "cookies"
					refresh_token_response_method = "cookies"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_token_response_method": "cookies",
				"project_settings.refresh_token_response_method": "cookies",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_token_response_method": "response_body",
				"project_settings.refresh_token_response_method": "response_body",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_cookie_policy = "foo"
				}
			`),
			ExpectError: regexp.MustCompile(`value must be one of`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_cookie_policy = "strict"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_cookie_policy": "strict",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_cookie_policy = "lax"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_cookie_policy": "lax",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					enable_inactivity = true
					inactivity_time = "1 hour"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "4 weeks",
				"project_settings.enable_inactivity":        true,
				"project_settings.inactivity_time":          "1 hour",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_migration = {
						vendor = "foo"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_migration = {
						vendor = "okta"
						client_id = "foo"
						domain = "example.com"
						issuer = "bar"
						loginid_matched_attributes = [ "username", "email" ]
					}
				}
			`),
			ExpectError: regexp.MustCompile(`should not be set`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_migration = {
						vendor = "okta"
						client_id = "foo"
						issuer = "bar"
						loginid_matched_attributes = [ "username", "email" ]
					}
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_migration": map[string]any{
					"vendor":                     "okta",
					"client_id":                  "foo",
					"domain":                     "",
					"audience":                   "",
					"issuer":                     "bar",
					"loginid_matched_attributes": []string{"username", "email"},
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					session_migration = null
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.session_migration": map[string]any{
					"vendor":                     "",
					"client_id":                  "",
					"domain":                     "",
					"audience":                   "",
					"issuer":                     "",
					"loginid_matched_attributes": []string{},
				},
			}),
		},
	)
}
