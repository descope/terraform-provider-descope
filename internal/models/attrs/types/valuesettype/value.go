package valuesettype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value            = (*SetValueOf[basetypes.StringValue])(nil)
	_ basetypes.SetValuable = (*SetValueOf[basetypes.StringValue])(nil)
)

type SetValueOf[T attr.Value] struct {
	basetypes.SetValue
}

func (v SetValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(SetValueOf[T])
	if !ok {
		return false
	}
	return v.SetValue.Equal(other.SetValue)
}

func (v SetValueOf[T]) Type(ctx context.Context) attr.Type {
	return setTypeOf[T]{basetypes.SetType{ElemType: elementTypeOf[T](ctx)}}
}

func (v SetValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.SetValue.ToTerraformValue(ctx)
}

func (v SetValueOf[T]) IsEmpty() bool {
	return len(v.Elements()) == 0
}

func (v SetValueOf[T]) ToSlice(ctx context.Context) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := []T{}
	d := v.ElementsAs(ctx, &values, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return values, diags
}

func NewNullValue[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetNull(elementTypeOf[T](ctx))}
}

func NewUnknownValue[T attr.Value](ctx context.Context) SetValueOf[T] {
	return SetValueOf[T]{SetValue: basetypes.NewSetUnknown(elementTypeOf[T](ctx))}
}

func NewValue[T attr.Value](ctx context.Context, elements []attr.Value) (SetValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewSetValue(elementTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: v}, diags
}
