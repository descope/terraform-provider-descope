package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestExportDayTwo(t *testing.T) {
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
	preamble := fmt.Sprintf(harness.SeedPreamble, fmt.Sprintf("testacc-tfexport-daytwo-%d", time.Now().UnixNano()%1000000))
	seedA := preamble + `
resource "descope_role" "keeper" {
  project_id  = descope_project.test.id
  name        = "daytwo-keeper"
  description = "original"
}

resource "descope_otp_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "5 minutes"
}
`
	seedB := preamble + `
resource "descope_role" "keeper" {
  project_id  = descope_project.test.id
  name        = "daytwo-keeper"
  description = "evolved"
}

resource "descope_otp_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "9 minutes"
}

resource "descope_role" "added" {
  project_id = descope_project.test.id
  name       = "daytwo-added"
}
`

	if err := os.WriteFile(filepath.Join(seedDir, "main.tf"), []byte(seedA), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.Run(t, seedDir, env, "apply", "-auto-approve", "-input=false")
	t.Cleanup(func() {
		harness.Run(t, seedDir, env, "destroy", "-auto-approve", "-input=false")
	})
	projectID := strings.TrimSpace(string(harness.Run(t, seedDir, env, "output", "-raw", "project_id")))

	workspace := t.TempDir()
	if _, _, err := export.Run(ctx, client, projectID, workspace, export.Options{}); err != nil {
		t.Fatalf("day one export failed: %s", err)
	}
	harness.WriteVariables(t, workspace, projectID)
	changes := harness.PlanJSON(t, workspace, env)
	harness.AssertBenignPlan(t, changes)
	if t.Failed() {
		return
	}
	harness.Run(t, workspace, env, "apply", "-input=false", "plan.bin")
	harness.AssertNoopPlan(t, harness.PlanJSON(t, workspace, env))
	t.Logf("day one: adopted %d resources and settled", len(changes))

	if err := os.WriteFile(filepath.Join(seedDir, "main.tf"), []byte(seedB), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.Run(t, seedDir, env, "apply", "-auto-approve", "-input=false")

	if _, _, err := export.Run(ctx, client, projectID, workspace, export.Options{}); err != nil {
		t.Fatalf("day two export failed: %s", err)
	}
	harness.WriteVariables(t, workspace, projectID)
	managed := map[string]bool{}
	for _, line := range strings.Split(string(harness.Run(t, workspace, env, "state", "list")), "\n") {
		managed[strings.TrimSpace(line)] = true
	}
	importsFile := filepath.Join(workspace, "import.tf")
	content, err := os.ReadFile(importsFile)
	if err != nil {
		t.Fatal(err)
	}
	var kept []string
	pruned := 0
	for _, block := range strings.Split(string(content), "\n\n") {
		match := harness.ImportToPattern.FindStringSubmatch(block)
		if match != nil && managed[match[1]+"."+match[2]] {
			pruned++
			continue
		}
		kept = append(kept, block)
	}
	if err := os.WriteFile(importsFile, []byte(strings.Join(kept, "\n\n")), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("day two: pruned %d import blocks for already-managed resources", pruned)

	imported, updated := 0, 0
	for _, change := range harness.PlanJSON(t, workspace, env) {
		actions := strings.Join(change.Change.Actions, ",")
		if len(change.Change.Importing) > 0 {
			imported++
			if name, _ := change.Change.After["name"].(string); name != "daytwo-added" {
				t.Errorf("day two: unexpected import of %s (%s)", change.Address, name)
			}
		}
		switch actions {
		case "no-op":
		case "update":
			dataOnly := true
			for _, attribute := range harness.ChangedAttributes(change.Change.Before, change.Change.After) {
				if attribute == "data" {
					// writing a flow bumps metadata.componentsVersion, so the first re-export after an apply rewrites flow data files once
					continue
				}
				dataOnly = false
				if attribute != "description" && attribute != "expiration_time" {
					t.Errorf("day two: unexpected update to %s attribute %q", change.Address, attribute)
				}
			}
			if !dataOnly {
				updated++
			}
		default:
			t.Errorf("day two: unexpected plan actions %q for %s", actions, change.Address)
		}
	}
	if imported != 1 {
		t.Errorf("day two: expected exactly 1 newly adopted resource, got %d", imported)
	}
	if updated != 0 {
		t.Errorf("day two: expected no updated resources, got %d", updated)
	}

	authorization, err := os.ReadFile(filepath.Join(workspace, "authorization.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(authorization), `"evolved"`) {
		t.Errorf("day two: expected the evolved role description in the re-export:\n%s", authorization)
	}
	auth, err := os.ReadFile(filepath.Join(workspace, "auth.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(auth), `"9 minutes"`) {
		t.Errorf("day two: expected the evolved otp expiration in the re-export:\n%s", auth)
	}

	harness.Run(t, workspace, env, "apply", "-input=false", "plan.bin")
	harness.AssertNoopPlan(t, harness.PlanJSON(t, workspace, env))
	t.Logf("day two: adopted the new resource, applied the drift, and settled")
}
