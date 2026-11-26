package texttemplate_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTextTemplate(t *testing.T) {
	projectID := testacc.ProjectID(t)
	name := testacc.GenerateAlias(t)
	x := testacc.TextTemplate(t)
	testacc.Run(t,
		// a new template gets a server-assigned id
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `"
				body = "Tap to sign in: {{link}}"
			`),
			Check: x.Check(map[string]any{
				"id":         testacc.AttributeIsSet,
				"project_id": testacc.AttributeIsSet,
				"method":     "magiclink",
				"name":       name,
				"body":       "Tap to sign in: {{link}}",
			}),
		},
		// the name is mutable since the server-assigned id is the identity
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `-renamed"
				body = "Sign in: {{link}}"
			`),
			Check: x.Check(map[string]any{
				"name": name + "-renamed",
				"body": "Sign in: {{link}}",
			}),
		},
		// the System name is reserved for the built-in template
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "System"
				body = "body"
			`),
			ExpectError: regexp.MustCompile(`Invalid text template`),
		},
		// end on a valid configuration so the post-test destroy passes config validation
		resource.TestStep{
			Config: x.Block(`
				project_id = "` + projectID + `"
				method = "magiclink"
				name = "` + name + `"
				body = "Tap to sign in: {{link}}"
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
