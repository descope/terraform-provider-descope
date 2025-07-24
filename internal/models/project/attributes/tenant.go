package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TenantAttributeModifier = objattr.NewModifier[TenantAttributeModel]("ensures a suitable id is used", objattr.ModifierAllowNullState)

var TenantAttributeAttributes = map[string]schema.Attribute{
	"id":             stringattr.Optional(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20)),
	"name":           stringattr.Required(stringattr.StandardLenValidator),
	"type":           stringattr.Required(attributeTypeValidator),
	"select_options": strsetattr.Default(),
	"authorization":  objattr.Default(TenantAttributeAuthorizationDefault, TenantAttributeAuthorizationAttributes),
}

type TenantAttributeModel struct {
	AttributeModel
	Authorization objattr.Type[TenantAttributeAuthorizationModel] `tfsdk:"authorization"`
}

func (m *TenantAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := m.AttributeModel.Values(h)
	objattr.Get(m.Authorization, data, helpers.RootKey, h)
	return data
}

func (m *TenantAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.AttributeModel.SetValues(h, data)
	objattr.Set(&m.Authorization, data, helpers.RootKey, h)
}

func (m *TenantAttributeModel) Modify(h *helpers.Handler, _ *TenantAttributeModel) {
	m.AttributeModel.Modify(h)
}

// Widget Authorization

var TenantAttributeAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(),
}

var TenantAttributeAuthorizationDefault = &TenantAttributeAuthorizationModel{
	ViewPermissions: strsetattr.Empty(),
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
