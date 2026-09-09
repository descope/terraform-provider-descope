package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEnchantedLinkSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.EnchantedLinkSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                 testacc.AttributeIsSet,
				"project_id":         testacc.AttributeIsSet,
				"disabled":           false,
				"expiration_time":    "3 minutes",
				"redirect_url":       "",
				"email_connector_id": "Descope",
				"text_connector_id":  "Descope",
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				expiration_time = "10 minutes"
				redirect_url = "https://example.com/enchantedlink"
			`),
			Check: m.Check(map[string]any{
				"expiration_time": "10 minutes",
				"redirect_url":    "https://example.com/enchantedlink",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"expiration_time":    "3 minutes",
				"redirect_url":       "",
				"email_connector_id": "Descope",
				"text_connector_id":  "Descope",
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

func TestEnchantedLinkSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	x := testacc.TextTemplate(t)
	m := testacc.EnchantedLinkSettings(t)
	testacc.Run(t,
		// a custom email connector with active email and text templates selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "enchantedlink"
				name = "`+name+`"
				subject = "Sign in"
				html_body = "Follow the link in this email to sign in"
			`) + x.Block(`
				project_id = "`+projectID+`"
				method = "enchantedlink"
				name = "`+name+`-text"
				body = "Tap to sign in: {{link}}"
			`) + m.Block(`
				project_id = "`+projectID+`"
				email_connector_id = `+c.Path()+`.id
				email_template_id = `+e.Path()+`.id
				text_template_id = `+x.Path()+`.id
			`),
			Check: m.Check(map[string]any{
				"email_connector_id": testacc.AttributeIsSet,
				"email_template_id":  testacc.AttributeIsSet,
				"text_template_id":   testacc.AttributeIsSet,
			}),
		},
	)
}
