package settings_test

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/settings"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				"role_policies.#":         0,
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
				"min_length":         8,
				"non_alphanumeric":   true,
				"expiration":         false,
				"lock":               false,
				"enforce_strength":   "none",
				"mask_errors":        true,
				"email_connector_id": "",
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
		// a role policy must list at least one role
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = [] }]
			`),
			ExpectError: regexp.MustCompile(`set\s+must\s+contain\s+at\s+least\s+1\s+elements`),
		},
		// a role policy can't require a minimum length above the maximum password length
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = ["admin"], min_length = 65 }]
			`),
			ExpectError: regexp.MustCompile(`must\s+be\s+between\s+0\s+and\s+64`),
		},
		// a role policy can't list more than 100 roles
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = [for i in range(101) : "role-${i}"] }]
			`),
			ExpectError: regexp.MustCompile(`set\s+must\s+contain\s+at\s+least\s+1\s+elements\s+and\s+at\s+most\s+100\s+elements`),
		},
		// a role name in a role policy can't be longer than 100 characters
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = ["` + strings.Repeat("r", 101) + `"] }]
			`),
			ExpectError: regexp.MustCompile(`string\s+length\s+must\s+be\s+between\s+1\s+and\s+100`),
		},
		// a role policy can't disallow more than 1024 characters
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = ["admin"], disallowed_characters = "` + strings.Repeat("x", 1025) + `" }]
			`),
			ExpectError: regexp.MustCompile(`string\s+length\s+must\s+be\s+at\s+most\s+1024`),
		},
		// there can't be more than 50 role policies
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [for i in range(51) : { roles = ["admin"] }]
			`),
			ExpectError: regexp.MustCompile(`list\s+must\s+contain\s+at\s+most\s+50\s+elements`),
		},
		// the roles of a role policy must exist in the project
		resource.TestStep{
			Config: m.Block(`
				project_id = "` + projectID + `"
				role_policies = [{ roles = ["` + testacc.GenerateAlias(t) + `"] }]
			`),
			ExpectError: regexp.MustCompile(`Role\s+not\s+found`),
		},
	)
}

func TestPasswordSettingsRolePolicies(t *testing.T) {
	projectID := testacc.ProjectID(t)
	project := `project_id = "` + projectID + `"`
	admin := testacc.Role(t)
	admin.ID = "admin"
	support := testacc.Role(t)
	support.ID = "support"
	roles := admin.Config(project) + support.Config(project)
	adminName := admin.Name
	admin.Name += "-renamed"
	renamedRoles := admin.Config(project) + support.Config(project)
	m := testacc.PasswordSettings(t)
	testacc.Run(t,
		// a policy for each role, referencing the roles by name
		resource.TestStep{
			Config: roles + m.Block(`
				project_id = "`+projectID+`"
				role_policies = [
					{
						roles = [`+admin.Path()+`.name]
						min_length = 16
						non_alphanumeric = true
						disallowed_characters = "<>"
						disallow_email_match = true
						enforce_strength = "strong"
					},
					{
						roles = [`+support.Path()+`.name]
						uppercase = true
						number = true
					},
				]
			`),
			Check: m.Check(map[string]any{
				"role_policies.#":                       2,
				"role_policies.0.roles":                 []string{adminName},
				"role_policies.0.min_length":            16,
				"role_policies.0.lowercase":             false,
				"role_policies.0.non_alphanumeric":      true,
				"role_policies.0.disallowed_characters": "<>",
				"role_policies.0.disallow_email_match":  true,
				"role_policies.0.enforce_strength":      "strong",
				"role_policies.1.roles":                 []string{support.Name},
				"role_policies.1.min_length":            0,
				"role_policies.1.uppercase":             true,
				"role_policies.1.number":                true,
				"role_policies.1.any_letter":            false,
				"role_policies.1.disallowed_characters": "",
				"role_policies.1.enforce_strength":      "none",
			}),
		},
		// a renamed role stays referenced, and both roles share one policy with updated rules
		resource.TestStep{
			Config: renamedRoles + m.Block(`
				project_id = "`+projectID+`"
				role_policies = [
					{
						roles = [`+admin.Path()+`.name, `+support.Path()+`.name]
						min_length = 20
						lowercase = true
						any_letter = true
						enforce_strength = "very_strong"
					},
				]
			`),
			Check: m.Check(map[string]any{
				"role_policies.#":                       1,
				"role_policies.0.roles":                 []string{admin.Name, support.Name},
				"role_policies.0.min_length":            20,
				"role_policies.0.lowercase":             true,
				"role_policies.0.uppercase":             false,
				"role_policies.0.any_letter":            true,
				"role_policies.0.non_alphanumeric":      false,
				"role_policies.0.disallowed_characters": "",
				"role_policies.0.disallow_email_match":  false,
				"role_policies.0.enforce_strength":      "very_strong",
			}),
		},
		// listing the roles of a policy in another order plans no changes
		resource.TestStep{
			Config: renamedRoles + m.Block(`
				project_id = "`+projectID+`"
				role_policies = [
					{
						roles = [`+support.Path()+`.name, `+admin.Path()+`.name]
						min_length = 20
						lowercase = true
						any_letter = true
						enforce_strength = "very_strong"
					},
				]
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
		},
		resource.TestStep{
			ResourceName:      m.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(m.Path(), "project_id"),
		},
		// removing the attribute clears all role policies
		resource.TestStep{
			Config: renamedRoles + m.Block(`
				project_id = "`+projectID+`"
			`),
			Check: m.Check(map[string]any{
				"role_policies.#": 0,
			}),
		},
	)
}

func TestPasswordSettingsRolePoliciesValues(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &settings.PasswordSettingsModel{RolePolicies: listattr.Value([]*settings.PasswordRolePolicyModel{strictRolePolicy(), defaultRolePolicy()})}
	assert.JSONEq(t, `[
		{"roleNames": ["Admin", "Support"], "minLength": 16, "lowercase": true, "uppercase": true, "number": true, "nonAlphanumeric": true, "anyLetter": true, "disallowedCharacters": "<>", "disallowEmailMatch": true, "enablePasswordStrength": true, "passwordStrengthScore": 3},
		{"roleNames": ["Auditor"], "minLength": 0, "lowercase": false, "uppercase": false, "number": false, "nonAlphanumeric": false, "anyLetter": false, "disallowedCharacters": "", "disallowEmailMatch": false, "enablePasswordStrength": false, "passwordStrengthScore": 0}
	]`, rolePoliciesJSON(t, m.Values(h)))

	m.RolePolicies = listattr.Empty[settings.PasswordRolePolicyModel]()
	assert.JSONEq(t, `[]`, rolePoliciesJSON(t, m.Values(h)))

	require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
}

func TestPasswordSettingsRolePoliciesSetValues(t *testing.T) {
	strict := `{"roleNames": ["Admin", "Support"], "minLength": 16, "lowercase": true, "uppercase": true, "number": true, "nonAlphanumeric": true, "anyLetter": true, "disallowedCharacters": "<>", "disallowEmailMatch": true, "enablePasswordStrength": true, "passwordStrengthScore": 3}`
	weakest := defaultRolePolicy()
	weakest.EnforceStrength = stringattr.Value("very_weak")
	cases := []struct {
		name     string
		current  []*settings.PasswordRolePolicyModel
		response string
		expected []*settings.PasswordRolePolicyModel
	}{
		{
			name:     "every key with the roles in another order",
			response: `{"rolePolicies": [` + strict + `]}`,
			expected: []*settings.PasswordRolePolicyModel{strictRolePolicy()},
		},
		{
			name:     "omitted keys are defaults",
			current:  []*settings.PasswordRolePolicyModel{strictRolePolicy(), defaultRolePolicy()},
			response: `{"rolePolicies": [` + strict + `, {"roleNames": ["Auditor"]}]}`,
			expected: []*settings.PasswordRolePolicyModel{strictRolePolicy(), defaultRolePolicy()},
		},
		{
			name:     "omitted keys override the current values",
			current:  []*settings.PasswordRolePolicyModel{strictRolePolicy()},
			response: `{"rolePolicies": [{"roleNames": ["Auditor"]}]}`,
			expected: []*settings.PasswordRolePolicyModel{defaultRolePolicy()},
		},
		{
			name:     "omitted strength score",
			response: `{"rolePolicies": [{"roleNames": ["Auditor"], "enablePasswordStrength": true}]}`,
			expected: []*settings.PasswordRolePolicyModel{weakest},
		},
		{
			name:     "null policies",
			current:  []*settings.PasswordRolePolicyModel{strictRolePolicy()},
			response: `{"rolePolicies": null}`,
			expected: []*settings.PasswordRolePolicyModel{},
		},
		{
			name:     "absent policies",
			current:  []*settings.PasswordRolePolicyModel{strictRolePolicy()},
			response: `{}`,
			expected: []*settings.PasswordRolePolicyModel{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var diags diag.Diagnostics
			h := helpers.NewHandler(context.Background(), &diags)
			var data map[string]any
			require.NoError(t, json.Unmarshal([]byte(c.response), &data))

			m := &settings.PasswordSettingsModel{RolePolicies: listattr.Value(c.current)}
			m.SetValues(h, data)

			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			expected := listattr.Value(c.expected)
			assert.True(t, expected.Equal(m.RolePolicies), "expected %s, got %s", expected, m.RolePolicies)
		})
	}
}

func rolePoliciesJSON(t *testing.T, data map[string]any) string {
	policies, ok := data["rolePolicies"]
	require.True(t, ok, "expected the rolePolicies key in %v", data)
	b, err := json.Marshal(policies)
	require.NoError(t, err)
	return string(b)
}

func strictRolePolicy() *settings.PasswordRolePolicyModel {
	return &settings.PasswordRolePolicyModel{
		Roles:                strsetattr.Value([]string{"Support", "Admin"}),
		MinLength:            intattr.Value(16),
		Lowercase:            boolattr.Value(true),
		Uppercase:            boolattr.Value(true),
		Number:               boolattr.Value(true),
		NonAlphanumeric:      boolattr.Value(true),
		AnyLetter:            boolattr.Value(true),
		DisallowedCharacters: stringattr.Value("<>"),
		DisallowEmailMatch:   boolattr.Value(true),
		EnforceStrength:      stringattr.Value("strong"),
	}
}

func defaultRolePolicy() *settings.PasswordRolePolicyModel {
	return &settings.PasswordRolePolicyModel{
		Roles:                strsetattr.Value([]string{"Auditor"}),
		MinLength:            intattr.Value(0),
		Lowercase:            boolattr.Value(false),
		Uppercase:            boolattr.Value(false),
		Number:               boolattr.Value(false),
		NonAlphanumeric:      boolattr.Value(false),
		AnyLetter:            boolattr.Value(false),
		DisallowedCharacters: stringattr.Value(""),
		DisallowEmailMatch:   boolattr.Value(false),
		EnforceStrength:      stringattr.Value("none"),
	}
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
				email_connector_id = `+c.Path()+`.id
				email_template_id = `+e.Path()+`.id
			`),
			Check: m.Check(map[string]any{
				"email_connector_id": testacc.AttributeIsSet,
				"email_template_id":  testacc.AttributeIsSet,
			}),
		},
	)
}
