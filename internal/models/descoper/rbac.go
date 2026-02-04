package descoper

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var RBacValidator = objattr.NewValidator[RBacModel]("must have is_company_admin set or at least one role assignment")

var RBacAttributes = map[string]schema.Attribute{
	"is_company_admin": boolattr.Default(false),
	"project_roles":    listattr.Default[DescoperProjectRoleModel](DescoperProjectRoleAttributes),
	"tag_roles":        listattr.Default[DescoperTagRoleModel](DescoperTagRoleAttributes),
}

type RBacModel struct {
	IsCompanyAdmin boolattr.Type                           `tfsdk:"is_company_admin"`
	ProjectRoles   listattr.Type[DescoperProjectRoleModel] `tfsdk:"project_roles"`
	TagRoles       listattr.Type[DescoperTagRoleModel]     `tfsdk:"tag_roles"`
}

func (m *RBacModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.IsCompanyAdmin, data, "isCompanyAdmin")
	listattr.Get(m.ProjectRoles, data, "projects", h)
	listattr.Get(m.TagRoles, data, "tags", h)
	return data
}

func (m *RBacModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.IsCompanyAdmin, data, "isCompanyAdmin")
	listattr.Set(&m.ProjectRoles, data, "projects", h)
	listattr.Set(&m.TagRoles, data, "tags", h)
}

func (m *RBacModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.IsCompanyAdmin, m.ProjectRoles, m.TagRoles) {
		return // skip validation if there are unknown values
	}

	isCompanyAdmin := m.IsCompanyAdmin.ValueBool()
	hasOtherRoles := !m.TagRoles.IsEmpty() || !m.ProjectRoles.IsEmpty()

	if isCompanyAdmin && hasOtherRoles {
		h.Conflict("The rbac attribute cannot have both is_company_admin together with project_roles or tag_roles")
	} else if !isCompanyAdmin && !hasOtherRoles {
		h.Missing("The rbac attribute must have is_company_admin set to true or at least one role in tag_roles or project_roles")
	}
}
