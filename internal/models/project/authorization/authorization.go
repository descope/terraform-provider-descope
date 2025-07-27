package authorization

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var systemPermissions = []string{
	"Impersonate",
	"User Admin",
	"SSO Admin",
}

var AuthorizationValidator = objattr.NewValidator[AuthorizationModel]("must have unique role and permission names")

var AuthorizationModifier = objattr.NewModifier[AuthorizationModel]("maintains permission and role identifiers between plan changes")

var AuthorizationAttributes = map[string]schema.Attribute{
	"roles":       listattr.Default[RoleModel](RoleAttributes),
	"permissions": listattr.Default[PermissionModel](PermissionAttributes),
}

type AuthorizationModel struct {
	Roles       listattr.Type[RoleModel]       `tfsdk:"roles"`
	Permissions listattr.Type[PermissionModel] `tfsdk:"permissions"`
}

func (m *AuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Roles, data, "roles", h)
	listattr.Get(m.Permissions, data, "permissions", h)
	return data
}

func (m *AuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatching(&m.Roles, data, "roles", "name", h)
	listattr.SetMatching(&m.Permissions, data, "permissions", "name", h)
}

func (m *AuthorizationModel) CollectReferences(h *helpers.Handler) {
	for v := range listattr.Iterator(m.Roles, h) {
		h.Refs.Add(helpers.RoleReferenceKey, "", v.ID.ValueString(), v.Name.ValueString())
	}
}

func (m *AuthorizationModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Permissions, m.Roles) {
		return // skip validation if there are unknown values
	}

	permissions := map[string]int{}
	roles := map[string]int{}

	for _, n := range systemPermissions {
		permissions[n] = 1
	}

	for p := range listattr.Iterator(m.Permissions, h) {
		name := p.Name.ValueString()
		permissions[name] += 1

		if slices.Contains(systemPermissions, name) {
			h.Invalid("The permission '%s' is a system permission and is already defined", name)
			return
		}
	}

	for r := range listattr.Iterator(m.Roles, h) {
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

func (m *AuthorizationModel) Modify(h *helpers.Handler, state *AuthorizationModel) {
	listattr.ModifyMatching(h, &m.Roles, state.Roles)
	listattr.ModifyMatching(h, &m.Permissions, state.Permissions)
}
