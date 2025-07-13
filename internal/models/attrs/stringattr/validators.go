package stringattr

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var TimeUnitValidator = stringvalidator.OneOf("seconds", "minutes", "hours", "days", "weeks")

var StandardLenValidator = stringvalidator.LengthAtMost(254)

var FlowIDValidator = stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "must only contain alphanumeric, underscore or hyphen characters")

var OTPValidator = stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]{6}$`), "must be a 6 digit code")

var NonEmptyValidator validator.String = &nonEmptyValidator{}

func JSONValidator(required ...string) validator.String {
	return &jsonValidator{required: required}
}

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

// JSON

type jsonValidator struct {
	required []string
}

func (v jsonValidator) Description(_ context.Context) string {
	return "must be valid JSON and have all requirements satisfied"
}

func (v jsonValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) == 0 {
		return // we let the Required/Optional/Default attribute handle the empty value case
	}

	m := map[string]any{}
	if err := json.Unmarshal([]byte(value), &m); err != nil {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid Attribute Value", fmt.Sprintf("Attribute %s must be valid JSON", req.Path)))
		return
	}

	for _, field := range v.required {
		if _, ok := m[field]; !ok {
			resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid Attribute Value", fmt.Sprintf("Attribute %s is expected to be valid JSON with a '%s' field", req.Path, field)))
			return
		}
	}
}
