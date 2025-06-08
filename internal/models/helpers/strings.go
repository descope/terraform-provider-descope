package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetStringSlice(data map[string]any, key string) []string {
	var strs []string
	if objects, ok := data[key].([]any); ok {
		for i := range objects {
			if s, ok := objects[i].(string); ok {
				strs = append(strs, s)
			}
		}
	}
	return strs
}

func GetStringMap(data map[string]any, key string) map[string]string {
	result := map[string]string{}
	if m, ok := data[key].(map[string]any); ok {
		for k, v := range m {
			if s, ok := v.(string); ok {
				result[k] = s
			}
		}
	}
	return result
}

func ConvertTerraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}

func ConvertTerraformStringMapToStringMap(m map[string]types.String) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		result[k] = v.ValueString()
	}
	return result
}
