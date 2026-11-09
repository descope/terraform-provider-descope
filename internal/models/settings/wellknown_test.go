package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWellKnownSettings(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.WellKnownSettings(t)
	testacc.Run(t,
		// create with defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"id":                        testacc.AttributeIsSet,
				"enabled":                   false,
				"hidden_fields":             []string{},
				"removed_values.#":          0,
				"scopes_supported.extra":    []string{},
				"scopes_supported.excluded": []string{},
				"claims_supported.extra":    []string{},
				"claims_supported.excluded": []string{},
			}),
		},
		// set every field
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				enabled = true
				hidden_fields = ["end_session_endpoint", "userinfo_endpoint"]
				removed_values = [
					{ field = "grant_types_supported", values = ["client_credentials"] },
					{ field = "token_endpoint_auth_methods_supported", values = ["client_secret_post"] },
				]
				scopes_supported = { extra = ["billing"], excluded = ["phone"] }
				claims_supported = { extra = ["department"] }
			`),
			Check: o.Check(map[string]any{
				"enabled":                   true,
				"hidden_fields":             []string{"end_session_endpoint", "userinfo_endpoint"},
				"removed_values.#":          2,
				"scopes_supported.extra":    []string{"billing"},
				"scopes_supported.excluded": []string{"phone"},
				"claims_supported.extra":    []string{"department"},
				"claims_supported.excluded": []string{},
			}, resource.TestCheckTypeSetElemNestedAttrs(o.Path(), "removed_values.*", map[string]string{
				"field":    "grant_types_supported",
				"values.#": "1",
				"values.0": "client_credentials",
			})),
		},
		// turn it off while keeping the customization
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				enabled = false
				hidden_fields = ["end_session_endpoint", "userinfo_endpoint"]
				removed_values = [
					{ field = "grant_types_supported", values = ["client_credentials"] },
					{ field = "token_endpoint_auth_methods_supported", values = ["client_secret_post"] },
				]
				scopes_supported = { extra = ["billing"], excluded = ["phone"] }
				claims_supported = { extra = ["department"] }
			`),
			Check: o.Check(map[string]any{
				"enabled":          false,
				"hidden_fields":    []string{"end_session_endpoint", "userinfo_endpoint"},
				"removed_values.#": 2,
			}),
		},
		// import using <project_id>/<id> (id == project_id for this singleton)
		resource.TestStep{
			ResourceName:      o.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id"),
		},
		// a field that cannot be hidden is rejected by the backend
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
				hidden_fields = ["issuer"]
			`),
			ExpectError: regexp.MustCompile(`cannot be hidden`),
		},
		// restore the shared project to its defaults
		resource.TestStep{
			Config: o.Block(`
				project_id = "` + projectID + `"
			`),
			Check: o.Check(map[string]any{
				"enabled":                false,
				"hidden_fields":          []string{},
				"removed_values.#":       0,
				"scopes_supported.extra": []string{},
			}),
		},
	)
}
