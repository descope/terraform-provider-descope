package objtype

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types"
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

func NewObjectTypeOf[T any](ctx context.Context) (objectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := types.AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return objectTypeOf[T]{}, diags
	}

	return objectTypeOf[T]{basetypes.ObjectType{AttrTypes: m}}, diags
}

func NewObjectTypeOfMust[T any](ctx context.Context) objectTypeOf[T] {
	return helpers.Must(NewObjectTypeOf[T](ctx))
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
	return fmt.Sprintf("ObjectTypeOf[%T]", zero)
}

func (t objectTypeOf[T]) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	if in.IsNull() {
		return NullValue[T](ctx), diags
	}
	if in.IsUnknown() {
		return UnknownValue[T](ctx), diags
	}

	m, d := types.AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	v, d := basetypes.NewObjectValue(m, in.Attributes())
	diags.Append(d...)
	if diags.HasError() {
		return UnknownValue[T](ctx), diags
	}

	value := ObjectValueOf[T]{
		ObjectValue: v,
	}

	return value, diags
}

func (t objectTypeOf[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) { // TODO
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

func (t objectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ObjectValueOf[T]{}
}

func (t objectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return ObjectTypeNewObjectPtr[T](ctx)
}

func (t objectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	return NullValue[T](ctx), diags
}

func (t objectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, ptr any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := ptr.(*T); ok {
		v, d := Value(ctx, v)
		diags.Append(d...)
		return v, diags
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: want %T, got %T", (*T)(nil), ptr)))
	return nil, diags
}

func ObjectTypeNewObjectPtr[T any](ctx context.Context) (*T, diag.Diagnostics) {
	t := new(T)

	diags := nullObjectFields(ctx, t)
	if diags.HasError() {
		return nil, diags
	}

	return t, diags
}
