package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AttributesAttributes = map[string]schema.Attribute{
	"tenant":     listattr.Default[TenantAttributeModel](TenantAttributeAttributes, TenantAttributeModifier),
	"user":       listattr.Default[UserAttributeModel](UserAttributeAttributes, UserAttributeModifier),
	"access_key": listattr.Default[AccessKeyAttributeModel](AccessKeyAttributeAttributes, AccessKeyAttributeModifier),
}

type AttributesModel struct {
	Tenant    listattr.Type[TenantAttributeModel]    `tfsdk:"tenant"`
	User      listattr.Type[UserAttributeModel]      `tfsdk:"user"`
	AccessKey listattr.Type[AccessKeyAttributeModel] `tfsdk:"access_key"`
}

func (m *AttributesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Tenant, data, "tenant", h)
	listattr.Get(m.User, data, "user", h)
	listattr.Get(m.AccessKey, data, "accessKey", h)
	return data
}

func (m *AttributesModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.Tenant, data, "tenant", "displayName", h)
	listattr.SetMatchingNames(&m.User, data, "user", "displayName", h)
	listattr.SetMatchingNames(&m.AccessKey, data, "accessKey", "displayName", h)
}
