package apps_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppImportTypeMismatch(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := testacc.OIDCApp(t)
	s := testacc.SAMLApp(t)
	project := `project_id = "` + projectID + `"`
	config := o.Config(project, `deletion_protection = false`) + s.Config(project,
		`deletion_protection = false`,
		`manual_configuration = {
			acs_url = "https://sp.example.com/acs"
			entity_id = "sp-entity"
		}`,
	)

	testacc.Run(t,
		resource.TestStep{
			Config: config,
		},
		resource.TestStep{
			Config:            config,
			ResourceName:      s.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id", "id"),
			ExpectError:       regexp.MustCompile(`is of type "oidc",\s+not\s+"saml"`),
		},
	)
}
