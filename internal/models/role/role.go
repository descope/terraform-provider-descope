package role

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single role in a Descope project. Roles reference permissions by name, and the referenced permissions must exist when the role is applied.",
	Attributes:          RoleAttributes,
}

var RoleAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"project_id":  stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":        stringattr.Required(stringvalidator.LengthBetween(1, 100)),
	"description": stringattr.Default("", stringvalidator.LengthAtMost(1024)),
	"permissions": strsetattr.Default(stringvalidator.LengthBetween(1, 100)),
	"default":     boolattr.Default(false),
	"private":     boolattr.Default(false),
}

type RoleModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	ProjectID   stringattr.Type `tfsdk:"project_id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Permissions strsetattr.Type `tfsdk:"permissions"`
	Default     boolattr.Type   `tfsdk:"default"`
	Private     boolattr.Type   `tfsdk:"private"`
}

func (m *RoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	strsetattr.Get(m.Permissions, data, "permissionNames", h)
	boolattr.Get(m.Default, data, "default")
	boolattr.Get(m.Private, data, "private")
	return data
}

func (m *RoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	strsetattr.Set(&m.Permissions, data, "permissionNames", h)
	boolattr.Set(&m.Default, data, "default")
	boolattr.Set(&m.Private, data, "private")
}

func (m *RoleModel) GetID() stringattr.Type {
	return m.ID
}

func (m *RoleModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *RoleModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
