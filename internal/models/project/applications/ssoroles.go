package applications

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// emitSSOAppRoles only writes permissions/roles into the app payload when at
// least one is configured. The backend project-update endpoint does not yet
// accept these keys on every SSO app entry, so an unconditional empty list
// is rejected with a 400.
func emitSSOAppRoles(h *helpers.Handler, data map[string]any, permissions listattr.Type[SSOAppPermissionModel], roles listattr.Type[SSOAppRoleModel]) {
	hasPerms := false
	for range listattr.Iterator(permissions, h) {
		hasPerms = true
		break
	}
	hasRoles := false
	for range listattr.Iterator(roles, h) {
		hasRoles = true
		break
	}
	if !hasPerms && !hasRoles {
		return
	}
	listattr.Get(permissions, data, "permissions", h)
	listattr.Get(roles, data, "roles", h)
}

// Permission

var SSOAppPermissionAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description": stringattr.Optional(stringattr.StandardLenValidator),
}

type SSOAppPermissionModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
}

func (m *SSOAppPermissionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	return data
}

func (m *SSOAppPermissionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
}

func (m *SSOAppPermissionModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SSOAppPermissionModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SSOAppPermissionModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Role

var SSOAppRoleAttributes = map[string]schema.Attribute{
	"id":            stringattr.IdentifierMatched(),
	"name":          stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description":   stringattr.Optional(stringattr.StandardLenValidator),
	"permissions":   strsetattr.Optional(),
	"role_mappings": strsetattr.Optional(),
}

type SSOAppRoleModel struct {
	ID           stringattr.Type `tfsdk:"id"`
	Name         stringattr.Type `tfsdk:"name"`
	Description  stringattr.Type `tfsdk:"description"`
	Permissions  strsetattr.Type `tfsdk:"permissions"`
	RoleMappings strsetattr.Type `tfsdk:"role_mappings"`
}

func (m *SSOAppRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")

	// snapshot stores role.permissions as a list of permission objects keyed by name;
	// the backend resolves names to ids on import. sort for deterministic payload order
	permNames := []string{}
	for name := range strsetattr.Iterator(m.Permissions, h) {
		permNames = append(permNames, name)
	}
	slices.Sort(permNames)
	perms := make([]any, 0, len(permNames))
	for _, name := range permNames {
		perms = append(perms, map[string]any{"name": name})
	}
	data["permissions"] = perms

	// roleMappings are project role identifiers; pass them through as raw strings
	mappingNames := []string{}
	for v := range strsetattr.Iterator(m.RoleMappings, h) {
		mappingNames = append(mappingNames, v)
	}
	slices.Sort(mappingNames)
	mappings := make([]any, 0, len(mappingNames))
	for _, v := range mappingNames {
		mappings = append(mappings, v)
	}
	data["roleMappings"] = mappings

	return data
}

func (m *SSOAppRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")

	if raw, ok := data["permissions"].([]any); ok {
		names := make([]any, 0, len(raw))
		for _, item := range raw {
			if perm, ok := item.(map[string]any); ok {
				if name, ok := perm["name"].(string); ok && name != "" {
					names = append(names, name)
				}
			}
		}
		strsetattr.Set(&m.Permissions, map[string]any{"permissions": names}, "permissions", h)
	}

	if raw, ok := data["roleMappings"].([]any); ok {
		values := make([]any, 0, len(raw))
		for _, item := range raw {
			if v, ok := item.(string); ok && v != "" {
				values = append(values, v)
			}
		}
		strsetattr.Set(&m.RoleMappings, map[string]any{"roleMappings": values}, "roleMappings", h)
	}
}

func (m *SSOAppRoleModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SSOAppRoleModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SSOAppRoleModel) SetID(id stringattr.Type) {
	m.ID = id
}
