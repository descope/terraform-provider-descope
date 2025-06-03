package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/iancoleman/strcase"
)

var AttributesAttributes = map[string]schema.Attribute{
	"tenant": listattr.Default[TenantAttributeModel](TenantAttributeAttributes),
	"user":   listattr.Default[UserAttributeModel](UserAttributeAttributes),
}

var AttributesDefault = &AttributesModel{
	Tenant: listattr.Empty[TenantAttributeModel](),
	User:   listattr.Empty[UserAttributeModel](),
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
	if m.Tenant.IsUnknown() {
		listattr.Set(&m.Tenant, data, "tenant", h)
	}
	if m.User.IsUnknown() {
		listattr.Set(&m.User, data, "user", h)
	}
}

// Tenant Attributes

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

var TenantAttributeAttributes = map[string]schema.Attribute{
	"name":           stringattr.Required(stringvalidator.LengthAtMost(20)),
	"type":           stringattr.Required(attributeTypeValidator),
	"select_options": strsetattr.Default(),
	"authorization":  objattr.Default[TenantAttributeAuthorizationModel](nil, TenantAttributeAuthorizationAttributes),
}

type TenantAttributeModel struct {
	Name          stringattr.Type                                 `tfsdk:"name"`
	Type          stringattr.Type                                 `tfsdk:"type"`
	SelectOptions strsetattr.Type                                 `tfsdk:"select_options"`
	Authorization objattr.Type[TenantAttributeAuthorizationModel] `tfsdk:"authorization"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.Authorization, data, helpers.RootKey, h)
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	data["name"] = strcase.ToLowerCamel(m.Name.ValueString())

	options := []map[string]any{}
	for option := range strsetattr.Iterator(m.SelectOptions, h) {
		options = append(options, map[string]any{
			"label": option,
			"value": option,
		})
	}
	data["options"] = options

	return data
}

func (m *TenantAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.Authorization, data, helpers.RootKey, h)
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")

	options := []string{}
	if vs, ok := data["options"].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					options = append(options, option)
				}
			}
		}
	}
	m.SelectOptions = strsetattr.Value(options)
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

var UserAttributeAttributes = map[string]schema.Attribute{
	"name":                 stringattr.Required(stringvalidator.LengthAtMost(20)),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       strsetattr.Default(),
	"widget_authorization": objattr.Default[UserAttributeAuthorizationModel](nil, UserAttributeWidgetAuthorizationAttributes),
}

type UserAttributeModel struct {
	Name                stringattr.Type                               `tfsdk:"name"`
	Type                stringattr.Type                               `tfsdk:"type"`
	SelectOptions       strsetattr.Type                               `tfsdk:"select_options"`
	WidgetAuthorization objattr.Type[UserAttributeAuthorizationModel] `tfsdk:"widget_authorization"`
}

func (m *UserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	data["name"] = strcase.ToLowerCamel(m.Name.ValueString())

	options := []map[string]any{}
	for option := range strsetattr.Iterator(m.SelectOptions, h) {
		options = append(options, map[string]any{
			"label": option,
			"value": option,
		})
	}
	data["options"] = options
	return data
}

func (m *UserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")

	options := []string{}
	if vs, ok := data["options"].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					options = append(options, option)
				}
			}
		}
	}
	m.SelectOptions = strsetattr.Value(options)
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
