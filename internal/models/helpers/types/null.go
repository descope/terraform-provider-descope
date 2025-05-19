package types

import (
	"context"
	"iter"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func NullOutObjectPtrFields[T any](ctx context.Context, t *T) diag.Diagnostics {
	var diags diag.Diagnostics
	val := reflect.ValueOf(t)
	if val.Type().Elem().Kind() != reflect.Struct {
		return diags
	}

	val = val.Elem()
	for field := range structFields(val.Type()) {
		fieldVal := val.FieldByIndex(field.Index)
		if !fieldVal.CanInterface() {
			continue
		}

		attrValue, err := nullValueOf(ctx, fieldVal.Interface())
		if err != nil {
			diags.Append(diag.NewErrorDiagnostic("attr.Type.ValueFromTerraform", err.Error()))
			return diags
		}

		if attrValue == nil {
			continue
		}

		fieldVal.Set(reflect.ValueOf(attrValue))
	}

	return diags
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

func structFields(typ reflect.Type) iter.Seq[reflect.StructField] {
	return func(yield func(reflect.StructField) bool) {
		for i := range typ.NumField() {
			field := typ.Field(i)

			if field.Anonymous {
				fieldIndexSequence := []int{i}
				for v := range structFields(field.Type) {
					v.Index = append(fieldIndexSequence, v.Index...)
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
