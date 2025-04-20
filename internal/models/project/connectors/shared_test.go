package connectors_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectorsShared(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"smtp": [
						{
							name = "My SMTP Connector"
							description = ""
							server = {
								host = "example.com"
								port = 587
							}
							sender = {
								email = "foo@bar.com"
								name = "Foo Bar"
							}
							authentication = {
								username = "foo"
								password = "bar"
							}
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.smtp.#":                1,
				"connectors.smtp.0.id":             testacc.AttributeMatchesPattern(`^(CI|MP)`),
				"connectors.smtp.0.name":           "My SMTP Connector",
				"connectors.smtp.0.description":    "",
				"connectors.smtp.0.use_static_ips": false,
			}),
		},
	)
}
