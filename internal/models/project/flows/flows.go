package flows

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

var FlowsValidator = mapvalidator.KeysAre(stringattr.FlowIDValidator)

func EnsureFlowIDs(m mapattr.Type[FlowModel], data map[string]any, key string, h *helpers.Handler) {
	values := data
	if key != helpers.RootKey {
		values, _ = data[key].(map[string]any)
	}

	for flowID := range mapattr.Iterator(m, h) {
		if v, ok := values[flowID].(map[string]any); ok {
			if valueID, _ := v["flowId"].(string); valueID != "" && valueID != flowID {
				h.Warn("Possible flow mismatch", "The '%s' flow data specifies a different flowId '%s'. You can update the flow data to use the same flowId or ignore this warning to use the '%s' flowId.", flowID, valueID, flowID)
			}
			values["flowId"] = flowID
		}
	}
}
