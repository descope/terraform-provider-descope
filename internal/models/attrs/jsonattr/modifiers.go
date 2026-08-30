package jsonattr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateWhenEquivalent keeps the state value in the plan when the configured document differs from it only in formatting. The framework
// applies the semantic equality of custom types when reading and after applying but never while planning, so without this a diff shows.
func UseStateWhenEquivalent() planmodifier.String {
	return useStateWhenEquivalentModifier{}
}

type useStateWhenEquivalentModifier struct{}

func (m useStateWhenEquivalentModifier) Description(_ context.Context) string {
	return "the existing value is kept when the configured JSON only differs in formatting"
}

func (m useStateWhenEquivalentModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useStateWhenEquivalentModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if equivalent(req.StateValue.ValueString(), req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}
