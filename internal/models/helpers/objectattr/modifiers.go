package objectattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func NewModifier[T any, M modifiableModel[T]](description string) planmodifier.Object {
	return &objectModifier[T, M]{description: description}
}

// Model

type modifiableModel[T any] interface {
	helpers.Model[T]
	Modify(h *helpers.Handler, state *T, config *T)
}

// Implementation

type objectModifier[T any, M modifiableModel[T]] struct {
	description string
}

func (v *objectModifier[T, M]) Description(_ context.Context) string {
	return v.description
}

func (v *objectModifier[T, M]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *objectModifier[T, M]) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() || req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	plan := helpers.ModelFromObject[T, M](ctx, req.PlanValue, &resp.Diagnostics)
	state := helpers.ModelFromObject[T, M](ctx, req.StateValue, &resp.Diagnostics)
	config := helpers.ModelFromObject[T, M](ctx, req.ConfigValue, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics, helpers.ReferencesMap{})
	plan.Modify(handler, state, config)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.PlanValue = helpers.ObjectFromModel[T, M](ctx, plan, req.PlanValue.AttributeTypes(ctx), &resp.Diagnostics)
}
