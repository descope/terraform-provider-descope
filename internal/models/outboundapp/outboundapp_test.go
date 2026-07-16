package outboundapp_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_OutboundApp_acceptance(t *testing.T) {
	projectID := os.Getenv("DESCOPE_OUTBOUND_APP_TEST_PROJECT_ID")
	clientID := os.Getenv("DESCOPE_OUTBOUND_APP_TEST_CLIENT_ID")
	clientSecret := os.Getenv("DESCOPE_OUTBOUND_APP_TEST_CLIENT_SECRET")
	authorizationURL := os.Getenv("DESCOPE_OUTBOUND_APP_TEST_AUTHORIZATION_URL")
	tokenURL := os.Getenv("DESCOPE_OUTBOUND_APP_TEST_TOKEN_URL")
	if projectID == "" || clientID == "" || clientSecret == "" || authorizationURL == "" || tokenURL == "" {
		t.Skip("set the DESCOPE_OUTBOUND_APP_TEST_* environment variables to run")
	}
	app := &testacc.Resource{Type: "outbound_app", ID: "test", Name: testacc.GenerateAlias(t)}
	config := outboundAppConfig(outboundAppTestConfig{
		app: app, projectID: projectID, clientID: clientID, clientSecret: clientSecret,
		authorizationURL: authorizationURL, tokenURL: tokenURL, description: "created by Terraform",
	})
	updated := outboundAppConfig(outboundAppTestConfig{
		app: app, projectID: projectID, clientID: clientID, clientSecret: clientSecret,
		authorizationURL: authorizationURL, tokenURL: tokenURL, description: "updated by Terraform",
	})

	testacc.RunIsolated(
		t,
		resource.TestStep{
			Config: config,
			Check: app.Check(map[string]any{
				"id":          testacc.AttributeIsSet,
				"project_id":  projectID,
				"name":        app.Name,
				"description": "created by Terraform",
			}),
		},
		resource.TestStep{
			Config: updated,
			Check:  app.Check(map[string]any{"description": "updated by Terraform"}),
		},
		resource.TestStep{
			ResourceName:      app.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(app.Path(), "project_id", "id"),
			ImportStateVerify: true,
		},
	)
}

type outboundAppTestConfig struct {
	app              *testacc.Resource
	projectID        string
	clientID         string
	clientSecret     string
	authorizationURL string
	tokenURL         string
	description      string
}

func outboundAppConfig(config outboundAppTestConfig) string {
	return config.app.Config(
		fmt.Sprintf("project_id = %q", config.projectID),
		fmt.Sprintf("description = %q", config.description),
		fmt.Sprintf("client_id = %q", config.clientID),
		fmt.Sprintf("client_secret = %q", config.clientSecret),
		fmt.Sprintf("authorization_url = %q", config.authorizationURL),
		fmt.Sprintf("token_url = %q", config.tokenURL),
	)
}
