package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SESAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key_id": stringattr.Required(),
	"secret":        stringattr.SecretRequired(),
	"region":        stringattr.Required(),
	"endpoint":      stringattr.Default(""),
	"sender":        objattr.Required[SenderFieldModel](SenderFieldAttributes),
}

// Model

type SESModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AccessKeyId stringattr.Type                `tfsdk:"access_key_id"`
	Secret      stringattr.Type                `tfsdk:"secret"`
	Region      stringattr.Type                `tfsdk:"region"`
	Endpoint    stringattr.Type                `tfsdk:"endpoint"`
	Sender      objattr.Type[SenderFieldModel] `tfsdk:"sender"`
}

func (m *SESModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "ses"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SESModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *SESModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.Region, c, "region")
	stringattr.Get(m.Endpoint, c, "endpoint")
	objattr.Get(m.Sender, c, helpers.RootKey, h)
	return c
}

func (m *SESModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.AccessKeyId, c, "accessKeyId")
	stringattr.Nil(&m.Secret)
	stringattr.Set(&m.Region, c, "awsSNSRegion")
	stringattr.Set(&m.Endpoint, c, "awsEndpoint")
	objattr.Set(&m.Sender, c, helpers.RootKey, h)
}

// Matching

func (m *SESModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SESModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SESModel) SetID(id stringattr.Type) {
	m.ID = id
}
