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
	ensureReferences(data, "connectors", "connector", helpers.ConnectorReferenceKey, h)
	ensureReferences(data, "roles", "role", helpers.RoleReferenceKey, h)
}

func getFlowData(data stringattr.Type, _ *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		panic("Invalid flow data after validation: " + err.Error())
	}
	return m
}

func ensureReferences(data map[string]any, key string, entity string, ref string, h *helpers.Handler) {
	references, _ := data["references"].(map[string]any)
	if names, ok := references[key].(map[string]any); ok {
		for name := range names {
			if r := h.Refs.Get(ref, name); r == nil {
				flowID, _ := data["flowId"].(string)
				h.Error("Unknown "+entity+" reference", "The flow %s requires a %s named '%s' to be defined", flowID, entity, name)
			}
		}
	}
}
