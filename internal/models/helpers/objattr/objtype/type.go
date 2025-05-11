package objtype

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.ObjectTypable = (*objectTypeOf[struct{}])(nil)
	_ NestedObjectType        = (*objectTypeOf[struct{}])(nil)
)

type NestedObjectType interface {
	attr.Type

	// NewObjectPtr returns a new, empty value as an object pointer (Go *struct).
	NewObjectPtr(context.Context) (any, diag.Diagnostics)

	// NullValue returns a Null Value.
	NullValue(context.Context) (attr.Value, diag.Diagnostics)

	// ValueFromObjectPtr returns a Value given an object pointer (Go *struct).
	ValueFromObjectPtr(context.Context, any) (attr.Value, diag.Diagnostics)
}

// objectTypeOf is the attribute type of an ObjectValueOf.
type objectTypeOf[T any] struct {
	basetypes.ObjectType
}

func newObjectTypeOf[T any](ctx context.Context) (objectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := attributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return objectTypeOf[T]{}, diags
	}

	return objectTypeOf[T]{basetypes.ObjectType{AttrTypes: m}}, diags
}

func NewObjectTypeOf[T any](ctx context.Context) objectTypeOf[T] {
	return diagsMust(newObjectTypeOf[T](ctx))
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
		return NewObjectValueOfNull[T](ctx), diags
	}
	if in.IsUnknown() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	m, d := attributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	v, d := basetypes.NewObjectValue(m, in.Attributes())
	diags.Append(d...)
	if diags.HasError() {
		return NewObjectValueOfUnknown[T](ctx), diags
	}

	value := ObjectValueOf[T]{
		ObjectValue: v,
	}

	return value, diags
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

func (t objectTypeOf[T]) ValueType(ctx context.Context) attr.Value {
	return ObjectValueOf[T]{}
}

func (t objectTypeOf[T]) NewObjectPtr(ctx context.Context) (any, diag.Diagnostics) {
	return objectTypeNewObjectPtr[T](ctx)
}

func (t objectTypeOf[T]) NullValue(ctx context.Context) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	return NewObjectValueOfNull[T](ctx), diags
}

func (t objectTypeOf[T]) ValueFromObjectPtr(ctx context.Context, ptr any) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	if v, ok := ptr.(*T); ok {
		v, d := NewObjectValueOf(ctx, v)
		diags.Append(d...)
		return v, diags
	}

	diags.Append(diag.NewErrorDiagnostic("Invalid pointer value", fmt.Sprintf("incorrect type: want %T, got %T", (*T)(nil), ptr)))
	return nil, diags
}

func objectTypeNewObjectPtr[T any](ctx context.Context) (*T, diag.Diagnostics) {
	var diags diag.Diagnostics

	t := new(T)
	diags.Append(NullOutObjectPtrFields(ctx, t)...)
	if diags.HasError() {
		return nil, diags
	}

	return t, diags
}

// NullOutObjectPtrFields sets all applicable fields of the specified object pointer to their null values.
func NullOutObjectPtrFields[T any](ctx context.Context, t *T) diag.Diagnostics {
	var diags diag.Diagnostics
	val := reflect.ValueOf(t)
	typ := val.Type().Elem()

	if typ.Kind() != reflect.Struct {
		return diags
	}

	val = val.Elem()

	for i := 0; i < typ.NumField(); i++ {
		val := val.Field(i)
		if !val.CanInterface() {
			continue
		}

		attrValue, err := NullValueOf(ctx, val.Interface())

		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("attr.Type.ValueFromTerraform", err.Error()))
			return diags
		}

		if attrValue == nil {
			continue
		}

		val.Set(reflect.ValueOf(attrValue))
	}

	return diags
}
