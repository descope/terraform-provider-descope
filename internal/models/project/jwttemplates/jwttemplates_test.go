package jwttemplates_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestJWTTemplates(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = []
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
							auth_schema = "tenantOnly"
							empty_claim_policy = "delete"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":                    1,
				"jwt_templates.user_templates.0.id":                 testacc.AttributeHasPrefix("JT"),
				"jwt_templates.user_templates.0.name":               "foo",
				"jwt_templates.user_templates.0.description":        "a",
				"jwt_templates.user_templates.0.auto_tenant_claim":  false,
				"jwt_templates.user_templates.0.conformance_issuer": false,
				"jwt_templates.user_templates.0.enforce_issuer":     false,
				"jwt_templates.user_templates.0.auth_schema":        "tenantOnly",
				"jwt_templates.user_templates.0.empty_claim_policy": "delete",
				"jwt_templates.user_templates.0.template":           "{}",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
			jwt_templates = {
				user_templates = [
					{
						name = "foo"
						description = "a"
						template = "{}"
						auth_schema = "tenantOnly"
						empty_claim_policy = "delete"
						auto_tenant_claim = true
						conformance_issuer = true
						enforce_issuer = true
						exclude_permission_claim = true
						override_subject_claim = "custom-sub"
						add_jti_claim = true
					}
				]
			}
		`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":                          1,
				"jwt_templates.user_templates.0.id":                       testacc.AttributeHasPrefix("JT"),
				"jwt_templates.user_templates.0.name":                     "foo",
				"jwt_templates.user_templates.0.description":              "a",
				"jwt_templates.user_templates.0.auto_tenant_claim":        true,
				"jwt_templates.user_templates.0.conformance_issuer":       true,
				"jwt_templates.user_templates.0.enforce_issuer":           true,
				"jwt_templates.user_templates.0.exclude_permission_claim": true,
				"jwt_templates.user_templates.0.override_subject_claim":   "custom-sub",
				"jwt_templates.user_templates.0.add_jti_claim":            true,
				"jwt_templates.user_templates.0.auth_schema":              "tenantOnly",
				"jwt_templates.user_templates.0.empty_claim_policy":       "delete",
				"jwt_templates.user_templates.0.template":                 "{}",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
						}
					]
					access_key_templates = [
						{
							name = "foo"
							description = "b"
							template = "{}"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Conflicting Attribute Values`),
		},
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
						}
					]
					access_key_templates = [
						{
							name = "bar"
							description = "b"
							template = "{}"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":            1,
				"jwt_templates.user_templates.0.id":         testacc.AttributeHasPrefix("JT"),
				"jwt_templates.user_templates.0.name":       "foo",
				"jwt_templates.access_key_templates.#":      1,
				"jwt_templates.access_key_templates.0.id":   testacc.AttributeHasPrefix("JT"),
				"jwt_templates.access_key_templates.0.name": "bar",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					user_jwt_template = "qux"
				}
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
						}
					]
					access_key_templates = [
						{
							name = "bar"
							description = "b"
							template = "{}"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Unknown JWT template reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					user_jwt_template = "bar"
				}
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
						}
					]
					access_key_templates = [
						{
							name = "bar"
							description = "b"
							template = "{}"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid JWT template reference`),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					user_jwt_template = "foo"
				}
				jwt_templates = {
					user_templates = [
						{
							name = "foo",
							description = "a"
							template = "{}"
						}
					]
					access_key_templates = [
						{
							name = "bar"
							description = "b"
							template = "{}"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.user_jwt_template": "foo",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				jwt_templates = {
					user_templates = [
						{
							name = "foo"
							description = "a"
							template = "{}"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":       1,
				"jwt_templates.access_key_templates.#": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					user_jwt_template = "foo"
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":       1,
				"jwt_templates.access_key_templates.#": 0,
			}),
		},
		resource.TestStep{ // XXX project_settings is optional, so we can't just remove it as a way to clear user_jwt_template
			Config: p.Config(`
				project_settings = {
					user_jwt_template = ""
				}
				jwt_templates = {
					user_templates = []
				}
			`),
			Check: p.Check(map[string]any{
				"jwt_templates.user_templates.#":       0,
				"jwt_templates.access_key_templates.#": 0,
			}),
		},
	)
}
