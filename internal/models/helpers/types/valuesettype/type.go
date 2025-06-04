package valuesettype

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
	_ attr.Type                = (*setTypeOf[basetypes.StringValue])(nil)
	_ attr.TypeWithElementType = (*setTypeOf[basetypes.StringValue])(nil)
	_ basetypes.SetTypable     = (*setTypeOf[basetypes.StringValue])(nil)
)

var StringSetType = setTypeOf[basetypes.StringValue]{basetypes.SetType{ElemType: basetypes.StringType{}}}

type setTypeOf[T attr.Value] struct {
	basetypes.SetType
}

func (t setTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(setTypeOf[T])
	if !ok {
		return false
	}
	return t.SetType.Equal(other.SetType)
}

func (t setTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("setTypeOf[%T]", zero)
}

func (t setTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return SetValueOf[T]{}
}

func (t setTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NewNullValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), diags
	}

	typ := helpers.AttrTypeOf[T](ctx)
	v, d := basetypes.NewSetValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return SetValueOf[T]{SetValue: v}, diags
}

func (t setTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	setValue, ok := attrValue.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	setValuable, diags := t.ValueFromSet(ctx, setValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
	}

	return setValuable, nil
}
