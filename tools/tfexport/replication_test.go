package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestExportReplication(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	ctx := context.Background()
	env := harness.Env(t)
	client := infra.NewClient("tfexport-test", os.Getenv("DESCOPE_MANAGEMENT_KEY"), os.Getenv("DESCOPE_BASE_URL"))

	seedDir := t.TempDir()
	flow, err := os.ReadFile(filepath.Join("testdata", "roundtrip-flow.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(seedDir, "roundtrip-flow.json"), flow, 0o644); err != nil {
		t.Fatal(err)
	}
	nameA := fmt.Sprintf("testacc-tfexport-replica-a-%d", time.Now().UnixNano()%1000000)
	seed := fmt.Sprintf(harness.SeedPreamble, nameA) + `
resource "descope_smtp_connector" "mail" {
  project_id   = descope_project.test.id
  name         = "replica-smtp"
  host         = "smtp.example.com"
  port         = 587
  username     = "mailer"
  password     = "not-a-real-password"
  sender_email = "auth@example.com"
}

resource "descope_email_template" "magic" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "replica-template"
  subject    = "Sign in"
  html_body  = "Follow the link"
}

resource "descope_magiclink_settings" "main" {
  project_id        = descope_project.test.id
  email_template_id = descope_email_template.magic.id
  email_connector_id = descope_smtp_connector.mail.id
}

resource "descope_permission" "use" {
  project_id = descope_project.test.id
  name       = "replica.use"
}

resource "descope_role" "user" {
  project_id  = descope_project.test.id
  name        = "replica-user"
  permissions = [descope_permission.use.name]
}

resource "descope_flow" "passkeys" {
  project_id = descope_project.test.id
  flow_id    = "replica-passkeys"
  data       = file("${path.module}/roundtrip-flow.json")
}

resource "descope_jwt_template" "claims" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "replica-jwt"
  template   = jsonencode({ replica = true })
}
`
	if err := os.WriteFile(filepath.Join(seedDir, "main.tf"), []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.Run(t, seedDir, env, "apply", "-auto-approve", "-input=false")
	t.Cleanup(func() {
		harness.Run(t, seedDir, env, "destroy", "-auto-approve", "-input=false")
	})
	projectA := strings.TrimSpace(string(harness.Run(t, seedDir, env, "output", "-raw", "project_id")))

	exportA := t.TempDir()
	if _, _, err := export.Run(ctx, client, projectA, exportA, export.Options{}); err != nil {
		t.Fatalf("export of project A failed: %s", err)
	}

	projectDir := t.TempDir()
	nameB := fmt.Sprintf("testacc-tfexport-replica-b-%d", time.Now().UnixNano()%1000000)
	if err := os.WriteFile(filepath.Join(projectDir, "main.tf"), fmt.Appendf(nil, harness.SeedPreamble, nameB), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.Run(t, projectDir, env, "apply", "-auto-approve", "-input=false")
	t.Cleanup(func() {
		harness.Run(t, projectDir, env, "destroy", "-auto-approve", "-input=false")
	})
	projectB := strings.TrimSpace(string(harness.Run(t, projectDir, env, "output", "-raw", "project_id")))

	if err := os.Remove(filepath.Join(exportA, "import.tf")); err != nil {
		t.Fatal(err)
	}
	stripBuiltinBlocks(t, exportA)
	retargetProject(t, exportA, projectB)
	harness.WriteVariables(t, exportA, projectB)
	harness.Run(t, exportA, env, "apply", "-auto-approve", "-input=false")
	t.Cleanup(func() {
		// unmanage B's resources before the project itself is destroyed
		_ = os.RemoveAll(filepath.Join(exportA, "terraform.tfstate"))
	})
	t.Logf("applied project A's export onto project B as a pure create")

	harness.AssertNoopPlan(t, harness.PlanJSON(t, exportA, env))

	exportB := t.TempDir()
	if _, _, err := export.Run(ctx, client, projectB, exportB, export.Options{}); err != nil {
		t.Fatalf("export of project B failed: %s", err)
	}
	a, b := exportSignature(t, exportA, nameA), exportSignature(t, exportB, nameB)
	for key, countA := range a {
		if b[key] != countA {
			t.Errorf("replication: %s appears %d times in A's export but %d in B's", key, countA, b[key])
		}
	}
	for key, countB := range b {
		if _, ok := a[key]; !ok {
			t.Errorf("replication: %s appears only in B's export (%d times)", key, countB)
		}
	}
	t.Logf("replication verified: %d distinct resource signatures match across projects", len(a))
}

var projectReference = regexp.MustCompile(`descope_project\.[a-z0-9_]+\.id`)

func retargetProject(t *testing.T, exportDir, projectID string) {
	t.Helper()
	if err := os.Remove(filepath.Join(exportDir, "project.tf")); err != nil {
		t.Fatal(err)
	}
	paths, err := filepath.Glob(filepath.Join(exportDir, "*.tf"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		retargeted := projectReference.ReplaceAllString(string(content), fmt.Sprintf("%q", projectID))
		if err := os.WriteFile(path, []byte(retargeted), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

var builtinLabels = regexp.MustCompile(`(?s)resource "(descope_oidc_app" "oidc_default_application|descope_jwt_template" "default_[a-z0-9_]*)".*?\n\}\n`)

var builtinReferences = regexp.MustCompile(`(?m)^\s*[a-z0-9_]+\s*= descope_[a-z0-9_]+\.(default_[a-z0-9_]*|oidc_default_application)\.[a-z_]+\n`)

func stripBuiltinBlocks(t *testing.T, exportDir string) {
	t.Helper()
	removed := 0
	for _, name := range []string{"apps.tf", "templates.tf", "auth.tf"} {
		path := filepath.Join(exportDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		stripped := builtinLabels.ReplaceAllString(string(content), "")
		stripped = builtinReferences.ReplaceAllString(stripped, "")
		removed += (len(content) - len(stripped)) / 80
		if err := os.WriteFile(path, []byte(stripped), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("stripped built-in entity blocks before replication (~%d lines)", removed)
}

func exportSignature(t *testing.T, exportDir, projectName string) map[string]int {
	t.Helper()
	// singletons are labelled by the project they configure, so the label legitimately differs between the two projects
	projectLabel := strings.ReplaceAll(projectName, "-", "_")
	signature := map[string]int{}
	configs, _ := harness.ExportFiles(t, exportDir)
	for name, content := range configs {
		if name == "import.tf" {
			continue
		}
		for _, match := range harness.ResourceHeaderPattern.FindAllStringSubmatch(content, -1) {
			if match[2] == "oidc_default_application" || strings.HasPrefix(match[2], "default_") {
				continue // built-in entities exist in every project and are not replicated
			}
			if match[1] == "descope_project" {
				continue // the project is the target being replicated into, not one of the replicated entities
			}
			label := match[2]
			if label == projectLabel {
				label = "singleton"
			}
			signature[match[1]+"."+label]++
		}
	}
	return signature
}
