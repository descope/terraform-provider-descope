package settype

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
	_ attr.Type                = (*setNestedObjectTypeOf[struct{}])(nil)
	_ attr.TypeWithElementType = (*setNestedObjectTypeOf[struct{}])(nil)
	_ basetypes.SetTypable     = (*setNestedObjectTypeOf[struct{}])(nil)
)

type setNestedObjectTypeOf[T any] struct {
	basetypes.SetType
}

func (t setNestedObjectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(setNestedObjectTypeOf[T])
	if !ok {
		return false
	}
	return t.SetType.Equal(other.SetType)
}

func (t setNestedObjectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("setNestedObjectTypeOf[%T]", zero)
}

func (t setNestedObjectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return SetNestedObjectValueOf[T]{}
}

func (t setNestedObjectTypeOf[T]) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
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

	v, d := basetypes.NewSetValue(typ, in.Elements())
	diags.Append(d...)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return SetNestedObjectValueOf[T]{SetValue: v}, diags
}

func (t setNestedObjectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func NewType[T any](ctx context.Context) setNestedObjectTypeOf[T] {
	typ := helpers.Must(objtype.NewTypeMaybe[T](ctx))
	return setNestedObjectTypeOf[T]{SetType: basetypes.SetType{ElemType: typ}}
}
