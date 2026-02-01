package managementkey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TagRoleAttributes = map[string]schema.Attribute{
	"tags":  strsetattr.Required(),
	"roles": strsetattr.Required(),
}

type TagRoleModel struct {
	Tags  strsetattr.Type `tfsdk:"tags"`
	Roles strsetattr.Type `tfsdk:"roles"`
}

func (m *TagRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.Tags, data, "tags", h)
	strsetattr.Get(m.Roles, data, "roles", h)
	return data
}

func (m *TagRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.Tags, data, "tags", h)
	strsetattr.Set(&m.Roles, data, "roles", h)
}
