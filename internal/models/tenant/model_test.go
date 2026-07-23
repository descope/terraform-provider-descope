package tenant

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestModelSetTenantPreservesConfiguredParentWhenReadOmitsIt(t *testing.T) {
	model := Model{Parent: types.StringValue("parent-1")}

	model.SetTenant(t.Context(), &infra.Tenant{ID: "tenant-1", Name: "Tenant"})

	require.Equal(t, "parent-1", model.Parent.ValueString())
}
