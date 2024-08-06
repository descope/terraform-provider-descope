package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RekognitionAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"access_key_id":     stringattr.Required(),
	"secret_access_key": stringattr.SecretRequired(),
	"collection_id":     stringattr.Required(),
}

// Model

type RekognitionModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	CollectionID    types.String `tfsdk:"collection_id"`
}

func (m *RekognitionModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "rekognition"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *RekognitionModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *RekognitionModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyID, c, "accessKeyId")
	stringattr.Get(m.SecretAccessKey, c, "secretAccessKey")
	stringattr.Get(m.CollectionID, c, "collectionId")
	return c
}

// Matching

func (m *RekognitionModel) GetName() types.String {
	return m.Name
}

func (m *RekognitionModel) GetID() types.String {
	return m.ID
}

func (m *RekognitionModel) SetID(id types.String) {
	m.ID = id
}
