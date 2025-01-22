package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SalesforceMarketingCloudAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"subdomain":     stringattr.Required(),
	"client_id":     stringattr.Required(),
	"client_secret": stringattr.SecretRequired(),
	"scope":         stringattr.Default(""),
	"account_id":    stringattr.Default(""),
}

// Model

type SalesforceMarketingCloudModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Subdomain    types.String `tfsdk:"subdomain"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Scope        types.String `tfsdk:"scope"`
	AccountId    types.String `tfsdk:"account_id"`
}

func (m *SalesforceMarketingCloudModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "salesforce-marketing-cloud"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SalesforceMarketingCloudModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SalesforceMarketingCloudModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.Subdomain, c, "subdomain")
	stringattr.Get(m.ClientID, c, "clientId")
	stringattr.Get(m.ClientSecret, c, "clientSecret")
	stringattr.Get(m.Scope, c, "scope")
	stringattr.Get(m.AccountId, c, "accountId")
	return c
}

// Matching

func (m *SalesforceMarketingCloudModel) GetName() types.String {
	return m.Name
}

func (m *SalesforceMarketingCloudModel) GetID() types.String {
	return m.ID
}

func (m *SalesforceMarketingCloudModel) SetID(id types.String) {
	m.ID = id
}
