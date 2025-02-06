package settings

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
	_ basetypes.ObjectTypable  = (*objectTypeOf[struct{}])(nil)
	_ NestedObjectType         = (*objectTypeOf[struct{}])(nil)
	_ basetypes.ObjectValuable = (*ObjectValueOf[struct{}])(nil)
	_ NestedObjectValue        = (*ObjectValueOf[struct{}])(nil)
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

type NestedObjectValue interface {
	attr.Value

	// ToObjectPtr returns the value as an object pointer (Go *struct).
	ToObjectPtr(context.Context) (any, diag.Diagnostics)
}

// objectTypeOf is the attribute type of an ObjectValueOf.
type objectTypeOf[T any] struct {
	basetypes.ObjectType
}

func newObjectTypeOf[T any](ctx context.Context) (objectTypeOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := AttributeTypes[T](ctx)
	diags.Append(d...)
	if diags.HasError() {
		return objectTypeOf[T]{}, diags
	}

	return objectTypeOf[T]{basetypes.ObjectType{AttrTypes: m}}, diags
}

func NewObjectTypeOf[T any](ctx context.Context) objectTypeOf[T] {
	return DiagsMust(newObjectTypeOf[T](ctx))
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

	m, d := AttributeTypes[T](ctx)
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
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectNull(AttributeTypesMust[T](ctx))}
}

func NewObjectValueOfUnknown[T any](ctx context.Context) ObjectValueOf[T] {
	return ObjectValueOf[T]{ObjectValue: basetypes.NewObjectUnknown(AttributeTypesMust[T](ctx))}
}

func NewObjectValueOf[T any](ctx context.Context, t *T) (ObjectValueOf[T], diag.Diagnostics) {
	var diags diag.Diagnostics

	m, d := AttributeTypes[T](ctx)
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
	return DiagsMust(NewObjectValueOf[T](ctx, t))
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
			tfType = tftypes.Object{AttributeTypes: ApplyToAllValues(v.AttributeTypes(), func(attrType attr.Type) tftypes.Type {
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

// AttributeTypes returns a map of attribute types for the specified type T.
// T must be a struct and reflection is used to find exported fields of T with the `tfsdk` tag.
func AttributeTypes[T any](ctx context.Context) (map[string]attr.Type, diag.Diagnostics) {
	var diags diag.Diagnostics
	var t T
	val := reflect.ValueOf(t)
	typ := val.Type()

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		val = reflect.New(typ.Elem()).Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		diags.Append(diag.NewErrorDiagnostic("Invalid type", fmt.Sprintf("%T has unsupported type: %s", t, typ)))
		return nil, diags
	}

	attributeTypes := make(map[string]attr.Type)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields.
		}
		tag := field.Tag.Get(`tfsdk`)
		if tag == "-" {
			continue // Skip explicitly excluded fields.
		}
		if tag == "" {
			diags.Append(diag.NewErrorDiagnostic("Invalid type", fmt.Sprintf(`%T needs a struct tag for "tfsdk" on %s`, t, field.Name)))
			return nil, diags
		}

		if v, ok := val.Field(i).Interface().(attr.Value); ok {
			attributeTypes[tag] = v.Type(ctx)
		}
	}

	return attributeTypes, nil
}

func AttributeTypesMust[T any](ctx context.Context) map[string]attr.Type {
	return DiagsMust(AttributeTypes[T](ctx))
}

func DiagsMust[T any](x T, diags diag.Diagnostics) T {
	if diags.HasError() {
		panic("invalid type")
	}
	return x
}

func ApplyToAllValues[M ~map[K]V1, K comparable, V1, V2 any](m M, f func(V1) V2) map[K]V2 {
	n := make(map[K]V2, len(m))

	for k, v := range m {
		n[k] = f(v)
	}

	return n
}
