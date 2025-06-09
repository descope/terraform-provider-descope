package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var FlowAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(),
}

type FlowModel struct {
	Data stringattr.Type `tfsdk:"data"`
}

func (m *FlowModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	return getFlowData(m.Data, h)
}

func (m *FlowModel) SetValues(h *helpers.Handler, data map[string]any) {
	b, err := json.Marshal(data)
	if err != nil {
		h.Error("Invalid flow data", "Failed to parse JSON: %s", err.Error())
		return
	}
	m.Data = stringattr.Value(string(b))
}

func (m *FlowModel) Check(h *helpers.Handler) {
	data := getFlowData(m.Data, h)

	for _, field := range []string{"metadata", "contents"} {
		if data[field] == nil {
			h.Error("Invalid flow data", "Expected a JSON object with a %s field", field)
		}
	}

	references, ok := data["references"].(map[string]any)
	if !ok {
		return
	}

	if connectors, ok := references["connectors"].(map[string]any); ok {
		for name := range connectors {
			if ref := h.Refs.Get(helpers.ConnectorReferenceKey, name); ref == nil {
				flowID, _ := data["flowId"].(string)
				h.Error("Unknown connector reference", "The flow %s requires a connector named '%s' to be defined", flowID, name)
			}
		}
	}
}

func getFlowData(data stringattr.Type, h *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		h.Error("Invalid flow data", "Failed to parse JSON: %s", err.Error())
	}
	return m
}
