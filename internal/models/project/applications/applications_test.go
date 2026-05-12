package applications_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestApplications(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		// Sanity
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"applications.%": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
				}
			`),
			Check: p.Check(map[string]any{
				"applications.%": 3,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = []
				}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.#": 0,
				"applications.saml_applications.#": 0,
			}),
		},
		// OIDC
		resource.TestStep{
			Config: p.Config(`
				applications = {
					oidc_applications = [
						{
							name = "foo"
							description = "bar"
							logo = "https://example.com/logo.png"
							disabled = true

							login_page_url = "https://example.com/login"
							claims = ["email", "name"]
							force_authentication = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.#": 1,
				"applications.oidc_applications.0": map[string]any{
					"id":                   testacc.AttributeHasPrefix("SA"),
					"name":                 "foo",
					"description":          "bar",
					"logo":                 "https://example.com/logo.png",
					"disabled":             true,
					"login_page_url":       "https://example.com/login",
					"claims":               []string{"email", "name"},
					"force_authentication": true,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.#": 0,
				"applications.saml_applications.#": 0,
			}),
		},
		// SAML
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = [
						{
							name = "foo"
							description = "bar"
							logo = "https://example.com/logo.png"
							disabled = true

							login_page_url = "https://example.com/login"
							force_authentication = true
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Missing Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = [
						{
							name = "foo"
							description = "bar"
							logo = "https://example.com/logo.png"
							disabled = true

							login_page_url = "https://example.com/login"
							dynamic_configuration = {
								metadata_url = "https://example.com/metadata"
							}
							acs_allowed_callback_urls = ["https://example.com/bar", "https://example.com/foo"]
							subject_name_id_type = "email"
							default_relay_state = "https://example.com/relay"
							attribute_mapping = [
								{
									name = "a"
									value = "a"
								},
							]								
							force_authentication = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.saml_applications.#": 1,
				"applications.saml_applications.0": map[string]any{
					"id":             testacc.AttributeHasPrefix("SA"),
					"name":           "foo",
					"description":    "bar",
					"logo":           "https://example.com/logo.png",
					"disabled":       true,
					"login_page_url": "https://example.com/login",
					"dynamic_configuration": map[string]any{
						"metadata_url": "https://example.com/metadata",
					},
					"manual_configuration.%":    0,
					"acs_allowed_callback_urls": []string{"https://example.com/foo", "https://example.com/bar"},
					"subject_name_id_type":      "email",
					"default_relay_state":       "https://example.com/relay",
					"attribute_mapping": map[string]any{
						"#":       1,
						"0.name":  "a",
						"0.value": "a",
					},
					"force_authentication": true,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = [
						{
							name = "meh"
							description = "bar"
							logo = "https://example.com/logo.png"

							login_page_url = "https://example.com/login"
							manual_configuration = {
								acs_url = "https://example.com/acs"
								entity_id = "foo"
							}
							acs_allowed_callback_urls = ["https://example.com/foo", "https://example.com/bar"]
							subject_name_id_type = "email"
							default_relay_state = "https://example.com/relay"
							attribute_mapping = [
								{
									name = "a"
									value = "a"
								},
								{
									name = "c"
									value = "c"
								},
								{
									name = "b"
									value = "b"
								},
							]								
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.saml_applications.#": 1,
				"applications.saml_applications.0": map[string]any{
					"id":                      testacc.AttributeHasPrefix("SA"),
					"name":                    "meh",
					"description":             "bar",
					"logo":                    "https://example.com/logo.png",
					"disabled":                false,
					"login_page_url":          "https://example.com/login",
					"dynamic_configuration.%": 0,
					"manual_configuration": map[string]any{
						"acs_url":     "https://example.com/acs",
						"entity_id":   "foo",
						"certificate": "",
					},
					"acs_allowed_callback_urls": []string{"https://example.com/foo", "https://example.com/bar"},
					"subject_name_id_type":      "email",
					"default_relay_state":       "https://example.com/relay",
					"attribute_mapping": map[string]any{
						"#":       3,
						"0.name":  "a",
						"0.value": "a",
						"1.name":  "c",
						"1.value": "c",
						"2.name":  "b",
						"2.value": "b",
					},
					"force_authentication": false,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = [
						{
							name = "meh"
							description = "bar"
							logo = "https://example.com/logo.png"

							dynamic_configuration = {
								metadata_url = "https://example.com/metadata"
							}
							subject_name_id_type = ""
							default_signature_algorithm = "sha256"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.saml_applications.#": 1,
				"applications.saml_applications.0": map[string]any{
					"id":                          testacc.AttributeHasPrefix("SA"),
					"name":                        "meh",
					"description":                 "bar",
					"logo":                        "https://example.com/logo.png",
					"login_page_url":              "",
					"dynamic_configuration.%":     1,
					"manual_configuration.%":      0,
					"acs_allowed_callback_urls":   []string{},
					"subject_name_id_type":        "",
					"subject_name_id_format":      "",
					"default_relay_state":         "",
					"default_signature_algorithm": "sha256",
					"attribute_mapping.#":         0,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					saml_applications = [
						{
							name = "meh"
							description = "bar"
							logo = "https://example.com/logo.png"

							dynamic_configuration = {
								metadata_url = "https://example.com/metadata"
							}
							subject_name_id_type = ""
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.saml_applications.#": 1,
				"applications.saml_applications.0": map[string]any{
					"id":                          testacc.AttributeHasPrefix("SA"),
					"name":                        "meh",
					"description":                 "bar",
					"logo":                        "https://example.com/logo.png",
					"login_page_url":              "",
					"dynamic_configuration.%":     1,
					"manual_configuration.%":      0,
					"acs_allowed_callback_urls":   []string{},
					"subject_name_id_type":        "",
					"subject_name_id_format":      "",
					"default_relay_state":         "",
					"default_signature_algorithm": "",
					"attribute_mapping.#":         0,
				},
			}),
		},
		// WS-Fed
		resource.TestStep{
			Config: p.Config(`
				applications = {
					wsfed_applications = [
						{
							name = "foo"
							description = "bar"
							logo = "https://example.com/logo.png"
							disabled = true

							realm = "https://example.com/realm"
							reply_url = "https://example.com/reply"
							login_page_url = "https://example.com/login"
							logout_redirect_url = "https://example.com/logout"
							error_redirect_url = "https://example.com/error"
							attribute_mapping = [
								{
									name = "email"
									value = "user.email"
								},
							]
							groups_mapping = [
								{
									name = "admins"
									type = "basic"
									filter_type = "roles"
									value = "admin"
									roles = [
										{
											id = "r1"
											name = "Admin"
										},
									]
								},
							]
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.wsfed_applications.#": 1,
				"applications.wsfed_applications.0": map[string]any{
					"id":                   testacc.AttributeHasPrefix("SA"),
					"name":                 "foo",
					"description":          "bar",
					"logo":                 "https://example.com/logo.png",
					"disabled":             true,
					"realm":                "https://example.com/realm",
					"reply_url":            "https://example.com/reply",
					"login_page_url":       "https://example.com/login",
					"logout_redirect_url":  "https://example.com/logout",
					"error_redirect_url":   "https://example.com/error",
					"force_authentication": false,
					"attribute_mapping": map[string]any{
						"#":       1,
						"0.name":  "email",
						"0.value": "user.email",
					},
					"groups_mapping": map[string]any{
						"#":              1,
						"0.name":         "admins",
						"0.type":         "basic",
						"0.filter_type":  "roles",
						"0.value":        "admin",
						"0.roles.#":      1,
						"0.roles.0.id":   "r1",
						"0.roles.0.name": "Admin",
					},
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {
					wsfed_applications = [
						{
							name = "foo"
							description = "updated"
							logo = "https://example.com/logo.png"

							realm = "https://example.com/realm2"
							reply_url = "https://example.com/reply"
							login_page_url = "https://example.com/login"
							force_authentication = true
							logout_redirect_url = "https://example.com/logout-updated"
							error_redirect_url = "https://example.com/error-updated"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.wsfed_applications.#": 1,
				"applications.wsfed_applications.0": map[string]any{
					"id":                   testacc.AttributeHasPrefix("SA"),
					"name":                 "foo",
					"description":          "updated",
					"logo":                 "https://example.com/logo.png",
					"disabled":             false,
					"realm":                "https://example.com/realm2",
					"reply_url":            "https://example.com/reply",
					"login_page_url":       "https://example.com/login",
					"force_authentication": true,
					"logout_redirect_url":  "https://example.com/logout-updated",
					"error_redirect_url":   "https://example.com/error-updated",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				applications = {}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.#":  0,
				"applications.saml_applications.#":  0,
				"applications.wsfed_applications.#": 0,
			}),
		},
		// Permissions and roles on an OIDC app
		resource.TestStep{
			Config: p.Config(`
				applications = {
					oidc_applications = [
						{
							name = "fed-app"
							description = "with auth"

							permissions = [
								{ name = "read:data", description = "read data" },
								{ name = "write:data" },
							]
							roles = [
								{
									name = "reader"
									description = "read-only"
									permissions = ["read:data"]
								},
								{
									name = "writer"
									permissions = ["read:data", "write:data"]
								},
							]
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.#":               1,
				"applications.oidc_applications.0.name":          "fed-app",
				"applications.oidc_applications.0.permissions.#": 2,
				"applications.oidc_applications.0.roles.#":       2,
			}),
		},
		// Edit roles + permissions; remove writer role and write:data permission
		resource.TestStep{
			Config: p.Config(`
				applications = {
					oidc_applications = [
						{
							name = "fed-app"

							permissions = [
								{ name = "read:data" },
							]
							roles = [
								{
									name = "reader"
									permissions = ["read:data"]
								},
							]
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.0.permissions.#": 1,
				"applications.oidc_applications.0.roles.#":       1,
			}),
		},
		// Clear roles and permissions
		resource.TestStep{
			Config: p.Config(`
				applications = {
					oidc_applications = [
						{
							name = "fed-app"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"applications.oidc_applications.0.permissions.#": 0,
				"applications.oidc_applications.0.roles.#":       0,
			}),
		},
	)
}
