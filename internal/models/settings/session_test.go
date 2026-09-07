package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSessionSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.SessionSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                                  testacc.AttributeIsSet,
				"project_id":                          testacc.AttributeIsSet,
				"refresh_token_expiration":            "4 weeks",
				"refresh_token_rotation":              false,
				"refresh_token_response_method":       "response_body",
				"refresh_token_cookie_policy":         "none",
				"session_token_expiration":            "10 minutes",
				"step_up_token_expiration":            "10 minutes",
				"access_key_session_token_expiration": "10 minutes",
				"enable_inactivity":                   false,
				"inactivity_time":                     "12 minutes",
			}),
		},
		// update the settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				refresh_token_expiration = "2 weeks"
				refresh_token_rotation = true
				refresh_token_response_method = "cookies"
				refresh_token_cookie_policy = "strict"
				session_token_expiration = "15 minutes"
				step_up_token_expiration = "8 minutes"
				trusted_device_token_expiration = "52 weeks"
				session_token_response_method = "cookies"
				session_token_cookie_policy = "lax"
				enable_inactivity = true
				inactivity_time = "30 minutes"
			`),
			Check: m.Check(map[string]any{
				"refresh_token_expiration":        "2 weeks",
				"refresh_token_rotation":          true,
				"refresh_token_response_method":   "cookies",
				"refresh_token_cookie_policy":     "strict",
				"session_token_expiration":        "15 minutes",
				"step_up_token_expiration":        "8 minutes",
				"trusted_device_token_expiration": "52 weeks",
				"session_token_response_method":   "cookies",
				"session_token_cookie_policy":     "lax",
				"enable_inactivity":               true,
				"inactivity_time":                 "30 minutes",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"refresh_token_expiration":        "4 weeks",
				"refresh_token_rotation":          false,
				"refresh_token_response_method":   "response_body",
				"refresh_token_cookie_policy":     "none",
				"session_token_expiration":        "10 minutes",
				"session_token_response_method":   "response_body",
				"session_token_cookie_policy":     "none",
				"step_up_token_expiration":        "10 minutes",
				"trusted_device_token_expiration": "365 days",
				"enable_inactivity":               false,
				"inactivity_time":                 "12 minutes",
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
