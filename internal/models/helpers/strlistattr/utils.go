package strlistattr

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func anySliceToStringSlice(data map[string]any, key string) []string {
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

func equalStringSlicesIgnoringOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Contains(b, a[i]) {
			return false
		}
	}
	return true
}

func stringSliceToTerraformSlice(strs []string) []types.String {
	var result []types.String
	for i := range strs {
		result = append(result, types.StringValue(strs[i]))
	}
	return result
}

func terraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}
