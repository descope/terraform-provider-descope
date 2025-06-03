package listtype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value             = (*ListNestedObjectValueOf[struct{}])(nil)
	_ basetypes.ListValuable = (*ListNestedObjectValueOf[struct{}])(nil)
)

type ListNestedObjectValueOf[T any] struct {
	basetypes.ListValue
}

func (v ListNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

func (v ListNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewListNestedObjectTypeOfMust[T](ctx)
}

func (v ListNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.ListValue.ToTerraformValue(ctx)
}

func (v ListNestedObjectValueOf[T]) Values(ctx context.Context) ([]*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	slice := []*T{}
	for _, v := range v.Elements() {
		ptr, d := objtype.ObjectValueObjectPtr[T](ctx, v)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		slice = append(slice, ptr)
	}

	return slice, diags
}

func (v ListNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.ListValue.Elements()) == 0
}

func NullValue[T any](ctx context.Context) ListNestedObjectValueOf[T] {
	typ := objtype.NewObjectTypeOfMust[T](ctx)
	value := basetypes.NewListNull(typ)
	return ListNestedObjectValueOf[T]{ListValue: value}
}

func UnknownValue[T any](ctx context.Context) ListNestedObjectValueOf[T] {
	typ := objtype.NewObjectTypeOfMust[T](ctx)
	value := basetypes.NewListUnknown(typ)
	return ListNestedObjectValueOf[T]{ListValue: value}
}

func Value[T any](ctx context.Context, values []*T) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	elements := []attr.Value{}
	for _, v := range values {
		elements = append(elements, helpers.Must(objtype.Value(ctx, v)))
	}
	return ValueOf[T](ctx, elements)
}

func ValueOf[T any](ctx context.Context, elements []attr.Value) (ListNestedObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	value, d := basetypes.NewListValue(typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	return ListNestedObjectValueOf[T]{ListValue: value}, diags
}
