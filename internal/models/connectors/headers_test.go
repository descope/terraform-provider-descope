package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectorSecretHeaders(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
				`headers = [
					{ key = "X-Plain", value = "visible" },
					{ key = "X-Secret", value = "topsecret", secret = true },
				]`,
			),
			Check: c.Check(map[string]any{
				"headers.#":        "2",
				"headers.0.key":    "X-Plain",
				"headers.0.value":  "visible",
				"headers.0.secret": false,
				"headers.1.key":    "X-Secret",
				"headers.1.value":  "topsecret",
				"headers.1.secret": true,
			}),
		},
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
				`headers = [
					{ key = "X-Secret", value = "rotated", secret = true },
				]`,
			),
			Check: c.Check(map[string]any{
				"headers.#":        "1",
				"headers.0.key":    "X-Secret",
				"headers.0.value":  "rotated",
				"headers.0.secret": true,
			}),
		},
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
				`headers = []`,
			),
			Check: c.Check(map[string]any{"headers.#": "0"}),
		},
	)
}

func TestConnectorSecretHeadersDuplicateKey(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")

	testacc.Run(t, resource.TestStep{
		Config: c.Config(`project_id = "`+projectID+`"`,
			`base_url = "https://example.com"`,
			`headers = [
				{ key = "X-Dup", value = "one" },
				{ key = "X-Dup", value = "two", secret = true },
			]`,
		),
		ExpectError: regexp.MustCompile(`The key "X-Dup" is used more than once`),
	})
}
