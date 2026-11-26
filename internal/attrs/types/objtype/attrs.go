package objtype

import (
	"context"
	"fmt"
	"iter"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func attrTypesOf[T any](ctx context.Context) map[string]attr.Type {
	object, val, typ := pointerTypeOf[T]()

	result := map[string]attr.Type{}
	for field := range exportedStructFields(typ) {
		tag := field.Tag.Get(`tfsdk`)
		if tag == "-" {
			continue
		}
		if tag == "" {
			panic(fmt.Sprintf(`%T is missing a tfsdk tag on %s`, object, field.Name))
		}

		fieldVal := val.FieldByIndex(field.Index)
		if v, ok := fieldVal.Interface().(attr.Value); ok {
			result[tag] = v.Type(ctx)
		}
	}

	return result
}

func nullObjectOf[T any](ctx context.Context) *T {
	object, val, typ := pointerTypeOf[T]()

	for field := range structFields(typ) {
		fieldVal := val.FieldByIndex(field.Index)
		if !fieldVal.CanInterface() {
			continue
		}

		attrValue, err := nullValueOf(ctx, fieldVal.Interface())
		if err != nil {
			panic(fmt.Sprintf("failed to create null value for field %s of type %T: %s", fieldVal.Type().Name(), object, err.Error()))
		}
		if attrValue == nil {
			continue
		}

		fieldVal.Set(reflect.ValueOf(attrValue))
	}

	return object
}

func nullValueOf(ctx context.Context, v any) (attr.Value, error) {
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
			types := map[string]tftypes.Type{}
			for k, attrType := range v.AttributeTypes() {
				types[k] = attrType.TerraformType(ctx)
			}
			tfType = tftypes.Object{AttributeTypes: types}
		} else {
			tfType = tftypes.Object{}
		}
	default:
		return nil, nil
	}

	return attrType.ValueFromTerraform(ctx, tftypes.NewValue(tfType, nil))
}

func pointerTypeOf[T any]() (*T, reflect.Value, reflect.Type) {
	object := new(T)

	val := reflect.ValueOf(object)
	typ := val.Type()

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("%T cannot be used to create pointer type", object))
	}

	val = val.Elem()
	typ = val.Type()

	return object, val, typ
}

func exportedStructFields(typ reflect.Type) iter.Seq[reflect.StructField] {
	return func(yield func(reflect.StructField) bool) {
		for field := range structFields(typ) {
			if !field.IsExported() && !field.Anonymous {
				continue
			}
			if !yield(field) {
				return
			}
		}
	}
}

func structFields(typ reflect.Type) iter.Seq[reflect.StructField] {
	return func(yield func(reflect.StructField) bool) {
		for i := range typ.NumField() {
			field := typ.Field(i)
			if field.Anonymous {
				indexSequence := []int{i}
				for v := range structFields(field.Type) {
					v.Index = append(indexSequence, v.Index...)
					if !yield(v) {
						return
					}
				}
				continue
			}
			if !yield(field) {
				return
			}
		}
	}
}
