package authorization

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/setattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AuthorizationValidator = objattr.NewValidator[AuthorizationModel]("must have unique role and permission names")

var AuthorizationAttributes = map[string]schema.Attribute{
	"roles":       setattr.Optional[RoleModel](RoleAttributes),
	"permissions": setattr.Optional[PermissionModel](PermissionAttributes),
}

type AuthorizationModel struct {
	Roles       setattr.Type[RoleModel]       `tfsdk:"roles"`
	Permissions setattr.Type[PermissionModel] `tfsdk:"permissions"`
}

func (m *AuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	setattr.Get(m.Roles, data, "roles", h)
	setattr.Get(m.Permissions, data, "permissions", h)
	return data
}

func (m *AuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	setattr.Set(&m.Roles, data, "roles", h)
	setattr.Set(&m.Permissions, data, "permissions", h)
}

func (m *AuthorizationModel) CollectReferences(h *helpers.Handler) {
	for v := range setattr.Iterator(m.Roles, h) {
		h.Refs.Add(helpers.RoleReferenceKey, "", v.ID.ValueString(), v.Name.ValueString())
	}
}

func (m *AuthorizationModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Permissions, m.Roles) {
		return // skip validation if there are unknown values
	}

	permissions := map[string]int{}
	roles := map[string]int{}

	for p := range setattr.Iterator(m.Permissions, h) {
		name := p.Name.ValueString()
		permissions[name] += 1
	}

	for r := range setattr.Iterator(m.Roles, h) {
		name := r.Name.ValueString()
		roles[name] += 1

		for p := range strsetattr.Iterator(r.Permissions, h) {
			if count := permissions[p]; count == 0 {
				h.Error("Missing Permission", "The role '%s' references a permission '%s' that doesn't exist", name, p)
			}
		}
	}

	for k, v := range permissions {
		if v > 1 {
			h.Error("Permission names must be unique", "The permission name '%s' is used %d times", k, v)
		}
	}

	for k, v := range roles {
		if v > 1 {
			h.Error("Role names must be unique", "The role name '%s' is used %d times", k, v)
		}
	}
}
