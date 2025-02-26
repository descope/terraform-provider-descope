package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var LokaliseAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_token":            stringattr.SecretRequired(),
	"project_id":           stringattr.Required(),
	"team_id":              stringattr.Default(""),
	"card_id":              stringattr.Default(""),
	"translation_provider": stringattr.Default(""),
}

// Model

type LokaliseModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIToken            types.String `tfsdk:"api_token"`
	ProjectID           types.String `tfsdk:"project_id"`
	TeamID              types.String `tfsdk:"team_id"`
	CardID              types.String `tfsdk:"card_id"`
	TranslationProvider types.String `tfsdk:"translation_provider"`
}

func (m *LokaliseModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "lokalise"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *LokaliseModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.APIToken, c, "apiToken")
		stringattr.Set(&m.ProjectID, c, "projectId")
		stringattr.Set(&m.TeamID, c, "teamId")
		stringattr.Set(&m.CardID, c, "cardId")
		stringattr.Set(&m.TranslationProvider, c, "translationProvider")
	}
}

// Configuration

func (m *LokaliseModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIToken, c, "apiToken")
	stringattr.Get(m.ProjectID, c, "projectId")
	stringattr.Get(m.TeamID, c, "teamId")
	stringattr.Get(m.CardID, c, "cardId")
	stringattr.Get(m.TranslationProvider, c, "translationProvider")
	return c
}

// Matching

func (m *LokaliseModel) GetName() types.String {
	return m.Name
}

func (m *LokaliseModel) GetID() types.String {
	return m.ID
}

func (m *LokaliseModel) SetID(id types.String) {
	m.ID = id
}
