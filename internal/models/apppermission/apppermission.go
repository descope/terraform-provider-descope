package apppermission

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single permission scoped to a federated application. Permissions are referenced by id from `descope_app_role` resources.",
	Attributes:          AppPermissionAttributes,
}

var AppPermissionAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"project_id":  stringattr.Required(stringplanmodifier.RequiresReplace()),
	"app_id":      stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":        stringattr.Required(stringvalidator.LengthBetween(1, 100)),
	"description": stringattr.Default("", stringattr.StandardLenValidator),
}

type AppPermissionModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	ProjectID   stringattr.Type `tfsdk:"project_id"`
	AppID       stringattr.Type `tfsdk:"app_id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
}

func (m *AppPermissionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AppID, data, "appId")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	return data
}

func (m *AppPermissionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.AppID, data, "appId")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
}

func (m *AppPermissionModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AppPermissionModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *AppPermissionModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

func (m *AppPermissionModel) GetAppID() stringattr.Type {
	return m.AppID
}
