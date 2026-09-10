package strmapattr

import (
	"context"
	"fmt"
	"iter"
	"maps"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/attrs"
	"github.com/descope/terraform-provider-descope/internal/attrs/types/valuemaptype"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = valuemaptype.MapValueOf[types.String]

func Value(value map[string]string) Type {
	return valueOf(context.Background(), value)
}

func Empty() Type {
	return valueOf(context.Background(), map[string]string{})
}

func valueOf(ctx context.Context, value map[string]string) Type {
	return convertStringMapToValue(ctx, value)
}

func Required(extras ...any) schema.MapAttribute {
	return schema.MapAttribute{
		Required:    true,
		CustomType:  valuemaptype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
	}
}

func Optional(extras ...any) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuemaptype.NewType[types.String](context.Background()),
		ElementType:   types.StringType,
		Validators:    parseExtras(extras),
		PlanModifiers: []planmodifier.Map{helpers.UseValidStateForUnknown()},
	}
}

func Default(extras ...any) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuemaptype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
		Default:     mapdefault.StaticValue(Empty().MapValue),
	}
}

func Secret(extras ...any) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		Sensitive:   true,
		CustomType:  valuemaptype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
		Default:     mapdefault.StaticValue(Empty().MapValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToMap(h.Ctx))
	data[key] = attrs.ConvertTerraformMapToStringMap(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	m := attrs.GetStringMap(data, key)
	*s = convertStringMapToValue(h.Ctx, m)
}

func Nil(s *Type, h *helpers.Handler) {
	if s.IsUnknown() {
		*s = convertStringMapToValue(h.Ctx, map[string]string{})
	}
}

// For APIs that carry a string map as an array of key and value objects. Entries are sorted by key so that the same map
// always serializes to the same array.
func GetKeyValueList(s Type, data map[string]any, key string, h *helpers.Handler) {
	pairs := map[string]string{}
	for k, v := range Iterator(s, h) {
		pairs[k] = v
	}
	entries := []any{}
	for _, k := range slices.Sorted(maps.Keys(pairs)) {
		entries = append(entries, map[string]any{"key": k, "value": pairs[k]})
	}
	data[key] = entries
}

func SetKeyValueList(s *Type, data map[string]any, key string, h *helpers.Handler) {
	pairs := map[string]string{}
	if entries, ok := data[key].([]any); ok {
		for i := range entries {
			if entry, ok := entries[i].(map[string]any); ok {
				k, _ := entry["key"].(string)
				v, _ := entry["value"].(string)
				pairs[k] = v
			}
		}
	}
	*s = valueOf(h.Ctx, pairs)
}

func Iterator(s Type, h *helpers.Handler) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range s.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			if str, ok := v.(types.String); ok {
				if !yield(k, str.ValueString()) {
					break
				}
			}
		}
	}
}

func parseExtras(extras []any) []validator.Map {
	var validators []validator.Map
	for _, e := range extras {
		matched := false
		if v, ok := e.(validator.Map); ok {
			matched = true
			validators = append(validators, v)
		}
		if v, ok := e.(validator.String); ok {
			matched = true
			validators = append(validators, mapvalidator.ValueStringsAre(v))
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in attribute", e))
		}
	}
	return validators
}

func convertStringMapToValue(ctx context.Context, m map[string]string) Type {
	elements := map[string]attr.Value{}
	for k, v := range m {
		elements[k] = types.StringValue(v)
	}
	return helpers.Require(valuemaptype.NewValue[types.String](ctx, elements))
}
