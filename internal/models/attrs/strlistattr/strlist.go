package strlistattr

import (
	"context"
	"fmt"
	"iter"

	"github.com/descope/terraform-provider-descope/internal/models/attrs"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/valuelisttype"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = valuelisttype.ListValueOf[types.String]

func Value(value []string) Type {
	return valueOf(context.Background(), value)
}

func Empty() Type {
	return valueOf(context.Background(), []string{})
}

func valueOf(ctx context.Context, value []string) Type {
	return convertStringSliceToValue(ctx, value)
}

func Required(extras ...any) schema.ListAttribute {
	return schema.ListAttribute{
		Required:    true,
		CustomType:  valuelisttype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
	}
}

func Optional(extras ...any) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuelisttype.NewType[types.String](context.Background()),
		ElementType:   types.StringType,
		Validators:    parseExtras(extras),
		PlanModifiers: []planmodifier.List{helpers.UseValidStateForUnknown()},
	}
}

func Default(extras ...any) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuelisttype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
		Validators:  parseExtras(extras),
		Default:     listdefault.StaticValue(Empty().ListValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToSlice(h.Ctx))
	data[key] = attrs.ConvertTerraformSliceToStringSlice(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := attrs.GetStringSlice(data, key)
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

func parseExtras(extras []any) []validator.List {
	var validators []validator.List
	for _, e := range extras {
		matched := false
		if v, ok := e.(validator.List); ok {
			matched = true
			validators = append(validators, v)
		}
		if v, ok := e.(validator.String); ok {
			matched = true
			validators = append(validators, listvalidator.ValueStringsAre(v))
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
	return helpers.Require(valuelisttype.NewValue[types.String](ctx, elements))
}
