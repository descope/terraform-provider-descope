package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/iancoleman/strcase"
)

var AttributesModifier = objattr.NewModifier[AttributesModel]("maintains attribute order between plan changes")

var AttributesAttributes = map[string]schema.Attribute{
	"tenant": listattr.Default[TenantAttributeModel](TenantAttributeAttributes, TenantAttributeModifier),
	"user":   listattr.Default[UserAttributeModel](UserAttributeAttributes, UserAttributeModifier),
}

type AttributesModel struct {
	Tenant listattr.Type[TenantAttributeModel] `tfsdk:"tenant"`
	User   listattr.Type[UserAttributeModel]   `tfsdk:"user"`
}

func (m *AttributesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Tenant, data, "tenant", h)
	listattr.Get(m.User, data, "user", h)
	return data
}

func (m *AttributesModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatching(&m.Tenant, data, "tenant", "displayName", h)
	listattr.SetMatching(&m.User, data, "user", "displayName", h)
}

func (m *AttributesModel) Modify(h *helpers.Handler, state *AttributesModel) {
	listattr.ModifyMatching(h, &m.Tenant, state.Tenant)
	listattr.ModifyMatching(h, &m.User, state.User)
}

// Tenant Attributes

var TenantAttributeModifier = objattr.NewModifier[TenantAttributeModel]("ensures a suitable id is used", objattr.ModifierAllowNullState)

var TenantAttributeAttributes = map[string]schema.Attribute{
	"id":             stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":           stringattr.Required(stringattr.StandardLenValidator),
	"type":           stringattr.Required(attributeTypeValidator),
	"authorization":  objattr.Default[TenantAttributeAuthorizationModel](nil, TenantAttributeAuthorizationAttributes),
	"select_options": strsetattr.Default(),
}

type TenantAttributeModel struct {
	ID            stringattr.Type                                 `tfsdk:"id"`
	Name          stringattr.Type                                 `tfsdk:"name"`
	Type          stringattr.Type                                 `tfsdk:"type"`
	SelectOptions strsetattr.Type                                 `tfsdk:"select_options"`
	Authorization objattr.Type[TenantAttributeAuthorizationModel] `tfsdk:"authorization"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "name")
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	objattr.Get(m.Authorization, data, helpers.RootKey, h)
	getOptions(m.SelectOptions, data, "options", h)
	return data
}

func (m *TenantAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "name")
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	objattr.Set(&m.Authorization, data, helpers.RootKey, h)
	setOptions(&m.SelectOptions, data, "options", h)
}

func (m *TenantAttributeModel) Modify(h *helpers.Handler, _ *TenantAttributeModel) {
	if v := m.Name.ValueString(); v != "" && m.ID.IsUnknown() {
		id := strcase.ToLowerCamel(v)
		h.Log("Using name attribute to modify id attribute in tenant attribute to '%s'", id)
		m.ID = stringattr.Value(id)
	}
}

func (m *TenantAttributeModel) GetName() stringattr.Type {
	return m.Name
}

func (m *TenantAttributeModel) GetID() stringattr.Type {
	return m.ID
}

func (m *TenantAttributeModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Widget Authorization

var TenantAttributeAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(),
}

type TenantAttributeAuthorizationModel struct {
	ViewPermissions strsetattr.Type `tfsdk:"view_permissions"`
}

func (m *TenantAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	return data
}

func (m *TenantAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
}

// User Attributes

var UserAttributeModifier = objattr.NewModifier[UserAttributeModel]("must have a valid configuration", objattr.ModifierAllowNullState)

var UserAttributeAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":                 stringattr.Required(stringattr.StandardLenValidator),
	"type":                 stringattr.Required(attributeTypeValidator),
	"widget_authorization": objattr.Default[UserAttributeAuthorizationModel](nil, UserAttributeWidgetAuthorizationAttributes),
	"select_options":       strsetattr.Default(),
}

type UserAttributeModel struct {
	ID                  stringattr.Type                               `tfsdk:"id"`
	Name                stringattr.Type                               `tfsdk:"name"`
	Type                stringattr.Type                               `tfsdk:"type"`
	WidgetAuthorization objattr.Type[UserAttributeAuthorizationModel] `tfsdk:"widget_authorization"`
	SelectOptions       strsetattr.Type                               `tfsdk:"select_options"`
}

func (m *UserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "name")
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	getOptions(m.SelectOptions, data, "options", h)
	return data
}

func (m *UserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "name")
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
	setOptions(&m.SelectOptions, data, "options", h)
}

func (m *UserAttributeModel) Modify(h *helpers.Handler, _ *UserAttributeModel) {
	if v := m.Name.ValueString(); v != "" && m.ID.IsUnknown() {
		id := strcase.ToLowerCamel(v)
		h.Log("Using name attribute to modify id attribute in user attribute to '%s'", id)
		m.ID = stringattr.Value(id)
	}
}

func (m *UserAttributeModel) GetName() stringattr.Type {
	return m.Name
}

func (m *UserAttributeModel) GetID() stringattr.Type {
	return m.ID
}

func (m *UserAttributeModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Widget Authorization

var UserAttributeWidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(),
	"edit_permissions": strsetattr.Default(),
}

type UserAttributeAuthorizationModel struct {
	ViewPermissions strsetattr.Type `tfsdk:"view_permissions"`
	EditPermissions strsetattr.Type `tfsdk:"edit_permissions"`
}

func (m *UserAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Get(m.EditPermissions, data, "editPermissions", h)
	return data
}

func (m *UserAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Set(&m.EditPermissions, data, "editPermissions", h)
}

// Shared

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

func getOptions(s strsetattr.Type, data map[string]any, key string, h *helpers.Handler) {
	options := []map[string]any{}
	for option := range strsetattr.Iterator(s, h) {
		options = append(options, map[string]any{"label": option, "value": option})
	}
	data[key] = options
}

func setOptions(s *strsetattr.Type, data map[string]any, key string, _ *helpers.Handler) {
	result := []string{}
	if vs, ok := data[key].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					result = append(result, option)
				}
			}
		}
	}
	*s = strsetattr.Value(result)
}
