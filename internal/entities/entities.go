package entities

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Interface representation for req.Plan and req.State in resource operations
type entitySource interface {
	Get(context.Context, any) diag.Diagnostics
}

// Interface representation for resp.State in resource operations
type entityTarget interface {
	Set(context.Context, any) diag.Diagnostics
}

// Loads data from the source Terraform plan or state into the model object
func load[T any](ctx context.Context, source entitySource, model *T, diagnostics *diag.Diagnostics) {
	diags := source.Get(ctx, model)
	diagnostics.Append(diags...)
}

// Saves data from the model object to the target Terraform state
func save[T any](ctx context.Context, target entityTarget, model *T, diagnostics *diag.Diagnostics) {
	diags := target.Set(ctx, model)
	diagnostics.Append(diags...)
}
