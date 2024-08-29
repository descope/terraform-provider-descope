package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GoogleCloudTranslationAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"project_id":           stringattr.Required(),
	"service_account_json": stringattr.SecretRequired(),
}

// Model

type GoogleCloudTranslationModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	ProjectID          types.String `tfsdk:"project_id"`
	ServiceAccountJSON types.String `tfsdk:"service_account_json"`
}

func (m *GoogleCloudTranslationModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "google-cloud-translation"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *GoogleCloudTranslationModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *GoogleCloudTranslationModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.ProjectID, c, "projectId")
	stringattr.Get(m.ServiceAccountJSON, c, "serviceAccountJSON")
	return c
}

// Matching

func (m *GoogleCloudTranslationModel) GetName() types.String {
	return m.Name
}

func (m *GoogleCloudTranslationModel) GetID() types.String {
	return m.ID
}

func (m *GoogleCloudTranslationModel) SetID(id types.String) {
	m.ID = id
}
