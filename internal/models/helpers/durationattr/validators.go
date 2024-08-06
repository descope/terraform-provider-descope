package durationattr

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MinimumValue(duration string) validator.String {
	return &durationValidator{minimum: duration}
}

var formatValidator validator.String = &durationValidator{}

type durationValidator struct {
	minimum string
}

func (v durationValidator) Description(_ context.Context) string {
	return fmt.Sprintf("must be a number between 0 and 1000 followed by a space and one of the valid time units: %s", units)
}

func (v durationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v durationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	seconds, ok := getSeconds(request.ConfigValue.ValueString())
	if !ok {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(request.Path, v.Description(ctx), request.ConfigValue.String()))
		return
	}

	if v.minimum != "" {
		minSeconds, ok := getSeconds(v.minimum)
		if !ok {
			response.Diagnostics.Append(validatordiag.BugInProviderDiagnostic("Invalid value for minimum duration: " + v.minimum))
		}
		if seconds < minSeconds {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(request.Path, "must be at least "+v.minimum, request.ConfigValue.String()))
		}
	}
}

func getSeconds(s string) (seconds int64, ok bool) {
	seconds, unit, ok := parseString(s)
	if !ok {
		return // failed
	}
	if unit == "seconds" {
		return // no need to change
	}
	if seconds *= 60; unit == "minutes" {
		return
	}
	if seconds *= 60; unit == "hours" {
		return
	}
	if seconds *= 24; unit == "days" {
		return
	}
	seconds *= 7 // weeks
	return
}
