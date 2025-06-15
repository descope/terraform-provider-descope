package attrs

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Returns a slice of Go strings from a specific key in a map.
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

// Returns a map of Go strings from a specific key in a map.
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

// Converts a slice of Terraform strings to a slice of Go strings.
func ConvertTerraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}

// Converts a map of Terraform strings to a map of Go strings.
func ConvertTerraformMapToStringMap(m map[string]types.String) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		result[k] = v.ValueString()
	}
	return result
}
