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
	_ attr.Value            = (*MapNestedObjectValueOf[struct{}])(nil)
	_ basetypes.MapValuable = (*MapNestedObjectValueOf[struct{}])(nil)
)

type MapNestedObjectValueOf[T any] struct {
	basetypes.MapValue
}

func (v MapNestedObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(MapNestedObjectValueOf[T])
	if !ok {
		return false
	}
	return v.MapValue.Equal(other.MapValue)
}

func (v MapNestedObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewType[T](ctx)
}

func (v MapNestedObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.MapValue.ToTerraformValue(ctx)
}

func (v MapNestedObjectValueOf[T]) IsEmpty() bool {
	return len(v.Elements()) == 0
}

func (v MapNestedObjectValueOf[T]) ToMap(ctx context.Context) (map[string]*T, diag.Diagnostics) {
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

func NewNullValue[T any](ctx context.Context) MapNestedObjectValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewMapNull(typ)
	return MapNestedObjectValueOf[T]{MapValue: value}
}

func NewUnknownValue[T any](ctx context.Context) MapNestedObjectValueOf[T] {
	typ := objtype.NewType[T](ctx)
	value := basetypes.NewMapUnknown(typ)
	return MapNestedObjectValueOf[T]{MapValue: value}
}

func NewValue[T any](ctx context.Context, elements map[string]*T) (MapNestedObjectValueOf[T], diag.Diagnostics) {
	values := map[string]attr.Value{}
	for k, v := range elements {
		values[k] = helpers.Must(objtype.NewValue(ctx, v))
	}
	return NewValueWith[T](ctx, values)
}

func NewValueWith[T any](ctx context.Context, elements map[string]attr.Value) (MapNestedObjectValueOf[T], diag.Diagnostics) {
	typ := objtype.NewType[T](ctx)
	value, diags := basetypes.NewMapValue(typ, elements)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}
	return MapNestedObjectValueOf[T]{MapValue: value}, diags
}
