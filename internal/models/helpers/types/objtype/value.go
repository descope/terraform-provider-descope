package objtype

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
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
	return NewObjectTypeOfMust[T](ctx)
}

func (v ObjectValueOf[T]) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	if v.IsNull() {
		return tftypes.NewValue(v.Type(ctx).TerraformType(ctx), nil), nil
	}
	return v.ObjectValue.ToTerraformValue(ctx)
}

func (v ObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToPtr(ctx)
}

func (v ObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return ObjectValueObjectPtr[T](ctx, v)
}

func (v ObjectValueOf[T]) ToPtrMust(ctx context.Context) *T {
	return types.Must(ObjectValueObjectPtr[T](ctx, v))
}

func (v ObjectValueOf[T]) IsSet() bool {
	return !v.IsNull() && !v.IsUnknown()
}

func ObjectValueObjectPtr[T any](ctx context.Context, val attr.Value) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	ptr, d := ObjectTypeNewObjectPtr[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	diags.Append(val.(ObjectValueOf[T]).As(ctx, ptr, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	return ptr, diags
}

func ObjectValueObjectPtrMust[T any](ctx context.Context, val attr.Value) *T {
	return types.Must(ObjectValueObjectPtr[T](ctx, val))
}

func NewObjectValueOfNull[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectNull(types.AttributeTypesMust[T](ctx))}
}

func NewObjectValueOfUnknown[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectUnknown(types.AttributeTypesMust[T](ctx))}
}

func NewObjectValueOf[T any](ctx context.Context, t *T) (ObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := types.AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewObjectValueFrom(ctx, m, t)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	return ObjectValueOf[T]{ObjectValue: v}, diags
}

func NewObjectValueOfMust[T any](ctx context.Context, t *T) ObjectValueOf[T] {
	return types.Must(NewObjectValueOf(ctx, t))
}
