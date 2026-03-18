package inboundapp_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestInboundApp(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
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
				"approved_callback_urls.#":         "0",
				"permissions_scopes.#":             "0",
				"attributes_scopes.#":              "0",
				"connections_scopes.#":             "0",
				"audience_whitelist.#":             "0",
				"force_add_all_authorization_info": "false",
				"default_audience":                 "",
			}),
		},
		// Test update with description and callback URLs
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				description = "Updated description"
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
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							template = "{}"
						}
					]
				}
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+p.Path()+`.jwt_templates.user_templates.0.id
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
		// Test updating the JWT template in the project and reflecting the new ID in the inbound app
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "bar"
							template = "{}"
						}
					]
				}
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+p.Path()+`.jwt_templates.user_templates.0.id
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
		// Destroy resource
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
			`),
			Destroy: true,
		},
	)
}
