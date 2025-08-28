package widgets

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var WidgetAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(stringattr.JSONValidator("widgetId", "metadata", "screens")),
}

type WidgetModel struct {
	Data stringattr.Type `tfsdk:"data"`
}

func (m *WidgetModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	return getWidgetData(m.Data, h)
}

func (m *WidgetModel) SetValues(h *helpers.Handler, data map[string]any) {
	if m.Data.ValueString() != "" {
		return // We do not currently update the widget data if it's already set because it might be different after apply
	}

	b, err := json.Marshal(data)
	if err != nil {
		h.Error("Unexpected widget data", "Failed to parse JSON: %s", err.Error())
		return
	}
	m.Data = stringattr.Value(string(b))
}

func (m *WidgetModel) Check(h *helpers.Handler) {
	data := getWidgetData(m.Data, h)

	references, _ := data["references"].(map[string]any)
	if connectors, ok := references["connectors"].(map[string]any); ok {
		for name := range connectors {
			if ref := h.Refs.Get(helpers.ConnectorReferenceKey, name); ref == nil {
				widgetID, _ := data["widgetId"].(string)
				h.Error("Unknown connector reference", "The widget %s requires a connector named '%s' to be defined", widgetID, name)
			}
		}
	}
}

func getWidgetData(data stringattr.Type, _ *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		panic("Invalid widget data after validation: " + err.Error())
	}
	return m
}
