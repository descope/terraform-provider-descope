package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
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
	// all attribute values are specified in the configuration
}

// Tenant Attributes

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

var TenantAttributeAttributes = map[string]schema.Attribute{
	"name":           stringattr.Required(),
	"type":           stringattr.Required(attributeTypeValidator),
	"select_options": listattr.StringOptional(),
}

type TenantAttributeModel struct {
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	SelectOptions []string     `tfsdk:"select_options"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
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
	// all attribute values are specified in the configuration
}

// User Attributes

var UserAttributeAttributes = map[string]schema.Attribute{
	"name":                 stringattr.Required(),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       listattr.StringOptional(),
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
	// all attribute values are specified in the configuration
}

// Widget Authorization

var UserAttributeWidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": listattr.StringOptional(),
	"edit_permissions": listattr.StringOptional(),
}

type UserAttributeAuthorizationModel struct {
	ViewPermissions []string `tfsdk:"view_permissions"`
	EditPermissions []string `tfsdk:"edit_permissions"`
}

func (m *UserAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	return map[string]any{
		"viewPermissions": m.ViewPermissions,
		"editPermissions": m.EditPermissions,
	}
}

func (m *UserAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all attribute values are specified in the configuration
}
