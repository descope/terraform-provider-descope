package objattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type ModifierOptions int

const (
	ModifierAllowNullState ModifierOptions = iota
)

func NewModifier[T any, M modifiableModel[T]](description string, options ...ModifierOptions) planmodifier.Object {
	modifier := &objectModifier[T, M]{description: description}
	for i := range options {
		if options[i] == ModifierAllowNullState {
			modifier.allowNullState = true
		}
	}
	return modifier
}

// Model

type modifiableModel[T any] interface {
	helpers.Model[T]
	Modify(h *helpers.Handler, state *T)
}

// Implementation

type objectModifier[T any, M modifiableModel[T]] struct {
	description    string
	allowNullState bool
}

func (v *objectModifier[T, M]) Description(_ context.Context) string {
	return v.description
}

func (v *objectModifier[T, M]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *objectModifier[T, M]) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.PlanValue.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	plan := modelFromObject[T, M](ctx, req.PlanValue, &resp.Diagnostics)

	var state M
	if !req.StateValue.IsNull() && !req.State.Raw.IsNull() {
		state = modelFromObject[T, M](ctx, req.StateValue, &resp.Diagnostics)
	} else if !v.allowNullState {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	plan.Modify(handler, state)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = valueOf(ctx, plan).ObjectValue
}
