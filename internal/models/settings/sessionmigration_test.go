package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSessionMigration(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.SessionMigration(t)
	testacc.Run(t,
		// defaults only, migration disabled
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"vendor":     "",
			}),
		},
		// vendor specific attributes are required
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				vendor = "okta"
				client_id = "0oa1b2c3d4e5f6g7h8i9"
			`),
			ExpectError: regexp.MustCompile(`issuer attribute is required`),
		},
		// enable okta session migration
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				vendor = "okta"
				client_id = "0oa1b2c3d4e5f6g7h8i9"
				issuer = "https://dev-123456.okta.com/oauth2/default"
				api_token = "00aBcDeFgHiJkLmNoPqRsTuVwXyZ"
				loginid_matched_attributes = ["email"]
				user_sync_type = "jit"
				user_mapping = [
					{ external_key = "email", descope_key = "email" },
					{ external_key = "given_name", descope_key = "givenName" },
				]
			`),
			Check: m.Check(map[string]any{
				"vendor":                      "okta",
				"client_id":                   "0oa1b2c3d4e5f6g7h8i9",
				"issuer":                      "https://dev-123456.okta.com/oauth2/default",
				"loginid_matched_attributes":  []string{"email"},
				"user_sync_type":              "jit",
				"user_mapping.#":              "2",
				"user_mapping.0.external_key": "email",
				"user_mapping.0.descope_key":  "email",
				"user_mapping.1.external_key": "given_name",
				"user_mapping.1.descope_key":  "givenName",
			}),
		},
		// switch to auth0, clearing the sync config by omitting it
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				vendor = "auth0"
				client_id = "abcdefghij"
				domain = "dev-123456.auth0.com"
				loginid_matched_attributes = ["email"]
			`),
			Check: m.Check(map[string]any{
				"vendor":         "auth0",
				"client_id":      "abcdefghij",
				"domain":         "dev-123456.auth0.com",
				"issuer":         "",
				"user_sync_type": "",
				"user_mapping.#": "0",
			}),
		},
		// removing all attributes disables migration
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"vendor":    "",
				"client_id": "",
				"domain":    "",
			}),
		},
		resource.TestStep{
			ResourceName:      m.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(m.Path(), "project_id"),
		},
	)
}
