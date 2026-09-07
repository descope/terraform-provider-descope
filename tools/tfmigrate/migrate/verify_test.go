package migrate

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestVerifyPlan(t *testing.T) {
	const (
		accessKeyAddress = "descope_access_key.migrated"
		settingsAddress  = "descope_project_settings.migrated"
	)

	tests := []struct {
		name         string
		fixture      string
		prepare      func(*Manifest)
		wantOK       bool
		wantChecked  int
		wantImported int
		wantFindings []Finding
		wantError    string
		wantReport   string
	}{
		{
			name:         "accepts all no-op imports",
			fixture:      "accepts_imports",
			wantOK:       true,
			wantChecked:  2,
			wantImported: 2,
			wantReport: "adoption plan verification: PASS\n" +
				"checked 2 managed changes, 2 imports\n",
		},
		{
			name:         "accepts unrelated managed no-op and skips data changes",
			fixture:      "accepts_unrelated_noop",
			wantOK:       true,
			wantChecked:  3,
			wantImported: 2,
		},
		{
			name:         "accepts imported secret update",
			fixture:      "accepts_secret_update",
			wantOK:       true,
			wantChecked:  2,
			wantImported: 2,
		},
		{
			name:         "accepts imported secret after-unknown update",
			fixture:      "accepts_secret_after_unknown",
			wantOK:       true,
			wantChecked:  2,
			wantImported: 2,
		},
		{
			name:         "rejects create",
			fixture:      "rejects_create",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedCreate, Address: accessKeyAddress, Actions: []string{"create"},
				Detail: "create is not allowed during adoption",
			}},
			wantReport: "adoption plan verification: FAIL\n" +
				"checked 2 managed changes, 2 imports\n" +
				"  unexpected_create descope_access_key.migrated: create is not allowed during adoption [actions: create]\n",
		},
		{
			name:         "rejects delete",
			fixture:      "rejects_delete",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedDelete, Address: accessKeyAddress, Actions: []string{"delete"},
				Detail: "delete is not allowed during adoption",
			}},
		},
		{
			name:         "rejects replacement",
			fixture:      "rejects_replace",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedReplace, Address: accessKeyAddress, Actions: []string{"delete", "create"},
				Detail: "replacement is not allowed during adoption",
			}},
		},
		{
			name:         "rejects forget",
			fixture:      "rejects_forget",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingForgottenObject, Address: accessKeyAddress, Actions: []string{"forget"},
				Detail: "forget is not allowed during adoption",
			}},
		},
		{
			name:         "rejects read",
			fixture:      "rejects_read",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedRead, Address: accessKeyAddress, Actions: []string{"read"},
				Detail: "read is not allowed during adoption",
			}},
		},
		{
			name:         "rejects imported update of non-secret attributes",
			fixture:      "rejects_nonsecret_update",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedUpdate, Address: accessKeyAddress, Actions: []string{"update"},
				Detail: "update changes non-secret attributes: name, project_id",
			}},
		},
		{
			name:         "rejects update without importing",
			fixture:      "rejects_update_without_import",
			wantChecked:  3,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedUpdate, Address: "descope_access_key.unrelated", Actions: []string{"update"},
				Detail: "update is not part of an import",
			}},
		},
		{
			name:         "rejects missing manifest import",
			fixture:      "rejects_missing_import",
			wantChecked:  2,
			wantImported: 1,
			wantFindings: []Finding{{
				Kind: FindingMissingImport, Address: accessKeyAddress,
				Detail: "manifest entry is not imported by the plan",
			}},
		},
		{
			name:         "rejects import outside manifest",
			fixture:      "rejects_unexpected_import",
			wantChecked:  3,
			wantImported: 3,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedImport, Address: "descope_flow.unreviewed", Actions: []string{"no-op"},
				Detail: "plan imports an address that is not in the migration manifest",
			}},
		},
		{
			name:         "rejects wrong import id",
			fixture:      "rejects_wrong_import_id",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedImport, Address: accessKeyAddress, Actions: []string{"no-op"},
				Detail: `plan imports id "project/wrong-access-key", manifest requires "project/access-key"`,
			}},
		},
		{
			name:         "rejects provider mismatch",
			fixture:      "rejects_provider_mismatch",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingProviderMismatch, Address: accessKeyAddress, Actions: []string{"no-op"},
				Detail: `plan uses provider "registry.terraform.io/example/descope", manifest requires "registry.terraform.io/descope/descope"`,
			}},
		},
		{
			name:         "rejects errored plan",
			fixture:      "rejects_errored",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingPlanErrored, Detail: "Terraform reported an errored plan",
			}},
		},
		{
			name:         "rejects unknown action",
			fixture:      "rejects_unknown_action",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnknownAction, Address: accessKeyAddress, Actions: []string{"migrate"},
				Detail: "unrecognized action \"migrate\"",
			}},
		},
		{
			name:         "rejects empty imported update",
			fixture:      "rejects_empty_secret_update",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind: FindingUnexpectedUpdate, Address: accessKeyAddress, Actions: []string{"update"},
				Detail: "update has no changed attributes",
			}},
		},
		{
			name:         "sorts findings",
			fixture:      "rejects_sorted_findings",
			wantChecked:  2,
			wantImported: 1,
			wantFindings: []Finding{
				{Kind: FindingMissingImport, Address: settingsAddress, Detail: "manifest entry is not imported by the plan"},
				{Kind: FindingPlanErrored, Detail: "Terraform reported an errored plan"},
				{
					Kind: FindingUnexpectedCreate, Address: accessKeyAddress, Actions: []string{"create"},
					Detail: "create is not allowed during adoption",
				},
			},
		},
		{
			name:         "rejects incomplete plan and hidden operations",
			fixture:      "rejects_incomplete_operations",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{
				{Kind: FindingActionInvocation, Detail: "plan contains 1 action invocation(s)"},
				{Kind: FindingDeferredChange, Detail: "plan contains 1 deferred change(s)"},
				{Kind: FindingIncompletePlan, Detail: "Terraform plan completeness was not confirmed; use Terraform 1.8 or newer and do not use -target"},
			},
		},
		{
			name:         "rejects plan without completeness metadata",
			fixture:      "rejects_missing_complete",
			wantChecked:  2,
			wantImported: 2,
			wantFindings: []Finding{{
				Kind:   FindingIncompletePlan,
				Detail: "Terraform plan completeness was not confirmed; use Terraform 1.8 or newer and do not use -target",
			}},
		},
		{
			name:        "errors on unsupported plan format",
			fixture:     "errors_unsupported_format",
			wantError:   "expected major version 1",
			wantChecked: 0,
		},
		{
			name:    "errors when manifest is not ready",
			fixture: "errors_manifest_not_ready",
			prepare: func(manifest *Manifest) {
				manifest.Ready = false
			},
			wantError: "manifest that is not ready",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			planJSON := readPlanFixture(t, test.fixture)
			planSnapshot := append([]byte(nil), planJSON...)
			manifest := testManifest()
			if test.prepare != nil {
				test.prepare(manifest)
			}
			manifestSnapshot, err := json.Marshal(manifest)
			if err != nil {
				t.Fatalf("marshaling manifest snapshot: %v", err)
			}

			result, err := VerifyPlan(planJSON, manifest)
			if test.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantError) {
					t.Fatalf("VerifyPlan() error = %v, want an error containing %q", err, test.wantError)
				}
				if result != nil {
					t.Fatalf("VerifyPlan() result = %#v, want nil on error", result)
				}
			} else {
				if err != nil {
					t.Fatalf("VerifyPlan() unexpected error: %v", err)
				}
				if result.OK != test.wantOK {
					t.Errorf("OK = %v, want %v", result.OK, test.wantOK)
				}
				if result.Checked != test.wantChecked {
					t.Errorf("Checked = %d, want %d", result.Checked, test.wantChecked)
				}
				if result.Imported != test.wantImported {
					t.Errorf("Imported = %d, want %d", result.Imported, test.wantImported)
				}
				if !reflect.DeepEqual(result.Findings, test.wantFindings) {
					t.Errorf("Findings = %#v, want %#v", result.Findings, test.wantFindings)
				}
				if test.wantReport != "" && result.Report() != test.wantReport {
					t.Errorf("Report() = %q, want %q", result.Report(), test.wantReport)
				}
			}

			if !bytes.Equal(planJSON, planSnapshot) {
				t.Error("VerifyPlan mutated the plan input")
			}
			manifestAfter, marshalErr := json.Marshal(manifest)
			if marshalErr != nil {
				t.Fatalf("marshaling manifest after verification: %v", marshalErr)
			}
			if !bytes.Equal(manifestAfter, manifestSnapshot) {
				t.Error("VerifyPlan mutated the manifest")
			}
		})
	}
}

func TestVerifyPlanNestedSecretPaths(t *testing.T) {
	manifest := &Manifest{
		Ready:          true,
		TargetProvider: ProviderPin{Source: ProviderSource},
		Entries: []Entry{{
			ImportID:         "project/connector",
			Address:          "descope_http_connector.example",
			SecretAttributes: []string{"authentication.basic.password"},
		}},
	}
	plan := func(after map[string]any) []byte {
		value := map[string]any{
			"format_version": "1.2",
			"complete":       true,
			"resource_changes": []any{map[string]any{
				"address":       "descope_http_connector.example",
				"mode":          "managed",
				"provider_name": "registry.terraform.io/descope/descope",
				"change": map[string]any{
					"actions":   []string{"update"},
					"importing": map[string]any{"id": "project/connector"},
					"before": map[string]any{"authentication": map[string]any{
						"basic": map[string]any{"username": "user", "password": nil},
					}},
					"after":         map[string]any{"authentication": map[string]any{"basic": after}},
					"after_unknown": map[string]any{},
				},
			}},
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return encoded
	}

	accepted, err := VerifyPlan(plan(map[string]any{"username": "user", "password": "secret"}), manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !accepted.OK {
		t.Fatalf("secret-only nested update rejected: %s", accepted.Report())
	}

	rejected, err := VerifyPlan(plan(map[string]any{"username": "different", "password": "secret"}), manifest)
	if err != nil {
		t.Fatal(err)
	}
	if rejected.OK || len(rejected.Findings) != 1 || !strings.Contains(rejected.Findings[0].Detail, "authentication.basic.username") {
		t.Fatalf("non-secret nested update accepted: %#v", rejected.Findings)
	}
	mapManifest := &Manifest{
		Ready:          true,
		TargetProvider: ProviderPin{Source: ProviderSource},
		Entries: []Entry{{
			Address:          "descope_http_connector.headers",
			ImportID:         "project/headers",
			SecretAttributes: []string{"secret_headers"},
		}},
	}
	mapPlan, err := json.Marshal(map[string]any{
		"format_version": "1.2",
		"complete":       true,
		"resource_changes": []any{map[string]any{
			"address":       "descope_http_connector.headers",
			"mode":          "managed",
			"provider_name": officialProviderSource,
			"change": map[string]any{
				"actions":   []string{"update"},
				"importing": map[string]any{"id": "project/headers"},
				"before":    map[string]any{"secret_headers": map[string]any{"Authorization": nil}},
				"after":     map[string]any{"secret_headers": map[string]any{"Authorization": "secret"}},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	mapResult, err := VerifyPlan(mapPlan, mapManifest)
	if err != nil {
		t.Fatal(err)
	}
	if !mapResult.OK {
		t.Fatalf("sensitive map update rejected: %s", mapResult.Report())
	}
}

func testManifest() *Manifest {
	return &Manifest{
		Ready:          true,
		TargetProvider: ProviderPin{Source: "descope/descope"},
		Entries: []Entry{
			{
				Address:          "descope_access_key.migrated",
				ImportID:         "project/access-key",
				SecretAttributes: []string{"client_secret"},
			},
			{Address: "descope_project_settings.migrated", ImportID: "project"},
		},
	}
}

func readPlanFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "plans", name+".json"))
	if err != nil {
		t.Fatalf("reading plan fixture %s: %v", name, err)
	}
	return data
}
