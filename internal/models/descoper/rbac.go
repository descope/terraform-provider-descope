package descoper

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/setattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var RBacValidator = objattr.NewValidator[RBacModel]("must have is_company_admin set or at least one role assignment")

var RBacAttributes = map[string]schema.Attribute{
	"is_company_admin": boolattr.Default(false),
	"project_roles":    setattr.Default[DescoperProjectRoleModel](DescoperProjectRoleAttributes),
	"tag_roles":        setattr.Default[DescoperTagRoleModel](DescoperTagRoleAttributes),
}

type RBacModel struct {
	IsCompanyAdmin boolattr.Type                          `tfsdk:"is_company_admin"`
	ProjectRoles   setattr.Type[DescoperProjectRoleModel] `tfsdk:"project_roles"`
	TagRoles       setattr.Type[DescoperTagRoleModel]     `tfsdk:"tag_roles"`
}

func (m *RBacModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.IsCompanyAdmin, data, "isCompanyAdmin")
	setattr.Get(m.ProjectRoles, data, "projects", h)
	setattr.Get(m.TagRoles, data, "tags", h)
	return data
}

func (m *RBacModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.IsCompanyAdmin, data, "isCompanyAdmin")
	setattr.Set(&m.ProjectRoles, data, "projects", h)
	setattr.Set(&m.TagRoles, data, "tags", h)
}

func (m *RBacModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.IsCompanyAdmin, m.ProjectRoles, m.TagRoles) {
		return
	}

	isCompanyAdmin := m.IsCompanyAdmin.ValueBool()
	hasOtherRoles := !m.TagRoles.IsEmpty() || !m.ProjectRoles.IsEmpty()

	if isCompanyAdmin && hasOtherRoles {
		h.Conflict("The rbac attribute cannot have both is_company_admin together with project_roles or tag_roles")
	} else if !isCompanyAdmin && !hasOtherRoles {
		h.Missing("The rbac attribute must have is_company_admin set to true or at least one role in tag_roles or project_roles")
	}
}
