package authentication_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthentication(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"authentication": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {}
			`),
			Check: p.Check(map[string]any{
				"authentication.%": 9,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						redirect_url = "1"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`The redirectUrl field must be a valid URL`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						disabled = true
						redirect_url = "https://example.com"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.magic_link": map[string]any{
					"disabled":        true,
					"redirect_url":    "https://example.com",
					"expiration_time": "3 minutes",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "2000 seconds"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`space and one of the valid time units`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "1 second"
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					magic_link = {
						expiration_time = "5 minutes"
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.magic_link": map[string]any{
					"disabled":        false,
					"redirect_url":    "https://example.com",
					"expiration_time": "5 minutes",
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					oauth = {
						custom = {
							apple = {
							}
						}
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Reserved OAuth Provider Name`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					oauth = {
						system = {
							apple = {
								allowed_grant_types = ["authorization_code", "implicit"]
							}
						}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.oauth.system.apple.allowed_grant_types": []string{"authorization_code", "implicit"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					oauth = {
						system = {
							apple = {
								client_id = "id"
							}
						}
					}
				}
			`),
			ExpectError: regexp.MustCompile(`Missing Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					oauth = {
						custom = {
							mobile_ios = {
								allowed_grant_types = ["authorization_code", "implicit"]
								client_id = "id"
								client_secret = "secret"
								authorization_endpoint = "https://auth.com"
								token_endpoint = "https://token.com"
								user_info_endpoint = "https://user.com"
								use_client_assertion = true
							}
						}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.oauth.custom.%": 1,
				"authentication.oauth.custom.mobile_ios": map[string]any{
					"allowed_grant_types":    []string{"authorization_code", "implicit"},
					"client_id":              "id",
					"client_secret":          testacc.AttributeIsSet,
					"authorization_endpoint": "https://auth.com",
					"token_endpoint":         "https://token.com",
					"user_info_endpoint":     "https://user.com",
					"use_client_assertion":   true,
				},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					sso = {
						sso_suite_settings = {
							hide_saml = true
							hide_oidc = true
						}
					}
				}
			`),
			ExpectError: regexp.MustCompile("The attributes hide_oidc and hide_saml cannot both be true"),
		},
		resource.TestStep{
			Config: p.Config(`
				authentication = {
					sso = {
						sso_suite_settings = {
							style_id = "koko"
							hide_saml = true
						}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.sso.sso_suite_settings": map[string]any{
					"style_id":  "koko",
					"hide_saml": true,
				},
			}),
		},
	)
}
