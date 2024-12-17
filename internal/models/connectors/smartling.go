package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SmartlingAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"user_identifier": stringattr.Required(),
	"user_secret":     stringattr.SecretRequired(),
	"account_uid":     stringattr.Required(),
}

// Model

type SmartlingModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	UserIdentifier types.String `tfsdk:"user_identifier"`
	UserSecret     types.String `tfsdk:"user_secret"`
	AccountUID     types.String `tfsdk:"account_uid"`
}

func (m *SmartlingModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "smartling"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SmartlingModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SmartlingModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.UserIdentifier, c, "userIdentifier")
	stringattr.Get(m.UserSecret, c, "userSecret")
	stringattr.Get(m.AccountUID, c, "accountUid")
	return c
}

// Matching

func (m *SmartlingModel) GetName() types.String {
	return m.Name
}

func (m *SmartlingModel) GetID() types.String {
	return m.ID
}

func (m *SmartlingModel) SetID(id types.String) {
	m.ID = id
}
