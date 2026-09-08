package outboundapp_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOutboundApp(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OutboundApp(t)
	testacc.RunWithDestroyCheck(t, "descope_outbound_app",
		// create with only the required fields, so the defaults are pinned
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
			`),
			Check: a.Check(map[string]any{
				"id":                         testacc.AttributeIsSet,
				"project_id":                 testacc.AttributeIsSet,
				"name":                       a.Name,
				"description":                "",
				"tenant_id":                  "",
				"logo":                       "",
				"app_type":                   "oauth",
				"client_id":                  "",
				"discovery_url":              "",
				"authorization_url":          "",
				"authorization_url_params.#": "0",
				"token_url":                  "",
				"token_url_params.#":         "0",
				"revocation_url":             "",
				"default_scopes.#":           "0",
				"default_redirect_url":       "",
				"callback_domain":            "",
				"pkce":                       false,
				"access_type":                "",
				"prompt.#":                   "0",
			}),
		},
		// the tenant cannot be changed after creation, and it must fail at plan time rather than
		// destroying the app, because deleting one cascades to every token stored against it
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				tenant_id = "T2abcdefghijklmnopqrstuvwxyz"
			`),
			ExpectError: regexp.MustCompile(`Immutable Attribute Changed`),
		},
		// populate every field, including both url param lists and the secret
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				description = "calendar access"
				logo = "https://example.com/logo.png"
				app_type = "oauth"
				client_id = "client-abc"
				client_secret = "secret-one"
				discovery_url = "https://accounts.example.com/.well-known/openid-configuration"
				authorization_url = "https://accounts.example.com/authorize"
				authorization_url_params = [
					{ key = "audience", value = "https://api.example.com" },
					{ key = "login_hint", value = "user@example.com" },
				]
				token_url = "https://oauth2.example.com/token"
				token_url_params = [
					{ key = "resource", value = "https://api.example.com" },
				]
				revocation_url = "https://oauth2.example.com/revoke"
				default_scopes = ["openid", "email", "calendar.read"]
				default_redirect_url = "https://app.example.com/oauth/callback"
				callback_domain = "app.example.com"
				pkce = true
				access_type = "offline"
				prompt = ["consent", "select_account"]
			`),
			Check: a.Check(map[string]any{
				"description":                    "calendar access",
				"logo":                           "https://example.com/logo.png",
				"app_type":                       "oauth",
				"client_id":                      "client-abc",
				"client_secret":                  "secret-one",
				"discovery_url":                  "https://accounts.example.com/.well-known/openid-configuration",
				"authorization_url":              "https://accounts.example.com/authorize",
				"authorization_url_params.#":     "2",
				"authorization_url_params.0.key": "audience",
				"authorization_url_params.1.key": "login_hint",
				"token_url":                      "https://oauth2.example.com/token",
				"token_url_params.#":             "1",
				"token_url_params.0.value":       "https://api.example.com",
				"revocation_url":                 "https://oauth2.example.com/revoke",
				"default_scopes":                 []string{"openid", "email", "calendar.read"},
				"default_redirect_url":           "https://app.example.com/oauth/callback",
				"callback_domain":                "app.example.com",
				"pkce":                           true,
				"access_type":                    "offline",
				"prompt":                         []string{"consent", "select_account"},
			}),
		},
		// Omitting client_secret must be accepted and must not plan a change: SecretOptional leaves it
		// null, stringattr.Get drops null keys, and the backend reads an absent secret as "keep the
		// existing one". Whether the server truly kept it is not observable here, because the API never
		// returns a secret - that half of the contract is the backend's, not this provider's.
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				description = "calendar access, updated"
				app_type = "oauth"
				client_id = "client-abc"
				authorization_url = "https://accounts.example.com/authorize"
				token_url = "https://oauth2.example.com/token"
				default_redirect_url = "https://app.example.com/oauth/callback"
				pkce = true
				access_type = "offline"
			`),
			Check: a.Check(map[string]any{
				"description":   "calendar access, updated",
				"client_secret": testacc.AttributeIsNotSet,
				"access_type":   "offline",
			}),
		},
		// rotating the secret in place must not require a replacement
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				description = "calendar access, updated"
				app_type = "oauth"
				client_id = "client-abc"
				client_secret = "secret-two"
				authorization_url = "https://accounts.example.com/authorize"
				token_url = "https://oauth2.example.com/token"
				default_redirect_url = "https://app.example.com/oauth/callback"
				pkce = true
				access_type = "offline"
			`),
			Check: a.Check(map[string]any{
				"id":            testacc.AttributeIsSet,
				"client_secret": "secret-two",
			}),
		},
		// removing the url params and scopes must actually clear them rather than leave the old values
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				app_type = "oauth"
				client_id = "client-abc"
				client_secret = "secret-two"
				authorization_url = "https://accounts.example.com/authorize"
				token_url = "https://oauth2.example.com/token"
			`),
			Check: a.Check(map[string]any{
				"authorization_url_params.#": "0",
				"token_url_params.#":         "0",
				"default_scopes.#":           "0",
				"description":                "",
				"pkce":                       false,
				"access_type":                "",
			}),
		},
		// the client secret is never returned by the backend, so it cannot participate in an import
		resource.TestStep{
			ResourceName:            a.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"client_secret"},
		},
	)
}

func TestOutboundAppAPIKeyType(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OutboundApp(t)
	testacc.RunWithDestroyCheck(t, "descope_outbound_app",
		// an apikey app runs no OAuth flow, so it needs none of the endpoint fields
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				app_type = "apikey"
				description = "static token app"
			`),
			Check: a.Check(map[string]any{
				"app_type":    "apikey",
				"description": "static token app",
				"pkce":        false,
			}),
		},
	)
}

func TestOutboundAppInvalidValues(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OutboundApp(t)

	// every closed set is rejected at plan time rather than as an apply time error from the backend
	testacc.Run(t,
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				app_type = "workflow"
			`),
			ExpectError: regexp.MustCompile(`Attribute app_type value must be one of`),
		},
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				access_type = "forever"
			`),
			ExpectError: regexp.MustCompile(`Attribute access_type value must be one of`),
		},
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				prompt = ["consent", "shout"]
			`),
			ExpectError: regexp.MustCompile(`value must be one of`),
		},
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				default_redirect_url = "not-a-url"
			`),
			ExpectError: regexp.MustCompile(`(?i)url`),
		},
		// Creating against a tenant that does not exist must fail with the backend's own message rather
		// than succeeding at project level. The regex is deliberately loose because the backend's wording
		// contains a typo ("requested tenant no found") that may be corrected later.
		resource.TestStep{
			Config: a.Config(`
				project_id = "` + projectID + `"
				tenant_id = "T2doesnotexist"
			`),
			ExpectError: regexp.MustCompile(`(?i)tenant`),
		},
	)
}
