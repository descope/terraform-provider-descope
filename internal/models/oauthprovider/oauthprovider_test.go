package oauthprovider_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOAuthProviderSystem(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.OAuthProvider(t)
	testacc.Run(t,
		// configure a built-in system provider with bring-your-own credentials
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "google"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
			`),
			Check: o.Check(map[string]any{
				"id":         "google",
				"project_id": testacc.AttributeIsSet,
				"disabled":   false,
				"client_id":  "my-client-id",
			}),
		},
		// a system provider cannot set custom-only endpoint fields
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "google"
				authorization_endpoint = "https://auth.example.com"
			`),
			ExpectError: regexp.MustCompile(`Reserved Attribute`),
		},
		// a system provider with client_id must also supply client_secret
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "google"
				client_id = "my-client-id"
			`),
			ExpectError: regexp.MustCompile(`Missing Attribute Value`),
		},
		// the id is the provider identity on the backend, so an import verifies it cleanly
		resource.TestStep{
			ResourceName:            o.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(o.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"client_secret"},
		},
		// validation runs at plan time, so the destroy that ends the test needs a valid configuration
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "google"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
			`),
		},
	)
}

func TestOAuthProviderCustom(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.OAuthProvider(t)
	testacc.Run(t,
		// a fully-specified custom provider
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "custom_idp"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
				allowed_grant_types = ["authorization_code"]
				authorization_endpoint = "https://auth.example.com"
				token_endpoint = "https://token.example.com"
				user_info_endpoint = "https://userinfo.example.com"
			`),
			Check: o.Check(map[string]any{
				"id":                     "custom_idp",
				"authorization_endpoint": "https://auth.example.com",
				"token_endpoint":         "https://token.example.com",
				"user_info_endpoint":     "https://userinfo.example.com",
			}),
		},
		// a different provider id forces a replacement, since an in-place update would keep the unset endpoint values from the prior state
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "custom_idp_partial"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
			`),
			ExpectError: regexp.MustCompile(`Invalid Custom OAuth Provider`),
		},
		// a custom provider cannot use the apple-only key generator
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "custom_idp"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
				allowed_grant_types = ["authorization_code"]
				authorization_endpoint = "https://auth.example.com"
				token_endpoint = "https://token.example.com"
				user_info_endpoint = "https://userinfo.example.com"
				apple_key_generator = {
					key_id      = "KEYID"
					team_id     = "TEAMID1234"
					private_key = "-----BEGIN PRIVATE KEY-----\nMOCK\n-----END PRIVATE KEY-----"
				}
			`),
			ExpectError: regexp.MustCompile(`Reserved Attribute`),
		},
		// validation runs at plan time, so the destroy that ends the test needs a valid configuration
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				id = "custom_idp"
				client_id = "my-client-id"
				client_secret = "my-client-secret"
				allowed_grant_types = ["authorization_code"]
				authorization_endpoint = "https://auth.example.com"
				token_endpoint = "https://token.example.com"
				user_info_endpoint = "https://userinfo.example.com"
			`),
		},
	)
}

func TestOAuthSettingsAndProviders(t *testing.T) {
	projectID := testacc.ProjectID(t)
	settings := testacc.OAuthSettings(t)
	google := &testacc.Resource{Type: "oauth_provider", ID: "google"}
	custom := &testacc.Resource{Type: "oauth_provider", ID: "custom"}

	config := settings.Block(`project_id = "`+projectID+`"`) +
		google.Block(`
			project_id = "`+projectID+`"
			id = "github"
			client_id = "gh-client-id"
			client_secret = "gh-client-secret"
		`) +
		custom.Block(`
			project_id = "`+projectID+`"
			id = "custom_idp"
			client_id = "my-client-id"
			client_secret = "my-client-secret"
			allowed_grant_types = ["authorization_code"]
			authorization_endpoint = "https://auth.example.com"
			token_endpoint = "https://token.example.com"
			user_info_endpoint = "https://userinfo.example.com"
		`)

	testacc.Run(t,
		resource.TestStep{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				settings.Check(map[string]any{"disabled": false}),
				google.Check(map[string]any{"id": "github"}),
				custom.Check(map[string]any{"id": "custom_idp"}),
			),
		},
	)
}
