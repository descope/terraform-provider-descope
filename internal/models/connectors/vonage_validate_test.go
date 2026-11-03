package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVonageConnectorSMSOnly(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "vonage_connector")

	testacc.RunWithDestroyCheck(t, "descope_vonage_connector", resource.TestStep{
		Config: c.Config(`project_id = "`+projectID+`"`,
			`api_key = "sms-key"`,
			`api_secret = "sms-secret"`,
			`sender = "Descope"`,
		),
		Check: c.Check(map[string]any{
			"api_key":          "sms-key",
			"api_secret":       "sms-secret",
			"sender":           "Descope",
			"application_id":   "",
			"from_phone_voice": "",
		}),
	})
}

func TestVonageConnectorVoiceOnly(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "vonage_connector")

	testacc.RunWithDestroyCheck(t, "descope_vonage_connector", resource.TestStep{
		Config: c.Config(`project_id = "`+projectID+`"`,
			`application_id = "voice-app"`,
			`private_key = "voice-key"`,
			`from_phone_voice = "+15555550100"`,
		),
		Check: c.Check(map[string]any{
			"api_key":          "",
			"sender":           "",
			"application_id":   "voice-app",
			"private_key":      "voice-key",
			"from_phone_voice": "+15555550100",
		}),
	})
}

func TestVonageConnectorNoCompleteChannel(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "vonage_connector")

	testacc.Run(t, resource.TestStep{
		Config: c.Config(`project_id = "`+projectID+`"`,
			`api_key = "sms-key"`,
			`api_secret = "sms-secret"`,
			`application_id = "voice-app"`,
			`private_key = "voice-key"`,
		),
		ExpectError: regexp.MustCompile(`requires\s+at\s+least\s+one\s+complete\s+channel`),
	})
}
