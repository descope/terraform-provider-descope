package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
