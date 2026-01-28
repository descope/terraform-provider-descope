package managementkey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ReBacValidator = objattr.NewValidator[ReBacModel]("must have at least one role assignment")

var ReBacAttributes = map[string]schema.Attribute{
	"company_roles": strsetattr.Default(),
	"project_roles": listattr.Default[ProjectRoleModel](ProjectRoleAttributes),
	"tag_roles":     listattr.Default[TagRoleModel](TagRoleAttributes),
}

type ReBacModel struct {
	CompanyRoles strsetattr.Type                 `tfsdk:"company_roles"`
	ProjectRoles listattr.Type[ProjectRoleModel] `tfsdk:"project_roles"`
	TagRoles     listattr.Type[TagRoleModel]     `tfsdk:"tag_roles"`
}

func (m *ReBacModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.CompanyRoles, data, "companyRoles", h)
	listattr.Get(m.ProjectRoles, data, "projectRoles", h)
	listattr.Get(m.TagRoles, data, "tagRoles", h)
	return data
}

func (m *ReBacModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.CompanyRoles, data, "companyRoles", h)
	listattr.Set(&m.ProjectRoles, data, "projectRoles", h)
	listattr.Set(&m.TagRoles, data, "tagRoles", h)
}

func (m *ReBacModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.CompanyRoles, m.ProjectRoles, m.TagRoles) {
		return // skip validation if there are unknown values
	}

	hasCompanyRoles := !m.CompanyRoles.IsEmpty()
	hasOtherRoles := !m.ProjectRoles.IsEmpty() || !m.TagRoles.IsEmpty()

	if hasCompanyRoles && hasOtherRoles {
		h.Conflict("The rebac attribute cannot have both company_roles and project_roles/tag_roles")
	} else if !hasCompanyRoles && !hasOtherRoles {
		h.Missing("The rebac attribute must have at least one role in company_roles or in project_roles/tag_roles")
	}
}
