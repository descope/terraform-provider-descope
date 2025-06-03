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
	return NewSetNestedObjectTypeOfMust[T](ctx)
}

func (v SetNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.SetValue.ToTerraformValue(ctx)
}

func (v SetNestedObjectValueOf[T]) Values(ctx context.Context) ([]*T, diag.Diagnostics) {
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

func (v SetNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.SetValue.Elements()) == 0
}

func NullValue[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	typ := objtype.NewObjectTypeOfMust[T](ctx)
	value := basetypes.NewSetNull(typ)
	return SetNestedObjectValueOf[T]{SetValue: value}
}

func UnknownValue[T any](ctx context.Context) SetNestedObjectValueOf[T] {
	typ := objtype.NewObjectTypeOfMust[T](ctx)
	value := basetypes.NewSetUnknown(typ)
	return SetNestedObjectValueOf[T]{SetValue: value}
}

func Value[T any](ctx context.Context, values []*T) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	elements := []attr.Value{}
	for _, v := range values {
		elements = append(elements, helpers.Must(objtype.Value(ctx, v)))
	}
	return ValueOf[T](ctx, elements)
}

func ValueOf[T any](ctx context.Context, elements []attr.Value) (SetNestedObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	value, d := basetypes.NewSetValue(typ, elements)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{SetValue: value}, diags
}
