package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	DescopeConnector = "Descope"
	DescopeTemplate  = "System"
)

// Pointer receiver interface for model objects.
type Model[T any] interface {
	Values(*Handler) map[string]any
	SetValues(*Handler, map[string]any)
	*T
}

// Creates a Terraform object from a model object.
func ModelFromObject[T any, M Model[T]](ctx context.Context, object types.Object, diagnostics *diag.Diagnostics) M {
	result := new(T)
	diags := object.As(ctx, result, basetypes.ObjectAsOptions{})
	diagnostics.Append(diags...)
	return result
}

// Creates a model object from a Terraform object.
func ObjectFromModel[T any, M Model[T]](ctx context.Context, value *T, types map[string]attr.Type, diagnostics *diag.Diagnostics) types.Object {
	result, diags := basetypes.NewObjectValueFrom(ctx, types, value)
	diagnostics.Append(diags...)
	return result
}
