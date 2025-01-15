package strlistattr

import (
	"fmt"
	"strings"

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

func Get(s []string, data map[string]any, key string) {
	data[key] = s
}

func Set(s *[]string, data map[string]any, key string) {
	if v, ok := data[key].([]any); ok {
		*s = []string{}
		if len(v) > 0 {
			for i := range v {
				str, ok := v[i].(string)
				if !ok {
					panic(fmt.Sprintf("unexpected value of type %T in string list: %s", v[i], key))
				}
				*s = append(*s, str)
			}
		}
	}
}

func GetCommaSeparated(s []string, data map[string]any, key string) {
	data[key] = strings.Join(s, ",")
}

func SetCommaSeparated(s *[]string, data map[string]any, key string) {
	if v, _ := data[key].(string); v != "" {
		*s = strings.Split(v, ",")
	}
}
