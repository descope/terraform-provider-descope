package stringattr

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var TimeUnitValidator = stringvalidator.OneOf("seconds", "minutes", "hours", "days", "weeks")

var StandardLenValidator = stringvalidator.LengthAtMost(254)

var FlowIDValidator = stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "must only contain alphanumeric, underscore or hyphen characters")

var NonEmptyValidator validator.String = &nonEmptyValidator{}

// Non-Empty

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
