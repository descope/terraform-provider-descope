package managementkey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ProjectRoleAttributes = map[string]schema.Attribute{
	"project_ids": strsetattr.Required(),
	"roles":       strsetattr.Required(),
}

type ProjectRoleModel struct {
	ProjectIDs strsetattr.Type `tfsdk:"project_ids"`
	Roles      strsetattr.Type `tfsdk:"roles"`
}

func (m *ProjectRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ProjectIDs, data, "projectIds", h)
	strsetattr.Get(m.Roles, data, "roles", h)
	return data
}

func (m *ProjectRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ProjectIDs, data, "projectIds", h)
	strsetattr.Set(&m.Roles, data, "roles", h)
}
