package inboundapp_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestInboundAppDeletionProtection(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: a.Check(map[string]any{
				"deletion_protection": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
			`),
			Destroy:     true,
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				non_confidential_client = true
			`),
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		// The flag change must be applied on its own before the replacement is allowed
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
			`),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
				non_confidential_client = true
			`),
			Check: a.Check(map[string]any{
				"non_confidential_client": "true",
				"deletion_protection":     "false",
			}),
		},
	)
}

func TestInboundApp(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	j := testacc.JWTTemplate(t)
	j2 := testacc.JWTTemplate(t)
	j2.ID = "other"
	testacc.Run(t,
		// Test basic creation with required fields only
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
			`),
			Check: a.Check(map[string]any{
				"id":                               testacc.AttributeIsSet,
				"project_id":                       testacc.AttributeIsSet,
				"name":                             a.Name,
				"description":                      "",
				"non_confidential_client":          "false",
				"client_id":                        testacc.AttributeIsSet,
				"client_secret":                    testacc.AttributeIsSet,
				"force_pkce":                       "false",
				"approved_callback_urls.#":         "0",
				"permissions_scopes.#":             "0",
				"attributes_scopes.#":              "0",
				"connections_scopes.#":             "0",
				"audience_whitelist.#":             "0",
				"force_add_all_authorization_info": "false",
				"force_dpop":                       "false",
				"default_audience":                 "",
			}),
		},
		// Test update with description and callback URLs
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				description = "Updated description"
				force_pkce = true
				approved_callback_urls = ["https://example.com/callback"]
				permissions_scopes = [
					{
						name = "openid"
						description = "Foo"
					},
					{
						name = "email"
						description = "Access email"
						optional = true
					},
				]
			`),
			Check: a.Check(map[string]any{
				"description":                      "Updated description",
				"force_pkce":                       "true",
				"approved_callback_urls.#":         "1",
				"approved_callback_urls.0":         "https://example.com/callback",
				"permissions_scopes.#":             "2",
				"permissions_scopes.0.name":        "openid",
				"permissions_scopes.1.name":        "email",
				"permissions_scopes.1.description": "Access email",
				"permissions_scopes.1.optional":    "true",
			}),
		},
		// Test update with session settings
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
			}),
		},
		// Test session settings with a user JWT template
		resource.TestStep{
			Config: p.Config() + j.Config(`
				project_id = `+p.Path()+`.id
				type = "user"
				template = "{}"
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+j.Path()+`.id
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
				"session_settings.user_template_id":             testacc.AttributeHasPrefix("JT"),
			}),
		},
		resource.TestStep{
			Config: p.Config() + j2.Config(`
				project_id = `+p.Path()+`.id
				type = "user"
				template = "{}"
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+j2.Path()+`.id
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
				"session_settings.user_template_id":             testacc.AttributeHasPrefix("JT"),
			}),
		},
		// Test import with composite ID
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		},
		// Disable the default deletion protection so the resource can be destroyed
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
			`),
		},
		// Destroy resource
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
			`),
			Destroy: true,
		},
	)
}
