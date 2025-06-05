package valuemaptype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Value            = (*MapValueOf[basetypes.StringValue])(nil)
	_ basetypes.MapValuable = (*MapValueOf[basetypes.StringValue])(nil)
)

type MapValueOf[T attr.Value] struct {
	basetypes.MapValue
}

func (v MapValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapValueOf[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v MapValueOf[T]) Type(ctx context.Context) attr.Type {
	return mapTypeOf[T]{basetypes.MapType{ElemType: helpers.AttrTypeOf[T](ctx)}}
}

func (v MapValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v MapValueOf[T]) IsEmpty() bool {
	return len(v.Elements()) == 0
}

func (v MapValueOf[T]) ToMap(ctx context.Context) (map[string]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := map[string]T{}
	d := v.ElementsAs(ctx, &values, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return values, diags
}

func NewMapValue[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapNull(helpers.AttrTypeOf[T](ctx))}
}

func NewUnknownValue[T attr.Value](ctx context.Context) MapValueOf[T] {
	return MapValueOf[T]{MapValue: basetypes.NewMapUnknown(helpers.AttrTypeOf[T](ctx))}
}

func NewValue[T attr.Value](ctx context.Context, elements map[string]attr.Value) (MapValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	v, d := basetypes.NewMapValue(helpers.AttrTypeOf[T](ctx), elements)
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: v}, diags
}
