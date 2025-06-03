package maptype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                = (*mapNestedObjectTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*mapNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.MapTypable     = (*mapNestedObjectTypeOf[struct{}])(nil)
)

type mapNestedObjectTypeOf[T any] struct {
	basetypes.MapType
}

func (t mapNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(mapNestedObjectTypeOf[T])
	if !ok {
		return false
	}
	return t.MapType.Equal(other.MapType)
}

func (t mapNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("mapNestedObjectTypeOf[%T]", zero)
}

func (t mapNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return MapNestedObjectValueOf[T]{}
}

func (t mapNestedObjectTypeOf[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewNullValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), diags
	}

	typ, d := objtype.NewTypeMaybe[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	v, d := basetypes.NewMapValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return MapNestedObjectValueOf[T]{MapValue: v}, diags
}

func (t mapNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func NewType[T any](ctx context.Context) mapNestedObjectTypeOf[T] {
	typ := helpers.Must(objtype.NewTypeMaybe[T](ctx))
	return mapNestedObjectTypeOf[T]{MapType: basetypes.MapType{ElemType: typ}}
}
