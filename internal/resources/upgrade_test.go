package resources_test

import (
	"context"
	"os"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestProjectUpgradeStateFromV03(t *testing.T) {
	raw, err := os.ReadFile("testdata/project-v0.3.16.json")
	require.NoError(t, err)

	state, diags := upgradeProjectState(t, raw)

	require.Len(t, diags, 1)
	require.Equal(t, tfprotov6.DiagnosticSeverityWarning, diags[0].Severity)
	require.Equal(t, "Project Configuration No Longer Managed", diags[0].Summary)
	require.Contains(t, diags[0].Detail, "guides/upgrading-from-v0.3")
	require.Equal(t, "P3JkQChdCY9nB5McOvoW72dH7i9Z", stringValue(t, state["id"]))
	require.Equal(t, "testacc-local-upgradefixture", stringValue(t, state["name"]))
	require.Equal(t, "", stringValue(t, state["environment"]))
	require.True(t, state["deletion_protection"].IsNull())
	var tags []tftypes.Value
	require.NoError(t, state["tags"].As(&tags))
	require.Len(t, tags, 1)
	require.Equal(t, "fixture", stringValue(t, tags[0]))
}

func TestProjectUpgradeStateFromPrerelease(t *testing.T) {
	raw := []byte(`{"id":"P1","name":"prerelease","environment":"production","deletion_protection":false,"tags":["a"]}`)

	state, diags := upgradeProjectState(t, raw)

	require.Empty(t, diags)
	require.Equal(t, "P1", stringValue(t, state["id"]))
	require.Equal(t, "prerelease", stringValue(t, state["name"]))
	require.Equal(t, "production", stringValue(t, state["environment"]))
	var protection bool
	require.NoError(t, state["deletion_protection"].As(&protection))
	require.False(t, protection)
}

func upgradeProjectState(t *testing.T, raw []byte) (map[string]tftypes.Value, []*tfprotov6.Diagnostic) {
	ctx := context.Background()
	server, err := providerserver.NewProtocol6WithError(provider.NewDescopeProvider("test")())()
	require.NoError(t, err)

	resp, err := server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
		TypeName: "descope_project",
		Version:  0,
		RawState: &tfprotov6.RawState{JSON: raw},
	})
	require.NoError(t, err)
	for _, d := range resp.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, d.Detail)
	}

	value, err := resp.UpgradedState.Unmarshal(project.Schema.Type().TerraformType(ctx))
	require.NoError(t, err)
	state := map[string]tftypes.Value{}
	require.NoError(t, value.As(&state))
	return state, resp.Diagnostics
}

func stringValue(t *testing.T, value tftypes.Value) string {
	var s string
	require.NoError(t, value.As(&s))
	return s
}
