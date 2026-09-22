package connectors_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestConnectorSecretHeaders(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
				`headers = { "X-Plain" = "visible" }`,
				`secret_headers = { "X-Api-Key" = "topsecret" }`,
			),
			Check: c.Check(map[string]any{
				"headers.X-Plain":          "visible",
				"secret_headers.X-Api-Key": "topsecret",
			}),
		},
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
				`headers = { "X-Plain" = "changed" }`,
				`secret_headers = { "X-Api-Key" = "rotated" }`,
			),
			Check: c.Check(map[string]any{
				"headers.X-Plain":          "changed",
				"secret_headers.X-Api-Key": "rotated",
			}),
		},
		resource.TestStep{
			Config: c.Config(project,
				`base_url = "https://example.com"`,
			),
			Check: c.Check(map[string]any{
				"headers.%":        "0",
				"secret_headers.%": "0",
			}),
		},
	)
}

func TestConnectorSecretHeadersDuplicateKey(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")

	testacc.Run(t, resource.TestStep{
		Config: c.Config(`project_id = "`+projectID+`"`,
			`base_url = "https://example.com"`,
			`headers = { "X-Dup" = "plain" }`,
			`secret_headers = { "X-Dup" = "secret" }`,
		),
		ExpectError: regexp.MustCompile(`must not both set the "X-Dup" key`),
	})
}

func TestConnectorSecretHeadersImport(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")
	config := c.Config(`project_id = "`+projectID+`"`,
		`base_url = "https://example.com"`,
		`headers = { "X-Plain" = "visible" }`,
		`secret_headers = { "X-Api-Key" = "topsecret" }`,
	)

	testacc.Run(t,
		resource.TestStep{
			Config: config,
		},
		resource.TestStep{
			Config:            config,
			ResourceName:      c.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(c.Path(), "project_id", "id"),
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 {
					return fmt.Errorf("expected a single imported state, got %d", len(states))
				}
				attrs := states[0].Attributes
				if v := attrs["headers.X-Plain"]; v != "visible" {
					return fmt.Errorf("expected the plain header to be imported, got %q", v)
				}
				if v := attrs["secret_headers.%"]; v != "0" {
					return fmt.Errorf("expected no secret headers to be imported, got %s entries", v)
				}
				return nil
			},
		},
	)
}
