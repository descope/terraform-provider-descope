package export

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/emit"
	"github.com/stretchr/testify/require"
)

func TestAdoptionWarningsFlagInboundApps(t *testing.T) {
	address, err := emit.ParseProjectAddress("descope_project.main")
	require.NoError(t, err)
	plan := &emit.Plan{
		ProjectAddress: address,
		Resources: []emit.Resource{
			{Type: "descope_inbound_app", Label: "partner", EntityID: "TPA123"},
			{Type: "descope_role", Label: "admin", EntityID: "ROL456"},
			{Type: "descope_inbound_app", Label: "mobile", EntityID: "TPA789"},
		},
	}

	warnings := adoptionWarnings(plan)

	require.Len(t, warnings, 2)
	require.Contains(t, warnings[0].Text, "descope_inbound_app.partner")
	require.Contains(t, warnings[0].Text, "already")
	require.Contains(t, warnings[1].Text, "descope_inbound_app.mobile")
	require.False(t, warnings[0].Lossy)
}

func TestAdoptionWarningsOnlyWhenAdopting(t *testing.T) {
	plan := &emit.Plan{Resources: []emit.Resource{{Type: "descope_inbound_app", Label: "partner", EntityID: "TPA123"}}}

	require.Empty(t, adoptionWarnings(plan))
}
