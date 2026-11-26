package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type defaulterModel struct {
	value bool
}

func (m *defaulterModel) DeletionProtectionDefault(_ context.Context) bool {
	return m.value
}

type noopSource struct{}

func (s noopSource) Get(_ context.Context, _ any) diag.Diagnostics {
	return nil
}

func (s noopSource) GetAttribute(_ context.Context, _ path.Path, _ any) diag.Diagnostics {
	return nil
}

func TestResolveDeletionProtection(t *testing.T) {
	testCases := []struct {
		name       string
		flag       types.Bool
		model      any
		protected  bool
		viaDefault bool
	}{
		{"explicit true", types.BoolValue(true), &struct{}{}, true, false},
		{"explicit false", types.BoolValue(false), &defaulterModel{value: true}, false, false},
		{"explicit true ignores defaulter", types.BoolValue(true), &defaulterModel{value: false}, true, false},
		{"null without defaulter", types.BoolNull(), &struct{}{}, false, false},
		{"null with protected default", types.BoolNull(), &defaulterModel{value: true}, true, true},
		{"null with unprotected default", types.BoolNull(), &defaulterModel{value: false}, false, true},
		{"unknown with protected default", types.BoolUnknown(), &defaulterModel{value: true}, true, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var diags diag.Diagnostics
			protected, viaDefault := resolveDeletionProtection(t.Context(), tc.flag, noopSource{}, tc.model, &diags)
			require.False(t, diags.HasError())
			assert.Equal(t, tc.protected, protected)
			assert.Equal(t, tc.viaDefault, viaDefault)
		})
	}
}

type flagSource struct {
	flag types.Bool
}

func (s flagSource) Get(_ context.Context, _ any) diag.Diagnostics {
	return nil
}

func (s flagSource) GetAttribute(_ context.Context, _ path.Path, target any) diag.Diagnostics {
	if flag, ok := target.(*types.Bool); ok {
		*flag = s.flag
	}
	return nil
}

func TestCheckRemovalProtection(t *testing.T) {
	testCases := []struct {
		name   string
		flag   types.Bool
		model  any
		errors int
	}{
		{"explicit true", types.BoolValue(true), &struct{}{}, 1},
		{"explicit false", types.BoolValue(false), &defaulterModel{value: true}, 0},
		{"null without defaulter", types.BoolNull(), &struct{}{}, 0},
		{"null with protected default", types.BoolNull(), &defaulterModel{value: true}, 1},
		{"null with unprotected default", types.BoolNull(), &defaulterModel{value: false}, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var diags diag.Diagnostics
			checkRemovalProtection(t.Context(), flagSource{flag: tc.flag}, tc.model, "test", &diags)
			assert.Equal(t, tc.errors, diags.ErrorsCount())
		})
	}
}

var replaceTestSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{Required: true},
		"pid":  schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
		"flag": schema.BoolAttribute{Optional: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()}},
	},
}

func replaceTestValue(name, pid string, flag bool) tftypes.Value {
	objectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name": tftypes.String,
		"pid":  tftypes.String,
		"flag": tftypes.Bool,
	}}
	return tftypes.NewValue(objectType, map[string]tftypes.Value{
		"name": tftypes.NewValue(tftypes.String, name),
		"pid":  tftypes.NewValue(tftypes.String, pid),
		"flag": tftypes.NewValue(tftypes.Bool, flag),
	})
}

func TestIsPlannedReplace(t *testing.T) {
	testCases := []struct {
		name    string
		plan    tftypes.Value
		state   tftypes.Value
		replace bool
	}{
		{"no changes", replaceTestValue("a", "p1", false), replaceTestValue("a", "p1", false), false},
		{"non-trigger change", replaceTestValue("b", "p1", false), replaceTestValue("a", "p1", false), false},
		{"string trigger change", replaceTestValue("a", "p2", false), replaceTestValue("a", "p1", false), true},
		{"bool trigger change", replaceTestValue("a", "p1", true), replaceTestValue("a", "p1", false), true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.ModifyPlanRequest{
				Plan:   tfsdk.Plan{Raw: tc.plan, Schema: replaceTestSchema},
				State:  tfsdk.State{Raw: tc.state, Schema: replaceTestSchema},
				Config: tfsdk.Config{Raw: tc.plan, Schema: replaceTestSchema},
			}
			assert.Equal(t, tc.replace, isPlannedReplace(t.Context(), replaceTestSchema, req))
		})
	}
}

var immutableTestSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"name":       schema.StringAttribute{Required: true},
		"project_id": schema.StringAttribute{Required: true},
		"app_id":     schema.StringAttribute{Required: true},
	},
}

func immutableTestValue(name, projectID, appID any) tftypes.Value {
	objectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name":       tftypes.String,
		"project_id": tftypes.String,
		"app_id":     tftypes.String,
	}}
	return tftypes.NewValue(objectType, map[string]tftypes.Value{
		"name":       tftypes.NewValue(tftypes.String, name),
		"project_id": tftypes.NewValue(tftypes.String, projectID),
		"app_id":     tftypes.NewValue(tftypes.String, appID),
	})
}

func TestCheckImmutableAttributes(t *testing.T) {
	testCases := []struct {
		name   string
		plan   tftypes.Value
		state  tftypes.Value
		errors int
	}{
		{"no changes", immutableTestValue("a", "p1", "s1"), immutableTestValue("a", "p1", "s1"), 0},
		{"other change", immutableTestValue("b", "p1", "s1"), immutableTestValue("a", "p1", "s1"), 0},
		{"project_id change", immutableTestValue("a", "p2", "s1"), immutableTestValue("a", "p1", "s1"), 1},
		{"app_id change", immutableTestValue("a", "p1", "s2"), immutableTestValue("a", "p1", "s1"), 1},
		{"both change", immutableTestValue("a", "p2", "s2"), immutableTestValue("a", "p1", "s1"), 2},
		{"null in state", immutableTestValue("a", "p1", "s1"), immutableTestValue("a", nil, nil), 0},
		{"unknown in plan", immutableTestValue("a", tftypes.UnknownValue, "s1"), immutableTestValue("a", "p1", "s1"), 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.ModifyPlanRequest{
				Plan:  tfsdk.Plan{Raw: tc.plan, Schema: immutableTestSchema},
				State: tfsdk.State{Raw: tc.state, Schema: immutableTestSchema},
			}
			resp := resource.ModifyPlanResponse{}
			checkImmutableAttributes(t.Context(), immutableTestSchema, "test", req, &resp)
			assert.Equal(t, tc.errors, resp.Diagnostics.ErrorsCount())
		})
	}
}
