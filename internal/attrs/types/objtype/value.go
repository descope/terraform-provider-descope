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
	_ attr.Value               = (*ObjectValueOf[struct{}])(nil)
	_ basetypes.ObjectValuable = (*ObjectValueOf[struct{}])(nil)
)

type ObjectValueOf[T any] struct {
	basetypes.ObjectValue
}

func (v ObjectValueOf[T]) Equal(o attr.Value) bool {
	other, ok := o.(ObjectValueOf[T])
	if !ok {
		return false
	}
	return v.ObjectValue.Equal(other.ObjectValue)
}

func (v ObjectValueOf[T]) Type(ctx context.Context) attr.Type {
	return NewType[T](ctx)
}

func (v ObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.ObjectValue.ToTerraformValue(ctx)
}

func (v ObjectValueOf[T]) IsSet() bool {
	return !v.IsNull() && !v.IsUnknown()
}

func (v ObjectValueOf[T]) ToObject(ctx context.Context) (*T, diag.Diagnostics) {
	return NewObjectWith[T](ctx, v)
}

func NewNullValue[T any](ctx context.Context) ObjectValueOf[T] {
	value := basetypes.NewObjectNull(attrTypesOf[T](ctx))
	return ObjectValueOf[T]{ObjectValue: value}
}

func NewUnknownValue[T any](ctx context.Context) ObjectValueOf[T] {
	value := basetypes.NewObjectUnknown(attrTypesOf[T](ctx))
	return ObjectValueOf[T]{ObjectValue: value}
}

func NewValue[T any](ctx context.Context, object *T) (ObjectValueOf[T], diag.Diagnostics) {
	value, diags := basetypes.NewObjectValueFrom(ctx, attrTypesOf[T](ctx), object)
	if diags.HasError() {
		return NewUnknownValue[T](ctx), diags
	}

	return ObjectValueOf[T]{ObjectValue: value}, diags
}

func NewObjectWith[T any](ctx context.Context, value attr.Value) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics
	v, ok := value.(ObjectValueOf[T])
	if !ok {
		var zero ObjectValueOf[T]
		diags.AddError("Invalid Object Value", fmt.Sprintf("Expected value of type %T, got %T", zero, value))
		return nil, diags
	}

	object := nullObjectOf[T](ctx)
	d := v.As(ctx, object, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return object, diags
}
