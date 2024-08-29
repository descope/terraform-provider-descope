package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var HIBPAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
}

// Model

type HIBPModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (m *HIBPModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "hibp"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *HIBPModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *HIBPModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	return c
}

// Matching

func (m *HIBPModel) GetName() types.String {
	return m.Name
}

func (m *HIBPModel) GetID() types.String {
	return m.ID
}

func (m *HIBPModel) SetID(id types.String) {
	m.ID = id
}
