package role_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRole(t *testing.T) {
	projectID := testacc.ProjectID(t)
	p := testacc.Permission(t)
	r := testacc.Role(t)
	project := `project_id = "` + projectID + `"`

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
		Config: p.Config(project) + r.Config(project, `permissions = [`+p.Path()+`.name]`),
		Check: r.Check(map[string]any{
			"id":          captureID,
			"project_id":  projectID,
			"name":        r.Name,
			"description": "",
			"permissions": []string{p.Name},
			"default":     false,
			"private":     false,
		}),
	}

	r.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: p.Config(project) + r.Config(project,
			`description = "updated"`,
			`permissions = []`,
			`default = true`,
			`private = true`,
		),
		Check: r.Check(map[string]any{
			"id":          sameID,
			"name":        r.Name,
			"description": "updated",
			"permissions": []string{},
			"default":     true,
			"private":     true,
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      r.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(r.Path(), "project_id", "id"),
	}

	testacc.RunWithDestroyCheck(t, "descope_role", createStep, updateStep, importStep)
}

func TestRoleOutOfBandDeletion(t *testing.T) {
	projectID := testacc.ProjectID(t)
	r := testacc.Role(t)
	config := r.Config(`project_id = "` + projectID + `"`)

	var roleID string
	testacc.Run(t,
		resource.TestStep{
			Config: config,
			Check: r.Check(map[string]any{
				"name": r.Name,
				"id": func(s string) error {
					roleID = s
					return nil
				},
			}),
		},
		resource.TestStep{
			PreConfig: func() {
				testacc.OutOfBandPost(t, projectID, "/v1/mgmt/role/delete", map[string]any{"id": roleID})
			},
			Config: config,
			Check: r.Check(map[string]any{
				"name": r.Name,
				"id": func(s string) error {
					if s == roleID {
						return fmt.Errorf("expected the role to be re-created with a new id, got the same id %s", s)
					}
					return nil
				},
			}),
		},
		resource.TestStep{
			Config:        config,
			ResourceName:  r.Path(),
			ImportState:   true,
			ImportStateId: projectID + "/Rbogusbogusbogusbogusbogusb",
			ExpectError:   regexp.MustCompile(`Error reading role`),
		},
	)
}
