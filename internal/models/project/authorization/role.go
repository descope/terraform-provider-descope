package authorization

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var RoleAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description": stringattr.Optional(stringattr.StandardLenValidator),
	"permissions": strsetattr.Optional(),
}

type RoleModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Permissions strsetattr.Type `tfsdk:"permissions"`
}

func (m *RoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	strsetattr.Get(m.Permissions, data, "permissions", h)

	// use the name as a lookup key to set the role reference or existing id
	roleName := m.Name.ValueString()
	if ref := h.Refs.Get(helpers.RoleReferenceKey, roleName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for role '%s' to: %s", roleName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown role reference", "No role named '%s' was defined", roleName)
	}

	return data
}

func (m *RoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	strsetattr.Set(&m.Permissions, data, "permissions", h)
}

// Matching

func (m *RoleModel) GetName() stringattr.Type {
	return m.Name
}

func (m *RoleModel) GetID() stringattr.Type {
	return m.ID
}

func (m *RoleModel) SetID(id stringattr.Type) {
	m.ID = id
}
