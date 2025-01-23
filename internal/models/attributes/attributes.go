package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/iancoleman/strcase"
)

var AttributesAttributes = map[string]schema.Attribute{
	"tenant": listattr.Optional(TenantAttributeAttributes),
	"user":   listattr.Optional(UserAttributeAttributes),
}

type AttributesModel struct {
	Tenant []*TenantAttributeModel `tfsdk:"tenant"`
	User   []*UserAttributeModel   `tfsdk:"user"`
}

func (m *AttributesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Tenant, data, "tenant", h)
	listattr.Get(m.User, data, "user", h)
	return data
}

func (m *AttributesModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.Tenant = []*TenantAttributeModel{}
	listattr.Set(&m.Tenant, data, "tenant", h)
	m.User = []*UserAttributeModel{}
	listattr.Set(&m.User, data, "user", h)
}

// Tenant Attributes

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

var TenantAttributeAttributes = map[string]schema.Attribute{
	"name":           stringattr.Required(stringvalidator.LengthAtMost(20)),
	"type":           stringattr.Required(attributeTypeValidator),
	"select_options": strlistattr.Optional(),
	"authorization":  objectattr.Optional(TenantAttributeAuthorizationAttributes),
}

type TenantAttributeModel struct {
	Name          types.String                       `tfsdk:"name"`
	Type          types.String                       `tfsdk:"type"`
	SelectOptions []string                           `tfsdk:"select_options"`
	Authorization *TenantAttributeAuthorizationModel `tfsdk:"authorization"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	if m.Authorization != nil {
		data = m.Authorization.Values(h)
	}
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	data["name"] = strcase.ToLowerCamel(m.Name.ValueString())
	options := []map[string]any{}
	for _, o := range m.SelectOptions {
		options = append(options, map[string]any{
			"label": o,
			"value": o,
		})
	}
	data["options"] = options

	return data
}

func (m *TenantAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.Authorization = utils.ZVL(m.Authorization)
	m.Authorization.SetValues(h, data)
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	if vs, ok := data["options"].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					m.SelectOptions = append(m.SelectOptions, option)
				}
			}
		}
	}
}

// Widget Authorization

var TenantAttributeAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strlistattr.Optional(),
}

type TenantAttributeAuthorizationModel struct {
	ViewPermissions []string `tfsdk:"view_permissions"`
}

func (m *TenantAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strlistattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	return data
}

func (m *TenantAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strlistattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
}

// User Attributes

var UserAttributeAttributes = map[string]schema.Attribute{
	"name":                 stringattr.Required(stringvalidator.LengthAtMost(20)),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       strlistattr.Optional(),
	"widget_authorization": objectattr.Optional(UserAttributeWidgetAuthorizationAttributes),
}

type UserAttributeModel struct {
	Name                types.String                     `tfsdk:"name"`
	Type                types.String                     `tfsdk:"type"`
	SelectOptions       []string                         `tfsdk:"select_options"`
	WidgetAuthorization *UserAttributeAuthorizationModel `tfsdk:"widget_authorization"`
}

func (m *UserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	if m.WidgetAuthorization != nil {
		data = m.WidgetAuthorization.Values(h)
	}
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	data["name"] = strcase.ToLowerCamel(m.Name.ValueString())
	options := []map[string]any{}
	for _, o := range m.SelectOptions {
		options = append(options, map[string]any{
			"label": o,
			"value": o,
		})
	}
	data["options"] = options

	return data
}

func (m *UserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.WidgetAuthorization = utils.ZVL(m.WidgetAuthorization)
	m.WidgetAuthorization.SetValues(h, data)
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	if vs, ok := data["options"].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					m.SelectOptions = append(m.SelectOptions, option)
				}
			}
		}
	}
}

// Widget Authorization

var UserAttributeWidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strlistattr.Optional(),
	"edit_permissions": strlistattr.Optional(),
}

type UserAttributeAuthorizationModel struct {
	ViewPermissions []string `tfsdk:"view_permissions"`
	EditPermissions []string `tfsdk:"edit_permissions"`
}

func (m *UserAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strlistattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	strlistattr.Get(m.EditPermissions, data, "editPermissions", h)
	return data
}

func (m *UserAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strlistattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
	strlistattr.Set(&m.EditPermissions, data, "editPermissions", h)
}
