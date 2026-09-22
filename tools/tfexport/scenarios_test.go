package main

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func roundTripScenarios(t *testing.T) []harness.Scenario {
	t.Helper()
	scenarios := harness.LoadScenarios(t, "testdata/scenarios")
	for i := range scenarios {
		if hook, ok := scenarioHooks[scenarios[i].Name]; ok {
			hook(&scenarios[i])
		}
	}
	for name := range scenarioHooks {
		if !slices.ContainsFunc(scenarios, func(s harness.Scenario) bool { return s.Name == name }) {
			t.Errorf("hooks defined for unknown scenario %q", name)
		}
	}
	return scenarios
}

var scenarioHooks = map[string]func(s *harness.Scenario){
	"flows-and-styles": func(s *harness.Scenario) {
		s.Files = []string{"roundtrip-flow.json", "roundtrip-styles.json", "roundtrip-styles-v2.json"}
	},
	"access-keys": func(s *harness.Scenario) {
		s.ExpectMissing = []string{"descope_access_key"}
	},
	"console-settings": func(s *harness.Scenario) {
		s.Steps[0].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			err := client.Post(context.Background(), projectID, "/v1/mgmt/magiclink/settings", map[string]any{
				"expirationTime":     7,
				"expirationTimeUnit": "minutes",
			})
			if err != nil {
				t.Errorf("out-of-band settings update failed: %s", err)
			}
		}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "auth.tf"))
			if err != nil {
				t.Fatalf("reading exported auth.tf: %s", err)
			}
			if !strings.Contains(string(content), `expiration_time = "7 minutes"`) {
				t.Errorf("expected the console-edited expiration to be exported, got:\n%s", content)
			}
		}
	},
	"out-of-band": func(s *harness.Scenario) {
		s.Steps[0].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			_, err := client.PostData(context.Background(), projectID, "/v1/mgmt/role/create", map[string]any{
				"name":        "roundtrip-console",
				"description": "Created out-of-band",
			})
			if err != nil {
				t.Errorf("out-of-band role creation failed: %s", err)
			}
		}
	},
	"file-extraction-edges": func(s *harness.Scenario) {
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "templates", "big_template.html"))
			if err != nil {
				t.Fatalf("expected the large body to be extracted to a file: %s", err)
			}
			for _, want := range []string{"${not_interpolated}", "%{ nope }", "ünïcödé 🎯", "filler paragraph number 39"} {
				if !strings.Contains(string(content), want) {
					t.Errorf("extracted body is missing %q", want)
				}
			}
		}
	},
	"settings-slices": func(s *harness.Scenario) {
		s.Steps[1].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "auth.tf"))
			if err != nil {
				t.Fatal(err)
			}
			// the sibling slices of the shared endpoint must be untouched by the session settings write
			for _, want := range []string{
				`invite_expiration`,
				`"2 weeks"`,
				`app_url`,
				`"https://app.example.com"`,
				`"45 minutes"`,
			} {
				if !strings.Contains(string(content), want) {
					t.Errorf("expected exported settings to contain %s, sibling slice may have been clobbered:\n%s", want, content)
				}
			}
		}
	},
	"settings-console-slices": func(s *harness.Scenario) {
		s.Steps[0].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			err := client.Post(context.Background(), projectID, "/v1/mgmt/project/settings", map[string]any{
				"refreshTokenExpiration":     2,
				"refreshTokenExpirationUnit": "weeks",
			})
			if err != nil {
				t.Errorf("out-of-band project settings update failed: %s", err)
			}
		}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "auth.tf"))
			if err != nil {
				t.Fatal(err)
			}
			// the console edit must appear and the terraform-managed sibling slice must have survived it (hclwrite aligns the = signs)
			for _, want := range []string{`refresh_token_expiration\s+= "2 weeks"`, `invite_expiration\s+= "3 weeks"`, `vendor`} {
				if !regexp.MustCompile(want).Match(content) {
					t.Errorf("expected exported settings to match %s:\n%s", want, content)
				}
			}
		}
	},
	"rename-and-dangling": func(s *harness.Scenario) {
		s.ExpectMissing = []string{"descope_email_template", "descope_magiclink_settings"}
		s.Steps[0].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			ctx := context.Background()
			err := client.Post(ctx, projectID, "/v1/mgmt/permission/update", map[string]any{
				"name":    "roundtrip.old",
				"newName": "roundtrip.renamed",
			})
			if err != nil {
				t.Errorf("out-of-band permission rename failed: %s", err)
				return
			}
			files, err := discover.Snapshot(ctx, client, projectID)
			if err != nil {
				t.Errorf("snapshot failed: %s", err)
				return
			}
			settings, _ := files["auth/magiclink.json"].(map[string]any)
			templates, _ := settings["emailTemplates"].([]any)
			for _, entry := range templates {
				template, _ := entry.(map[string]any)
				if template["name"] == "dangle-template" {
					id, _ := template["id"].(string)
					if err := client.Del(ctx, projectID, "/v1/mgmt/template", map[string]string{"type": "email", "method": "magiclink", "id": id}); err != nil {
						t.Errorf("out-of-band template deletion failed: %s", err)
					}
					return
				}
			}
			t.Error("dangle-template not found in snapshot")
		}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			authorization, err := os.ReadFile(filepath.Join(exportDir, "authorization.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(authorization), "roundtrip.renamed") {
				t.Errorf("expected the renamed permission in the export:\n%s", authorization)
			}
			t.Logf("authorization.tf after rename:\n%s", authorization)
			if auth, err := os.ReadFile(filepath.Join(exportDir, "auth.tf")); err == nil {
				t.Logf("auth.tf after template deletion:\n%s", auth)
			}
		}
	},
	"pagination": func(s *harness.Scenario) {
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "authorization.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if roles := strings.Count(string(content), `resource "descope_role"`); roles != 35 {
				t.Errorf("expected 35 exported roles, got %d - discovery may be paginated", roles)
			}
			if permissions := strings.Count(string(content), `resource "descope_permission"`); permissions != 25 {
				t.Errorf("expected 25 exported permissions, got %d - discovery may be paginated", permissions)
			}
			misc, err := os.ReadFile(filepath.Join(exportDir, "misc.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if lists := strings.Count(string(misc), `resource "descope_list"`); lists != 10 {
				t.Errorf("expected 10 exported lists, got %d - discovery may be paginated", lists)
			}
		}
	},
	"widget-takeover": func(s *harness.Scenario) {
		s.Files = []string{"roundtrip-widget.json"}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "flows.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(content), `resource "descope_widget"`) {
				t.Errorf("expected the materialized widget in the export:\n%s", content)
			}
			if _, err := os.Stat(filepath.Join(exportDir, "widgets", "user_profile.json")); err != nil {
				t.Errorf("expected the widget data extracted to a file: %s", err)
			}
		}
	},
	"polarity-flips": func(s *harness.Scenario) {
		s.ExpectMissing = []string{"descope_oauth_settings", "descope_otp_settings"}
	},
	"flow-variants": func(s *harness.Scenario) {
		s.Files = []string{"roundtrip-flow.json", "roundtrip-flow-signin.json"}
		s.Steps[0].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			if err := client.Del(context.Background(), projectID, "/v1/mgmt/flow", map[string]string{"id": "sign-in"}); err != nil {
				t.Errorf("out-of-band flow deletion failed: %s", err)
			}
		}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "flows.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(content), `"sign-in"`) {
				t.Errorf("expected the deleted sign-in flow to be absent from the export:\n%s", content)
			}
			if !strings.Contains(string(content), `"passkeys-copy"`) {
				t.Errorf("expected the passkeys-copy flow in the export:\n%s", content)
			}
		}
	},
	"connector-rename": func(s *harness.Scenario) {
		s.ExpectMissing = []string{"descope_smtp_connector", "descope_magiclink_settings"}
		s.Steps[1].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			files, err := discover.Snapshot(context.Background(), client, projectID)
			if err != nil {
				t.Errorf("snapshot failed: %s", err)
				return
			}
			for path, value := range files {
				connector, ok := value.(map[string]any)
				if !ok || !strings.HasPrefix(path, "connectors/") || connector["name"] != "mail-beta" {
					continue
				}
				id, _ := connector["id"].(string)
				if err := client.Del(context.Background(), projectID, "/v1/mgmt/connector", map[string]string{"id": id}); err != nil {
					t.Errorf("out-of-band connector deletion failed: %s", err)
				}
				return
			}
			t.Error("mail-beta connector not found in snapshot")
		}
		s.Steps[0].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			content, err := os.ReadFile(filepath.Join(exportDir, "auth.tf"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(content), "descope_smtp_connector.mail_alpha.id") {
				t.Errorf("expected the settings to reference the mail_alpha connector:\n%s", content)
			}
		}
		s.Steps[1].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			if content, err := os.ReadFile(filepath.Join(exportDir, "auth.tf")); err == nil {
				t.Logf("auth.tf after connector deletion:\n%s", content)
			} else {
				t.Logf("no auth.tf after connector deletion (settings all default)")
			}
		}
	},
	"lifecycle-marathon": func(s *harness.Scenario) {
		s.Steps[2].Mutate = func(t *testing.T, client *infra.Client, projectID string) {
			t.Helper()
			ctx := context.Background()
			if _, err := client.PostData(ctx, projectID, "/v1/mgmt/role/create", map[string]any{
				"name": "marathon-console", "description": "console drift",
			}); err != nil {
				t.Errorf("out-of-band role creation failed: %s", err)
			}
			if err := client.Post(ctx, projectID, "/v1/mgmt/otp/settings", map[string]any{
				"expirationTime": 8, "expirationTimeUnit": "minutes",
			}); err != nil {
				t.Errorf("out-of-band settings update failed: %s", err)
			}
		}
		s.Steps[2].Verify = func(t *testing.T, exportDir string) {
			t.Helper()
			authorization, _ := os.ReadFile(filepath.Join(exportDir, "authorization.tf"))
			if !strings.Contains(string(authorization), "marathon-console") {
				t.Errorf("expected the console-created role in the export:\n%s", authorization)
			}
			auth, _ := os.ReadFile(filepath.Join(exportDir, "auth.tf"))
			if !regexp.MustCompile(`expiration_time\s+= "8 minutes"`).Match(auth) {
				t.Errorf("expected the console-edited otp expiration in the export:\n%s", auth)
			}
		}
	},
}
