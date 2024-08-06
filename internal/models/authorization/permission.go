package authorization

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var PermissionAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
}

type PermissionModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (m *PermissionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	return data
}

func (m *PermissionModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all permission values are specified in the configuration
}