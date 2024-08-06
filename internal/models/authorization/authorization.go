package authorization

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AuthorizationValidator = objectattr.NewValidator[AuthorizationModel]("must have unique role and permission names")

var AuthorizationAttributes = map[string]schema.Attribute{
	"roles":       listattr.Optional(RoleAttributes),
	"permissions": listattr.Optional(PermissionAttributes),
}

type AuthorizationModel struct {
	Roles       []*RoleModel       `tfsdk:"roles"`
	Permissions []*PermissionModel `tfsdk:"permissions"`
}

func (m *AuthorizationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Roles, data, "roles", h)
	listattr.Get(m.Permissions, data, "permissions", h)
	return data
}

func (m *AuthorizationModel) SetValues(h *helpers.Handler, data map[string]any) {
	roles, permissions := m.getAuthorizationIDs(data)

	for _, role := range m.Roles {
		name := role.Name.ValueString()
		id, found := roles[name]
		if found {
			value := types.StringValue(id)
			if !role.ID.Equal(value) {
				h.Log("Setting new ID '" + id + "' for role named '" + name + "'")
				role.ID = value
			} else {
				h.Log("Keeping existing ID '" + id + "' for role named '" + name + "'")
			}
		} else {
			h.Error("Role not found", "Expected to find role to match with '"+name+"'")
		}
	}

	for _, permission := range m.Permissions {
		name := permission.Name.ValueString()
		id, found := permissions[name]
		if found {
			value := types.StringValue(id)
			if !permission.ID.Equal(value) {
				h.Log("Setting new ID '" + id + "' for permission named '" + name + "'")
				permission.ID = value
			} else {
				h.Log("Keeping existing ID '" + id + "' for permission named '" + name + "'")
			}
		} else {
			h.Error("Permission not found", "Expected to find permission to match with '"+name+"'")
		}
	}
}

func (m *AuthorizationModel) References(ctx context.Context) helpers.ReferencesMap {
	refs := helpers.ReferencesMap{}
	for _, v := range m.Roles {
		refs.Add(helpers.RoleReferenceKey, "", v.ID.ValueString(), v.Name.ValueString())
	}
	return refs
}

func (m *AuthorizationModel) getAuthorizationIDs(data map[string]any) (roles, permissions map[string]string) {
	roles = map[string]string{}
	permissions = map[string]string{}

	ps, _ := data["permissions"].([]any)
	for _, v := range ps {
		if p, ok := v.(map[string]any); ok {
			id, _ := p["id"].(string)
			name, _ := p["name"].(string)
			permissions[name] = id
		}
	}

	rs, _ := data["roles"].([]any)
	for _, v := range rs {
		if r, ok := v.(map[string]any); ok {
			id, _ := r["id"].(string)
			name, _ := r["name"].(string)
			roles[name] = id
		}
	}

	return
}

func (m *AuthorizationModel) Validate(h *helpers.Handler) {
	permissions := map[string]int{}
	for _, v := range m.Permissions {
		permissions[v.Name.ValueString()] += 1
	}
	roles := map[string]int{}
	for _, v := range m.Roles {
		roles[v.Name.ValueString()] += 1
		for _, p := range v.Permissions {
			if value := permissions[p]; value == 0 {
				h.Error("Permission doesn't exist", "The role '%s' references a permission '%s' that doesn't exist", v.Name.ValueString(), p)
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
