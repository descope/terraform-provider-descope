package durationattr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MinimumValue(duration string) validator.String {
	return &durationValidator{minimum: duration}
}

func MaximumValue(duration string) validator.String {
	return &durationValidator{maximum: duration}
}

var formatValidator validator.String = &durationValidator{}

type durationValidator struct {
	minimum string
	maximum string
}

func (v durationValidator) Description(_ context.Context) string {
	return fmt.Sprintf("must be a number between 0 and 1000 followed by a space and one of the valid time units: %s", strings.Join(units, ", "))
}

func (v durationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v durationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	tflog.Trace(ctx, "Validating duration", map[string]any{"path": request.Path.String()})
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
	if v.maximum != "" {
		maxSeconds, ok := getSeconds(v.maximum)
		if !ok {
			response.Diagnostics.Append(validatordiag.BugInProviderDiagnostic("Invalid value for maximum duration: " + v.maximum))
		}
		if seconds > maxSeconds {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(request.Path, "must be at most "+v.maximum, request.ConfigValue.String()))
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
