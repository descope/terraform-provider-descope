package objtype

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.Type               = (*objectTypeOf[struct{}])(nil)
	_ basetypes.ObjectTypable = (*objectTypeOf[struct{}])(nil)
)

type objectTypeOf[T any] struct {
	basetypes.ObjectType
}

func (t objectTypeOf[T]) Equal(o attr.Type) bool {
	other, ok := o.(objectTypeOf[T])
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t objectTypeOf[T]) String() string {
	var zero T
	return fmt.Sprintf("objectTypeOf[%T]", zero)
}

func (t objectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ObjectValueOf[T]{}
}

func (t objectTypeOf[T]) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	if in.IsNull() {
		return NewNullValue[T](ctx), nil
	}
	if in.IsUnknown() {
		return NewUnknownValue[T](ctx), nil
	}

	objectValue, diags := basetypes.NewObjectValue(attrTypesOf[T](ctx), in.Attributes())
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return ObjectValueOf[T]{ObjectValue: objectValue}, diags
}

func (t objectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	objectValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	objectValuable, diags := t.ValueFromObject(ctx, objectValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ObjectValue to ObjectValuable: %v", diags)
	}

	return objectValuable, nil
}

func NewType[T any](ctx context.Context) objectTypeOf[T] {
	return objectTypeOf[T]{basetypes.ObjectType{AttrTypes: attrTypesOf[T](ctx)}}
}
