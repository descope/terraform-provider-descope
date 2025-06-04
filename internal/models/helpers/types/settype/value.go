package settype

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
	_ attr.Value            = (*SetNestedObjectValueOf[struct{}])(nil)
	_ basetypes.SetValuable = (*SetNestedObjectValueOf[struct{}])(nil)
)

type SetNestedObjectValueOf[T any] struct {
	basetypes.SetValue
}

func (v SetNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.SetValue.Equal(other.SetValue)
}

func (v SetNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewType[T](ctx)
}

func (v SetNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.SetValue.ToTerraformValue(ctx)
}

func (v SetNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.Elements()) == 0
}

func (v SetNestedObjectValueOf[T]) ToSlice(ctx context.Context) ([]*T, diag.Diagnostics) {
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

func NewNullValue[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewSetNull(typ)
	return SetNestedObjectValueOf[T]{SetValue: value}
}

func NewUnknownValue[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewSetUnknown(typ)
	return SetNestedObjectValueOf[T]{SetValue: value}
}

func NewValue[T any](ctx context.Context, values []*T) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	elements := []attr.Value{}
	for _, v := range values {
		elements = append(elements, helpers.Require(objtype.NewValue(ctx, v)))
	}
	return NewValueWith[T](ctx, elements)
}

func NewValueWith[T any](ctx context.Context, elements []attr.Value) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	typ := objtype.NewType[T](ctx)
	value, diags := basetypes.NewSetValue(typ, elements)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}
	return SetNestedObjectValueOf[T]{SetValue: value}, diags
}
