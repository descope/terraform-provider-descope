package harness

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/read"
)

func RequireEnv(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}
}

type Step struct {
	Config string
	Mutate func(t *testing.T, client *infra.Client, projectID string)
	Verify func(t *testing.T, exportDir string)
}

type Scenario struct {
	Name          string
	Files         []string // testdata files copied next to main.tf for file() references
	ExpectMissing []string // resource types the exporter is known to not discover yet
	Steps         []Step
}

const SeedPreamble = `
terraform {
  required_providers {
    descope = {
      source = "descope/descope"
    }
  }
}

provider "descope" {
}

resource "descope_project" "test" {
  name                = %q
  deletion_protection = false
}

output "project_id" {
  value = descope_project.test.id
}
`

func RunScenario(t *testing.T, scenario Scenario) {
	ctx := context.Background()
	env := Env(t)
	client := infra.NewClient("tfexport-test", os.Getenv("DESCOPE_MANAGEMENT_KEY"), os.Getenv("DESCOPE_BASE_URL"))

	seedDir := t.TempDir()
	projectName := fmt.Sprintf("testacc-tfexport-%s-%d", scenario.Name, time.Now().UnixNano()%1000000)
	preamble := fmt.Sprintf(SeedPreamble, projectName)

	for _, file := range scenario.Files {
		content, err := os.ReadFile(filepath.Join("testdata", file))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(seedDir, file), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	// destroys are retried: parallel connector deletion fails spuriously and each pass deletes more
	destroy := func() {
		for attempt := 1; ; attempt++ {
			output, err := Try(seedDir, env, "destroy", "-auto-approve", "-input=false")
			if err == nil {
				return
			}
			if attempt == 3 {
				t.Errorf("terraform destroy failed after %d attempts: %s\n%s", attempt, err, output)
				return
			}
			t.Logf("retrying destroy (attempt %d failed): %s", attempt, err)
		}
	}
	destroyed := false
	t.Cleanup(func() {
		if !destroyed {
			destroy()
		}
	})

	for i, step := range scenario.Steps {
		stepName := fmt.Sprintf("step%d", i+1)
		if err := os.WriteFile(filepath.Join(seedDir, "main.tf"), []byte(preamble+step.Config), 0o644); err != nil {
			t.Fatal(err)
		}
		Run(t, seedDir, env, "apply", "-auto-approve", "-input=false")

		output := Run(t, seedDir, env, "output", "-raw", "project_id")
		projectID := strings.TrimSpace(string(output))
		if !strings.HasPrefix(projectID, "P") {
			t.Fatalf("%s: unexpected project id %q", stepName, projectID)
		}

		if step.Mutate != nil {
			step.Mutate(t, client, projectID)
			if t.Failed() {
				return
			}
		}

		exportDir := t.TempDir()
		count, warnings, err := export.Run(ctx, client, projectID, exportDir, export.Options{})
		if err != nil {
			t.Fatalf("%s: export failed: %s", stepName, err)
		}
		t.Logf("%s: exported %d resources with %d warnings", stepName, count, len(warnings))
		for _, warning := range warnings {
			t.Logf("%s: warning: %s", stepName, warning.Text)
			if strings.Contains(warning.Text, read.UnreadableWarning) {
				t.Errorf("%s: %s", stepName, warning.Text)
			}
		}

		assertCoverage(t, stepName, seedDir, exportDir, env, scenario.ExpectMissing)
		assertExportIntegrity(t, stepName, exportDir)
		assertFormatting(t, stepName, exportDir, env)

		secondDir := t.TempDir()
		if _, _, err := export.Run(ctx, client, projectID, secondDir, export.Options{}); err != nil {
			t.Fatalf("%s: second export failed: %s", stepName, err)
		}
		assertExportDeterminism(t, stepName, exportDir, secondDir)

		if step.Verify != nil {
			step.Verify(t, exportDir)
			if t.Failed() {
				return
			}
		}

		WriteVariables(t, exportDir, projectID)
		changes := PlanJSON(t, exportDir, env)
		planSummary(t, stepName, changes)
		if imported := AssertBenignPlan(t, changes); imported != count {
			t.Errorf("%s: expected %d imported resources, got %d", stepName, count, imported)
		}
		if step.Mutate == nil {
			seeded := seedStateResources(t, seedDir, env)
			assertStateConsistency(t, stepName, seeded, changes, scenario.ExpectMissing)
			assertExportedValues(t, stepName, seeded, changes, scenario.ExpectMissing)
		} else {
			t.Logf("%s: consistency skipped: the step mutated the project out-of-band, so seed state and server intentionally diverge", stepName)
		}
		if t.Failed() {
			return // don't apply an unexpected plan, and skip the remaining steps
		}

		if os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP_NOAPPLY") == "" {
			// the apply is retried once: concurrent writes can fail spuriously on backend version conflicts
			if output, err := Try(exportDir, env, "apply", "-input=false", "plan.bin"); err != nil {
				t.Logf("%s: exported apply failed, re-planning and retrying once (backend write contention):\n%s", stepName, output)
				retry := PlanJSON(t, exportDir, env)
				AssertBenignPlan(t, retry)
				if t.Failed() {
					return
				}
				Run(t, exportDir, env, "apply", "-input=false", "plan.bin")
			}
			AssertNoopPlan(t, PlanJSON(t, exportDir, env))
			if t.Failed() {
				return
			}
		}
	}

	destroy()
	destroyed = true
}

var resourceBlockPattern = regexp.MustCompile(`(?m)^resource "(descope_[a-z0-9_]+)" `)

func assertCoverage(t *testing.T, stepName, seedDir, exportDir string, env []string, expectMissing []string) {
	t.Helper()

	seeded := map[string]bool{}
	for _, line := range strings.Split(string(Run(t, seedDir, env, "state", "list")), "\n") {
		resourceType, _, ok := strings.Cut(strings.TrimSpace(line), ".")
		if ok && strings.HasPrefix(resourceType, "descope_") {
			seeded[resourceType] = true
		}
	}

	exported := map[string]bool{}
	entries, err := os.ReadDir(exportDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(exportDir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range resourceBlockPattern.FindAllStringSubmatch(string(content), -1) {
			exported[match[1]] = true
		}
	}

	for resourceType := range seeded {
		if !exported[resourceType] && !slices.Contains(expectMissing, resourceType) {
			t.Errorf("%s: discovery gap: seeded %s does not appear in the export", stepName, resourceType)
		}
	}
}
