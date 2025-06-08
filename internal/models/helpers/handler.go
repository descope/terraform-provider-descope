package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// A wrapper struct for the context, diags and references to pass around to model function calls.
type Handler struct {
	Ctx         context.Context
	Diagnostics *diag.Diagnostics
	Refs        ReferencesMap
}

// NewHandler creates a new Handler instance that wraps a context, an empty references map, and a
// reference to the diagnostics collection of a response of an entity object.
func NewHandler(ctx context.Context, diags *diag.Diagnostics) *Handler {
	return &Handler{
		Ctx:         ctx,
		Diagnostics: diags,
		Refs:        ReferencesMap{},
	}
}

func (h *Handler) Log(format string, a ...any) {
	tflog.Info(h.Ctx, fmt.Sprintf(format, a...))
}

func (h *Handler) Warn(summary string, format string, a ...any) {
	h.Diagnostics.AddWarning(summary, fmt.Sprintf(format, a...))
}

func (h *Handler) Error(summary string, format string, a ...any) {
	h.Diagnostics.AddError(summary, fmt.Sprintf(format, a...))
}

func (h *Handler) Invalid(format string, a ...any) {
	if !h.Diagnostics.HasError() {
		h.Diagnostics.AddError("Invalid Attribute Value", fmt.Sprintf(format, a...))
	}
}

func (h *Handler) Missing(format string, a ...any) {
	if !h.Diagnostics.HasError() {
		h.Diagnostics.AddError("Missing Attribute Value", fmt.Sprintf(format, a...))
	}
}

func (h *Handler) Conflict(format string, a ...any) {
	if !h.Diagnostics.HasError() {
		h.Diagnostics.AddError("Conflicting Attribute Values", fmt.Sprintf(format, a...))
	}
}
