package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ClearAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"project_id": stringattr.Required(),
	"api_key":    stringattr.SecretRequired(),
}

// Model

type ClearModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	ProjectID types.String `tfsdk:"project_id"`
	APIKey    types.String `tfsdk:"api_key"`
}

func (m *ClearModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "clear"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *ClearModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *ClearModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.ProjectID, c, "projectId")
	stringattr.Get(m.APIKey, c, "apiKey")
	return c
}

// Matching

func (m *ClearModel) GetName() types.String {
	return m.Name
}

func (m *ClearModel) GetID() types.String {
	return m.ID
}

func (m *ClearModel) SetID(id types.String) {
	m.ID = id
}
