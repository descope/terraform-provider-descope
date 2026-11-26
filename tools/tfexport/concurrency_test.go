package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestConcurrencyStress(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	env := harness.Env(t)
	seedDir := t.TempDir()
	preamble := fmt.Sprintf(harness.SeedPreamble, fmt.Sprintf("testacc-tfexport-stress-%d", time.Now().UnixNano()%1000000))

	config := func(pass int) string {
		var blocks []string
		for i := 0; i < 6; i++ {
			blocks = append(blocks, fmt.Sprintf(`
resource "descope_generic_sms_gateway_connector" "stress_%d" {
  project_id = descope_project.test.id
  name       = "stress-sms-%d"
  post_url   = "https://stress.example.com/sms/%d/%d"
  sender     = "stress-%d"
}`, i, i, i, pass, pass))
		}
		blocks = append(blocks, fmt.Sprintf(`
resource "descope_magiclink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "%d minutes"
}

resource "descope_otp_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "%d minutes"
}

resource "descope_enchantedlink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "%d minutes"
}

resource "descope_password_settings" "main" {
  project_id = descope_project.test.id
  min_length = %d
}

resource "descope_invite_settings" "main" {
  project_id        = descope_project.test.id
  invite_expiration = "%d weeks"
}

resource "descope_session_settings" "main" {
  project_id               = descope_project.test.id
  session_token_expiration = "%d minutes"
}

resource "descope_totp_settings" "main" {
  project_id    = descope_project.test.id
  service_label = "stress %d"
}`, 5+pass, 6+pass, 7+pass, 8+pass, 2+pass, 9+pass, pass))
		return preamble + strings.Join(blocks, "\n")
	}

	before := backendConflictCount(t)
	const passes = 3
	for pass := 0; pass < passes; pass++ {
		if err := os.WriteFile(filepath.Join(seedDir, "main.tf"), []byte(config(pass)), 0o644); err != nil {
			t.Fatal(err)
		}
		harness.Run(t, seedDir, env, "apply", "-auto-approve", "-input=false", "-parallelism=20")
		t.Logf("stress pass %d applied", pass+1)
	}
	harness.Run(t, seedDir, env, "destroy", "-auto-approve", "-input=false", "-parallelism=20")

	if after := backendConflictCount(t); before >= 0 && after >= 0 {
		t.Logf("stress telemetry: %d backend version-conflict errors across %d full-parallel passes", after-before, passes)
	}
}

func backendConflictCount(t *testing.T) int {
	t.Helper()
	docker, err := exec.LookPath("docker")
	if err != nil {
		return -1
	}
	output, err := exec.Command(docker, "logs", "--since", "24h", "euw1-managementservice-1").CombinedOutput()
	if err != nil {
		return -1
	}
	return strings.Count(string(output), "E013008")
}
