package list_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestList(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	var listID string
	captureID := func(s string) error {
		listID = s
		return nil
	}
	sameID := func(s string) error {
		if s != listID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", listID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: l.Config(project, `description = "a texts list"`, `texts = ["alpha", "beta"]`),
		Check: l.Check(map[string]any{
			"id":          captureID,
			"project_id":  projectID,
			"name":        l.Name,
			"description": "a texts list",
			"texts":       []string{"alpha", "beta"},
			"ips":         []string{},
			"json":        "",
		}),
	}

	updateStep := resource.TestStep{
		Config: l.Config(project, `description = "updated"`, `texts = ["gamma"]`),
		Check: l.Check(map[string]any{
			"id":          sameID,
			"description": "updated",
			"texts":       []string{"gamma"},
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      l.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(l.Path(), "project_id", "id"),
	}

	testacc.RunWithDestroyCheck(t, "descope_list", createStep, updateStep, importStep)
}

func TestListJSON(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t, resource.TestStep{
		Config: l.Config(project, `json = jsonencode({ foo = "bar", enabled = true })`),
		Check: l.Check(map[string]any{
			"name":  l.Name,
			"json":  `{"enabled":true,"foo":"bar"}`,
			"texts": []string{},
			"ips":   []string{},
		}),
	})
}

func TestListIPs(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t, resource.TestStep{
		Config: l.Config(project, `ips = ["10.0.0.1", "192.168.0.0/24"]`),
		Check: l.Check(map[string]any{
			"name": l.Name,
			"ips":  []string{"10.0.0.1", "192.168.0.0/24"},
		}),
	})
}

func TestListSwitchKind(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config: l.Config(project, `texts = ["alpha"]`),
			Check:  l.Check(map[string]any{"texts": []string{"alpha"}, "json": ""}),
		},
		resource.TestStep{
			Config: l.Config(project, `json = jsonencode({ foo = "bar" })`),
			Check:  l.Check(map[string]any{"json": `{"foo":"bar"}`, "texts": []string{}}),
		},
	)
}

func TestListDuplicatesCollapse(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t, resource.TestStep{
		Config: l.Config(project, `texts = ["dup", "dup", "uniq"]`),
		Check:  l.Check(map[string]any{"texts": []string{"dup", "uniq"}}),
	})
}

func TestListValidation(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config:      l.Config(project),
			ExpectError: regexp.MustCompile(`Exactly one of the 'texts', 'ips' or 'json' attributes`),
		},
		resource.TestStep{
			Config:      l.Config(project, `texts = ["alpha"]`, `json = jsonencode({})`),
			ExpectError: regexp.MustCompile(`Exactly one of the 'texts', 'ips' or 'json' attributes`),
		},
		resource.TestStep{
			Config:      l.Config(project, `ips = ["not-an-ip"]`),
			ExpectError: regexp.MustCompile(`must be an IP address or CIDR range`),
		},
		resource.TestStep{
			Config:      l.Config(project, `json = jsonencode(["array-not-object"])`),
			ExpectError: regexp.MustCompile(`must be valid JSON`),
		},
		resource.TestStep{
			Config:      l.Config(project, `json = jsonencode(null)`),
			ExpectError: regexp.MustCompile(`must be valid JSON`),
		},
	)
}

// A document written with indentation must survive an apply, see TestJWTTemplateIndentedConfig.
func TestListIndentedConfig(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	project := `project_id = "` + projectID + `"`
	indented := "{\n  \"foo\": \"bar\"\n}"

	testacc.Run(t,
		resource.TestStep{
			Config: l.Config(project, `json = "{\n  \"foo\": \"bar\"\n}"`),
			Check:  l.Check(map[string]any{"json": indented}),
		},
		resource.TestStep{
			Config: l.Config(project, `json = jsonencode({ foo = "bar" })`),
			Check:  l.Check(map[string]any{"json": indented}),
		},
	)
}

func TestListOutOfBandDeletion(t *testing.T) {
	projectID := testacc.ProjectID(t)
	l := testacc.List(t)
	config := l.Config(`project_id = "`+projectID+`"`, `texts = ["alpha"]`)

	var listID string
	testacc.Run(t,
		resource.TestStep{
			Config: config,
			Check: l.Check(map[string]any{
				"name": l.Name,
				"id": func(s string) error {
					listID = s
					return nil
				},
			}),
		},
		resource.TestStep{
			PreConfig: func() {
				testacc.OutOfBandPost(t, projectID, "/v1/mgmt/list/delete", map[string]any{"id": listID})
			},
			Config: config,
			Check: l.Check(map[string]any{
				"name": l.Name,
				"id": func(s string) error {
					if s == listID {
						return fmt.Errorf("expected the list to be re-created with a new id, got the same id %s", s)
					}
					return nil
				},
			}),
		},
		resource.TestStep{
			Config:        config,
			ResourceName:  l.Path(),
			ImportState:   true,
			ImportStateId: projectID + "/does-not-exist-xyz",
			ExpectError:   regexp.MustCompile(`Error reading list`),
		},
	)
}
