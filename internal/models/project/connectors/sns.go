package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SNSValidator = objattr.NewValidator[SNSModel]("must have a valid configuration")

var SNSAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key_id":      stringattr.SecretRequired(),
	"secret":             stringattr.SecretRequired(),
	"region":             stringattr.Required(),
	"endpoint":           stringattr.Default(""),
	"origination_number": stringattr.Default(""),
	"sender_id":          stringattr.Default(""),
	"entity_id":          stringattr.Default(""),
	"template_id":        stringattr.Default(""),

	// Deprecated fields
	"organization_number": stringattr.Renamed("organization_number", "origination_number"),
}

// Model

type SNSModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccessKeyId       stringattr.Type `tfsdk:"access_key_id"`
	Secret            stringattr.Type `tfsdk:"secret"`
	Region            stringattr.Type `tfsdk:"region"`
	Endpoint          stringattr.Type `tfsdk:"endpoint"`
	OriginationNumber stringattr.Type `tfsdk:"origination_number"`
	SenderID          stringattr.Type `tfsdk:"sender_id"`
	EntityID          stringattr.Type `tfsdk:"entity_id"`
	TemplateID        stringattr.Type `tfsdk:"template_id"`

	// Deprecated fields
	OrganizationNumber stringattr.Type `tfsdk:"organization_number"`
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

func (m *SNSModel) Validate(h *helpers.Handler) {
	if m.OrganizationNumber.ValueString() != "" && m.OriginationNumber.ValueString() != "" {
		h.Conflict("The organization_number field has been renamed to origination_number, please use only origination_number going forward")
	}
}

// Configuration

func (m *SNSModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.Region, c, "awsSNSRegion")
	stringattr.Get(m.Endpoint, c, "awsEndpoint")
	stringattr.Get(m.OriginationNumber, c, "originationNumber")
	stringattr.Get(m.SenderID, c, "senderId")
	stringattr.Get(m.EntityID, c, "entityId")
	stringattr.Get(m.TemplateID, c, "templateId")

	// Deprecated fields
	if m.OriginationNumber.ValueString() == "" {
		stringattr.Get(m.OrganizationNumber, c, "originationNumber")
	}
	return c
}

func (m *SNSModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Nil(&m.AccessKeyId)
	stringattr.Nil(&m.Secret)
	stringattr.Set(&m.Region, c, "awsSNSRegion")
	stringattr.Set(&m.Endpoint, c, "awsEndpoint")
	stringattr.Set(&m.OriginationNumber, c, "originationNumber")
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
