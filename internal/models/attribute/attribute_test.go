package attribute_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestUserAttribute(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.UserAttribute(t)
	testacc.Run(t,
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccusertier"
				name = "Tier"
				type = "string"
				widget_authorization = {
					view_permissions = ["Tenant Admin"]
				}
			`),
			Check: a.Check(map[string]any{
				"id":                                    "tfaccusertier",
				"project_id":                            projectID,
				"name":                                  "Tier",
				"type":                                  "string",
				"widget_authorization.view_permissions": []string{"Tenant Admin"},
			}),
		},
		// updating the display name is an in-place update (id and type are unchanged)
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccusertier"
				name = "Tier Updated"
				type = "string"
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(a.Path(), plancheck.ResourceActionUpdate)},
			},
			Check: a.Check(map[string]any{"name": "Tier Updated"}),
		},
		// changing the type forces a replacement (the backend rejects an in-place type change)
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccusertier"
				name = "Tier Updated"
				type = "number"
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(a.Path(), plancheck.ResourceActionDestroyBeforeCreate)},
			},
			Check: a.Check(map[string]any{"type": "number"}),
		},
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		},
	)
}

func TestUserAttributeSelectOptions(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.UserAttribute(t)
	testacc.Run(t,
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccuserdept"
				name = "Department"
				type = "singleselect"
				select_options = ["gold", "silver"]
			`),
			Check: a.Check(map[string]any{
				"type":           "singleselect",
				"select_options": []string{"gold", "silver"},
			}),
		},
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccuserdept"
				name = "Department"
				type = "singleselect"
				select_options = ["gold", "silver", "bronze"]
			`),
			Check: a.Check(map[string]any{
				"select_options": []string{"gold", "silver", "bronze"},
			}),
		},
	)
}

func TestTenantAttribute(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.TenantAttribute(t)
	testacc.Run(t,
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfacctenantreg"
				name = "Region"
				type = "string"
				authorization = {
					view_permissions = ["Tenant Admin"]
				}
			`),
			Check: a.Check(map[string]any{
				"id":                             "tfacctenantreg",
				"name":                           "Region",
				"type":                           "string",
				"authorization.view_permissions": []string{"Tenant Admin"},
			}),
		},
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		},
	)
}

func TestAccessKeyAttribute(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.AccessKeyAttribute(t)
	testacc.Run(t,
		resource.TestStep{
			Config: a.Block(`
				project_id = "` + projectID + `"
				id = "tfaccakenv"
				name = "Environment"
				type = "string"
			`),
			Check: a.Check(map[string]any{
				"id":   "tfaccakenv",
				"name": "Environment",
				"type": "string",
			}),
		},
	)
}
