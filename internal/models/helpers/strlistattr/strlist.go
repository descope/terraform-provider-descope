package strlistattr

import (
	"context"
	"iter"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/valuelisttype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
	return convertStringSliceToTerraformValue(ctx, value)
}

func Required(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Required:    true,
		CustomType:  valuelisttype.StringListType,
		ElementType: types.StringType,
		Validators:  validators,
	}
}

func Optional(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    valuelisttype.StringListType,
		ElementType:   types.StringType,
		Validators:    validators,
		PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
	}
}

func Default(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  valuelisttype.StringListType,
		ElementType: types.StringType,
		Validators:  validators,
		Default:     listdefault.StaticValue(Empty().ListValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToSlice(h.Ctx))
	data[key] = helpers.ConvertTerraformSliceToStringSlice(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := helpers.GetStringSlice(data, key)
	*s = convertStringSliceToTerraformValue(h.Ctx, values)
}

func GetCommaSeparated(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := helpers.Require(s.ToSlice(h.Ctx))
	data[key] = strings.Join(helpers.ConvertTerraformSliceToStringSlice(values), ",")
}

func SetCommaSeparated(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := helpers.GetCommaSeparatedStringSlice(data, key)
	*s = convertStringSliceToTerraformValue(h.Ctx, values)
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

func convertStringSliceToTerraformValue(ctx context.Context, values []string) Type {
	var elements []attr.Value
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return helpers.Require(valuelisttype.NewValue[types.String](ctx, elements))
}
