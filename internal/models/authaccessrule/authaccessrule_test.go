package authaccessrule_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthAccessRule(t *testing.T) {
	projectID := testacc.ProjectID(t)
	r := testacc.AuthAccessRule(t)
	project := `project_id = "` + projectID + `"`

	var ruleID string
	captureID := func(s string) error {
		ruleID = s
		return nil
	}
	sameID := func(s string) error {
		if s != ruleID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", ruleID, s)
		}
		return nil
	}

	testacc.RunWithDestroyCheck(t, "descope_auth_access_rule",
		// create with defaults
		resource.TestStep{
			Config: r.Config(project, `effect = "allow"`),
			Check: r.Check(map[string]any{
				"id":           captureID,
				"project_id":   projectID,
				"name":         r.Name,
				"description":  "",
				"enabled":      true,
				"effect":       "allow",
				"conditions.#": 0,
				"expire_time":  0,
			}),
		},
		// update every attribute in place
		resource.TestStep{
			Config: r.Config(project,
				`description = "block risky sign ins"`,
				`enabled = false`,
				`effect = "block"`,
				`expire_time = 2000000000`,
				`conditions = [
					{ key = "request.ip", operator = "inCidr", values = ["10.0.0.0/8"] },
					{ key = "request.country", operator = "notIn", values = ["US", "IL"] },
					{ key = "user.email", operator = "equal", values = ["a@example.com"] },
					{ key = "request.country", operator = "notEqual", values = ["{{jwtClaims.loginCountry}}"] },
				]`,
			),
			Check: r.Check(map[string]any{
				"id":                    sameID,
				"description":           "block risky sign ins",
				"enabled":               false,
				"effect":                "block",
				"expire_time":           2000000000,
				"conditions.#":          4,
				"conditions.0.key":      "request.ip",
				"conditions.0.operator": "inCidr",
				"conditions.0.values.0": "10.0.0.0/8",
				"conditions.1.values.#": 2,
				"conditions.1.values.0": "US",
				"conditions.1.values.1": "IL",
				"conditions.2.operator": "equal",
				"conditions.2.values.0": "a@example.com",
				"conditions.3.operator": "notEqual",
				"conditions.3.values.0": "{{jwtClaims.loginCountry}}",
			}),
		},
		// clear the conditions
		resource.TestStep{
			Config: r.Config(project, `effect = "block"`, `conditions = []`),
			Check: r.Check(map[string]any{
				"id":           sameID,
				"enabled":      true,
				"conditions.#": 0,
			}),
		},
		// import
		resource.TestStep{
			ResourceName:      r.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(r.Path(), "project_id", "id"),
		},
	)
}

func TestAuthAccessRuleInvalid(t *testing.T) {
	projectID := testacc.ProjectID(t)
	r := testacc.AuthAccessRule(t)
	project := `project_id = "` + projectID + `"`

	testacc.Run(t,
		// operators that take one operand reject several values
		resource.TestStep{
			Config:      r.Config(project, `effect = "block"`, `conditions = [{ key = "request.ip", operator = "equal", values = ["1.1.1.1", "2.2.2.2"] }]`),
			ExpectError: regexp.MustCompile(`requires exactly one value`),
		},
		// CIDR operators only apply to the request IP
		resource.TestStep{
			Config:      r.Config(project, `effect = "block"`, `conditions = [{ key = "user.email", operator = "inCidr", values = ["10.0.0.0/8"] }]`),
			ExpectError: regexp.MustCompile(`only supported with the request.ip key`),
		},
		// expiry is limited to epoch seconds that fit in 32 bits
		resource.TestStep{
			Config:      r.Config(project, `effect = "block"`, `expire_time = 4102444800`),
			ExpectError: regexp.MustCompile(`must be between 0 and 2147483647`),
		},
		// values must already be in canonical form
		resource.TestStep{
			Config:      r.Config(project, `effect = "block"`, `conditions = [{ key = "user.email", operator = "in", values = ["Mixed@Example.com"] }]`),
			ExpectError: regexp.MustCompile(`not in canonical form`),
		},
	)
}
