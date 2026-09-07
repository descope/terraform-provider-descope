package accesskey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var AccessKeyAttributes = map[string]schema.Attribute{
	"id":                stringattr.Identifier(),
	"project_id":        stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":              stringattr.Required(),
	"description":       stringattr.Default("", stringattr.StandardLenValidator),
	"status":            stringattr.Default("active", stringvalidator.OneOf("active", "inactive")),
	"expire_time":       intattr.Default(0, int64planmodifier.RequiresReplace()),
	"bound_user_id":     stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"roles":             strlistattr.Default(stringattr.NonEmptyValidator),
	"tenants":           listattr.Default[AccessKeyTenantModel](AccessKeyTenantAttributes),
	"custom_claims":     jsonattr.Default("{}"),
	"custom_attributes": jsonattr.Default("{}"),
	"permitted_ips":     strlistattr.Default(),
	"client_id":         stringattr.Identifier(),
	"created_time":      intattr.Generated(),
	"created_by":        stringattr.Generated(),
	"cleartext":         stringattr.SecretGenerated(false),
}

var Schema = schema.Schema{
	Attributes: AccessKeyAttributes,
}

type AccessKeyModel struct {
	ID               stringattr.Type                     `tfsdk:"id"`
	ProjectID        stringattr.Type                     `tfsdk:"project_id"`
	Name             stringattr.Type                     `tfsdk:"name"`
	Description      stringattr.Type                     `tfsdk:"description"`
	Status           stringattr.Type                     `tfsdk:"status"`
	ExpireTime       intattr.Type                        `tfsdk:"expire_time"`
	BoundUserID      stringattr.Type                     `tfsdk:"bound_user_id"`
	Roles            strlistattr.Type                    `tfsdk:"roles"`
	Tenants          listattr.Type[AccessKeyTenantModel] `tfsdk:"tenants"`
	CustomClaims     jsonattr.Type                       `tfsdk:"custom_claims"`
	CustomAttributes jsonattr.Type                       `tfsdk:"custom_attributes"`
	PermittedIPs     strlistattr.Type                    `tfsdk:"permitted_ips"`
	ClientID         stringattr.Type                     `tfsdk:"client_id"`
	CreatedTime      intattr.Type                        `tfsdk:"created_time"`
	CreatedBy        stringattr.Type                     `tfsdk:"created_by"`
	Cleartext        stringattr.Type                     `tfsdk:"cleartext"`
}

func (m *AccessKeyModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Status, data, "status")
	intattr.Get(m.ExpireTime, data, "expireTime")
	stringattr.Get(m.BoundUserID, data, "boundUserId")
	strlistattr.Get(m.Roles, data, "roleNames", h)
	listattr.Get(m.Tenants, data, "keyTenants", h)
	jsonattr.Get(m.CustomClaims, data, "customClaims")
	jsonattr.Get(m.CustomAttributes, data, "customAttributes")
	strlistattr.Get(m.PermittedIPs, data, "permittedIps", h)

	if m.ID.ValueString() == "" && m.Status.ValueString() == "inactive" {
		h.Invalid("Cannot set status to inactive when creating a new access key")
	}

	if !helpers.HasUnknownValues(m.Roles, m.Tenants) {
		if !m.Roles.IsEmpty() && !m.Tenants.IsEmpty() {
			h.Conflict("The roles attribute cannot be set when tenants is set; specify the roles within each tenant object instead")
		}
	}

	return data
}

func (m *AccessKeyModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Status, data, "status")
	intattr.Set(&m.ExpireTime, data, "expireTime")
	stringattr.Set(&m.BoundUserID, data, "boundUserId")
	strlistattr.Set(&m.Roles, data, "roleNames", h)
	listattr.Set(&m.Tenants, data, "keyTenants", h)
	jsonattr.Set(&m.CustomClaims, data, "customClaims", jsonattr.SkipIfAlreadySet)
	jsonattr.Set(&m.CustomAttributes, data, "customAttributes", jsonattr.SkipIfAlreadySet)
	strlistattr.Set(&m.PermittedIPs, data, "permittedIps", h)
	stringattr.Set(&m.ClientID, data, "clientId")
	intattr.Set(&m.CreatedTime, data, "createdTime")
	stringattr.Set(&m.CreatedBy, data, "createdBy")
	stringattr.Set(&m.Cleartext, data, "cleartext")
}

func (m *AccessKeyModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AccessKeyModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *AccessKeyModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
