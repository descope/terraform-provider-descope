package strlistattr

import (
	"context"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/strlisttype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getStringSliceValue(data map[string]any, key string) []string {
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

func getCommaSeparatedStringSliceValue(data map[string]any, key string) []string {
	var strs []string
	if v, _ := data[key].(string); v != "" {
		strs = strings.Split(v, ",")
	}
	return strs
}

func terraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}

func stringSliceToStringListValue(ctx context.Context, values []string) Type {
	var elements []attr.Value
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return strlisttype.NewListValueOfMust[types.String](ctx, elements)
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
