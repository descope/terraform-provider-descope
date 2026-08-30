package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_admin_portal is the project-level admin portal configuration singleton (id = project_id).

var AdminPortalSchema = schema.Schema{
	MarkdownDescription: "Manages the project's admin portal configuration. This is a singleton resource, and its id is always the project ID.",
	Attributes:          AdminPortalAttributes,
}

var AdminPortalAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"enabled":    boolattr.Default(false),
	"style_id":   stringattr.Default(""),
	"widgets":    listattr.Default[AdminPortalWidgetModel](AdminPortalWidgetAttributes),
}

type AdminPortalModel struct {
	ID        stringattr.Type                       `tfsdk:"id"`
	ProjectID stringattr.Type                       `tfsdk:"project_id"`
	Enabled   boolattr.Type                         `tfsdk:"enabled"`
	StyleID   stringattr.Type                       `tfsdk:"style_id"`
	Widgets   listattr.Type[AdminPortalWidgetModel] `tfsdk:"widgets"`
}

func (m *AdminPortalModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	stringattr.Get(m.StyleID, data, "styleId")
	listattr.Get(m.Widgets, data, "widgets", h)
	return data
}

func (m *AdminPortalModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	stringattr.Set(&m.StyleID, data, "styleId")
	listattr.Set(&m.Widgets, data, "widgets", h)
}

func (m *AdminPortalModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Enabled, m.Widgets) {
		return
	}
	if m.Enabled.ValueBool() && m.Widgets.IsEmpty() {
		h.Missing("admin_portal must have at least one widget when enabled")
	}
}

func (m *AdminPortalModel) GetID() stringattr.Type        { return m.ID }
func (m *AdminPortalModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *AdminPortalModel) GetProjectID() stringattr.Type { return m.ProjectID }

var AdminPortalWidgetAttributes = map[string]schema.Attribute{
	"widget_id": stringattr.Required(),
	"type":      stringattr.Required(),
}

type AdminPortalWidgetModel struct {
	WidgetID stringattr.Type `tfsdk:"widget_id"`
	Type     stringattr.Type `tfsdk:"type"`
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
