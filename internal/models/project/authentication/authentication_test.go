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
				"authentication.magic_link.disabled":        true,
				"authentication.magic_link.redirect_url":    "https://example.com",
				"authentication.magic_link.expiration_time": "3 minutes",
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
				"authentication.magic_link.disabled":        false,
				"authentication.magic_link.redirect_url":    "https://example.com",
				"authentication.magic_link.expiration_time": "5 minutes",
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
							}
						}
					}
				}
			`),
			Check: p.Check(map[string]any{
				"authentication.oauth.custom.%":                                 1,
				"authentication.oauth.custom.mobile_ios.allowed_grant_types":    []string{"authorization_code", "implicit"},
				"authentication.oauth.custom.mobile_ios.client_id":              "id",
				"authentication.oauth.custom.mobile_ios.client_secret":          testacc.AttributeIsSet,
				"authentication.oauth.custom.mobile_ios.authorization_endpoint": "https://auth.com",
				"authentication.oauth.custom.mobile_ios.token_endpoint":         "https://token.com",
				"authentication.oauth.custom.mobile_ios.user_info_endpoint":     "https://user.com",
			}),
		},
	)
}
