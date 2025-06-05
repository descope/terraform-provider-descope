package maptype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                = (*mapTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*mapTypeOf[struct{}])(nil)
	_ basetypes.MapTypable     = (*mapTypeOf[struct{}])(nil)
)

type mapTypeOf[T any] struct {
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
	return fmt.Sprintf("mapNestedObjectTypeOf[%T]", zero)
}

func (t mapTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return MapValueOf[T]{}
}

func (t mapTypeOf[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNullValue[T](ctx), nil
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), nil
	}

	setValue, diags := basetypes.NewMapValue(objtype.NewType[T](ctx), in.Elements())
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return MapValueOf[T]{MapValue: setValue}, diags
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

func NewType[T any](ctx context.Context) mapTypeOf[T] {
	return mapTypeOf[T]{MapType: basetypes.MapType{ElemType: objtype.NewType[T](ctx)}}
}
