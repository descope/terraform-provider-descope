package resources_test

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/flow"
	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestReplacingAFlowWarnsAboutOverwriting(t *testing.T) {
	ctx := context.Background()
	server, err := providerserver.NewProtocol6WithError(provider.NewDescopeProvider("test")())()
	require.NoError(t, err)

	flowType := flow.Schema.Type().TerraformType(ctx)
	value := func(id any, flowID string) *tfprotov6.DynamicValue {
		v, err := tfprotov6.NewDynamicValue(flowType, tftypes.NewValue(flowType, map[string]tftypes.Value{
			"id":         tftypes.NewValue(tftypes.String, id),
			"project_id": tftypes.NewValue(tftypes.String, "P1"),
			"flow_id":    tftypes.NewValue(tftypes.String, flowID),
			"data":       tftypes.NewValue(tftypes.String, `{"metadata":{},"contents":{}}`),
		}))
		require.NoError(t, err)
		return &v
	}

	resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "descope_flow",
		PriorState:       value("old-flow", "old-flow"),
		ProposedNewState: value("old-flow", "new-flow"),
		Config:           value(nil, "new-flow"),
	})
	require.NoError(t, err)

	summaries := []string{}
	for _, d := range resp.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, d.Detail)
		summaries = append(summaries, d.Summary)
	}
	require.Contains(t, summaries, "Existing Configuration Might Be Overwritten")
}
