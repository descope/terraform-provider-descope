package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAWSEndUserMessagingConnectorAuth(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "aws_end_user_messaging_connector")
	project := `project_id = "` + projectID + `"`

	testacc.RunWithDestroyCheck(t, "descope_aws_end_user_messaging_connector",
		// credentials auth without a secret access key is rejected
		resource.TestStep{
			Config: c.Config(project,
				`aws_region = "us-east-1"`,
				`access_key_id = "AKIAEXAMPLE"`,
				`notify_configuration_id = "notify-1"`,
			),
			ExpectError: regexp.MustCompile(`The secret_access_key field is required`),
		},
		// minimal connector with credentials auth and every optional field at its default
		resource.TestStep{
			Config: c.Config(project,
				`aws_region = "us-east-1"`,
				`access_key_id = "AKIAEXAMPLE"`,
				`secret_access_key = "example-secret"`,
				`notify_configuration_id = "notify-1"`,
			),
			Check: c.Check(map[string]any{
				"id":                      testacc.AttributeIsSet,
				"auth_type":               "credentials",
				"aws_region":              "us-east-1",
				"notify_configuration_id": "notify-1",
				"default_template_id":     "",
				"template_expiry_minutes": 0,
			}),
		},
		// switching to assume role drops the access keys
		resource.TestStep{
			Config: c.Config(project,
				`aws_region = "us-west-2"`,
				`auth_type = "assumeRole"`,
				`role_arn = "arn:aws:iam::123456789012:role/descope-eum"`,
				`external_id = "descope-eum-external-id"`,
				`notify_configuration_id = "notify-1"`,
			),
			Check: c.Check(map[string]any{
				"auth_type":         "assumeRole",
				"aws_region":        "us-west-2",
				"role_arn":          "arn:aws:iam::123456789012:role/descope-eum",
				"external_id":       "descope-eum-external-id",
				"access_key_id":     testacc.AttributeIsNotSet,
				"secret_access_key": testacc.AttributeIsNotSet,
			}),
		},
	)
}
