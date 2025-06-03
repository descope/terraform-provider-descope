package valuelisttype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
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
	return newListTypeOf[T](ctx)
}

func (v ListValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.ListValue.ToTerraformValue(ctx)
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

func (v ListValueOf[T]) ToSliceMust(ctx context.Context) []T {
	return helpers.Must(v.ToSlice(ctx))
}

func (v ListValueOf[T]) IsEmpty() bool {
	return len(v.ListValue.Elements()) == 0
}

func NewListValueOfNull[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListNull(types.NewAttrTypeOf[T](ctx))}
}

func NewListValueOfUnknown[T attr.Value](ctx context.Context) ListValueOf[T] {
	return ListValueOf[T]{ListValue: basetypes.NewListUnknown(types.NewAttrTypeOf[T](ctx))}
}

func NewListValueOf[T attr.Value](ctx context.Context, elements []attr.Value) (ListValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewListValue(types.NewAttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
}

func NewListValueOfMust[T attr.Value](ctx context.Context, elements []attr.Value) ListValueOf[T] {
	return helpers.Must(NewListValueOf[T](ctx, elements))
}
