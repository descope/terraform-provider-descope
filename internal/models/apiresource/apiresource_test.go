package apiresource_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAPIResource(t *testing.T) {
	projectID := testacc.ProjectID(t)
	r := testacc.APIResource(t)
	project := `project_id = "` + projectID + `"`

	createStep := resource.TestStep{
		Config: r.Config(project, `uri = "https://`+r.Name+`.example.com"`, `
			permissions_scopes = [{ name = "read", description = "Read access" }]
			id_jag_target = {
				authorization_server = "https://as.example.com"
				continuation_trust = [{
					issuer          = "https://ras.example.com"
					jwks_uri        = "https://ras.example.com/jwks"
					sign_algorithm  = "ES256"
					actor_client_id = "gateway-1"
				}]
			}
		`),
		Check: r.Check(map[string]any{
			"project_id":                                 projectID,
			"name":                                       r.Name,
			"permissions_scopes.#":                       "1",
			"id_jag_target.authorization_server":         "https://as.example.com",
			"id_jag_target.continuation_trust.#":         "1",
			"id_jag_target.continuation_trust.0.issuer":  "https://ras.example.com",
			"id_jag_target.continuation_trust.0.enabled": "true",
		}),
	}

	// Disables the first entry and adds a second one
	updateStep := resource.TestStep{
		Config: r.Config(project, `uri = "https://`+r.Name+`.example.com"`, `
			id_jag_target = {
				authorization_server = "https://as.example.com"
				continuation_trust = [{
					issuer          = "https://ras.example.com"
					jwks_uri        = "https://ras.example.com/jwks"
					sign_algorithm  = "ES256"
					actor_client_id = "gateway-1"
					enabled         = false
				}, {
					issuer          = "https://ras.example.com"
					jwks_uri        = "https://ras.example.com/jwks"
					sign_algorithm  = "RS256"
					actor_client_id = "gateway-2"
				}]
			}
		`),
		Check: r.Check(map[string]any{
			"permissions_scopes.#":                               "0",
			"id_jag_target.continuation_trust.#":                 "2",
			"id_jag_target.continuation_trust.0.enabled":         "false",
			"id_jag_target.continuation_trust.1.actor_client_id": "gateway-2",
		}),
	}

	// Removing the target clears it on the backend
	clearStep := resource.TestStep{
		Config: r.Config(project, `uri = "https://`+r.Name+`.example.com"`),
		Check:  r.Check(map[string]any{"id_jag_target": testacc.AttributeIsNotSet}),
	}

	importStep := resource.TestStep{
		ResourceName:      r.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(r.Path(), "project_id", "id"),
	}

	testacc.RunWithDestroyCheck(t, "descope_api_resource", createStep, updateStep, clearStep, importStep)
}

func TestAPIResourceValidation(t *testing.T) {
	r := testacc.APIResource(t)

	testacc.Run(t, resource.TestStep{
		Config: r.Config(`project_id = "P123"`, `uri = "https://x.example.com"`, `
			id_jag_target = {
				authorization_server = "https://as.example.com"
				continuation_trust = [{ issuer = "https://ras.example.com", jwks_uri = "https://ras.example.com/jwks", sign_algorithm = "HS256", actor_client_id = "a" }]
			}
		`),
		ExpectError: regexp.MustCompile(`sign_algorithm`),
	})
}
