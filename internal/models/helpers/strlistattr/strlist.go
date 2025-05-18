package strlistattr

import (
	"context"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/strlisttype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

func Get(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsNull() || s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = terraformSliceToStringSlice(values)
}

func Set(s *Type, data map[string]any, key string, h *helpers.Handler) {
	// if len(*s) > 0 { TODO
	// 	current := terraformSliceToStringSlice(*s)
	// 	if !equalStringSlicesIgnoringOrder(current, values) {
	// 		h.Mismatch("Mismatched string array value in '%s' key: received [%s], expected [%s]", key, strings.Join(values, ","), strings.Join(current, ","))
	// 	}
	// 	return
	// }
	*s = stringSliceToStringListValue(h.Ctx, anySliceToStringSlice(data, key))
}

func GetCommaSeparated(s Type, data map[string]any, key string, h *helpers.Handler) {
	if s.IsNull() || s.IsUnknown() {
		return
	}

	values := s.ToSliceMust(h.Ctx)
	data[key] = strings.Join(terraformSliceToStringSlice(values), ",")
}

func SetCommaSeparated(s *Type, data map[string]any, key string, h *helpers.Handler) {
	v, _ := data[key].(string)
	*s = stringSliceToStringListValue(h.Ctx, strings.Split(v, ","))
}
