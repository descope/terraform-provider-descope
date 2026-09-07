package attribute

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var typeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date", "monthday")

type AttributeModel struct {
	ID            stringattr.Type `tfsdk:"id"`
	ProjectID     stringattr.Type `tfsdk:"project_id"`
	Name          stringattr.Type `tfsdk:"name"`
	Type          stringattr.Type `tfsdk:"type"`
	SelectOptions strsetattr.Type `tfsdk:"select_options"`
}

func (m *AttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "name")
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	strsetattr.Get(m.SelectOptions, data, "options", h)
	return data
}

func (m *AttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "name")
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	strsetattr.Set(&m.SelectOptions, data, "options", h)
}

func (m *AttributeModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AttributeModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *AttributeModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

var WidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtMost(254))),
	"edit_permissions": strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtMost(254))),
}

var widgetAuthorizationDefault = &widgetAuthorizationModel{
	ViewPermissions: strsetattr.Empty(),
	EditPermissions: strsetattr.Empty(),
}

type widgetAuthorizationModel struct {
	ViewPermissions strsetattr.Type `tfsdk:"view_permissions"`
	EditPermissions strsetattr.Type `tfsdk:"edit_permissions"`
}

func (m *widgetAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Get(m.EditPermissions, data, "editPermissions", h)
	return data
}

func (m *widgetAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Set(&m.EditPermissions, data, "editPermissions", h)
}
