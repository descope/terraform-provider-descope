package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMagicLinkSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.MagicLinkSettings(t)
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
				"expiration_time": "3 minutes",
				"redirect_url":    "",
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				expiration_time = "10 minutes"
				redirect_url = "https://example.com/magiclink"
			`),
			Check: m.Check(map[string]any{
				"expiration_time": "10 minutes",
				"redirect_url":    "https://example.com/magiclink",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"expiration_time": "3 minutes",
				"redirect_url":    "",
			}),
		},
		// email_service and text_service are server-populated on read, outside the import contract
		resource.TestStep{
			ResourceName:            m.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(m.Path(), "project_id"),
			ImportStateVerifyIgnore: []string{"email_service", "text_service"},
		},
	)
}

func TestMagicLinkSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	m := testacc.MagicLinkSettings(t)
	testacc.Run(t,
		// a custom email connector with an active template selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "magiclink"
				name = "`+name+`"
				subject = "Sign in"
				html_body = "Follow the link in this email to sign in"
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
				"text_template_id":           "",
			}),
		},
	)
}
