package outboundscim_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOutboundSCIMConfiguration(t *testing.T) {
	projectID := os.Getenv("DESCOPE_OUTBOUND_SCIM_TEST_PROJECT_ID")
	appID := os.Getenv("DESCOPE_OUTBOUND_SCIM_TEST_APP_ID")
	configuration := os.Getenv("DESCOPE_OUTBOUND_SCIM_TEST_CONFIGURATION")
	if projectID == "" || appID == "" || configuration == "" {
		t.Skip("set DESCOPE_OUTBOUND_SCIM_TEST_PROJECT_ID, DESCOPE_OUTBOUND_SCIM_TEST_APP_ID, and DESCOPE_OUTBOUND_SCIM_TEST_CONFIGURATION to run")
	}
	secrets := os.Getenv("DESCOPE_OUTBOUND_SCIM_TEST_SECRETS")

	config := renderOutboundSCIMConfig(outboundSCIMTestConfig{projectID: projectID, appID: appID, configuration: configuration, secrets: secrets, enabled: true})
	updated := renderOutboundSCIMConfig(outboundSCIMTestConfig{projectID: projectID, appID: appID, configuration: configuration, secrets: secrets, enabled: false})
	resourceName := "descope_outbound_scim_configuration.test"

	testacc.RunIsolated(
		t,
		resource.TestStep{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "id", appID),
				resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				resource.TestCheckResourceAttr(resourceName, "app_id", appID),
				resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
			),
		},
		resource.TestStep{
			Config: updated,
			Check:  resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		},
		resource.TestStep{
			ResourceName:      resourceName,
			ImportState:       true,
			ImportStateId:     projectID + "/" + appID,
			ImportStateVerify: true,
		},
	)
}

type outboundSCIMTestConfig struct {
	projectID     string
	appID         string
	configuration string
	secrets       string
	enabled       bool
}

func renderOutboundSCIMConfig(config outboundSCIMTestConfig) string {
	secretAttribute := ""
	if config.secrets != "" {
		secretAttribute = fmt.Sprintf("secrets = %q", config.secrets)
	}
	return fmt.Sprintf(`
resource "descope_outbound_scim_configuration" "test" {
  project_id    = %q
  app_id        = %q
  configuration = %q
  enabled       = %t
  %s
}
`, config.projectID, config.appID, config.configuration, config.enabled, secretAttribute)
}
