package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var UserAttributeModifier = objattr.NewModifier[UserAttributeModel]("ensures the user attribute has an id value", objattr.ModifierAllowNullState)

var UserAttributeAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":                 stringattr.Required(stringattr.StandardLenValidator),
	"type":                 stringattr.Required(attributeTypeValidator),
	"select_options":       strsetattr.Default(),
	"widget_authorization": objattr.Default(UserAttributeAuthorizationDefault, UserAttributeWidgetAuthorizationAttributes),
}

type UserAttributeModel struct {
	AttributeModel
	WidgetAuthorization objattr.Type[UserAttributeAuthorizationModel] `tfsdk:"widget_authorization"`
}

func (m *UserAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := m.AttributeModel.Values(h)
	if m.WidgetAuthorization.IsSet() {
		objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	}
	return data
}

func (m *UserAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.AttributeModel.SetValues(h, data)
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
}

func (m *UserAttributeModel) Modify(h *helpers.Handler, _ *UserAttributeModel) {
	m.AttributeModel.Modify(h)
}

// Widget Authorization

var UserAttributeWidgetAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(),
	"edit_permissions": strsetattr.Default(),
}

var UserAttributeAuthorizationDefault = &UserAttributeAuthorizationModel{
	ViewPermissions: strsetattr.Empty(),
	EditPermissions: strsetattr.Empty(),
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
