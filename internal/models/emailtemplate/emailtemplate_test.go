package emailtemplate_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEmailTemplate(t *testing.T) {
	projectID := testacc.ProjectID(t)
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	testacc.Run(t,
		resource.TestStep{
			Config: e.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `"
				subject = "Sign in"
				html_body = "Follow the link in this email to sign in"
			`),
			Check: e.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"method":     "magiclink",
				"name":       name,
				"subject":    "Sign in",
				"html_body":  "Follow the link in this email to sign in",
			}),
		},
		resource.TestStep{
			Config: e.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `-renamed"
				subject = "Sign in here"
				plain_text_body = "Follow the link to sign in"
				use_plain_text_body = true
			`),
			Check: e.Check(map[string]any{
				"name":                name + "-renamed",
				"subject":             "Sign in here",
				"plain_text_body":     "Follow the link to sign in",
				"use_plain_text_body": true,
			}),
		},
		resource.TestStep{
			Config: e.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "System"
				subject = "Sign in"
				html_body = "body"
			`),
			ExpectError: regexp.MustCompile(`Invalid email template`),
		},
		resource.TestStep{
			Config: e.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `"
				subject = "Sign in"
			`),
			ExpectError: regexp.MustCompile(`Missing Attribute Value`),
		},
		resource.TestStep{
			Config: e.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `"
				subject = "Sign in"
				html_body = "Follow the link in this email to sign in"
			`),
			Check: e.Check(map[string]any{"name": name}),
		},
		resource.TestStep{
			ResourceName:      e.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(e.Path(), "project_id", "method", "id"),
		},
	)
}

func TestMagicLinkTemplatesAndSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	names := testacc.GenerateAlias(t)
	settings := testacc.MagicLinkSettings(t)
	welcome := &testacc.Resource{Type: "email_template", ID: "welcome"}
	reminder := &testacc.Resource{Type: "email_template", ID: "reminder"}
	text := &testacc.Resource{Type: "text_template", ID: "text"}
	connector := testacc.NewResource(t, "generic_email_gateway_connector")

	templates := connector.Config(`
			project_id = "`+projectID+`"
			post_url = "https://mail.example.com/send"
		`) + welcome.Block(`
			project_id = "`+projectID+`"
			method = "magiclink"
			name = "`+names+`-welcome"
			subject = "Sign in"
			html_body = "Follow the link in this email to sign in"
		`) + reminder.Block(`
			project_id = "`+projectID+`"
			method = "magiclink"
			name = "`+names+`-reminder"
			subject = "Reminder"
			html_body = "You have a pending sign in link"
		`) + text.Block(`
			project_id = "`+projectID+`"
			method = "magiclink"
			name = "`+names+`-short"
			body = "Tap to sign in: {{link}}"
		`)

	testacc.Run(t,
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
				email_service = {
					connector_id = `+connector.Path()+`.id
				}
				email_template_id = `+welcome.Path()+`.id
			`),
			Check: resource.ComposeAggregateTestCheckFunc(
				welcome.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-welcome"}),
				reminder.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-reminder"}),
				text.Check(map[string]any{"id": testacc.AttributeIsSet, "name": names + "-short"}),
				settings.Check(map[string]any{"email_template_id": testacc.AttributeIsSet}),
			),
		},
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
				email_service = {
					connector_id = `+connector.Path()+`.id
				}
				email_template_id = `+reminder.Path()+`.id
			`),
			Check: settings.Check(map[string]any{"email_template_id": testacc.AttributeIsSet}),
		},
		resource.TestStep{
			Config: templates + settings.Block(`
				project_id = "`+projectID+`"
			`),
			Check: settings.Check(map[string]any{"email_template_id": ""}),
		},
	)
}
