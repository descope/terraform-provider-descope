package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ForterValidator = objectattr.NewValidator[ForterModel]("must have a valid configuration")

var ForterAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"site_id":             stringattr.Required(),
	"secret_key":          stringattr.SecretRequired(),
	"overrides":           boolattr.Default(false),
	"override_ip_address": stringattr.Default(""),
	"override_user_email": stringattr.Default(""),
}

// Model

type ForterModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	SiteID            types.String `tfsdk:"site_id"`
	SecretKey         types.String `tfsdk:"secret_key"`
	Overrides         types.Bool   `tfsdk:"overrides"`
	OverrideIPAddress types.String `tfsdk:"override_ip_address"`
	OverrideUserEmail types.String `tfsdk:"override_user_email"`
}

func (m *ForterModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "forter"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *ForterModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.SiteID, c, "siteId")
		stringattr.Set(&m.SecretKey, c, "secretKey")
		boolattr.Set(&m.Overrides, c, "overrides")
		stringattr.Set(&m.OverrideIPAddress, c, "overrideIpAddress")
		stringattr.Set(&m.OverrideUserEmail, c, "overrideUserEmail")
	}
}

func (m *ForterModel) Validate(h *helpers.Handler) {
	if !m.OverrideIPAddress.IsNull() && !m.Overrides.ValueBool() {
		h.Error("Invalid connector configuration", "The override_ip_address field cannot be used unless overrides is set to true")
	}
	if !m.OverrideUserEmail.IsNull() && !m.Overrides.ValueBool() {
		h.Error("Invalid connector configuration", "The override_user_email field cannot be used unless overrides is set to true")
	}
}

// Configuration

func (m *ForterModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.SiteID, c, "siteId")
	stringattr.Get(m.SecretKey, c, "secretKey")
	boolattr.Get(m.Overrides, c, "overrides")
	stringattr.Get(m.OverrideIPAddress, c, "overrideIpAddress")
	stringattr.Get(m.OverrideUserEmail, c, "overrideUserEmail")
	return c
}

// Matching

func (m *ForterModel) GetName() types.String {
	return m.Name
}

func (m *ForterModel) GetID() types.String {
	return m.ID
}

func (m *ForterModel) SetID(id types.String) {
	m.ID = id
}
