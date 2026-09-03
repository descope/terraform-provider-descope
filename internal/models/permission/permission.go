package permission

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var ReservedNames = []string{"Impersonate", "SSO Admin", "Super User", "User Admin"}

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single custom permission in a Descope project. Permissions are referenced by name from `descope_role` resources.",
	Attributes:          PermissionAttributes,
}

var PermissionAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"project_id":  stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":        stringattr.Required(stringvalidator.LengthBetween(1, 100), stringvalidator.NoneOf(ReservedNames...)),
	"description": stringattr.Default("", stringvalidator.LengthAtMost(1024)),
}

type PermissionModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	ProjectID   stringattr.Type `tfsdk:"project_id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
}

func (m *PermissionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	return data
}

func (m *PermissionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
}

func (m *PermissionModel) GetID() stringattr.Type {
	return m.ID
}

func (m *PermissionModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *PermissionModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
