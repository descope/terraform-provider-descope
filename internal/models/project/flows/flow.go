package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var FlowAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(stringattr.JSONValidator("metadata", "contents")),
}

type FlowModel struct {
	Data stringattr.Type `tfsdk:"data"`
}

func (m *FlowModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	return getFlowData(m.Data, h)
}

func (m *FlowModel) SetValues(h *helpers.Handler, data map[string]any) {
	if m.Data.ValueString() != "" {
		return // We do not currently update the flow data if it's already set because it might be different after apply
	}

	b, err := json.Marshal(data)
	if err != nil {
		h.Error("Unexpected flow data", "Failed to parse JSON: %s", err.Error())
		return
	}
	m.Data = stringattr.Value(string(b))
}

func (m *FlowModel) Check(h *helpers.Handler) {
	data := getFlowData(m.Data, h)

	references, _ := data["references"].(map[string]any)
	if connectors, ok := references["connectors"].(map[string]any); ok {
		for name := range connectors {
			if ref := h.Refs.Get(helpers.ConnectorReferenceKey, name); ref == nil {
				flowID, _ := data["flowId"].(string)
				h.Error("Unknown connector reference", "The flow %s requires a connector named '%s' to be defined", flowID, name)
			}
		}
	}
}

func getFlowData(data stringattr.Type, _ *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		panic("Invalid flow data after validation: " + err.Error())
	}
	return m
}
