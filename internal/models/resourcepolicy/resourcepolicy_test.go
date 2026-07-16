package resourcepolicy_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/resourcepolicy"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestModel_SetPolicy_roundtrips_policy_values(t *testing.T) {
	t.Parallel()

	// Given
	want := infra.ResourcePolicy{
		ApplicationID:      "AP1",
		ResourceID:         "RS1",
		UserAccessScopes:   []string{"read", "write"},
		ClientAccessScopes: []string{"admin"},
		AllUserScopes:      true,
	}
	var diagnostics diag.Diagnostics
	handler := helpers.NewHandler(context.Background(), &diagnostics)
	model := &resourcepolicy.Model{}

	// When
	model.SetPolicy(&want)
	got := model.Policy(handler)

	// Then
	require.False(t, diagnostics.HasError())
	require.Equal(t, want, got)
}

func TestSchema_requires_stable_reference_fields(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"project_id", "application_id", "resource_id"} {
		attribute, ok := resourcepolicy.Schema.Attributes[name].(schema.StringAttribute)
		require.True(t, ok, name)
		require.True(t, attribute.Required, name)
		require.NotEmpty(t, attribute.Validators, name)
		require.NotEmpty(t, attribute.PlanModifiers, name)
	}
}

func TestAccResourcePolicy(t *testing.T) {
	projectID := os.Getenv("DESCOPE_TESTACC_RESOURCE_POLICY_PROJECT_ID")
	applicationID := os.Getenv("DESCOPE_TESTACC_RESOURCE_POLICY_APPLICATION_ID")
	resourceID := os.Getenv("DESCOPE_TESTACC_RESOURCE_POLICY_RESOURCE_ID")
	acceptanceEnabled, _ := strconv.ParseBool(os.Getenv("TF_ACC"))
	if acceptanceEnabled && (projectID == "" || applicationID == "" || resourceID == "") {
		t.Fatal("DESCOPE_TESTACC_RESOURCE_POLICY_PROJECT_ID, DESCOPE_TESTACC_RESOURCE_POLICY_APPLICATION_ID, and DESCOPE_TESTACC_RESOURCE_POLICY_RESOURCE_ID must be set")
	}

	resourcePath := "descope_resource_policy.test"
	config := func(allClientScopes bool) string {
		return fmt.Sprintf(`
resource "descope_resource_policy" "test" {
  project_id          = %q
  application_id      = %q
  resource_id         = %q
  user_access_scopes  = ["read"]
  client_access_scopes = ["write"]
  all_client_scopes   = %t
}
`, projectID, applicationID, resourceID, allClientScopes)
	}

	testacc.Run(
		t,
		resource.TestStep{
			Config: config(false),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(resourcePath, "id", applicationID+"/"+resourceID),
				resource.TestCheckResourceAttr(resourcePath, "user_access_scopes.#", "1"),
				resource.TestCheckResourceAttr(resourcePath, "all_client_scopes", "false"),
			),
		},
		resource.TestStep{
			Config: config(true),
			Check:  resource.TestCheckResourceAttr(resourcePath, "all_client_scopes", "true"),
		},
		resource.TestStep{
			ResourceName:      resourcePath,
			ImportState:       true,
			ImportStateId:     projectID + "/" + applicationID + "/" + resourceID,
			ImportStateVerify: true,
		},
	)
}
