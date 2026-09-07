package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSSOSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.SSOSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                            testacc.AttributeIsSet,
				"project_id":                    testacc.AttributeIsSet,
				"disabled":                      false,
				"merge_users":                   false,
				"redirect_url":                  "",
				"allow_duplicate_domains":       false,
				"groups_priority":               false,
				"require_sso_domains":           false,
				"require_groups_attribute_name": false,
				"mandatory_user_attributes.#":   "0",
				"sso_suite_settings.style_id":   "",
				"sso_suite_settings.hide_saml":  false,
			}),
		},
		// update the plain settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				merge_users = true
				redirect_url = "https://example.com/sso"
				allow_duplicate_domains = true
				groups_priority = true
				require_sso_domains = true
				require_groups_attribute_name = true
				allow_merge_users_with_multiple_tenants = true
				mandatory_user_attributes = [
					{
						id = "email"
					},
					{
						id = "department"
						custom = true
					}
				]
				sso_suite_settings = {
					hide_scim = true
					hide_role_mapping = true
					hide_fga_mapping = true
					show_xaa = true
					support_email = "help@example.com"
					show_help_contact = true
				}
			`),
			Check: m.Check(map[string]any{
				"merge_users":                             true,
				"redirect_url":                            "https://example.com/sso",
				"allow_duplicate_domains":                 true,
				"groups_priority":                         true,
				"require_sso_domains":                     true,
				"require_groups_attribute_name":           true,
				"allow_merge_users_with_multiple_tenants": true,
				"mandatory_user_attributes.#":             "2",
				"mandatory_user_attributes.0.id":          "email",
				"mandatory_user_attributes.0.custom":      false,
				"mandatory_user_attributes.1.id":          "department",
				"mandatory_user_attributes.1.custom":      true,
				"sso_suite_settings.hide_scim":            true,
				"sso_suite_settings.hide_role_mapping":    true,
				"sso_suite_settings.hide_fga_mapping":     true,
				"sso_suite_settings.show_xaa":             true,
				"sso_suite_settings.support_email":        "help@example.com",
				"sso_suite_settings.show_help_contact":    true,
			}),
		},
		// conflicting sso suite settings are rejected
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				sso_suite_settings = {
					hide_saml = true
					hide_oidc = true
				}
			`),
			ExpectError: regexp.MustCompile(`hide_oidc and hide_saml cannot both be true`),
		},
		// the suite must offer either SSO or SCIM configuration
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				sso_suite_settings = {
					hide_sso = true
					hide_scim = true
				}
			`),
			ExpectError: regexp.MustCompile(`hide_sso and hide_scim cannot both be true`),
		},
		// the deprecated combined flag conflicts with the split pair even when values are false
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				sso_suite_settings = {
					hide_groups_mapping = true
					hide_role_mapping = false
				}
			`),
			ExpectError: regexp.MustCompile(`hide_groups_mapping attribute cannot be combined`),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"merge_users":                             false,
				"redirect_url":                            "",
				"allow_duplicate_domains":                 false,
				"groups_priority":                         false,
				"require_sso_domains":                     false,
				"require_groups_attribute_name":           false,
				"allow_merge_users_with_multiple_tenants": false,
				"mandatory_user_attributes.#":             "0",
				"sso_suite_settings.hide_scim":            false,
				"sso_suite_settings.hide_role_mapping":    false,
				"sso_suite_settings.hide_fga_mapping":     false,
				"sso_suite_settings.show_xaa":             false,
				"sso_suite_settings.support_email":        "",
				"sso_suite_settings.show_help_contact":    false,
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

func TestSSOSettingsTemplates(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "generic_email_gateway_connector")
	name := testacc.GenerateAlias(t)
	e := testacc.EmailTemplate(t)
	m := testacc.SSOSettings(t)
	testacc.Run(t,
		// a custom email connector with an active template selected by reference
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + e.Block(`
				project_id = "`+projectID+`"
				method = "sso"
				name = "`+name+`"
				subject = "SSO invitation"
				html_body = "Follow the link in this email to connect"
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
		// Dropping the reference resets the selection to the built-in System template; the settings update runs first, so the same apply can destroy it.
		resource.TestStep{
			Config: c.Config(`
				project_id = "`+projectID+`"
				post_url = "https://mail.example.com/send"
			`) + m.Block(`
				project_id = "`+projectID+`"
			`),
			Check: m.Check(map[string]any{"email_template_id": ""}),
		},
	)
}
