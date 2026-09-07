package permission_test

import (
	"fmt"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPermission(t *testing.T) {
	projectID := testacc.ProjectID(t)
	p := testacc.Permission(t)
	project := `project_id = "` + projectID + `"`

	var permissionID string
	captureID := func(s string) error {
		permissionID = s
		return nil
	}
	sameID := func(s string) error {
		if s != permissionID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", permissionID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: p.Config(project),
		Check: p.Check(map[string]any{
			"id":          captureID,
			"project_id":  projectID,
			"name":        p.Name,
			"description": "",
		}),
	}

	updateStep := resource.TestStep{
		Config: p.Config(project, `description = "updated"`),
		Check: p.Check(map[string]any{
			"id":          sameID,
			"description": "updated",
		}),
	}

	p.Name += "-renamed"
	renameStep := resource.TestStep{
		Config: p.Config(project, `description = "updated"`),
		Check: p.Check(map[string]any{
			"id":   sameID,
			"name": p.Name,
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      p.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(p.Path(), "project_id", "id"),
	}

	testacc.Run(t, createStep, updateStep, renameStep, importStep)
}
