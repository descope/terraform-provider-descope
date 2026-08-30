package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOTPSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.OTPSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":              testacc.AttributeIsSet,
				"project_id":      testacc.AttributeIsSet,
				"disabled":        false,
				"domain":          "",
				"expiration_time": "3 minutes",
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				domain = "example.com"
				expiration_time = "10 minutes"
			`),
			Check: m.Check(map[string]any{
				"domain":          "example.com",
				"expiration_time": "10 minutes",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"domain":          "",
				"expiration_time": "3 minutes",
			}),
		},
		// the messaging service blocks are server-populated on read, outside the import contract
		resource.TestStep{
			ResourceName:            m.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(m.Path(), "project_id"),
			ImportStateVerifyIgnore: []string{"email_service", "text_service", "voice_service"},
		},
	)
}

func TestOTPSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	v := testacc.VoiceTemplate(t)
	m := testacc.OTPSettings(t)
	testacc.Run(t,
		// a custom email connector with active email and voice templates selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "otp"
				name = "`+name+`"
				subject = "Your code"
				html_body = "Use the code in this email to sign in"
			`) + v.Block(`
				project_id = "`+projectID+`"
				method = "otp"
				name = "`+name+`-voice"
				body = "Your code is {{.code}}"
			`) + m.Block(`
				project_id = "`+projectID+`"
				email_service = {
					connector_id = `+c.Path()+`.id
				}
				email_template_id = `+e.Path()+`.id
				voice_template_id = `+v.Path()+`.id
			`),
			Check: m.Check(map[string]any{
				"email_service.connector_id": testacc.AttributeIsSet,
				"email_template_id":          testacc.AttributeIsSet,
				"voice_template_id":          testacc.AttributeIsSet,
				"text_template_id":           "",
			}),
		},
	)
}
