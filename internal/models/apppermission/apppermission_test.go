package apppermission_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppPermission(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	p := testacc.AppPermission(t)
	project := `project_id = "` + projectID + `"`
	app := a.Config(project, `deletion_protection = false`)
	appID := `app_id = ` + a.Path() + `.id`

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
		Config: app + p.Config(project, appID),
		Check: p.Check(map[string]any{
			"id":          captureID,
			"project_id":  projectID,
			"name":        p.Name,
			"description": "",
		}),
	}

	p.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: app + p.Config(project, appID, `description = "updated"`),
		Check: p.Check(map[string]any{
			"id":          sameID,
			"name":        p.Name,
			"description": "updated",
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      p.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(p.Path(), "project_id", "app_id", "id"),
	}

	immutableProjectStep := resource.TestStep{
		Config:      app + p.Config(`project_id = "P2aaaaaaaaaaaaaaaaaaaaaaaaaa"`, appID),
		ExpectError: regexp.MustCompile(`Immutable Attribute Changed`),
	}
	immutableAppStep := resource.TestStep{
		Config:      app + p.Config(project, `app_id = "SA2aaaaaaaaaaaaaaaaaaaaaaaaa"`),
		ExpectError: regexp.MustCompile(`Immutable Attribute Changed`),
	}

	testacc.Run(t, createStep, updateStep, importStep, immutableProjectStep, immutableAppStep)
}
