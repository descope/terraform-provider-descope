package strlisttype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
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

func newListTypeOf[T attr.Value](ctx context.Context) listTypeOf[T] {
	return listTypeOf[T]{basetypes.ListType{ElemType: types.NewAttrTypeOf[T](ctx)}}
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
	return fmt.Sprintf("ListTypeOf[%T]", zero)
}

func (t listTypeOf[T]) ValueFromList(ctx context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewListValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewListValue(types.NewAttrTypeOf[T](ctx), in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewListValueOfUnknown[T](ctx), diags
	}

	return ListValueOf[T]{ListValue: v}, diags
}

func (t listTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) { // TODO
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

func (t listTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ListValueOf[T]{}
}
