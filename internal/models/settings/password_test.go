package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPasswordSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.PasswordSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                      testacc.AttributeIsSet,
				"project_id":              testacc.AttributeIsSet,
				"disabled":                false,
				"min_length":              8,
				"lowercase":               true,
				"uppercase":               true,
				"number":                  true,
				"non_alphanumeric":        true,
				"any_letter":              false,
				"disallowed_characters":   "",
				"disallow_email_match":    false,
				"expiration":              false,
				"expiration_weeks":        20,
				"reuse":                   false,
				"reuse_amount":            10,
				"lock":                    false,
				"lock_attempts":           5,
				"temporary_lock":          false,
				"temporary_lock_attempts": 3,
				"temporary_lock_duration": "5 minutes",
				"enforce_strength":        "none",
				"mask_errors":             true,
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				min_length = 12
				non_alphanumeric = false
				any_letter = true
				disallowed_characters = "$%"
				disallow_email_match = true
				expiration = true
				expiration_weeks = 4
				lock = true
				lock_attempts = 7
				temporary_lock = true
				temporary_lock_duration = "10 minutes"
				enforce_strength = "strong"
				mask_errors = false
			`),
			Check: m.Check(map[string]any{
				"min_length":              12,
				"non_alphanumeric":        false,
				"any_letter":              true,
				"disallowed_characters":   "$%",
				"disallow_email_match":    true,
				"expiration":              true,
				"expiration_weeks":        4,
				"lock":                    true,
				"lock_attempts":           7,
				"temporary_lock":          true,
				"temporary_lock_duration": "10 minutes",
				"enforce_strength":        "strong",
				"mask_errors":             false,
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"min_length":       8,
				"non_alphanumeric": true,
				"expiration":       false,
				"lock":             false,
				"enforce_strength": "none",
				"mask_errors":      true,
			}),
		},
		// email_service is server-populated on read, so it is outside the import contract
		resource.TestStep{
			ResourceName:            m.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(m.Path(), "project_id"),
			ImportStateVerifyIgnore: []string{"email_service"},
		},
	)
}

func TestPasswordSettingsInvalid(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.PasswordSettings(t)

	testacc.Run(t,
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				temporary_lock_duration = "90 seconds"
			`),
			ExpectError: regexp.MustCompile(`must be a whole number of minutes`),
		},
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				temporary_lock_duration = "2 minutes"
			`),
			Check: m.Check(map[string]any{"temporary_lock_duration": "2 minutes"}),
		},
	)
}

func TestPasswordSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	m := testacc.PasswordSettings(t)
	testacc.Run(t,
		// a custom email connector with an active reset template selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "password"
				name = "`+name+`"
				subject = "Reset your password"
				html_body = "Follow the link in this email to reset your password"
			`) + m.Block(`
				project_id = "`+projectID+`"
				email_service = {
					connector_id = `+c.Path()+`.id
				}
				email_template_id = `+e.Path()+`.id
			`),
			Check: m.Check(map[string]any{
				"email_service.connector_id": testacc.AttributeIsSet,
				"email_template_id":          testacc.AttributeIsSet,
			}),
		},
	)
}
