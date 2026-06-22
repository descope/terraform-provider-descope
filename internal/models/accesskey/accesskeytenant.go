package accesskey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AccessKeyTenantAttributes = map[string]schema.Attribute{
	"tenant_id": stringattr.Required(),
	"roles":     strsetattr.Default(),
}

type AccessKeyTenantModel struct {
	TenantID stringattr.Type `tfsdk:"tenant_id"`
	Roles    strsetattr.Type `tfsdk:"roles"`
}

func (m *AccessKeyTenantModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.TenantID, data, "tenantId")
	strsetattr.Get(m.Roles, data, "roleNames", h)
	return data
}

func (m *AccessKeyTenantModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.TenantID, data, "tenantId")
	strsetattr.Set(&m.Roles, data, "roleNames", h)
}
