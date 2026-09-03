package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	m := testacc.ProjectSettings(t)
	testacc.Run(t,
		// defaults only
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"id":                                  testacc.AttributeIsSet,
				"project_id":                          testacc.AttributeIsSet,
				"app_url":                             "",
				"custom_domain":                       "",
				"approved_domains":                    []string{},
				"default_no_sso_apps":                 false,
				"tenant_user_isolation":               false,
				"allow_auth_hosting_iframe_embedding": false,
				"test_users_loginid_regexp":           "",
				"test_users_static_otp":               "",
				"test_users_verifier_regexp":          "",
			}),
		},
		// static otp and verifier regexp must be set together
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				test_users_static_otp = "123456"
			`),
			ExpectError: regexp.MustCompile(`set together`),
		},
		// update the settings fields
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				app_url = "https://app.example.com"
				approved_domains = ["example.com", "*.example.dev"]
				default_no_sso_apps = true
				allow_auth_hosting_iframe_embedding = true
				test_users_loginid_regexp = ".*@test\\.example\\.com"
				test_users_static_otp = "123456"
				test_users_verifier_regexp = ".*"
			`),
			Check: m.Check(map[string]any{
				"app_url":                             "https://app.example.com",
				"approved_domains":                    []string{"example.com", "*.example.dev"},
				"default_no_sso_apps":                 true,
				"allow_auth_hosting_iframe_embedding": true,
				"test_users_loginid_regexp":           ".*@test\\.example\\.com",
				"test_users_static_otp":               "123456",
				"test_users_verifier_regexp":          ".*",
			}),
		},
		// removing optional fields reverts them to their defaults
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
			`),
			Check: m.Check(map[string]any{
				"app_url":                             "",
				"approved_domains":                    []string{},
				"default_no_sso_apps":                 false,
				"allow_auth_hosting_iframe_embedding": false,
				"test_users_loginid_regexp":           "",
				"test_users_static_otp":               "",
				"test_users_verifier_regexp":          "",
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
