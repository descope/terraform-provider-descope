package jwttemplate_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestJWTTemplate(t *testing.T) {
	projectID := testacc.ProjectID(t)
	j := testacc.JWTTemplate(t)
	project := `project_id = "` + projectID + `"`

	var templateID string
	captureID := func(s string) error {
		if !strings.HasPrefix(s, "JT") {
			return fmt.Errorf("expected id to have JT prefix, got %s", s)
		}
		templateID = s
		return nil
	}
	sameID := func(s string) error {
		if s != templateID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", templateID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: j.Config(project, `type = "user"`, `template = jsonencode({ custom = "value" })`),
		Check: j.Check(map[string]any{
			"id":                       captureID,
			"project_id":               projectID,
			"name":                     j.Name,
			"description":              "",
			"type":                     "user",
			"issuer_type":              "legacy",
			"auth_schema":              "default",
			"empty_claim_policy":       "none",
			"auto_tenant_claim":        false,
			"conformance_issuer":       false,
			"enforce_issuer":           false,
			"exclude_permission_claim": false,
			"override_subject_claim":   false,
			"add_jti_claim":            false,
			"template":                 `{"custom":"value"}`,
		}),
	}

	j.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: j.Config(project, `type = "user"`, `description = "updated"`, `auth_schema = "tenantOnly"`,
			`empty_claim_policy = "delete"`, `add_jti_claim = true`, `template = jsonencode({ custom = "other", extra = 1 })`),
		Check: j.Check(map[string]any{
			"id":                 sameID,
			"name":               j.Name,
			"description":        "updated",
			"auth_schema":        "tenantOnly",
			"empty_claim_policy": "delete",
			"add_jti_claim":      true,
			"template":           `{"custom":"other","extra":1}`,
		}),
	}

	flipTypeStep := resource.TestStep{
		Config: j.Config(project, `type = "key"`, `description = "updated"`, `auth_schema = "tenantOnly"`,
			`empty_claim_policy = "delete"`, `add_jti_claim = true`, `template = jsonencode({ custom = "other", extra = 1 })`),
		Check: j.Check(map[string]any{
			"id":   sameID,
			"type": "key",
		}),
	}

	importStep := resource.TestStep{
		ResourceName:      j.Path(),
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: testacc.GenerateImportStateID(j.Path(), "project_id", "id"),
	}

	testacc.RunWithDestroyCheck(t, "descope_jwt_template", createStep, updateStep, flipTypeStep, importStep)
}

func TestJWTTemplateValidation(t *testing.T) {
	projectID := testacc.ProjectID(t)
	j := testacc.JWTTemplate(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		resource.TestStep{
			Config:      j.Config(project, `type = "user"`, `template = jsonencode([1, 2])`),
			ExpectError: regexp.MustCompile(`must be valid JSON`),
		},
		resource.TestStep{
			Config:      j.Config(project, `type = "user"`, `template = "null"`),
			ExpectError: regexp.MustCompile(`must be a JSON object`),
		},
		resource.TestStep{
			Config:      j.Config(project, `type = "user"`, `template = jsonencode({})`, `issuer_type = "federated"`, `conformance_issuer = true`),
			ExpectError: regexp.MustCompile(`cannot be enabled when 'issuer_type'`),
		},
	)
}

// The server returns the template compact, and overwriting the indented state with that would fail the apply as an inconsistent result.
func TestJWTTemplateIndentedConfig(t *testing.T) {
	projectID := testacc.ProjectID(t)
	j := testacc.JWTTemplate(t)
	project := `project_id = "` + projectID + `"`
	indented := "{\n  \"custom\": \"value\"\n}"

	testacc.Run(t,
		resource.TestStep{
			Config: j.Config(project, `type = "user"`, `template = "{\n  \"custom\": \"value\"\n}"`),
			Check:  j.Check(map[string]any{"template": indented}),
		},
		resource.TestStep{
			Config: j.Config(project, `type = "user"`, `template = jsonencode({ custom = "value" })`),
			Check:  j.Check(map[string]any{"template": indented}),
		},
	)
}
