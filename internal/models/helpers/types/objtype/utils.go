package objtype

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// attributeTypes returns a map of attribute types for the specified type T.
// T must be a struct and reflection is used to find exported fields of T with the `tfsdk` tag.
func attributeTypes[T any](ctx context.Context) (map[string]attr.Type, diag.Diagnostics) {
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
	for i := range typ.NumField() {
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

func attributeTypesMust[T any](ctx context.Context) map[string]attr.Type {
	return diagsMust(attributeTypes[T](ctx))
}

func diagsMust[T any](x T, diags diag.Diagnostics) T {
	if errs := diags.Errors(); len(errs) > 0 {
		panic(fmt.Sprintf("%s: %s", errs[0].Summary(), errs[0].Detail()))
	}
	return x
}

func applyToAllValues[M ~map[K]V1, K comparable, V1, V2 any](m M, f func(V1) V2) map[K]V2 {
	n := make(map[K]V2, len(m))

	for k, v := range m {
		n[k] = f(v)
	}

	return n
}
