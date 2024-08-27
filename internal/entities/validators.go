package entities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewValidator[T any, E validatableEntity[T]](constructor validatableEntityConstructor[T, E], description string) resource.ConfigValidator {
	return &entityValidator[T, E]{constructor: constructor, description: description}
}

// Entity

type validatableEntityConstructor[T any, E validatableEntity[T]] func(context.Context, entitySource, *diag.Diagnostics) E

type validatableEntity[T any] interface {
	Validate(context.Context)
	*T
}

// Implementation

type entityValidator[T any, E validatableEntity[T]] struct {
	constructor validatableEntityConstructor[T, E]
	description string
}

func (v *entityValidator[T, E]) Description(_ context.Context) string {
	return v.description
}

func (v *entityValidator[T, E]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *entityValidator[T, E]) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	e := v.constructor(ctx, req.Config, &resp.Diagnostics)
	e.Validate(ctx)
}
