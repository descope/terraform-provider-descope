package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SNSAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key_id":       stringattr.Required(),
	"secret":              stringattr.SecretRequired(),
	"region":              stringattr.Required(),
	"endpoint":            stringattr.Default(""),
	"organization_number": stringattr.Default(""),
	"sender_id":           stringattr.Default(""),
	"entity_id":           stringattr.Default(""),
	"template_id":         stringattr.Default(""),
}

// Model

type SNSModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyId        types.String `tfsdk:"access_key_id"`
	Secret             types.String `tfsdk:"secret"`
	Region             types.String `tfsdk:"region"`
	Endpoint           types.String `tfsdk:"endpoint"`
	OrganizationNumber types.String `tfsdk:"organization_number"`
	SenderID           types.String `tfsdk:"sender_id"`
	EntityID           types.String `tfsdk:"entity_id"`
	TemplateID         types.String `tfsdk:"template_id"`
}

func (m *SNSModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "sns"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SNSModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SNSModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.Region, c, "awsSNSRegion")
	stringattr.Get(m.Endpoint, c, "awsEndpoint")
	stringattr.Get(m.OrganizationNumber, c, "originationNumber")
	stringattr.Get(m.SenderID, c, "senderId")
	stringattr.Get(m.EntityID, c, "entityId")
	stringattr.Get(m.TemplateID, c, "templateId")
	return c
}

// Matching

func (m *SNSModel) GetName() types.String {
	return m.Name
}

func (m *SNSModel) GetID() types.String {
	return m.ID
}

func (m *SNSModel) SetID(id types.String) {
	m.ID = id
}
