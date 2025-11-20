package authorization

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/iancoleman/strcase"
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
	listattr.SetMatchingNames(&m.Roles, data, "roles", "name", h)
	listattr.SetMatchingNames(&m.Permissions, data, "permissions", "name", h)
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
	roleNames := map[string]int{}
	roleKeys := map[string]int{}

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
		roleNames[name] += 1

		if key := r.Key.ValueString(); key != "" {
			roleKeys[key] += 1
		}

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

	for k, v := range roleNames {
		if v > 1 {
			h.Error("Role names must be unique", "The role name '%s' is used %d times", k, v)
		}
	}

	for k, v := range roleKeys {
		if v > 1 {
			h.Error("Role keys must be unique", "The role key '%s' is used %d times", k, v)
		}
	}

	if len(roleKeys) > 0 && len(roleKeys) != len(roleNames) {
		h.Missing("The 'key' attribute must be set for all roles")
	}
}

func (m *AuthorizationModel) Modify(h *helpers.Handler, state *AuthorizationModel) {
	// we use this chance to warn the user about missing role keys
	planKeys := map[string]string{}
	for p := range listattr.MutatingIterator(&m.Roles, h) {
		name := p.Name.ValueString()
		if p.Key.ValueString() == "" {
			h.Warn("Missing Key Attribute In "+name+" Role", "The role '%s' is missing a value for the 'key' attribute. It's strongly recommended to set a unique value (e.g., '%s') as the value of the 'key' attribute in the Terraform plan to ensure user roles are maintained correctly in future plan changes. This will become an error in a future version of the provider.", name, strcase.ToSnake(name))
		} else {
			// keep the key->name mapping in the plan to compare against the state below
			planKeys[name] = p.Key.ValueString()
		}
	}

	// try to warn a about accidental role key changes
	for s := range listattr.Iterator(state.Roles, h) {
		name := s.Name.ValueString()
		stateKey := s.Key.ValueString()
		if planKey, ok := planKeys[name]; ok && stateKey != "" && stateKey != planKey {
			h.Warn("Role Key Modified", "The key for role '%s' has been modified in the plan from '%s' to '%s'. This may lead to unintended changes to user roles.", name, stateKey, planKey)
		}
	}

	listattr.ModifyMatchingKeysOrNames(h, &m.Roles, state.Roles)
	listattr.ModifyMatchingNames(h, &m.Permissions, state.Permissions)
}
