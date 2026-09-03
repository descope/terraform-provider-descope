package fgaschema_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFGASchema(t *testing.T) {
	projectID := testacc.ProjectID(t)
	s := testacc.FGASchema(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config: s.Block(project, `schema = "model AuthZ 1.0\n\ntype user\n"`),
			Check: s.Check(map[string]any{
				"id.==":      "project_id",
				"project_id": projectID,
				"schema":     testacc.AttributeIsSet,
			}),
		},
		resource.TestStep{
			Config: s.Block(project, `schema = "model AuthZ 1.0\n\ntype user\n\ntype doc\n  relation owner: user\n  permission view: owner\n"`),
			Check: s.Check(map[string]any{
				"id.==":  "project_id",
				"schema": testacc.AttributeIsSet,
			}),
		},
		resource.TestStep{
			Config: s.Block(project, `schema = ""`),
			Check: s.Check(map[string]any{
				"id.==":  "project_id",
				"schema": "",
			}),
		},
		resource.TestStep{
			ResourceName:      s.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateId:     projectID,
		},
	)
}

func TestFGASchemaInvalidPrefix(t *testing.T) {
	projectID := testacc.ProjectID(t)
	s := testacc.FGASchema(t)

	testacc.Run(t,
		resource.TestStep{
			Config:      s.Block(`project_id = "`+projectID+`"`, `schema = "type user"`),
			ExpectError: regexp.MustCompile(`must start with 'model AuthZ'`),
		},
	)
}
