package adminportal

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AdminPortalValidator = objattr.NewValidator[AdminPortalModel]("must not be empty and have at least one widget")

var AdminPortalWidgetAttributes = map[string]schema.Attribute{
	"widget_id": stringattr.Required(),
	"type":      stringattr.Required(),
}

type AdminPortalWidgetModel struct {
	WidgetID stringattr.Type `tfsdk:"widget_id"`
	Type     stringattr.Type `tfsdk:"type"`
}

var AdminPortalAttributes = map[string]schema.Attribute{
	"enabled":  boolattr.Default(false),
	"style_id": stringattr.Default(""),
	"widgets":  listattr.Default[AdminPortalWidgetModel](AdminPortalWidgetAttributes),
}

func (m *AdminPortalWidgetModel) Values(_ *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.WidgetID, data, "id")
	stringattr.Get(m.Type, data, "type")
	return data
}

func (m *AdminPortalWidgetModel) SetValues(_ *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.WidgetID, data, "id")
	stringattr.Set(&m.Type, data, "type")
}

type AdminPortalModel struct {
	Enabled boolattr.Type                         `tfsdk:"enabled"`
	StyleID stringattr.Type                       `tfsdk:"style_id"`
	Widgets listattr.Type[AdminPortalWidgetModel] `tfsdk:"widgets"`
}

func (m *AdminPortalModel) Values(h *helpers.Handler) map[string]any {
	config := map[string]any{}
	boolattr.Get(m.Enabled, config, "enabled")
	stringattr.Get(m.StyleID, config, "styleId")
	listattr.Get(m.Widgets, config, "widgets", h)
	return map[string]any{"config": config}
}

func (m *AdminPortalModel) SetValues(h *helpers.Handler, data map[string]any) {
	config, ok := data["config"].(map[string]any)
	if !ok {
		return
	}
	boolattr.Set(&m.Enabled, config, "enabled")
	stringattr.Set(&m.StyleID, config, "styleId")
	listattr.Set(&m.Widgets, config, "widgets", h)
}

func (m *AdminPortalModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Enabled, m.StyleID, m.Widgets) {
		return
	}

	if m.Enabled.ValueBool() && (m.Widgets.IsEmpty()) {
		h.Missing("admin_portal must have at least one widget when enabled")
	}
}
