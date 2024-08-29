package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TelesignAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"customer_id": stringattr.Required(),
	"api_key":     stringattr.SecretRequired(),
}

// Model

type TelesignModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	CustomerID types.String `tfsdk:"customer_id"`
	APIKey     types.String `tfsdk:"api_key"`
}

func (m *TelesignModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "telesign"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *TelesignModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *TelesignModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.CustomerID, c, "customerID")
	stringattr.Get(m.APIKey, c, "apiKey")
	return c
}

// Matching

func (m *TelesignModel) GetName() types.String {
	return m.Name
}

func (m *TelesignModel) GetID() types.String {
	return m.ID
}

func (m *TelesignModel) SetID(id types.String) {
	m.ID = id
}
