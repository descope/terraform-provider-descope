package valuesettype

import (
	"context"
	"fmt"

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

	v, d := basetypes.NewSetValue(elementTypeOf[T](ctx), in.Elements())
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

func NewType[T attr.Value](ctx context.Context) setTypeOf[T] {
	return setTypeOf[T]{basetypes.SetType{ElemType: elementTypeOf[T](ctx)}}
}

func elementTypeOf[T attr.Value](ctx context.Context) attr.Type {
	var zero T
	return zero.Type(ctx)
}
