package managementkey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
)

var ManagementKeyAttributes = map[string]schema.Attribute{
	"id":            stringattr.Identifier(),
	"name":          stringattr.Required(),
	"description":   stringattr.Default("", stringattr.StandardLenValidator),
	"status":        stringattr.Default("active", stringvalidator.OneOf("active", "inactive")),
	"expire_time":   intattr.Default(0, int64planmodifier.RequiresReplace()),
	"permitted_ips": strlistattr.Default(),
	"rebac":         objattr.Required[ReBacModel](ReBacAttributes, ReBacValidator, objectplanmodifier.RequiresReplace()),
	"cleartext":     stringattr.SecretComputed(),
}

type ManagementKeyModel struct {
	ID           stringattr.Type          `tfsdk:"id"`
	Name         stringattr.Type          `tfsdk:"name"`
	Description  stringattr.Type          `tfsdk:"description"`
	Status       stringattr.Type          `tfsdk:"status"`
	ExpireTime   intattr.Type             `tfsdk:"expire_time"`
	PermittedIPs strlistattr.Type         `tfsdk:"permitted_ips"`
	ReBac        objattr.Type[ReBacModel] `tfsdk:"rebac"`
	Cleartext    stringattr.Type          `tfsdk:"cleartext"`
}

func (m *ManagementKeyModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Status, data, "status")
	intattr.Get(m.ExpireTime, data, "expireTime")
	strlistattr.Get(m.PermittedIPs, data, "permittedIps", h)
	objattr.Get(m.ReBac, data, "reBac", h)
	if m.ID.ValueString() == "" && m.Status.ValueString() == "inactive" {
		h.Invalid("Cannot set status to inactive when creating a new management key")
	}
	return data
}

func (m *ManagementKeyModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Status, data, "status")
	intattr.Set(&m.ExpireTime, data, "expireTime")
	strlistattr.Set(&m.PermittedIPs, data, "permittedIps", h)
	objattr.Set(&m.ReBac, data, "reBac", h)
	stringattr.Set(&m.Cleartext, data, "cleartext")
}
