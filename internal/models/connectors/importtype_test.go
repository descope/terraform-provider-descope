package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectorImportTypeMismatch(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.NewResource(t, "abuseipdb_connector")
	b := testacc.NewResource(t, "alloy_connector")
	project := `project_id = "` + projectID + `"`
	config := a.Config(project, `api_key = "mhvece"`) + b.Config(project,
		`api_token = "mybopddv"`,
		`api_secret = "hgg666mus"`,
		`base_url = "https://api.alloy.co/v1"`,
	)

	testacc.Run(t,
		resource.TestStep{
			Config: config,
		},
		resource.TestStep{
			Config:            config,
			ResourceName:      b.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
			ExpectError:       regexp.MustCompile(`is of type "abuseipdb",\s+not\s+"alloy"`),
		},
	)
}
