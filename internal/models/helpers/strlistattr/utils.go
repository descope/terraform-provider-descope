package strlistattr

import (
	"context"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/strlisttype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getStringSlice(data map[string]any, key string) []string {
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

func getCommaSeparatedStringSlice(data map[string]any, key string) []string {
	var strs []string
	if v, _ := data[key].(string); v != "" {
		strs = strings.Split(v, ",")
	}
	return strs
}

func convertTerraformSliceToStringSlice(strs []types.String) []string {
	var result []string
	for i := range strs {
		result = append(result, strs[i].ValueString())
	}
	return result
}

func convertStringSliceToTerraformValue(ctx context.Context, values []string) Type {
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
	as := slices.Clone(a)
	bs := slices.Clone(b)
	slices.Sort(as)
	slices.Sort(bs)
	return slices.Equal(as, bs)
}
