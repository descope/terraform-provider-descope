package apps_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSAMLApp(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.SAMLApp(t)
	project := `project_id = "` + projectID + `"`

	var appID string
	captureID := func(s string) error {
		appID = s
		return nil
	}
	sameID := func(s string) error {
		if s != appID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", appID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`manual_configuration = {
				acs_url = "https://sp.example.com/acs"
				entity_id = "sp-entity"
			}`,
		),
		Check: a.Check(map[string]any{
			"id":                               captureID,
			"project_id":                       projectID,
			"name":                             a.Name,
			"disabled":                         false,
			"login_page_url":                   testacc.AttributeIsSet,
			"manual_configuration.acs_url":     "https://sp.example.com/acs",
			"manual_configuration.entity_id":   "sp-entity",
			"manual_configuration.certificate": "",
			"subject_name_id_type":             "",
			"default_signature_algorithm":      "",
		}),
	}

	a.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`manual_configuration = {
				acs_url = "https://sp.example.com/acs2"
				entity_id = "sp-entity-2"
			}`,
			`acs_allowed_callback_urls = ["https://sp.example.com/*"]`,
			`subject_name_id_type = "email"`,
			`default_signature_algorithm = "sha256"`,
			`default_relay_state = "https://sp.example.com/home"`,
			`attribute_mapping = [{ name = "email", value = "user.email" }]`,
			`force_authentication = true`,
		),
		Check: a.Check(map[string]any{
			"id":                             sameID,
			"name":                           a.Name,
			"manual_configuration.acs_url":   "https://sp.example.com/acs2",
			"manual_configuration.entity_id": "sp-entity-2",
			"acs_allowed_callback_urls":      []string{"https://sp.example.com/*"},
			"subject_name_id_type":           "email",
			"default_signature_algorithm":    "sha256",
			"default_relay_state":            "https://sp.example.com/home",
			"attribute_mapping.0.name":       "email",
			"attribute_mapping.0.value":      "user.email",
			"force_authentication":           true,
		}),
	}

	dynamicStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`dynamic_configuration = {
				metadata_url = "https://sp.example.com/metadata"
			}`,
		),
		Check: a.Check(map[string]any{
			"id":                                 sameID,
			"dynamic_configuration.metadata_url": "https://sp.example.com/metadata",
		}),
	}

	importStep := resource.TestStep{
		ResourceName:            a.Path(),
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateIdFunc:       testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		ImportStateVerifyIgnore: []string{"deletion_protection"},
	}

	testacc.RunWithDestroyCheck(t, "descope_saml_app", createStep, updateStep, dynamicStep, importStep)
}

func TestSAMLAppInvalid(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.SAMLApp(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config:      a.Config(project),
			ExpectError: regexp.MustCompile(`Missing Attribute Value`),
		},
		resource.TestStep{
			Config: a.Config(project,
				`dynamic_configuration = {
					metadata_url = "https://sp.example.com/metadata"
				}`,
				`manual_configuration = {
					acs_url = "https://sp.example.com/acs"
					entity_id = "sp-entity"
				}`,
			),
			ExpectError: regexp.MustCompile(`Conflicting Attribute Values`),
		},
	)
}

func TestSAMLAppOutOfBandDeletion(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.SAMLApp(t)
	config := a.Config(`project_id = "`+projectID+`"`,
		`deletion_protection = false`,
		`manual_configuration = {
			acs_url = "https://sp.example.com/acs"
			entity_id = "sp-entity"
		}`,
	)

	var appID string
	testacc.Run(t,
		resource.TestStep{
			Config: config,
			Check: a.Check(map[string]any{
				"name": a.Name,
				"id": func(s string) error {
					appID = s
					return nil
				},
			}),
		},
		resource.TestStep{
			PreConfig: func() {
				testacc.OutOfBandPost(t, projectID, "/v1/mgmt/sso/idp/app/delete", map[string]any{"id": appID})
			},
			Config: config,
			Check: a.Check(map[string]any{
				"name": a.Name,
				"id": func(s string) error {
					if s == appID {
						return fmt.Errorf("expected the app to be re-created with a new id, got the same id %s", s)
					}
					return nil
				},
			}),
		},
		resource.TestStep{
			Config:        config,
			ResourceName:  a.Path(),
			ImportState:   true,
			ImportStateId: projectID + "/bogusbogusbogusbogus",
			ExpectError:   regexp.MustCompile(`Error reading saml_app`),
		},
	)
}
