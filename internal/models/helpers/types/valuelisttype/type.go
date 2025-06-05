package valuelisttype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type                = (*listTypeOf[basetypes.StringValue])(nil)
	_ attr.TypeWithElementType = (*listTypeOf[basetypes.StringValue])(nil)
	_ basetypes.ListTypable    = (*listTypeOf[basetypes.StringValue])(nil)
)

var StringListType = listTypeOf[basetypes.StringValue]{basetypes.ListType{ElemType: basetypes.StringType{}}}

type listTypeOf[T attr.Value] struct {
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
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewNullValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), diags
	}

	typ := helpers.AttrTypeOf[T](ctx)
	v, d := basetypes.NewListValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
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
