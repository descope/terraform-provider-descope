package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SegmentAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"write_key": stringattr.SecretRequired(),
	"host":      stringattr.Default(""),
}

// Model

type SegmentModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	WriteKey types.String `tfsdk:"write_key"`
	Host     types.String `tfsdk:"host"`
}

func (m *SegmentModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "segment"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SegmentModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SegmentModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.WriteKey, c, "writeKey")
	stringattr.Get(m.Host, c, "host")
	return c
}

// Matching

func (m *SegmentModel) GetName() types.String {
	return m.Name
}

func (m *SegmentModel) GetID() types.String {
	return m.ID
}

func (m *SegmentModel) SetID(id types.String) {
	m.ID = id
}
