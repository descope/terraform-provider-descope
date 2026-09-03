package approle

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single role scoped to a federated application. Roles reference application permissions by id via `permission_ids`, and project-level roles by id via `role_mappings`.",
	Attributes:          AppRoleAttributes,
}

var AppRoleAttributes = map[string]schema.Attribute{
	"id":             stringattr.Identifier(),
	"project_id":     stringattr.Required(stringplanmodifier.RequiresReplace()),
	"app_id":         stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":           stringattr.Required(stringvalidator.LengthBetween(1, 100)),
	"description":    stringattr.Default("", stringattr.StandardLenValidator),
	"permission_ids": strsetattr.Default(),
	"role_mappings":  strsetattr.Default(),
}

type AppRoleModel struct {
	ID            stringattr.Type `tfsdk:"id"`
	ProjectID     stringattr.Type `tfsdk:"project_id"`
	AppID         stringattr.Type `tfsdk:"app_id"`
	Name          stringattr.Type `tfsdk:"name"`
	Description   stringattr.Type `tfsdk:"description"`
	PermissionIDs strsetattr.Type `tfsdk:"permission_ids"`
	RoleMappings  strsetattr.Type `tfsdk:"role_mappings"`
}

func (m *AppRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AppID, data, "appId")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	strsetattr.Get(m.PermissionIDs, data, "permissionIds", h)
	strsetattr.Get(m.RoleMappings, data, "roleMappings", h)
	return data
}

func (m *AppRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.AppID, data, "appId")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	strsetattr.Set(&m.PermissionIDs, data, "permissionIds", h)
	strsetattr.Set(&m.RoleMappings, data, "roleMappings", h)
}

func (m *AppRoleModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AppRoleModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *AppRoleModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

func (m *AppRoleModel) GetAppID() stringattr.Type {
	return m.AppID
}
