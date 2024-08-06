package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TraceableAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"secret_key": stringattr.SecretRequired(),
	"eu_region":  boolattr.Default(false),
}

// Model

type TraceableModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	SecretKey types.String `tfsdk:"secret_key"`
	EURegion  types.Bool   `tfsdk:"eu_region"`
}

func (m *TraceableModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "traceable"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *TraceableModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *TraceableModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.SecretKey, c, "secretKey")
	boolattr.Get(m.EURegion, c, "euRegion")
	return c
}

// Matching

func (m *TraceableModel) GetName() types.String {
	return m.Name
}

func (m *TraceableModel) GetID() types.String {
	return m.ID
}

func (m *TraceableModel) SetID(id types.String) {
	m.ID = id
}
