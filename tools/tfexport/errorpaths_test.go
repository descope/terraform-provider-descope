package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestExportErrorPaths(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	if managementKey == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	ctx := context.Background()
	client := infra.NewClient("tfexport-test", managementKey, os.Getenv("DESCOPE_BASE_URL"))

	t.Run("NonexistentProject", func(t *testing.T) {
		outDir := t.TempDir()
		_, _, err := export.Run(ctx, client, "P0000000000000000000000000000", outDir, export.Options{})
		if err == nil {
			t.Fatal("expected an error exporting a nonexistent project")
		}
		if !strings.Contains(err.Error(), "snapshot") {
			t.Errorf("expected a clear snapshot error, got: %s", err)
		}
		entries, _ := os.ReadDir(outDir)
		if len(entries) != 0 {
			t.Errorf("expected no partial output for a failed export, found %d entries", len(entries))
		}
	})

	t.Run("FilterMatchesNothing", func(t *testing.T) {
		projectID := os.Getenv("DESCOPE_TESTACC_PROJECT_ID")
		if projectID == "" {
			t.Skip("set DESCOPE_TESTACC_PROJECT_ID to run this test")
		}
		outDir := t.TempDir()
		count, _, err := export.Run(ctx, client, projectID, outDir, export.Options{Only: "no_such_resource_type"})
		if err != nil {
			t.Fatalf("export failed: %s", err)
		}
		if count != 0 {
			t.Errorf("expected 0 resources, got %d", count)
		}
		env := harness.Env(t)
		if output, err := harness.Try(outDir, env, "validate"); err != nil {
			t.Errorf("empty export does not validate: %s\n%s", err, output)
		}
	})
}
