package strlistattr

import (
	"context"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/strlisttype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Type = strlisttype.ListValueOf[types.String]

func Value(values []string) Type {
	return stringSliceToStringListValue(context.Background(), values)
}

func Required(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Required:    true,
		CustomType:  strlisttype.ListOfStringType,
		ElementType: types.StringType,
		Validators:  validators,
	}
}

func Optional(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:      true,
		Computed:      true,
		CustomType:    strlisttype.ListOfStringType,
		ElementType:   types.StringType,
		Validators:    validators,
		PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
	}
}

func Default(value []string, validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:    true,
		Computed:    true,
		CustomType:  strlisttype.ListOfStringType,
		ElementType: types.StringType,
		Validators:  validators,
		Default:     listdefault.StaticValue(Value(value).ListValue),
	}
}

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = terraformSliceToStringSlice(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := getStringSliceValue(data, key)

	if !s.IsEmpty() {
		current := terraformSliceToStringSlice(s.ToSliceMust(h.Ctx))
		if !equalStringSlicesIgnoringOrder(current, values) {
			h.Mismatch("Mismatched string array value in '%s' key: received [%s], expected [%s]", key, strings.Join(values, ","), strings.Join(current, ","))
		}
		return
	}

	*s = stringSliceToStringListValue(h.Ctx, values)
}

func GetCommaSeparated(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = strings.Join(terraformSliceToStringSlice(values), ",")
}

func SetCommaSeparated(s *Type, data map[string]any, key string, h *helpers.Handler) {
	values := getCommaSeparatedStringSliceValue(data, key)

	if !s.IsEmpty() {
		current := terraformSliceToStringSlice(s.ToSliceMust(h.Ctx))
		if !equalStringSlicesIgnoringOrder(current, values) {
			h.Mismatch("Mismatched comma separated string array value in '%s' key: received [%s], expected [%s]", key, strings.Join(values, ","), strings.Join(current, ","))
		}
		return
	}

	*s = stringSliceToStringListValue(h.Ctx, values)
}
