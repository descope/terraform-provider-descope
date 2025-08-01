package objattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewValidator[T any, M validatableModel[T]](description string) validator.Object {
	return &objectValidator[T, M]{description: description}
}

// Model

type validatableModel[T any] interface {
	helpers.Model[T]
	Validate(*helpers.Handler)
}

// Implementation

type objectValidator[T any, M validatableModel[T]] struct {
	description string
}

func (v *objectValidator[T, M]) Description(_ context.Context) string {
	return v.description
}

func (v *objectValidator[T, M]) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *objectValidator[T, M]) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	tflog.Debug(ctx, "Validating object", map[string]any{"path": req.Path.String()})
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	model := modelFromObject[T, M](ctx, req.ConfigValue, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	model.Validate(handler)
}
