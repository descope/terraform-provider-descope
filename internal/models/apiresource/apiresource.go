package apiresource

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages an API resource in a Descope project: an external API or MCP server that inbound applications request access to, protected by another authorization server.",
	Attributes:          APIResourceAttributes,
}

var accessValidator = stringvalidator.OneOf("", "allow", "deny", "allowViaPolicy")

var APIResourceAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"project_id":         stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":               stringattr.Required(stringattr.StandardLenValidator),
	"description":        stringattr.Default("", stringvalidator.LengthAtMost(1000)),
	"uri":                stringattr.Required(stringvalidator.LengthAtMost(3000)),
	"user_access":        stringattr.Default("", accessValidator),
	"client_access":      stringattr.Default("", accessValidator),
	"permissions_scopes": listattr.Default[inboundapp.ApplicationScopeModel](inboundapp.ApplicationScopeAttributes),
	"attributes_scopes":  listattr.Default[inboundapp.ApplicationScopeModel](inboundapp.ApplicationScopeAttributes),
	"connections_scopes": listattr.Default[inboundapp.ApplicationScopeModel](inboundapp.ApplicationScopeAttributes),
	"id_jag_target":      objattr.Default[IDJAGTargetModel](nil, IDJAGTargetAttributes),
}

type APIResourceModel struct {
	ID                stringattr.Type                                 `tfsdk:"id"`
	ProjectID         stringattr.Type                                 `tfsdk:"project_id"`
	Name              stringattr.Type                                 `tfsdk:"name"`
	Description       stringattr.Type                                 `tfsdk:"description"`
	URI               stringattr.Type                                 `tfsdk:"uri"`
	UserAccess        stringattr.Type                                 `tfsdk:"user_access"`
	ClientAccess      stringattr.Type                                 `tfsdk:"client_access"`
	PermissionsScopes listattr.Type[inboundapp.ApplicationScopeModel] `tfsdk:"permissions_scopes"`
	AttributesScopes  listattr.Type[inboundapp.ApplicationScopeModel] `tfsdk:"attributes_scopes"`
	ConnectionsScopes listattr.Type[inboundapp.ApplicationScopeModel] `tfsdk:"connections_scopes"`
	IDJAGTarget       objattr.Type[IDJAGTargetModel]                  `tfsdk:"id_jag_target"`
}

func (m *APIResourceModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{"type": "api"}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.URI, data, "uri")
	stringattr.Get(m.UserAccess, data, "userAccess")
	stringattr.Get(m.ClientAccess, data, "clientAccess")
	scopes := map[string]any{}
	listattr.Get(m.PermissionsScopes, scopes, "permissionsScopes", h)
	listattr.Get(m.AttributesScopes, scopes, "attributesScopes", h)
	listattr.Get(m.ConnectionsScopes, scopes, "connectionsScopes", h)
	data["scopes"] = scopes
	objattr.Get(m.IDJAGTarget, data, "idJagTarget", h)
	return data
}

func (m *APIResourceModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.URI, data, "uri")
	stringattr.Set(&m.UserAccess, data, "userAccess")
	stringattr.Set(&m.ClientAccess, data, "clientAccess")
	scopes, _ := data["scopes"].(map[string]any)
	listattr.Set(&m.PermissionsScopes, scopes, "permissionsScopes", h)
	listattr.Set(&m.AttributesScopes, scopes, "attributesScopes", h)
	listattr.Set(&m.ConnectionsScopes, scopes, "connectionsScopes", h)
	objattr.Set(&m.IDJAGTarget, data, "idJagTarget", h)
}

func (m *APIResourceModel) GetID() stringattr.Type {
	return m.ID
}

func (m *APIResourceModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *APIResourceModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
