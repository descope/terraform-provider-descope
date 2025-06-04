package valuelisttype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value             = (*ListValueOf[basetypes.StringValue])(nil)
	_ basetypes.ListValuable = (*ListValueOf[basetypes.StringValue])(nil)
)

type ListValueOf[T attr.Value] struct {
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
	return listTypeOf[T]{basetypes.ListType{ElemType: helpers.AttrTypeOf[T](ctx)}}
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

func (v ListValueOf[T]) ToSlice(ctx context.Context) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := []T{}
	d := v.ElementsAs(ctx, &values, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return values, diags
}

func NewNullValue[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListNull(helpers.AttrTypeOf[T](ctx))}
}

func NewUnknownValue[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListUnknown(helpers.AttrTypeOf[T](ctx))}
}

func NewValue[T attr.Value](ctx context.Context, elements []attr.Value) (ListValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewListValue(helpers.AttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
}
