package resources_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPolicyRule(t *testing.T) {
	project := testacc.Project(t)
	rule := &testacc.Resource{Type: "policy_rule", ID: "test", Name: testacc.GenerateAlias(t)}

	testacc.Run(
		t,
		resource.TestStep{
			Config: project.Config() + rule.Config(`
				project_id = `+project.Path()+`.id
				enabled = true
				rule_family = "resource_access"
				action_kinds = ["client_access"]
				effect = "permit"
				principal_type = "client"
				resource_targets = [{
					type = "api"
					ids = ["api-1"]
				}]
			`),
			Check: rule.Check(map[string]any{
				"id":      testacc.AttributeIsSet,
				"version": 1,
				"name":    rule.Name,
				"enabled": true,
			}),
		},
		resource.TestStep{
			Config: project.Config() + rule.Config(`
				project_id = `+project.Path()+`.id
				description = "updated"
				enabled = false
				rule_family = "resource_access"
				action_kinds = ["client_access"]
				effect = "permit"
				principal_type = "client"
				resource_targets = [{
					type = "api"
					ids = ["api-1"]
				}]
			`),
			Check: rule.Check(map[string]any{
				"description": "updated",
				"enabled":     false,
				"version":     2,
			}),
		},
		resource.TestStep{
			ResourceName:      rule.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(rule.Path(), "project_id", "id"),
			ImportStateVerify: true,
		},
	)
}
