package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TrueWhenSiblingTrue plans a boolean attribute as true whenever a sibling boolean
// attribute (by name, within the same parent object) is set to true in the
// configuration.
//
// Use this when the attribute's effective value is the OR of itself and that sibling —
// for example a deprecated flag that folds into it and disables it. The backend returns
// the folded (effective) value on read, so when the sibling is true this attribute is
// deterministically true. Planning it as known true keeps the plan both consistent with
// the apply result and idempotent; without it, a state-based modifier such as
// UseValidStateForUnknown would plan the stale prior value, yielding a "provider produced
// inconsistent result" error.
//
// When the sibling is unknown the folded result cannot be determined, so the attribute is
// planned as unknown. When the sibling is null or false the plan value is left untouched
// so any earlier modifier (e.g. UseValidStateForUnknown) still applies.
func TrueWhenSiblingTrue(siblingName string) trueWhenSiblingTrueModifier {
	return trueWhenSiblingTrueModifier{siblingName: siblingName}
}

type trueWhenSiblingTrueModifier struct {
	siblingName string
}

func (m trueWhenSiblingTrueModifier) Description(_ context.Context) string {
	return "the value is true when " + m.siblingName + " is set to true"
}

func (m trueWhenSiblingTrueModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m trueWhenSiblingTrueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	var sibling types.Bool
	siblingPath := req.Path.ParentPath().AtName(m.siblingName)
	if diags := req.Config.GetAttribute(ctx, siblingPath, &sibling); diags.HasError() {
		return
	}

	switch {
	case sibling.IsUnknown():
		// Cannot determine the folded result yet, so let apply resolve it.
		resp.PlanValue = types.BoolUnknown()
	case !sibling.IsNull() && sibling.ValueBool():
		// The sibling being true folds this attribute to true deterministically.
		resp.PlanValue = types.BoolValue(true)
	}
}
