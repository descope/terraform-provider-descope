package accesskey

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
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
	"custom_claims":     stringattr.Default("{}", stringattr.JSONValidator()),
	"custom_attributes": stringattr.Default("{}", stringattr.JSONValidator()),
	"permitted_ips":     strlistattr.Default(),
	"client_id":         stringattr.Identifier(),
	"created_time":      intattr.Generated(),
	"created_by":        stringattr.Generated(),
	"cleartext":         stringattr.SecretGenerated(false),
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
	CustomClaims     stringattr.Type                     `tfsdk:"custom_claims"`
	CustomAttributes stringattr.Type                     `tfsdk:"custom_attributes"`
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
	getJSONField(m.CustomClaims, data, "customClaims")
	getJSONField(m.CustomAttributes, data, "customAttributes")
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
	setJSONField(&m.CustomClaims, data, "customClaims")
	setJSONField(&m.CustomAttributes, data, "customAttributes")
	strlistattr.Set(&m.PermittedIPs, data, "permittedIps", h)
	stringattr.Set(&m.ClientID, data, "clientId")
	intattr.Set(&m.CreatedTime, data, "createdTime")
	stringattr.Set(&m.CreatedBy, data, "createdBy")
	stringattr.Set(&m.Cleartext, data, "cleartext")
}

func getJSONField(s stringattr.Type, data map[string]any, key string) {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(s.ValueString()), &m); err != nil {
		panic("Invalid JSON data after validation: " + err.Error())
	}
	data[key] = m
}

func setJSONField(s *stringattr.Type, data map[string]any, key string) {
	// We do not currently update the field data if it's already set because it might be slightly different after apply
	if s.ValueString() == "" {
		value := "{}"
		if v, ok := data[key].(map[string]any); ok {
			if b, err := json.MarshalIndent(v, "", "  "); err == nil {
				value = string(b)
			}
		}
		*s = stringattr.Value(value)
	}
}
