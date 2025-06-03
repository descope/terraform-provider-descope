package listtype

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
	_ attr.Type                = (*listNestedObjectTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*listNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.ListTypable    = (*listNestedObjectTypeOf[struct{}])(nil)
)

type listNestedObjectTypeOf[T any] struct {
	basetypes.ListType
}

func NewListNestedObjectTypeOf[T any](ctx context.Context) (listNestedObjectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	elemType, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return listNestedObjectTypeOf[T]{}, diags
	}

	return listNestedObjectTypeOf[T]{
		ListType: basetypes.ListType{ElemType: elemType},
	}, diags
}

func NewListNestedObjectTypeOfMust[T any](ctx context.Context) listNestedObjectTypeOf[T] {
	return types.Must(NewListNestedObjectTypeOf[T](ctx))
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
	return fmt.Sprintf("ListNestedObjectTypeOf[%T]", zero)
}

func (t listNestedObjectTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NullValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return UnknownValue[T](ctx), diags
	}

	typ, d := objtype.NewObjectTypeOf[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	v, d := basetypes.NewListValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	return ListNestedObjectValueOf[T]{ListValue: v}, diags
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
