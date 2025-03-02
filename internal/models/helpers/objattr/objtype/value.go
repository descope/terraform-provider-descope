package objtype

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ObjectValuable = (*ObjectValueOf[struct{}])(nil)
	_ NestedObjectValue        = (*ObjectValueOf[struct{}])(nil)
)

type NestedObjectValue interface {
	attr.Value

	// ToObjectPtr returns the value as an object pointer (Go *struct).
	ToObjectPtr(context.Context) (any, diag.Diagnostics)
}

// ObjectValueOf represents a Terraform Plugin Framework Object value whose corresponding Go type is the structure T.
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
	return NewObjectTypeOf[T](ctx)
}

func (v ObjectValueOf[T]) ToObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return v.ToPtr(ctx)
}

func (v ObjectValueOf[T]) ToPtr(ctx context.Context) (*T, diag.Diagnostics) {
	return objectValueObjectPtr[T](ctx, v)
}

func (v ObjectValueOf[T]) ToPtrMust(ctx context.Context) *T {
	return diagsMust(objectValueObjectPtr[T](ctx, v))
}

func objectValueObjectPtr[T any](ctx context.Context, val attr.Value) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	ptr, d := objectTypeNewObjectPtr[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	diags.Append(val.(ObjectValueOf[T]).ObjectValue.As(ctx, ptr, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	return ptr, diags
}

func NewObjectValueOfNull[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectNull(attributeTypesMust[T](ctx))}
}

func NewObjectValueOfUnknown[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectUnknown(attributeTypesMust[T](ctx))}
}

func NewObjectValueOf[T any](ctx context.Context, t *T) (ObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := attributeTypes[T](ctx)
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
	return diagsMust(NewObjectValueOf[T](ctx, t))
}

func NullValueOf(ctx context.Context, v any) (attr.Value, error) {
	var attrType attr.Type
	var tfType tftypes.Type

	switch v := v.(type) {
	case basetypes.BoolValuable:
		attrType = v.Type(ctx)
		tfType = tftypes.Bool
	case basetypes.Float64Valuable:
		attrType = v.Type(ctx)
		tfType = tftypes.Number
	case basetypes.Int64Valuable:
		attrType = v.Type(ctx)
		tfType = tftypes.Number
	case basetypes.StringValuable:
		attrType = v.Type(ctx)
		tfType = tftypes.String
	case basetypes.ListValuable:
		attrType = v.Type(ctx)
		if v, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.List{ElementType: v.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.List{}
		}
	case basetypes.SetValuable:
		attrType = v.Type(ctx)
		if v, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.Set{ElementType: v.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.Set{}
		}
	case basetypes.MapValuable:
		attrType = v.Type(ctx)
		if v, ok := attrType.(attr.TypeWithElementType); ok {
			tfType = tftypes.Map{ElementType: v.ElementType().TerraformType(ctx)}
		} else {
			tfType = tftypes.Map{}
		}
	case basetypes.ObjectValuable:
		attrType = v.Type(ctx)
		if v, ok := attrType.(attr.TypeWithAttributeTypes); ok {
			tfType = tftypes.Object{AttributeTypes: applyToAllValues(v.AttributeTypes(), func(attrType attr.Type) tftypes.Type {
				return attrType.TerraformType(ctx)
			})}
		} else {
			tfType = tftypes.Object{}
		}
	default:
		return nil, nil
	}

	return attrType.ValueFromTerraform(ctx, tftypes.NewValue(tfType, nil))
}
