package accesskey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var AccessKeyAttributes = map[string]schema.Attribute{
	"id":            stringattr.Identifier(),
	"name":          stringattr.Required(),
	"description":   stringattr.Default("", stringattr.StandardLenValidator),
	"status":        stringattr.Default("active", stringvalidator.OneOf("active", "inactive")),
	"expire_time":   intattr.Default(0, int64planmodifier.RequiresReplace()),
	"role_names":    strsetattr.Default(),
	"key_tenants":   listattr.Default[TenantModel](TenantAttributes),
	"permitted_ips": strlistattr.Default(),
	"user_id": schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Default:       stringdefault.StaticString(""),
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	},
	"cleartext":     stringattr.SecretComputed(),
	"client_id":     stringattr.Identifier(),
	"created_time":  intattr.Default(0, int64planmodifier.UseStateForUnknown()),
	"created_by":    stringattr.Identifier(),
}

type AccessKeyModel struct {
	ID           stringattr.Type                `tfsdk:"id"`
	Name         stringattr.Type                `tfsdk:"name"`
	Description  stringattr.Type                `tfsdk:"description"`
	Status       stringattr.Type                `tfsdk:"status"`
	ExpireTime   intattr.Type                   `tfsdk:"expire_time"`
	RoleNames    strsetattr.Type                `tfsdk:"role_names"`
	KeyTenants   listattr.Type[TenantModel]     `tfsdk:"key_tenants"`
	PermittedIPs strlistattr.Type               `tfsdk:"permitted_ips"`
	UserID       stringattr.Type                `tfsdk:"user_id"`
	Cleartext    stringattr.Type                `tfsdk:"cleartext"`
	ClientID     stringattr.Type                `tfsdk:"client_id"`
	CreatedTime  intattr.Type                   `tfsdk:"created_time"`
	CreatedBy    stringattr.Type                `tfsdk:"created_by"`
}

func (m *AccessKeyModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Status, data, "status")
	intattr.Get(m.ExpireTime, data, "expireTime")
	strsetattr.Get(m.RoleNames, data, "roleNames", h)
	listattr.Get(m.KeyTenants, data, "keyTenants", h)
	strlistattr.Get(m.PermittedIPs, data, "permittedIps", h)
	stringattr.Get(m.UserID, data, "userId")
	if m.ID.ValueString() == "" && m.Status.ValueString() == "inactive" {
		h.Invalid("Cannot set status to inactive when creating a new access key")
	}
	return data
}

func (m *AccessKeyModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Status, data, "status")
	intattr.Set(&m.ExpireTime, data, "expireTime")
	strsetattr.Set(&m.RoleNames, data, "roleNames", h)
	listattr.Set(&m.KeyTenants, data, "keyTenants", h)
	strlistattr.Set(&m.PermittedIPs, data, "permittedIps", h)
	stringattr.Set(&m.UserID, data, "userId")
	stringattr.Set(&m.Cleartext, data, "cleartext")
	stringattr.Set(&m.ClientID, data, "clientId")
	intattr.Set(&m.CreatedTime, data, "createdTime")
	stringattr.Set(&m.CreatedBy, data, "createdBy")
}
