package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseValidStateForUnknown behaves the same as the framework provided UseStateForUnknown
// used to work until version v1.15.1 of the framework, in that it uses the current state
// value if the plan value is Unknown, but unlike UseStateForUnknown it does so only if the
// state value is not Null. In other words, whereas UseStateForUnknown copies the state value
// into the plan value on every update, UseValidStateForUnknown only does so if the state
// value has some non-Null value.
//
// We use this modifier for Optional attributes instead of UseStateForUnknown so that attributes
// that plan writer did not set do not always appear in the diffs with a "known after apply" value.
// In return, we require the SetValues functions in every model to ensure all attributes have some
// value set and are not left as Unknown.
func UseValidStateForUnknown() useValidStateForUnknownModifier {
	return useValidStateForUnknownModifier{}
}

type useValidStateForUnknownModifier struct{}

func (m useValidStateForUnknownModifier) Description(_ context.Context) string {
	return "the value will be set to any existing state value"
}

func (m useValidStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useValidStateForUnknownModifier) ShouldModifyPlan(stateValue, planValue, configValue attr.Value) bool {
	// Do nothing if there is no valid state value (resource is being created or state value was null).
	if stateValue.IsNull() {
		return false
	}

	// Do nothing if there is a known planned value.
	if !planValue.IsUnknown() {
		return false
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if configValue.IsUnknown() {
		return false
	}

	return true
}

func (m useValidStateForUnknownModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyMap(_ context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}

func (m useValidStateForUnknownModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if m.ShouldModifyPlan(req.StateValue, req.PlanValue, req.ConfigValue) {
		resp.PlanValue = req.StateValue
	}
}
