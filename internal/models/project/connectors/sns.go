package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccessKeyId        stringattr.Type `tfsdk:"access_key_id"`
	Secret             stringattr.Type `tfsdk:"secret"`
	Region             stringattr.Type `tfsdk:"region"`
	Endpoint           stringattr.Type `tfsdk:"endpoint"`
	OrganizationNumber stringattr.Type `tfsdk:"organization_number"`
	SenderID           stringattr.Type `tfsdk:"sender_id"`
	EntityID           stringattr.Type `tfsdk:"entity_id"`
	TemplateID         stringattr.Type `tfsdk:"template_id"`
}

func (m *SNSModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "sns"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SNSModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
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

func (m *SNSModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.AccessKeyId, c, "accessKeyId")
	stringattr.Nil(&m.Secret)
	stringattr.Set(&m.Region, c, "awsSNSRegion")
	stringattr.Set(&m.Endpoint, c, "awsEndpoint")
	stringattr.Set(&m.OrganizationNumber, c, "originationNumber")
	stringattr.Set(&m.SenderID, c, "senderId")
	stringattr.Set(&m.EntityID, c, "entityId")
	stringattr.Set(&m.TemplateID, c, "templateId")
}

// Matching

func (m *SNSModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SNSModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SNSModel) SetID(id stringattr.Type) {
	m.ID = id
}
