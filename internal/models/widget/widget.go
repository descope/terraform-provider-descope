package widget

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single widget in a Descope project via its exported JSON representation, the same format produced by exporting a widget in the Descope console.",
	Attributes:          WidgetAttributes,
}

var WidgetAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"widget_id":  stringattr.Required(stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace()),
	"data":       jsonattr.Required("metadata", "screens"),
}

type WidgetModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	WidgetID  stringattr.Type `tfsdk:"widget_id"`
	Data      jsonattr.Type   `tfsdk:"data"`
}

func (m *WidgetModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	jsonattr.Get(m.Data, data, helpers.RootKey)
	widgetID := m.WidgetID.ValueString()
	if valueID, _ := data["widgetId"].(string); valueID != "" && valueID != widgetID {
		h.Warn("Possible widget mismatch", "The widget data specifies a different widgetId '%s'. You can update the widget data to use the same widgetId or ignore this warning to use the '%s' widgetId.", valueID, widgetID)
	}
	data["widgetId"] = widgetID
	return data
}

func (m *WidgetModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.WidgetID, data, "widgetId", stringattr.SkipIfAlreadySet)
	jsonattr.Set(&m.Data, data, helpers.RootKey, jsonattr.SkipIfAlreadySet)
}

func (m *WidgetModel) GetID() stringattr.Type        { return m.ID }
func (m *WidgetModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *WidgetModel) GetProjectID() stringattr.Type { return m.ProjectID }
