package connectors_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Stored secrets survive an import: omitting them plans a visible clear, supplying them re-sends them and settles.
func TestConnectorSecretsImport(t *testing.T) {
	projectID := testacc.ProjectID(t)
	c := testacc.NewResource(t, "http_connector")
	adopted := &testacc.Resource{Type: c.Type, ID: "adopted", Name: c.Name}
	base := []string{`project_id = "` + projectID + `"`, `base_url = "https://example.com"`}
	secrets := []string{
		`authentication = { bearer_token = "bearer-value" }`,
		`hmac_secret = "hmac-value"`,
		`secret_headers = { "X-Api-Key" = "header-value" }`,
	}
	withSecrets := c.Config(append(base, secrets...)...)
	withoutSecrets := c.Config(base...)

	ids := map[string]string{}
	adoptedConfig := `
		variable "connector_id" {}

		removed {
			from = ` + c.Path() + `
			lifecycle {
				destroy = false
			}
		}

		import {
			to = ` + adopted.Path() + `
			id = "` + projectID + `/${var.connector_id}"
		}
	` + adopted.Config(append(base, secrets...)...)
	adoptedVariables := config.Variables{"connector_id": lazyVariable(func() string { return ids["connector"] })}
	adoptedChecks := adopted.Check(map[string]any{
		"authentication.bearer_token": "bearer-value",
		"hmac_secret":                 "hmac-value",
		"secret_headers.X-Api-Key":    "header-value",
	})

	testacc.Run(t,
		resource.TestStep{
			Config: withSecrets,
			Check: c.Check(map[string]any{
				"id": func(s string) error {
					ids["connector"] = s
					return nil
				},
			}),
		},
		resource.TestStep{
			Config:             withoutSecrets,
			ResourceName:       c.Path(),
			ImportState:        true,
			ImportStateKind:    resource.ImportBlockWithID,
			ImportStateIdFunc:  testacc.GenerateImportStateID(c.Path(), "project_id", "id"),
			ExpectNonEmptyPlan: true,
			ImportPlanChecks: resource.ImportPlanChecks{
				PreApply: []plancheck.PlanCheck{expectImportedSecrets(c.Path(), nil)},
			},
		},
		resource.TestStep{
			Config:             withSecrets,
			ResourceName:       c.Path(),
			ImportState:        true,
			ImportStateKind:    resource.ImportBlockWithID,
			ImportStateIdFunc:  testacc.GenerateImportStateID(c.Path(), "project_id", "id"),
			ExpectNonEmptyPlan: true,
			ImportPlanChecks: resource.ImportPlanChecks{
				PreApply: []plancheck.PlanCheck{expectImportedSecrets(c.Path(), map[string]string{
					"bearer_token": "bearer-value",
					"hmac_secret":  "hmac-value",
					"X-Api-Key":    "header-value",
				})},
			},
		},
		resource.TestStep{
			Config:          adoptedConfig,
			ConfigVariables: adoptedVariables,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply:             []plancheck.PlanCheck{plancheck.ExpectResourceAction(adopted.Path(), plancheck.ResourceActionUpdate)},
				PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
			Check: adoptedChecks,
		},
		resource.TestStep{
			Config:          adoptedConfig,
			ConfigVariables: adoptedVariables,
			Check:           adoptedChecks,
		},
	)
}

type importedSecretsCheck struct {
	address  string
	expected map[string]string
}

func expectImportedSecrets(address string, expected map[string]string) plancheck.PlanCheck {
	return importedSecretsCheck{address: address, expected: expected}
}

func (c importedSecretsCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	var before, after map[string]string
	for _, rc := range req.Plan.ResourceChanges {
		if rc.Address == c.address && rc.Change != nil {
			before, after = secretValues(rc.Change.Before), secretValues(rc.Change.After)
		}
	}
	if before == nil {
		resp.Error = fmt.Errorf("no planned change for %s", c.address)
		return
	}
	for _, key := range []string{"bearer_token", "hmac_secret", "X-Api-Key"} {
		if before[key] == "" {
			resp.Error = fmt.Errorf("expected the import to record the stored %s, got an empty value", key)
			return
		}
		if after[key] != c.expected[key] {
			resp.Error = fmt.Errorf("expected %s to be planned as %q, got %q", key, c.expected[key], after[key])
			return
		}
	}
}

func secretValues(object any) map[string]string {
	var values struct {
		Authentication struct {
			BearerToken string `json:"bearer_token"`
		} `json:"authentication"`
		HMACSecret    string            `json:"hmac_secret"`
		SecretHeaders map[string]string `json:"secret_headers"`
	}
	b, _ := json.Marshal(object)
	_ = json.Unmarshal(b, &values)
	return map[string]string{
		"bearer_token": values.Authentication.BearerToken,
		"hmac_secret":  values.HMACSecret,
		"X-Api-Key":    values.SecretHeaders["X-Api-Key"],
	}
}

type lazyVariable func() string

func (v lazyVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v())
}
