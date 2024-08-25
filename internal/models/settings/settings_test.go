package settings_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSettings(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "3 weeks"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "3 weeks",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "1 day"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "1 day",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					refresh_token_expiration = "1 minute"
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
			`),
			Check: p.Check(map[string]any{
				"project_settings": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "1 days",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				project_settings = {
					domain = "example.com"
				}
			`),
			Check: p.Check(map[string]any{
				"project_settings.refresh_token_expiration": "1 days",
				"project_settings.domain":                   "example.com",
			}),
		},
	)
}
