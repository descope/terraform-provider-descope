package helpers

import (
	"strings"

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

func GetCommaSeparatedStringSlice(data map[string]any, key string) []string {
	var strs []string
	if v, _ := data[key].(string); v != "" {
		strs = strings.Split(v, ",")
	}
	return strs
}

func ConvertTerraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}
