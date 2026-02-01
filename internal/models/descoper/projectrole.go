package descoper

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var DescoperProjectRoleAttributes = map[string]schema.Attribute{
	"project_ids": strsetattr.Required(),
	"role":        stringattr.Required(stringvalidator.OneOf("admin", "developer", "support", "auditor")),
}

type DescoperProjectRoleModel struct {
	ProjectIDs strsetattr.Type `tfsdk:"project_ids"`
	Role       stringattr.Type `tfsdk:"role"`
}

func (m *DescoperProjectRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.ProjectIDs, data, "projectIds", h)
	stringattr.Get(m.Role, data, "role")
	return data
}

func (m *DescoperProjectRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.ProjectIDs, data, "projectIds", h)
	stringattr.Set(&m.Role, data, "role")
}
