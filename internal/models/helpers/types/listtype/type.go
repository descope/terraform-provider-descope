package listtype

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
	_ attr.Type                = (*listNestedObjectTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*listNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.ListTypable    = (*listNestedObjectTypeOf[struct{}])(nil)
)

type listNestedObjectTypeOf[T any] struct {
	basetypes.ListType
}

func (t listNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(listNestedObjectTypeOf[T])
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

func (t listNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("listNestedObjectTypeOf[%T]", zero)
}

func (t listNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListNestedObjectValueOf[T]{}
}

func (t listNestedObjectTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNullValue[T](ctx), nil
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), nil
	}

	listValue, diags := basetypes.NewListValue(objtype.NewType[T](ctx), in.Elements())
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return ListNestedObjectValueOf[T]{ListValue: listValue}, diags
}

func (t listNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	listValue, ok := attrValue.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	listValuable, diags := t.ValueFromList(ctx, listValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListValue to ListValuable: %v", diags)
	}

	return listValuable, nil
}

func NewType[T any](ctx context.Context) listNestedObjectTypeOf[T] {
	return listNestedObjectTypeOf[T]{ListType: basetypes.ListType{ElemType: objtype.NewType[T](ctx)}}
}
