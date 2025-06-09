package flows

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

var FlowsValidator = mapvalidator.KeysAre(stringattr.FlowIDValidator)

func Get(m mapattr.Type[FlowModel], h *helpers.Handler) map[string]any {
	data := map[string]any{}
	for flowID, flow := range mapattr.Iterator(m, h) {
		values := flow.Values(h)
		if valuesID, _ := values["flowId"].(string); valuesID != "" && valuesID != flowID {
			h.Warn("Possible flow mismatch", "The '%s' flow data specifies a different flowId '%s'. You can update the flow data to use the same flowId or ignore this warning to use the '%s' flowId.", flowID, valuesID, flowID)
		}
		values["flowId"] = flowID
		data[flowID] = values
	}
	return data
}

func Set(m *mapattr.Type[FlowModel], data map[string]any, key string, h *helpers.Handler) {
	values := data
	if key != helpers.RootKey {
		values, _ = data[key].(map[string]any)
	}

	flows := map[string]*FlowModel{}
	for flowID, v := range values {
		if flowData, ok := v.(map[string]any); ok {
			flow := &FlowModel{}
			flow.SetValues(h, flowData)
			flows[flowID] = flow
		}
	}
	*m = mapattr.Value(flows)
}
