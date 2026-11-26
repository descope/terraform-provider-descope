package main

import (
	"context"
	"os"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestExportGolden(t *testing.T) {
	// gated with the rest of the battery: it reads the shared acceptance project, which the provider's own suites mutate in parallel
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	projectID := os.Getenv("DESCOPE_TESTACC_PROJECT_ID")
	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	if projectID == "" || managementKey == "" {
		t.Skip("set DESCOPE_TESTACC_PROJECT_ID and DESCOPE_MANAGEMENT_KEY to run this test")
	}

	ctx := context.Background()
	outDir := t.TempDir()

	client := infra.NewClient("tfexport-test", managementKey, os.Getenv("DESCOPE_BASE_URL"))
	count, warnings, err := export.Run(ctx, client, projectID, outDir, export.Options{})
	if err != nil {
		t.Fatalf("export failed: %s", err)
	}
	t.Logf("exported %d resources with %d warnings", count, len(warnings))

	env := harness.Env(t)
	harness.WriteVariables(t, outDir, projectID)
	changes := harness.PlanJSON(t, outDir, env)
	if len(changes) == 0 {
		t.Fatal("expected resource changes in the plan")
	}
	if imported := harness.AssertBenignPlan(t, changes); imported != count {
		t.Errorf("expected %d imported resources, got %d", count, imported)
	}
}
