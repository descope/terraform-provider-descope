package emit

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func write(t *testing.T, plan *Plan) map[string]string {
	t.Helper()
	outDir := t.TempDir()
	if err := Write(context.Background(), outDir, plan); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		content, err := os.ReadFile(filepath.Join(outDir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files[entry.Name()] = string(content)
	}
	return files
}

func TestProjectIDReferencesTheProjectResource(t *testing.T) {
	plan := &Plan{
		ProjectID: "P123",
		Resources: []Resource{
			{
				Type:     "descope_project",
				Label:    "my_auth_project",
				EntityID: "P123",
				ImportID: "P123",
				Attrs:    []prune.Attr{{Name: "name", Value: types.StringValue("My Auth Project")}},
			},
			{
				Type:         "descope_role",
				Label:        "admin",
				EntityID:     "ROL456",
				ImportID:     "P123/ROL456",
				HasProjectID: true,
				Attrs:        []prune.Attr{{Name: "name", Value: types.StringValue("admin")}},
			},
		},
	}

	files := write(t, plan)

	if _, ok := files["project.tf"]; !ok {
		t.Errorf("expected the project in its own file, got %v", keys(files))
	}
	if want := "project_id = descope_project.my_auth_project.id"; !strings.Contains(files["authorization.tf"], want) {
		t.Errorf("expected %q in authorization.tf:\n%s", want, files["authorization.tf"])
	}
	if _, ok := files["variables.tf"]; ok {
		t.Errorf("expected no variables file at all, got:\n%s", files["variables.tf"])
	}
	if want := `id = "P123"`; !strings.Contains(files["import.tf"], want) {
		t.Errorf("expected the project import %q:\n%s", want, files["import.tf"])
	}
}

func TestProjectIDFallsBackToAVariable(t *testing.T) {
	plan := &Plan{
		ProjectID: "P123",
		Resources: []Resource{
			{
				Type:         "descope_role",
				Label:        "admin",
				EntityID:     "ROL456",
				ImportID:     "P123/ROL456",
				HasProjectID: true,
				Attrs:        []prune.Attr{{Name: "name", Value: types.StringValue("admin")}},
			},
		},
	}

	files := write(t, plan)

	if want := "project_id = var.project_id"; !strings.Contains(files["authorization.tf"], want) {
		t.Errorf("expected %q without an exported project:\n%s", want, files["authorization.tf"])
	}
	if want := `variable "project_id"`; !strings.Contains(files["variables.tf"], want) {
		t.Errorf("expected %q to be declared:\n%s", want, files["variables.tf"])
	}
}

func keys(files map[string]string) []string {
	var names []string
	for name := range files {
		names = append(names, name)
	}
	return names
}

func elementsResource(elements [][]prune.Attr) *Plan {
	return &Plan{
		ProjectID: "P123",
		Resources: []Resource{{
			Type: "descope_http_connector", Label: "hook", EntityID: "CON1", ImportID: "P123/CON1", HasProjectID: true,
			Attrs: []prune.Attr{{Name: "headers", IsElements: true, Elements: elements}},
		}},
	}
}

func TestElementVariablesAreNamedByIndexOnlyWhenAmbiguous(t *testing.T) {
	t.Run("one element needs no index", func(t *testing.T) {
		files := write(t, elementsResource([][]prune.Attr{
			{{Name: "value", Placeholder: true}},
		}))
		if want := `variable "hook_headers_value"`; !strings.Contains(files["variables.tf"], want) {
			t.Errorf("expected %q:\n%s", want, files["variables.tf"])
		}
	})

	t.Run("several elements get one variable each", func(t *testing.T) {
		files := write(t, elementsResource([][]prune.Attr{
			{{Name: "value", Placeholder: true}},
			{{Name: "value", Placeholder: true}},
		}))
		for _, want := range []string{`variable "hook_headers_0_value"`, `variable "hook_headers_1_value"`} {
			if !strings.Contains(files["variables.tf"], want) {
				t.Errorf("expected %q:\n%s", want, files["variables.tf"])
			}
		}
		if count := strings.Count(files["variables.tf"], `variable "hook_headers`); count != 2 {
			t.Errorf("expected exactly two header variables, got %d:\n%s", count, files["variables.tf"])
		}
	})
}

func TestSecretVariablesAreScopedOnlyWhenTheyCollide(t *testing.T) {
	plan := &Plan{
		ProjectID: "P123",
		Resources: []Resource{
			{
				Type:         "descope_outbound_app",
				Label:        "crm",
				EntityID:     "OA1",
				ImportID:     "P123/OA1",
				HasProjectID: true,
				Attrs:        []prune.Attr{{Name: "client_secret", Placeholder: true}},
			},
			{
				Type:         "descope_salesforce_connector",
				Label:        "crm",
				EntityID:     "CON1",
				ImportID:     "P123/CON1",
				HasProjectID: true,
				Attrs:        []prune.Attr{{Name: "client_secret", Placeholder: true}, {Name: "api_key", Placeholder: true}},
			},
		},
	}

	files := write(t, plan)

	for _, want := range []string{"outbound_app_crm_client_secret", "salesforce_connector_crm_client_secret", "crm_api_key"} {
		if !strings.Contains(files["variables.tf"], `variable "`+want+`"`) {
			t.Errorf("expected a %q variable:\n%s", want, files["variables.tf"])
		}
	}
	if strings.Contains(files["variables.tf"], `variable "crm_client_secret"`) {
		t.Errorf("expected the colliding name to be gone:\n%s", files["variables.tf"])
	}
	if want := "client_secret = var.outbound_app_crm_client_secret"; !strings.Contains(files["apps.tf"], want) {
		t.Errorf("expected %q in apps.tf:\n%s", want, files["apps.tf"])
	}
}

func TestProjectIDReferencesAnExistingProject(t *testing.T) {
	address, err := ParseProjectAddress(`descope_project.main["prod"]`)
	if err != nil {
		t.Fatal(err)
	}
	prefix, err := ParseModulePrefix("module.auth")
	if err != nil {
		t.Fatal(err)
	}
	plan := &Plan{
		ProjectID:      "P123",
		ProjectAddress: address,
		ImportPrefix:   prefix,
		Resources: []Resource{
			{
				Type:         "descope_role",
				Label:        "admin",
				EntityID:     "ROL456",
				ImportID:     "P123/ROL456",
				HasProjectID: true,
				Attrs:        []prune.Attr{{Name: "name", Value: types.StringValue("admin")}},
			},
		},
	}

	files := write(t, plan)

	if want := `project_id = descope_project.main["prod"].id`; !strings.Contains(files["authorization.tf"], want) {
		t.Errorf("expected %q in authorization.tf:\n%s", want, files["authorization.tf"])
	}
	for _, name := range []string{"provider.tf", "variables.tf", "project.tf"} {
		if _, ok := files[name]; ok {
			t.Errorf("expected no %s, got:\n%s", name, files[name])
		}
	}
	if want := "to = module.auth.descope_role.admin"; !strings.Contains(files["import.tf"], want) {
		t.Errorf("expected %q in import.tf:\n%s", want, files["import.tf"])
	}
}

func TestParseProjectAddress(t *testing.T) {
	for _, valid := range []string{"descope_project.main", "descope_project.main[0]", `descope_project.main["prod"]`} {
		if _, err := ParseProjectAddress(valid); err != nil {
			t.Errorf("expected %s to be valid: %v", valid, err)
		}
	}
	for _, invalid := range []string{"", "main", "descope_project", "descope_role.main", "descope_project.main.id", "module.auth.descope_project.main", "descope_project.main[0][1]", "descope_project.1main", "descope_project.main + 1"} {
		if _, err := ParseProjectAddress(invalid); err == nil {
			t.Errorf("expected %s to be invalid", invalid)
		}
	}
}

func TestParseModulePrefix(t *testing.T) {
	for _, valid := range []string{"module.auth", "module.auth[0]", `module.auth["prod"]`, "module.auth.module.descope", `module.auth["prod"].module.descope[1]`} {
		if _, err := ParseModulePrefix(valid); err != nil {
			t.Errorf("expected %s to be valid: %v", valid, err)
		}
	}
	for _, invalid := range []string{"", "module", "auth", "module.auth.descope", "module.auth.module", "module.auth[0][1]", "descope_project.main", "module.auth."} {
		if _, err := ParseModulePrefix(invalid); err == nil {
			t.Errorf("expected %s to be invalid", invalid)
		}
	}
}

func TestSecretMapKeysBecomeVariables(t *testing.T) {
	plan := &Plan{
		ProjectID: "P123",
		Resources: []Resource{
			{
				Type:         "descope_http_connector",
				Label:        "webhook",
				EntityID:     "CI1",
				ImportID:     "P123/CI1",
				HasProjectID: true,
				Attrs: []prune.Attr{{Name: "secret_headers", IsNested: true, Nested: []prune.Attr{
					{Name: "X-Api-Key", Placeholder: true},
					{Name: "X.Trace", Placeholder: true},
				}}},
			},
		},
	}

	files := write(t, plan)
	all := strings.Join([]string{files["connectors.tf"], files["main.tf"]}, "\n")

	for _, want := range []string{"X-Api-Key = var.webhook_secret_headers_x_api_key", `"X.Trace" = var.webhook_secret_headers_x_trace`} {
		if !strings.Contains(all, want) {
			t.Errorf("expected %q in the generated resource:\n%s", want, all)
		}
	}
	for _, want := range []string{"webhook_secret_headers_x_api_key", "webhook_secret_headers_x_trace"} {
		if !strings.Contains(files["variables.tf"], `variable "`+want+`"`) {
			t.Errorf("expected a %q variable:\n%s", want, files["variables.tf"])
		}
	}
}
