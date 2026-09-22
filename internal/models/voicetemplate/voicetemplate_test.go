package voicetemplate_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVoiceTemplate(t *testing.T) {
	projectID := testacc.ProjectID(t)
	name := testacc.GenerateAlias(t)
	x := testacc.VoiceTemplate(t)
	testacc.Run(t,
		// a new template gets a server-assigned id
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "otp"
				name = "` + name + `"
				body = "Your code is {{.code}}"
			`),
			Check: x.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"method":     "otp",
				"name":       name,
				"body":       "Your code is {{.code}}",
			}),
		},
		// the name is mutable since the server-assigned id is the identity
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "otp"
				name = "` + name + `-renamed"
				body = "The code is {{.code}}"
			`),
			Check: x.Check(map[string]any{
				"name": name + "-renamed",
				"body": "The code is {{.code}}",
			}),
		},
		// the System name is reserved for the built-in template
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "otp"
				name = "System"
				body = "body"
			`),
			ExpectError: regexp.MustCompile(`Invalid voice template`),
		},
		// end on a valid configuration so the post-test destroy passes config validation
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "otp"
				name = "` + name + `"
				body = "Your code is {{.code}}"
			`),
			Check: x.Check(map[string]any{"name": name}),
		},
		resource.TestStep{
			ResourceName:      x.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(x.Path(), "project_id", "method", "id"),
		},
	)
}

func TestOTPTemplatesAndSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	names := testacc.GenerateAlias(t)
	settings := testacc.OTPSettings(t)
	welcome := &testacc.Resource{Type: "email_template", ID: "welcome"}
	reminder := &testacc.Resource{Type: "email_template", ID: "reminder"}
	text := &testacc.Resource{Type: "text_template", ID: "text"}
	voice := &testacc.Resource{Type: "voice_template", ID: "voice"}

	templates := welcome.Block(`
			project_id = "`+projectID+`"
			method = "otp"
			name = "`+names+`-welcome"
			subject = "Your code"
			html_body = "Use the code in this email to sign in"
		`) + reminder.Block(`
			project_id = "`+projectID+`"
			method = "otp"
			name = "`+names+`-reminder"
			subject = "Reminder"
			html_body = "Your code is still valid"
		`) + text.Block(`
			project_id = "`+projectID+`"
			method = "otp"
			name = "`+names+`-short"
			body = "Code: {{code}}"
		`) + voice.Block(`
			project_id = "`+projectID+`"
			method = "otp"
			name = "`+names+`-voice"
			body = "Your code is {{code}}"
		`)

	testacc.Run(t,
		// all templates and the settings write the shared messaging row in one apply
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
				email_template_id = `+welcome.Path()+`.id
				text_template_id = `+text.Path()+`.id
				voice_template_id = `+voice.Path()+`.id
			`),
			Check: resource.ComposeAggregateTestCheckFunc(
				welcome.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-welcome"}),
				reminder.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-reminder"}),
				text.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-short"}),
				voice.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-voice"}),
				settings.Check(map[string]any{
					"email_template_id": testacc.AttributeIsSet,
					"text_template_id":  testacc.AttributeIsSet,
					"voice_template_id": testacc.AttributeIsSet,
				}),
			),
		},
		// the active selection can be repointed at another template
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
				email_template_id = `+reminder.Path()+`.id
			`),
			Check: settings.Check(map[string]any{
				"email_template_id": testacc.AttributeIsSet,
				"text_template_id":  "",
				"voice_template_id": "",
			}),
		},
		// dropping the references resets the selections to the built-in System template, so the same apply can also destroy the previously active templates
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
			`),
			Check: settings.Check(map[string]any{"email_template_id": ""}),
		},
	)
}
