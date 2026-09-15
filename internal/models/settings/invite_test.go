package settings_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestInviteSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.InviteSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                   testacc.AttributeIsSet,
				"project_id":           testacc.AttributeIsSet,
				"require_invitation":   false,
				"invite_url":           "",
				"add_magiclink_token":  false,
				"expire_invited_users": false,
				"invite_expiration":    "1 week",
				"send_email":           true,
				"send_text":            false,
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				require_invitation = true
				invite_url = "https://example.com/invite"
				add_magiclink_token = true
				expire_invited_users = true
				invite_expiration = "2 weeks"
				send_text = true
			`),
			Check: m.Check(map[string]any{
				"require_invitation":   true,
				"invite_url":           "https://example.com/invite",
				"add_magiclink_token":  true,
				"expire_invited_users": true,
				"invite_expiration":    "2 weeks",
				"send_email":           true,
				"send_text":            true,
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"require_invitation":   false,
				"invite_url":           "",
				"add_magiclink_token":  false,
				"expire_invited_users": false,
				"invite_expiration":    "1 week",
				"send_email":           true,
				"send_text":            false,
			}),
		},
		// invite_expiration unit is normalized to plural on read
		resource.TestStep{
			ResourceName:            m.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(m.Path(), "project_id"),
			ImportStateVerifyIgnore: []string{"invite_expiration"},
		},
	)
}

func TestInviteSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	m := testacc.InviteSettings(t)
	testacc.Run(t,
		// a custom email connector with an active template selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "invite"
				name = "`+name+`"
				subject = "You're invited"
				html_body = "Follow the link in this email to join"
			`) + m.Block(`
				project_id = "`+projectID+`"
				email_connector_id = `+c.Path()+`.id
				email_template_id = `+e.Path()+`.id
			`),
			Check: m.Check(map[string]any{
				"email_connector_id": testacc.AttributeIsSet,
				"email_template_id":  testacc.AttributeIsSet,
			}),
		},
		// dropping the references resets delivery and the template to the built-in ones
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + m.Block(`
				project_id = "`+projectID+`"
			`),
			Check: m.Check(map[string]any{
				"send_email":         true,
				"email_connector_id": "",
				"email_template_id":  "",
			}),
		},
	)
}
