package attribute

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var TenantAttributeSchema = schema.Schema{
	MarkdownDescription: "Manages a single tenant custom attribute definition in a Descope project.",
	Attributes:          TenantAttributeAttributes,
}

var TenantAttributeAttributes = map[string]schema.Attribute{
	"id":             stringattr.Required(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20), stringplanmodifier.RequiresReplace()),
	"project_id":     stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":           stringattr.Required(stringattr.StandardLenValidator),
	"type":           stringattr.Required(typeValidator, stringplanmodifier.RequiresReplace()),
	"select_options": strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1))),
	"authorization":  objattr.Default(tenantAttributeAuthorizationDefault, TenantAttributeAuthorizationAttributes),
}

type TenantAttributeModel struct {
	AttributeModel
	Authorization objattr.Type[tenantAttributeAuthorizationModel] `tfsdk:"authorization"`
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

var TenantAttributeAuthorizationAttributes = map[string]schema.Attribute{
	"view_permissions": strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtMost(254))),
}

var tenantAttributeAuthorizationDefault = &tenantAttributeAuthorizationModel{
	ViewPermissions: strsetattr.Empty(),
}

type tenantAttributeAuthorizationModel struct {
	ViewPermissions strsetattr.Type `tfsdk:"view_permissions"`
}

func (m *tenantAttributeAuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ViewPermissions, data, "viewPermissions", h)
	return data
}

func (m *tenantAttributeAuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ViewPermissions, data, "viewPermissions", h)
}
