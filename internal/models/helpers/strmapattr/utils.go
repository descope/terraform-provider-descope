package strmapattr

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/types/valuemaptype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getStringMap(data map[string]any, key string) map[string]string {
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

func convertStringMapToTerraformValue(ctx context.Context, m map[string]string) Type {
	elements := map[string]attr.Value{}
	for k, v := range m {
		elements[k] = types.StringValue(v)
	}
	return helpers.Require(valuemaptype.NewValue[types.String](ctx, elements))
}

func convertTerraformStringMapToStringMap(m map[string]types.String) map[string]string {
	result := map[string]string{}
	for k, v := range m {
		result[k] = v.ValueString()
	}
	return result
}
