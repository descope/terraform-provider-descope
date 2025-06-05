package maptype

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
	_ attr.Value            = (*MapValueOf[struct{}])(nil)
	_ basetypes.MapValuable = (*MapValueOf[struct{}])(nil)
)

type MapValueOf[T any] struct {
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
	return NewType[T](ctx)
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

func (v MapValueOf[T]) ToMap(ctx context.Context) (map[string]*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := map[string]*T{}
	for k, element := range v.Elements() {
		ptr, d := objtype.NewObjectWith[T](ctx, element)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		result[k] = ptr
	}

	return result, diags
}

func NewNullValue[T any](ctx context.Context) MapValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewMapNull(typ)
	return MapValueOf[T]{MapValue: value}
}

func NewUnknownValue[T any](ctx context.Context) MapValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewMapUnknown(typ)
	return MapValueOf[T]{MapValue: value}
}

func NewValue[T any](ctx context.Context, elements map[string]*T) (MapValueOf[T], diag.Diagnostics) {
	values := map[string]attr.Value{}
	for k, v := range elements {
		values[k] = helpers.Require(objtype.NewValue(ctx, v))
	}
	return NewValueWith[T](ctx, values)
}

func NewValueWith[T any](ctx context.Context, elements map[string]attr.Value) (MapValueOf[T], diag.Diagnostics) {
	typ := objtype.NewType[T](ctx)
	value, diags := basetypes.NewMapValue(typ, elements)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}
	return MapValueOf[T]{MapValue: value}, diags
}
