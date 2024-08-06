package authorization

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RoleAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"permissions": listattr.StringOptional(),
}

type RoleModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Permissions []string     `tfsdk:"permissions"`
}

func (m *RoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	data["permissions"] = m.Permissions

	// use the name as a lookup key to set the role reference or existing id
	roleName := m.Name.ValueString()
	if ref := h.Refs.Get(helpers.RoleReferenceKey, roleName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for role '%s' to: %s", roleName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown role reference", "No role named '"+roleName+"' was defined")
		data["id"] = m.ID.ValueString()
	}

	return data
}

func (m *RoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all roles values are specified in the configuration
}
