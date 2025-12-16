package listattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/listtype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// Creates a new modifier that matches entities in the list by name (does not match by list order like ModifyMatchingNames)
func NewModifierMatchingNames[T any, M helpers.NamedModel[T]](description string) planmodifier.List {
	return &modifierMatchingNames[T, M]{description: description}
}

// Implementation

type modifierMatchingNames[T any, M helpers.NamedModel[T]] struct {
	description string
}

func (v *modifierMatchingNames[T, M]) Description(_ context.Context) string {
	return v.description
}

func (v *modifierMatchingNames[T, M]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *modifierMatchingNames[T, M]) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() || req.Plan.Raw.IsNull() || req.StateValue.IsNull() || req.State.Raw.IsNull() {
		return
	}

	plan, diags := listtype.NewValueWith[T](ctx, req.PlanValue.Elements())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diags := listtype.NewValueWith[T](ctx, req.StateValue.Elements())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	h := helpers.NewHandler(ctx, &resp.Diagnostics)
	for p := range MutatingIterator(&plan, h) {
		var planned M = p
		for e := range Iterator(state, h) {
			var existing M = e
			if planned.GetName().Equal(existing.GetName()) {
				h.Log("Setting ID '%s' for %T named '%s' by matching name", existing.GetID().ValueString(), *planned, planned.GetName().ValueString())
				planned.SetID(existing.GetID())
				break
			}
		}
		if planned.GetID().ValueString() == "" {
			h.Log("No existing ID found for %T named '%s'", *planned, planned.GetName().ValueString())
		}
	}

	resp.PlanValue = plan.ListValue
}
