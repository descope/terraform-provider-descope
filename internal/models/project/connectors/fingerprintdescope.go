package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FingerprintDescopeAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"custom_domain": stringattr.Default(""),
}

// Model

type FingerprintDescopeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	CustomDomain types.String `tfsdk:"custom_domain"`
}

func (m *FingerprintDescopeModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "fingerprint-descope"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *FingerprintDescopeModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.CustomDomain, c, "customDomain")
	}
}

// Configuration

func (m *FingerprintDescopeModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.CustomDomain, c, "customDomain")
	return c
}

// Matching

func (m *FingerprintDescopeModel) GetName() types.String {
	return m.Name
}

func (m *FingerprintDescopeModel) GetID() types.String {
	return m.ID
}

func (m *FingerprintDescopeModel) SetID(id types.String) {
	m.ID = id
}
