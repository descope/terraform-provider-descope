package listtype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value             = (*ListValueOf[struct{}])(nil)
	_ basetypes.ListValuable = (*ListValueOf[struct{}])(nil)
)

type ListValueOf[T any] struct {
	basetypes.ListValue
}

func (v ListValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ListValueOf[T])
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

func (v ListValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewType[T](ctx)
}

func (v ListValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.ListValue.ToTerraformValue(ctx)
}

func (v ListValueOf[T]) IsEmpty() bool {
	return len(v.Elements()) == 0
}

func (v ListValueOf[T]) ToSlice(ctx context.Context) ([]*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := []*T{}
	for _, element := range v.Elements() {
		ptr, d := objtype.NewObjectWith[T](ctx, element)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		result = append(result, ptr)
	}

	return result, diags
}

func NewNullValue[T any](ctx context.Context) ListValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewListNull(typ)
	return ListValueOf[T]{ListValue: value}
}

func NewUnknownValue[T any](ctx context.Context) ListValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewListUnknown(typ)
	return ListValueOf[T]{ListValue: value}
}

func NewValue[T any](ctx context.Context, values []*T) (ListValueOf[T], diag.Diagnostics) {
	elements := []attr.Value{}
	for _, v := range values {
		elem, diags := objtype.NewValue(ctx, v)
		if diags.HasError() {
			return NewUnknownValue[T](ctx), diags
		}
		elements = append(elements, elem)
	}
	return NewValueWith[T](ctx, elements)
}

func NewValueWith[T any](ctx context.Context, elements []attr.Value) (ListValueOf[T], diag.Diagnostics) {
	typ := objtype.NewType[T](ctx)
	value, diags := basetypes.NewListValue(typ, elements)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}
	return ListValueOf[T]{ListValue: value}, diags
}
