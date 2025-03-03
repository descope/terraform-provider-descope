package flows

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
)

var FlowsValidator = mapvalidator.KeysAre(stringattr.FlowIDValidator)

type FlowsModel map[string]*FlowModel

func (m *FlowsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	for flowID, flow := range *m {
		values := flow.Values(h)
		if valuesID, _ := values["flowId"].(string); valuesID != "" && valuesID != flowID {
			h.Warn("Possible flow mismatch", "The '%s' flow data specifies a different flowId '%s'. You can update the flow data to use the same flowId or ignore this warning to use the '%s' flowId.", flowID, valuesID, flowID)
		}
		values["flowId"] = flowID
		data[flowID] = values
	}
	return data
}

func (m *FlowsModel) SetValues(h *helpers.Handler, data map[string]any) {
	if len(*m) == 0 {
		for flowID, v := range data {
			if flowData, ok := v.(map[string]any); ok {
				flow := &FlowModel{}
				flow.SetValues(h, flowData)
				(*m)[flowID] = flow
			}
		}
	}
}
