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

var AttributesAttributes = map[string]schema.Attribute{
	"tenant": listattr.Default[TenantAttributeModel](TenantAttributeAttributes, TenantAttributeValidator, TenantAttributeModifier),
	"user":   listattr.Default[UserAttributeModel](UserAttributeAttributes, UserAttributeValidator, UserAttributeModifier),
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
	if m.Tenant.IsEmpty() { // XXX test without this check
		listattr.Set(&m.Tenant, data, "tenant", h)
	}
	if m.User.IsEmpty() { // XXX test without this check
		listattr.Set(&m.User, data, "user", h)
	}
}

// Tenant Attributes

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

var TenantAttributeValidator = objattr.NewValidator[TenantAttributeModel]("must have a valid configuration")

var TenantAttributeModifier = objattr.NewModifier[TenantAttributeModel]("must have a valid configuration", objattr.ModifierAllowNullState)

var TenantAttributeAttributes = map[string]schema.Attribute{
	"display_name":   stringattr.Optional(stringattr.StandardLenValidator),
	"machine_name":   stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":           stringattr.Renamed("name", "display_name", stringvalidator.LengthAtMost(20)),
	"type":           stringattr.Required(attributeTypeValidator),
	"select_options": strsetattr.Default(),
	"authorization":  objattr.Default[TenantAttributeAuthorizationModel](nil, TenantAttributeAuthorizationAttributes),
}

type TenantAttributeModel struct {
	DisplayName   stringattr.Type                                 `tfsdk:"display_name"`
	MachineName   stringattr.Type                                 `tfsdk:"machine_name"`
	Name          stringattr.Type                                 `tfsdk:"name"`
	Type          stringattr.Type                                 `tfsdk:"type"`
	SelectOptions strsetattr.Type                                 `tfsdk:"select_options"`
	Authorization objattr.Type[TenantAttributeAuthorizationModel] `tfsdk:"authorization"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.Authorization, data, helpers.RootKey, h)
	stringattr.Get(m.DisplayName, data, "displayName")
	stringattr.Get(m.MachineName, data, "name")
	stringattr.Get(m.Type, data, "type")
	getOptions(m.SelectOptions, data, "options", h)
	return data
}

func (m *TenantAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.Authorization, data, helpers.RootKey, h)
	stringattr.Set(&m.DisplayName, data, "displayName")
	stringattr.Set(&m.MachineName, data, "name")
	stringattr.Set(&m.Type, data, "type")
	setOptions(&m.SelectOptions, data, "options", h)
}

func (m *TenantAttributeModel) Modify(h *helpers.Handler, _ *TenantAttributeModel) {
	if n := m.Name.ValueString(); n != "" && m.DisplayName.IsUnknown() && m.MachineName.IsUnknown() {
		c := strcase.ToLowerCamel(n)
		h.Log("Using name attribute to modify attributes in tenant attribute to set them to '%s' and '%s'", n, c)
		m.DisplayName = stringattr.Value(n)
		m.MachineName = stringattr.Value(c)
	} else if n := m.DisplayName.ValueString(); n != "" && m.MachineName.IsUnknown() {
		c := strcase.ToLowerCamel(n)
		h.Log("Using display_name attribute to modify machine_name attribute in tenant attribute to set it to '%s'", c)
		m.MachineName = stringattr.Value(c)
	}
}

func (m *TenantAttributeModel) Validate(h *helpers.Handler) {
	if n := m.Name.ValueString(); n != "" {
		c := strcase.ToLowerCamel(n)
		if n != c {
			h.Warn("The display_name attribute should be set to '%s' and the machine_name attribute should be set to '%s' to match the behavior of the deprecated name attribute from previous versions of the provider", n, c)
		}
		if v := m.DisplayName.ValueString(); v != "" {
			if v != n {
				h.Error("Unexpected Value Conflict", "The display_name attribute is expected to be set to same value '%s' as the deprecated name attribute, and the deprecated name attribute should be removed", n)
			} else {
				h.Error("Conflicting Attribute Values", "The deprecated name attribute should be removed once the display_name attribute is set")
			}
		}
		if v := m.MachineName.ValueString(); v != "" {
			if v != c {
				h.Error("Unexpected Value Conflict", "The machine_name attribute is expected to be set to '%s' to match the value that was generated from the value '%s' of the deprecated name attribute", c, n)
			} else {
				h.Error("Conflicting Attribute Values", "The deprecated name attribute should be removed once the machine_name attribute is set")
			}
		}
		return
	}

	if m.DisplayName.ValueString() == "" && !m.DisplayName.IsUnknown() {
		h.Invalid("The display_name attribute is required and must not be empty")
	}
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

var UserAttributeValidator = objattr.NewValidator[UserAttributeModel]("must have a valid configuration")

var UserAttributeModifier = objattr.NewModifier[UserAttributeModel]("must have a valid configuration", objattr.ModifierAllowNullState)

var UserAttributeAttributes = map[string]schema.Attribute{
	"display_name":         stringattr.Optional(stringattr.StandardLenValidator),
	"machine_name":         stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":                 stringattr.Renamed("name", "display_name", stringvalidator.LengthAtMost(20)),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       strsetattr.Default(),
	"widget_authorization": objattr.Default[UserAttributeAuthorizationModel](nil, UserAttributeWidgetAuthorizationAttributes),
}

type UserAttributeModel struct {
	DisplayName         stringattr.Type                               `tfsdk:"display_name"`
	MachineName         stringattr.Type                               `tfsdk:"machine_name"`
	Name                stringattr.Type                               `tfsdk:"name"`
	Type                stringattr.Type                               `tfsdk:"type"`
	SelectOptions       strsetattr.Type                               `tfsdk:"select_options"`
	WidgetAuthorization objattr.Type[UserAttributeAuthorizationModel] `tfsdk:"widget_authorization"`
}

func (m *UserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	stringattr.Get(m.DisplayName, data, "displayName")
	stringattr.Get(m.MachineName, data, "name")
	stringattr.Get(m.Type, data, "type")
	getOptions(m.SelectOptions, data, "options", h)
	return data
}

func (m *UserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
	stringattr.Set(&m.DisplayName, data, "displayName")
	stringattr.Set(&m.MachineName, data, "name")
	stringattr.Set(&m.Type, data, "type")
	setOptions(&m.SelectOptions, data, "options", h)
}

func (m *UserAttributeModel) Modify(h *helpers.Handler, _ *UserAttributeModel) {
	if n := m.Name.ValueString(); n != "" && m.DisplayName.IsUnknown() && m.MachineName.IsUnknown() {
		c := strcase.ToLowerCamel(n)
		h.Log("Using name attribute to modify attributes in user attribute to set them to '%s' and '%s'", n, c)
		m.DisplayName = stringattr.Value(n)
		m.MachineName = stringattr.Value(c)
	} else if n := m.DisplayName.ValueString(); n != "" && m.MachineName.IsUnknown() {
		c := strcase.ToLowerCamel(n)
		h.Log("Using display_name attribute to modify machine_name attribute in user attribute to set it to '%s'", c)
		m.MachineName = stringattr.Value(c)
	}
}

func (m *UserAttributeModel) Validate(h *helpers.Handler) {
	if n := m.Name.ValueString(); n != "" {
		c := strcase.ToLowerCamel(n)
		if n != c {
			h.Warn("Deprecated Attribute", "The display_name attribute should be set to '%s' and the machine_name attribute should be set to '%s' to match the behavior of the deprecated name attribute from previous versions of the provider", n, c)
		}
		if v := m.DisplayName.ValueString(); v != "" {
			if v != n {
				h.Error("Unexpected Value Conflict", "The display_name attribute is expected to be set to same value '%s' as the deprecated name attribute, and the deprecated name attribute should be removed", n)
			} else {
				h.Error("Conflicting Attribute Values", "The deprecated name attribute should be removed once the display_name attribute is set")
			}
		}
		if v := m.MachineName.ValueString(); v != "" {
			if v != c {
				h.Error("Unexpected Value Conflict", "The machine_name attribute is expected to be set to '%s' to match the value that was generated from the value '%s' of the deprecated name attribute", c, n)
			} else {
				h.Error("Conflicting Attribute Values", "The deprecated name attribute should be removed once the machine_name attribute is set")
			}
		}
		return
	}

	if m.DisplayName.ValueString() == "" && !m.DisplayName.IsUnknown() {
		h.Invalid("The display_name is required and must not be empty")
	}
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
