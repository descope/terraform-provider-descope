package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FlowValidator = objectattr.NewValidator[FlowModel]("must be valid JSON data")

var FlowAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(),
}

type FlowModel struct {
	Data types.String `tfsdk:"data"`
}

func (m *FlowModel) Values(h *helpers.Handler) map[string]any {
	data := getFlowData(m.Data, h)
	checkReferences(h, data)
	return data
}

func (m *FlowModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all styles values are specified in the configuration
}

func (m *FlowModel) Validate(h *helpers.Handler) {
	data := getFlowData(m.Data, h)
	for _, field := range []string{"metadata", "contents", "screens"} {
		if data[field] == nil {
			h.Error("Invalid flow data", "Expected a JSON object with a %s field", field)
		}
	}
}

func getFlowData(data types.String, h *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		h.Error("Invalid flow data", "Failed to parse JSON: %s", err.Error())
	}
	return m
}

func checkReferences(h *helpers.Handler, data map[string]any) {
	references, ok := data["references"].(map[string]any)
	if !ok {
		return
	}

	if connectors, ok := references["connectors"].(map[string]any); ok {
		for name := range connectors {
			if ref := h.Refs.Get(helpers.ConnectorReferenceKey, name); ref == nil {
				flowID, _ := data["flowId"].(string)
				h.Error("Unknown connector reference", "No connector named '%s' was defined but it's required for flow %s", name, flowID)
			}
		}
	}
}
