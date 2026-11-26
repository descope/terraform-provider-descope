package project_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"
)

func TestProjectUpgradeFromV03(t *testing.T) {
	cliConfig := filepath.Join(t.TempDir(), "terraformrc")
	if err := os.WriteFile(cliConfig, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	// dev_overrides in a developer's CLI configuration would replace the published v0.3.16 with the local build
	t.Setenv("TF_CLI_CONFIG_FILE", cliConfig)

	p := testacc.Project(t)
	ids := map[string]string{}
	capture := func(key string) func(string) error {
		return func(s string) error {
			ids[key] = s
			return nil
		}
	}
	sameProjectID := func(s string) error {
		if s != ids["project"] {
			return fmt.Errorf("expected the project id to be preserved as %s, got %s", ids["project"], s)
		}
		return nil
	}

	legacyConfig := `
		resource "descope_project" "test" {
			name = "` + p.Name + `"
			project_settings = {
				refresh_token_expiration = "3 weeks"
			}
			authorization = {
				roles = [{ name = "Tester", permissions = ["Tests"] }]
				permissions = [{ name = "Tests" }]
			}
			attributes = {
				user = [{ id = "tfaccupgradetier", name = "Tier", type = "string" }]
			}
			lists = [{ name = "Allowed", type = "texts", data = jsonencode(["alpha", "beta"]) }]
			connectors = {
				http = [{ name = "Webhook", base_url = "https://example.com" }]
			}
		}
	`

	trimmedConfig := `
		resource "descope_project" "test" {
			name = "` + p.Name + `"
		}
	`

	adoptedConfig := `
		variable "role_id" {}
		variable "permission_id" {}
		variable "list_id" {}
		variable "connector_id" {}

		resource "descope_project" "test" {
			name = "` + p.Name + `"
			deletion_protection = false
		}

		import {
			to = descope_session_settings.test
			id = descope_project.test.id
		}

		resource "descope_session_settings" "test" {
			project_id = descope_project.test.id
			refresh_token_expiration = "3 weeks"
		}

		import {
			to = descope_permission.test
			id = "${descope_project.test.id}/${var.permission_id}"
		}

		resource "descope_permission" "test" {
			project_id = descope_project.test.id
			name = "Tests"
		}

		import {
			to = descope_role.test
			id = "${descope_project.test.id}/${var.role_id}"
		}

		resource "descope_role" "test" {
			project_id = descope_project.test.id
			name = "Tester"
			permissions = [descope_permission.test.name]
		}

		import {
			to = descope_user_attribute.test
			id = "${descope_project.test.id}/tfaccupgradetier"
		}

		resource "descope_user_attribute" "test" {
			project_id = descope_project.test.id
			id = "tfaccupgradetier"
			name = "Tier"
			type = "string"
		}

		import {
			to = descope_list.test
			id = "${descope_project.test.id}/${var.list_id}"
		}

		resource "descope_list" "test" {
			project_id = descope_project.test.id
			name = "Allowed"
			texts = ["alpha", "beta"]
		}

		import {
			to = descope_http_connector.test
			id = "${descope_project.test.id}/${var.connector_id}"
		}

		resource "descope_http_connector" "test" {
			project_id = descope_project.test.id
			name = "Webhook"
			base_url = "https://example.com"
		}
	`

	adopted := []string{
		"descope_session_settings.test",
		"descope_permission.test",
		"descope_role.test",
		"descope_user_attribute.test",
		"descope_list.test",
	}
	importChecks := []plancheck.PlanCheck{}
	for _, addr := range adopted {
		importChecks = append(importChecks, plancheck.ExpectResourceAction(addr, plancheck.ResourceActionNoop))
	}
	importChecks = append(importChecks, plancheck.ExpectResourceAction("descope_http_connector.test", plancheck.ResourceActionUpdate))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testacc.PreCheck(t) },
		Steps: []resource.TestStep{
			// creates the project and its configuration with the published v0.3.16 provider
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"descope": {Source: "descope/descope", VersionConstraint: "0.3.16"},
				},
				Config: legacyConfig,
				Check: p.Check(map[string]any{
					"id":                   capture("project"),
					"lists.0.id":           capture("list"),
					"connectors.http.0.id": capture("connector"),
				}),
			},
			// upgrading and cutting the project down to its core attributes plans no changes
			{
				ProtoV6ProviderFactories: testacc.ProviderFactories,
				Config:                   trimmedConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: p.Check(map[string]any{
					"id":   sameProjectID,
					"name": p.Name,
				}),
			},
			// adopting the configuration with import blocks changes nothing but the http connector's new defaults
			{
				ProtoV6ProviderFactories: testacc.ProviderFactories,
				PreConfig:                func() { resolveAuthorizationIDs(t, ids) },
				Config:                   adoptedConfig,
				ConfigVariables: config.Variables{
					"role_id":       lazyVariable(func() string { return ids["role"] }),
					"permission_id": lazyVariable(func() string { return ids["permission"] }),
					"list_id":       lazyVariable(func() string { return ids["list"] }),
					"connector_id":  lazyVariable(func() string { return ids["connector"] }),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             importChecks,
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

type lazyVariable func() string

func (v lazyVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v())
}

func resolveAuthorizationIDs(t *testing.T, ids map[string]string) {
	ctx := context.Background()
	client := infra.NewClient("testacc", os.Getenv("DESCOPE_MANAGEMENT_KEY"), os.Getenv("DESCOPE_BASE_URL"))
	roles, err := client.PostData(ctx, ids["project"], "/v1/mgmt/role/search", map[string]any{})
	require.NoError(t, err)
	ids["role"] = findIDByName(t, roles["roles"], "Tester")
	permissions, err := client.Get(ctx, ids["project"], "/v1/mgmt/permission/all", nil)
	require.NoError(t, err)
	ids["permission"] = findIDByName(t, permissions["permissions"], "Tests")
}

func findIDByName(t *testing.T, entities any, name string) string {
	list, _ := entities.([]any)
	for _, entity := range list {
		if m, ok := entity.(map[string]any); ok && m["name"] == name {
			id, _ := m["id"].(string)
			return id
		}
	}
	require.Fail(t, "entity not found", "no entity named %s", name)
	return ""
}
