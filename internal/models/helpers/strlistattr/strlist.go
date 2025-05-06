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

func Get(s []types.String, data map[string]any, key string, _ *helpers.Handler) {
	data[key] = terraformSliceToStringSlice(s)
}

func Set(s *[]types.String, data map[string]any, key string, h *helpers.Handler) {
	values := anySliceToStringSlice(data, key)
	if len(*s) > 0 {
		current := terraformSliceToStringSlice(*s)
		if !equalStringSlicesIgnoringOrder(current, values) {
			h.Mismatch("Mismatched string array value in '%s' key: received [%s], expected [%s]", key, strings.Join(values, ","), strings.Join(current, ","))
		}
		return
	}
	*s = stringSliceToTerraformSlice(values)
}

func GetCommaSeparated(s []types.String, data map[string]any, key string) {
	data[key] = strings.Join(terraformSliceToStringSlice(s), ",")
}

func SetCommaSeparated(s *[]types.String, data map[string]any, key string) {
	if v, _ := data[key].(string); v != "" {
		*s = stringSliceToTerraformSlice(strings.Split(v, ","))
	}
}
