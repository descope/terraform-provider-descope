package strsetattr

import (
	"context"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/valuesettype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = valuesettype.SetValueOf[types.String]

func Value(value []string) Type {
	return valueOf(context.Background(), value)
}

func Empty() Type {
	return valueOf(context.Background(), []string{})
}

func valueOf(ctx context.Context, value []string) Type {
	return convertStringSliceToValue(ctx, value)
}

func Required(extras ...any) schema.SetAttribute {
	return schema.SetAttribute{
		Required:    true,
		CustomType:  valuesettype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
	}
}

func Optional(extras ...any) schema.SetAttribute {
	return schema.SetAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuesettype.NewType[types.String](context.Background()),
		ElementType:   types.StringType,
		Validators:    parseExtras(extras),
		PlanModifiers: []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
	}
}

func Default(extras ...any) schema.SetAttribute {
	return schema.SetAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuesettype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
		Default:     setdefault.StaticValue(Empty().SetValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToSlice(h.Ctx))
	strings := attrs.ConvertTerraformSliceToStringSlice(values)

	// sort string slice to prevent sporadic order changes in resource updates
	slices.Sort(strings)

	data[key] = strings
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := attrs.GetStringSlice(data, key)
	*s = convertStringSliceToValue(h.Ctx, values)
}

func GetCommaSeparated(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToSlice(h.Ctx))
	value := strings.Join(attrs.ConvertTerraformSliceToStringSlice(values), ",")

	data[key] = value
}

func SetCommaSeparated(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := []string{}
	if v, _ := data[key].(string); v != "" {
		values = strings.Split(v, ",")
	}
	*s = convertStringSliceToValue(h.Ctx, values)
}

func Iterator(l Type, h *helpers.Handler) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, v := range l.Elements() {
			if v.IsNull() || v.IsUnknown() {
				continue
			}

			if str, ok := v.(types.String); ok {
				if !yield(str.ValueString()) {
					break
				}
			}
		}
	}
}

func parseExtras(extras []any) []validator.Set {
	var validators []validator.Set
	for _, e := range extras {
		matched := false
		if v, ok := e.(validator.Set); ok {
			matched = true
			validators = append(validators, v)
		}
		if v, ok := e.(validator.String); ok {
			matched = true
			validators = append(validators, setvalidator.ValueStringsAre(v))
		}
		if !matched {
			panic(fmt.Sprintf("unexpected extra value of type %T in attribute", e))
		}
	}
	return validators
}

func convertStringSliceToValue(ctx context.Context, values []string) Type {
	var elements []attr.Value
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return helpers.Require(valuesettype.NewValue[types.String](ctx, elements))
}
