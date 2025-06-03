package types

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func AttrTypesOf[T any](ctx context.Context) map[string]attr.Type {
	var t T
	val := reflect.ValueOf(t)
	typ := val.Type()

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		val = reflect.New(typ.Elem()).Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%T has unsupported type: %s", t, typ))
	}

	attributeTypes := make(map[string]attr.Type)
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() {
			continue // Skip unexported fields.
		}
		tag := field.Tag.Get(`tfsdk`)
		if tag == "-" {
			continue // Skip explicitly excluded fields.
		}
		if tag == "" {
			panic(fmt.Sprintf(`%T is missing a tfsdk tag on %s`, t, field.Name))
		}

		if v, ok := val.Field(i).Interface().(attr.Value); ok {
			attributeTypes[tag] = v.Type(ctx)
		}
	}

	return attributeTypes
}

func AttrTypeOf[T attr.Value](ctx context.Context) attr.Type {
	var zero T
	return zero.Type(ctx)
}
