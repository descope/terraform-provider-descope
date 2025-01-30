package strlistattr

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Required(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
		Validators:  validators,
	}
}

func Optional(validators ...validator.List) schema.ListAttribute {
	return schema.ListAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.StringType,
		Default:     listdefault.StaticValue(types.ListNull(types.StringType)),
		Validators:  validators,
	}
}

func Get(s []string, data map[string]any, key string, _ *helpers.Handler) {
	data[key] = s
}

func Set(s *[]string, data map[string]any, key string, h *helpers.Handler) {
	if *s == nil {
		return
	}
	values := helpers.AnySliceToStringSlice(data, key)
	if len(*s) > 0 {
		if !helpers.EqualStringSliceElements(*s, values) {
			h.Mismatch("Mismatched string array value in '%s' key: received [%s], expected [%s]", key, strings.Join(values, ","), strings.Join(*s, ","))
		}
		return
	}
	*s = values
}

func GetCommaSeparated(s []string, data map[string]any, key string) {
	data[key] = strings.Join(s, ",")
}

func SetCommaSeparated(s *[]string, data map[string]any, key string) {
	if v, _ := data[key].(string); v != "" {
		*s = strings.Split(v, ",")
	}
}
