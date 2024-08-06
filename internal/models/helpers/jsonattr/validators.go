package jsonattr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewValidator(description string, validate func(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse)) validator.Dynamic {
	return &dynamicValidator{description: description, validate: validate}
}

type dynamicValidator struct {
	description string
	validate    func(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse)
}

func (v *dynamicValidator) Description(_ context.Context) string {
	return v.description
}

func (v *dynamicValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *dynamicValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	v.validate(ctx, req, resp)
}

// Type

var typeValidator = NewValidator("must be a valid JSON object", validateObject)

func validateObject(_ context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	d := req.ConfigValue
	if d.IsNull() || d.IsUnknown() || d.IsUnderlyingValueNull() || d.IsUnderlyingValueUnknown() {
		return
	}

	u := d.UnderlyingValue()
	if object, ok := u.(types.Object); ok {
		ensureValue(object, resp)
	} else if _, ok := u.(types.String); ok {
		resp.Diagnostics.AddError("Invalid JSON value", "Use jsondecode() to convert file contents or string literal to a JSON object")
	} else {
		resp.Diagnostics.AddError("Invalid JSON value", fmt.Sprintf("The property must be a JSON object, found %T", u))
	}
}

func ensureValue(value attr.Value, resp *validator.DynamicResponse) {
	switch v := value.(type) {
	case types.String, types.Bool, types.Number:
		return
	case types.Object:
		for _, e := range v.Attributes() {
			ensureValue(e, resp)
		}
	case types.Tuple:
		for _, e := range v.Elements() {
			ensureValue(e, resp)
		}
	default:
		resp.Diagnostics.AddError("Invalid JSON value", fmt.Sprintf("Unexpected value of type %T", v))
	}
}
