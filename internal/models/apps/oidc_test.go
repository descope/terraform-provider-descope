package apps_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/apps"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestOIDCAppClientSecretPlan(t *testing.T) {
	testCases := []struct {
		name       string
		clientType types.String
		config     types.String
		state      types.String
		plan       types.String
		expected   types.String
	}{
		{"config secret is kept", types.StringValue("confidential"), types.StringValue("imported"), types.StringNull(), types.StringValue("imported"), types.StringValue("imported")},
		{"existing secret is kept", types.StringValue("confidential"), types.StringNull(), types.StringValue("s3cr3t"), types.StringValue("s3cr3t"), types.StringValue("s3cr3t")},
		{"confidential create is generated", types.StringValue("confidential"), types.StringNull(), types.StringNull(), types.StringUnknown(), types.StringUnknown()},
		{"switch to confidential expects a new secret", types.StringValue("confidential"), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringUnknown()},
		{"unknown type expects a new secret", types.StringUnknown(), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringUnknown()},
		{"legacy type has no secret", types.StringValue(""), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringValue("")},
		{"public type has no secret", types.StringValue("public"), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringValue("")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plan := &apps.OIDCAppModel{ClientType: tc.clientType, ClientSecret: tc.plan}
			plan.ModifyPlan(nil, &apps.OIDCAppModel{ClientSecret: tc.config}, &apps.OIDCAppModel{ClientSecret: tc.state})
			if !plan.ClientSecret.Equal(tc.expected) {
				t.Errorf("expected planned client_secret to be %s, got %s", tc.expected, plan.ClientSecret)
			}
		})
	}
}

func TestOIDCAppClientIDPlan(t *testing.T) {
	testCases := []struct {
		name       string
		clientType types.String
		config     types.String
		state      types.String
		plan       types.String
		expected   types.String
	}{
		{"config id is kept", types.StringValue("confidential"), types.StringValue("pinned"), types.StringNull(), types.StringValue("pinned"), types.StringValue("pinned")},
		{"existing id is kept", types.StringValue("confidential"), types.StringNull(), types.StringValue("client-x"), types.StringValue("client-x"), types.StringValue("client-x")},
		{"switch to confidential expects a new id", types.StringValue("confidential"), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringUnknown()},
		{"switch to public expects a new id", types.StringValue("public"), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringUnknown()},
		{"unknown type expects a new id", types.StringUnknown(), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringUnknown()},
		{"legacy type has no id", types.StringValue(""), types.StringNull(), types.StringValue(""), types.StringValue(""), types.StringValue("")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plan := &apps.OIDCAppModel{ClientType: tc.clientType, ClientID: tc.plan}
			plan.ModifyPlan(nil, &apps.OIDCAppModel{ClientID: tc.config}, &apps.OIDCAppModel{ClientID: tc.state})
			if !plan.ClientID.Equal(tc.expected) {
				t.Errorf("expected planned client_id to be %s, got %s", tc.expected, plan.ClientID)
			}
		})
	}
}

func TestOIDCApp(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	project := `project_id = "` + projectID + `"`

	var appID, clientSecret string
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
	// raw check since attribute checks treat empty as not set: older backends generate a secret on any create, newer only for confidential apps
	captureSecret := func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[a.Path()]
		if !ok {
			return fmt.Errorf("resource %s not found in state", a.Path())
		}
		clientSecret = res.Primary.Attributes["client_secret"]
		return nil
	}
	confidentialSecret := func(s string) error {
		if s == "" {
			return fmt.Errorf("expected a client secret after the update to a confidential client type")
		}
		if clientSecret != "" && s != clientSecret {
			return fmt.Errorf("expected the client secret to be preserved in the state")
		}
		clientSecret = s
		return nil
	}

	createStep := resource.TestStep{
		Config: a.Config(project, `deletion_protection = false`),
		Check: a.Check(map[string]any{
			"id":                                  captureID,
			"project_id":                          projectID,
			"name":                                a.Name,
			"description":                         "",
			"logo":                                "",
			"disabled":                            false,
			"login_page_url":                      testacc.AttributeIsSet,
			"claims":                              []string{},
			"force_authentication":                false,
			"backchannel_logout_url":              "",
			"custom_idp_initiated_login_page_url": "",
			"client_id":                           "",
			"client_type":                         "",
			"approved_redirect_urls":              []string{},
			"force_pkce":                          false,
			"default_audience":                    "",
			"trusted_apps_audience":               "",
		}, captureSecret),
	}

	a.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`description = "updated"`,
			`claims = ["sub", "email"]`,
			`client_type = "confidential"`,
			`approved_redirect_urls = ["https://sp.example.com/callback"]`,
			`default_audience = "empty"`,
			`trusted_apps_audience = "appId"`,
			`backchannel_logout_url = "https://sp.example.com/logout"`,
			`custom_idp_initiated_login_page_url = "https://sp.example.com/login"`,
			`force_pkce = true`,
			`refresh_token_disabled = true`,
		),
		Check: a.Check(map[string]any{
			"id":                                  sameID,
			"name":                                a.Name,
			"description":                         "updated",
			"claims":                              []string{"sub", "email"},
			"client_id":                           testacc.AttributeIsSet,
			"client_secret":                       confidentialSecret,
			"client_type":                         "confidential",
			"approved_redirect_urls":              []string{"https://sp.example.com/callback"},
			"default_audience":                    "empty",
			"trusted_apps_audience":               "appId",
			"backchannel_logout_url":              "https://sp.example.com/logout",
			"custom_idp_initiated_login_page_url": "https://sp.example.com/login",
			"force_pkce":                          true,
			"refresh_token_disabled":              true,
		}),
	}

	importStep := resource.TestStep{
		ResourceName:            a.Path(),
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateIdFunc:       testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		ImportStateVerifyIgnore: []string{"deletion_protection"},
	}

	testacc.Run(t, createStep, updateStep, importStep)
}

func TestOIDCAppImportedClient(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	project := `project_id = "` + projectID + `"`
	appID := "tf-oidc-" + a.Name[len(a.Name)-8:] // the app id attribute is limited to 30 characters

	createStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`id = "`+appID+`"`,
			`client_id = "client-`+appID+`"`,
			`client_secret = "secret-value-for-`+appID+`"`,
			`client_type = "confidential"`,
		),
		Check: a.Check(map[string]any{
			"id":            appID,
			"client_id":     "client-" + appID,
			"client_secret": "secret-value-for-" + appID,
		}),
	}

	testacc.Run(t, createStep)
}

func TestOIDCAppDeletionProtection(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config: a.Config(project),
			Check: a.Check(map[string]any{
				"id":                  testacc.AttributeIsSet,
				"deletion_protection": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config:      a.Config(project),
			Destroy:     true,
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		resource.TestStep{
			Config: a.Config(project, `deletion_protection = false`),
		},
	)
}
