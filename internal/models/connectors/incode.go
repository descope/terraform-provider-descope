package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var IncodeAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_key": stringattr.SecretRequired(),
	"api_url": stringattr.Required(),
	"flow_id": stringattr.Required(),
}

// Model

type IncodeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey types.String `tfsdk:"api_key"`
	ApiUrl types.String `tfsdk:"api_url"`
	FlowId types.String `tfsdk:"flow_id"`
}

func (m *IncodeModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "incode"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *IncodeModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *IncodeModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.ApiUrl, c, "apiUrl")
	stringattr.Get(m.FlowId, c, "flowId")
	return c
}

// Matching

func (m *IncodeModel) GetName() types.String {
	return m.Name
}

func (m *IncodeModel) GetID() types.String {
	return m.ID
}

func (m *IncodeModel) SetID(id types.String) {
	m.ID = id
}
