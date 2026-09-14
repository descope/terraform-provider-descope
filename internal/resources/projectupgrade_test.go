package resources

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectStateUpgrade(t *testing.T) {
	protected := true
	testCases := []struct {
		name           string
		rawState       string
		diagnostic     string
		detailContains []string
		expected       *projectUpgradeExpectedState
	}{
		{
			name:     "clean legacy state",
			rawState: `{"id":"project-id","name":"Example","environment":"production","tags":["one","two"],"deletion_protection":true}`,
			expected: &projectUpgradeExpectedState{
				id:                 "project-id",
				name:               "Example",
				environment:        "production",
				tags:               []string{"one", "two"},
				deletionProtection: &protected,
			},
		},
		{
			name:           "populated connectors",
			rawState:       `{"id":"project-id","name":"Example","connectors":{"http":{}}}`,
			diagnostic:     "Descope Project State Requires Migration",
			detailContains: []string{"connectors", "tfmigrate"},
		},
		{
			name:           "several populated ownership attributes",
			rawState:       `{"id":"project-id","name":"Example","authorization":{"roles":[]},"flows":[{}],"connectors":{"http":{}}}`,
			diagnostic:     "Descope Project State Requires Migration",
			detailContains: []string{"authorization, connectors, flows"},
		},
		{
			name:     "empty ownership attributes",
			rawState: `{"id":"project-id","name":"Example","environment":null,"tags":null,"connectors":null,"flows":{},"lists":[]}`,
			expected: &projectUpgradeExpectedState{
				id:   "project-id",
				name: "Example",
			},
		},
		{
			name:           "unknown populated attribute",
			rawState:       `{"id":"project-id","name":"Example","future_ownership":{"id":"owned"}}`,
			diagnostic:     "Descope Project State Requires Migration",
			detailContains: []string{"future_ownership"},
		},
		{
			name:       "malformed JSON",
			rawState:   `{"id":`,
			diagnostic: "Invalid Descope Project State",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			req := resource.UpgradeStateRequest{RawState: &tfprotov6.RawState{JSON: []byte(tc.rawState)}}
			resp := &resource.UpgradeStateResponse{State: tfsdk.State{
				Schema: project.Schema,
				Raw:    tftypes.NewValue(project.Schema.Type().TerraformType(ctx), nil),
			}}

			resourceUnderTest, ok := NewProjectResource().(*projectResource)
			require.True(t, ok)
			upgrader, ok := resourceUnderTest.UpgradeState(ctx)[0]
			require.True(t, ok)
			require.NotNil(t, upgrader.StateUpgrader)
			upgrader.StateUpgrader(ctx, req, resp)

			if tc.diagnostic != "" {
				require.Equal(t, 1, resp.Diagnostics.ErrorsCount())
				require.Len(t, resp.Diagnostics, 1)
				assert.Equal(t, tc.diagnostic, resp.Diagnostics[0].Summary())
				for _, text := range tc.detailContains {
					assert.Contains(t, resp.Diagnostics[0].Detail(), text)
				}
				assert.True(t, resp.State.Raw.IsNull())
				return
			}

			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			assertProjectUpgradeState(ctx, t, resp.State, *tc.expected)
		})
	}
}

type projectUpgradeExpectedState struct {
	id                 string
	name               string
	environment        string
	tags               []string
	deletionProtection *bool
}

func assertProjectUpgradeState(ctx context.Context, t *testing.T, state tfsdk.State, expected projectUpgradeExpectedState) {
	t.Helper()

	var actual project.ProjectModel
	diagnostics := state.Get(ctx, &actual)
	require.False(t, diagnostics.HasError(), diagnostics)
	assert.Equal(t, expected.id, actual.ID.ValueString())
	assert.Equal(t, expected.name, actual.Name.ValueString())
	if expected.environment == "" {
		assert.True(t, actual.Environment.IsNull())
	} else {
		assert.Equal(t, expected.environment, actual.Environment.ValueString())
	}
	if expected.deletionProtection == nil {
		assert.True(t, actual.DeletionProtection.IsNull())
	} else {
		assert.Equal(t, *expected.deletionProtection, actual.DeletionProtection.ValueBool())
	}
	if expected.tags == nil {
		assert.True(t, actual.Tags.IsNull())
	} else {
		var tags []string
		diagnostics = actual.Tags.ElementsAs(ctx, &tags, false)
		require.False(t, diagnostics.HasError(), diagnostics)
		assert.ElementsMatch(t, expected.tags, tags)
	}
}
