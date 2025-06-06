package strmapattr

import (
	"context"
	"iter"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/valuemaptype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
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
	return convertStringMapToTerraformValue(ctx, value)
}

func Required(validators ...validator.Map) schema.MapAttribute {
	return schema.MapAttribute{
		Required:    true,
		CustomType:  valuemaptype.StringMapType,
		ElementType: types.StringType,
		Validators:  validators,
	}
}

func Optional(validators ...validator.Map) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuemaptype.StringMapType,
		ElementType:   types.StringType,
		Validators:    validators,
		PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
	}
}

func Default(validators ...validator.Map) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuemaptype.StringMapType,
		ElementType: types.StringType,
		Validators:  validators,
		Default:     mapdefault.StaticValue(Empty().MapValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToMap(h.Ctx))
	data[key] = convertTerraformStringMapToStringMap(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	m := getStringMap(data, key)

	if !s.IsEmpty() {
		values := helpers.Require(s.ToMap(h.Ctx))
		current := convertTerraformStringMapToStringMap(values)
		if !maps.Equal(current, m) {
			h.Mismatch("Mismatched string map value in '%s' key", key)
		}
		return
	}

	*s = convertStringMapToTerraformValue(h.Ctx, m)
}

func Nil(s *Type, h *helpers.Handler) {
	if s.IsUnknown() {
		*s = convertStringMapToTerraformValue(h.Ctx, map[string]string{})
	}
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
