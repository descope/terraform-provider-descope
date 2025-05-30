package maptype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
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

func NewMapNestedObjectTypeOf[T any](ctx context.Context) (mapNestedObjectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	elemType, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return mapNestedObjectTypeOf[T]{}, diags
	}

	return mapNestedObjectTypeOf[T]{MapType: basetypes.MapType{ElemType: elemType}}, diags
}

func NewMapNestedObjectTypeOfMust[T any](ctx context.Context) mapNestedObjectTypeOf[T] {
	return types.Must(NewMapNestedObjectTypeOf[T](ctx))
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
	return fmt.Sprintf("MapNestedObjectTypeOf[%T]", zero)
}

func (t mapNestedObjectTypeOf[T]) ValueFromMap(ctx context.Context, in basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewMapNestedObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewMapNestedObjectValueOfUnknown[T](ctx), diags
	}

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewMapNestedObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewMapValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewMapNestedObjectValueOfUnknown[T](ctx), diags
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

func (t mapNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return MapNestedObjectValueOf[T]{}
}

func (t mapNestedObjectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objtype.ObjectTypeNewObjectPtr[T](ctx)
}

func (t mapNestedObjectTypeOf[T]) NewObjectMap(ctx context.Context, len int) (any, diag.Diagnostics) {
	return nestedObjectTypeNewObjectSlice[T](ctx, len)
}

func (t mapNestedObjectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	return NewMapNestedObjectValueOfNull[T](ctx), diags
}

func (t mapNestedObjectTypeOf[T]) ValueFromObjectMap(ctx context.Context, m any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := m.(map[string]*T); ok {
		v, d := NewMapNestedObjectValueOfMap(ctx, v)
		diags.Append(d...)
		return v, d
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid map value", fmt.Sprintf("incorrect type: want %T, got %T", (*map[string]T)(nil), m)))
	return nil, diags
}

func nestedObjectTypeNewObjectSlice[T any](_ context.Context, len int) (map[string]*T, diag.Diagnostics) {
	var diags diag.Diagnostics
	return make(map[string]*T, len), diags
}
