package stringattr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var TimeUnitValidator = stringvalidator.OneOf("seconds", "minutes", "hours", "days", "weeks")

// Non-Empty

func NonEmptyValidator() validator.String {
	return &nonEmptyValidator{}
}

type nonEmptyValidator struct {
}

func (v nonEmptyValidator) Description(_ context.Context) string {
	return "string must not be empty"
}

func (v nonEmptyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v nonEmptyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) == 0 {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Empty Attribute Value", fmt.Sprintf("Attribute %s must not be empty", req.Path)))
	}
}
