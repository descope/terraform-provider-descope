package approle_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAppRole(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	p := testacc.AppPermission(t)
	pr := testacc.Role(t)
	r := testacc.AppRole(t)
	project := `project_id = "` + projectID + `"`
	app := a.Config(project, `deletion_protection = false`)
	appID := `app_id = ` + a.Path() + `.id`

	var roleID string
	captureID := func(s string) error {
		roleID = s
		return nil
	}
	sameID := func(s string) error {
		if s != roleID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", roleID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: app + p.Config(project, appID) + pr.Config(project) + r.Config(project, appID,
			`permission_ids = [`+p.Path()+`.id]`,
			`role_mappings = [`+pr.Path()+`.id]`,
		),
		Check: r.Check(map[string]any{
			"id":               captureID,
			"project_id":       projectID,
			"name":             r.Name,
			"description":      "",
			"permission_ids.#": "1",
			"role_mappings.#":  "1",
		}),
	}

	r.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: app + p.Config(project, appID) + pr.Config(project) + r.Config(project, appID,
			`description = "updated"`,
			`permission_ids = []`,
			`role_mappings = []`,
		),
		Check: r.Check(map[string]any{
			"id":             sameID,
			"name":           r.Name,
			"description":    "updated",
			"permission_ids": []string{},
			"role_mappings":  []string{},
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      r.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(r.Path(), "project_id", "app_id", "id"),
	}

	testacc.Run(t, createStep, updateStep, importStep)
}

func TestAppRoleInvalidPermission(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.OIDCApp(t)
	r := testacc.AppRole(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t, resource.TestStep{
		Config: a.Config(project, `deletion_protection = false`) + r.Config(project,
			`app_id = `+a.Path()+`.id`,
			`permission_ids = ["SAP0000000000000000000000000"]`,
		),
		ExpectError: regexp.MustCompile(`references unknown permissions`),
	})
}
