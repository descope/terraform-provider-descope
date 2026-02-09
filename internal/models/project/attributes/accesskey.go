package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AccessKeyAttributeModifier = objattr.NewModifier[AccessKeyAttributeModel]("ensures the access key attribute has an id value", objattr.ModifierAllowNullState)

var AccessKeyAttributeAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":                 stringattr.Required(stringattr.StandardLenValidator),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       strsetattr.Default(),
	"widget_authorization": objattr.Default(AccessKeyAttributeAuthorizationDefault, AccessKeyAttributeWidgetAuthorizationAttributes),
}

type AccessKeyAttributeModel struct {
	AttributeModel
	WidgetAuthorization objattr.Type[AccessKeyAttributeAuthorizationModel] `tfsdk:"widget_authorization"`
}

func (m *AccessKeyAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := m.AttributeModel.Values(h)
	if m.WidgetAuthorization.IsSet() {
		objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	}
	return data
}

func (m *AccessKeyAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.AttributeModel.SetValues(h, data)
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
}

func (m *AccessKeyAttributeModel) Modify(h *helpers.Handler, _ *AccessKeyAttributeModel) {
	m.AttributeModel.Modify(h)
}

// Widget Authorization

var AccessKeyAttributeWidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(),
	"edit_permissions": strsetattr.Default(),
}

var AccessKeyAttributeAuthorizationDefault = &AccessKeyAttributeAuthorizationModel{
	ViewPermissions: strsetattr.Empty(),
	EditPermissions: strsetattr.Empty(),
}

type AccessKeyAttributeAuthorizationModel struct {
	ViewPermissions strsetattr.Type `tfsdk:"view_permissions"`
	EditPermissions strsetattr.Type `tfsdk:"edit_permissions"`
}

func (m *AccessKeyAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Get(m.EditPermissions, data, "editPermissions", h)
	return data
}

func (m *AccessKeyAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
	strsetattr.Set(&m.EditPermissions, data, "editPermissions", h)
}
