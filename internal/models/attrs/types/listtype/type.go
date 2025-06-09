package listtype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/objtype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                = (*listTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*listTypeOf[struct{}])(nil)
	_ basetypes.ListTypable    = (*listTypeOf[struct{}])(nil)
)

type listTypeOf[T any] struct {
	basetypes.ListType
}

func (t listTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(listTypeOf[T])
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

func (t listTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("listTypeOf[%T]", zero)
}

func (t listTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListValueOf[T]{}
}

func (t listTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
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

	return ListValueOf[T]{ListValue: listValue}, diags
}

func (t listTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func NewType[T any](ctx context.Context) listTypeOf[T] {
	return listTypeOf[T]{ListType: basetypes.ListType{ElemType: objtype.NewType[T](ctx)}}
}
