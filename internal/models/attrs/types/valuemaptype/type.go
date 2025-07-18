package valuemaptype

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                = (*mapTypeOf[basetypes.StringValue])(nil)
	_ attr.TypeWithElementType = (*mapTypeOf[basetypes.StringValue])(nil)
	_ basetypes.MapTypable     = (*mapTypeOf[basetypes.StringValue])(nil)
)

type mapTypeOf[T attr.Value] struct {
	basetypes.MapType
}

func (t mapTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(mapTypeOf[T])
	if !ok {
		return false
	}
	return t.MapType.Equal(other.MapType)
}

func (t mapTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("mapTypeOf[%T]", zero)
}

func (t mapTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return MapValueOf[T]{}
}

func (t mapTypeOf[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewMapValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), diags
	}

	v, d := basetypes.NewMapValue(elementTypeOf[T](ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: v}, diags
}

func (t mapTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.MapType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	mapValue, ok := attrValue.(basetypes.MapValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	mapValuable, diags := t.ValueFromMap(ctx, mapValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting MapValue to MapValuable: %v", diags)
	}

	return mapValuable, nil
}

func NewType[T attr.Value](ctx context.Context) mapTypeOf[T] {
	return mapTypeOf[T]{basetypes.MapType{ElemType: elementTypeOf[T](ctx)}}
}

func elementTypeOf[T attr.Value](ctx context.Context) attr.Type {
	var zero T
	return zero.Type(ctx)
}
