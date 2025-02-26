package connectors

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SESAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key_id": stringattr.Required(),
	"secret":        stringattr.SecretRequired(),
	"region":        stringattr.Required(),
	"endpoint":      stringattr.Default(""),
	"sender":        objectattr.Required(SenderFieldAttributes),
}

// Model

type SESModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyId types.String      `tfsdk:"access_key_id"`
	Secret      types.String      `tfsdk:"secret"`
	Region      types.String      `tfsdk:"region"`
	Endpoint    types.String      `tfsdk:"endpoint"`
	Sender      *SenderFieldModel `tfsdk:"sender"`
}

func (m *SESModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "ses"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SESModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SESModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.Region, c, "region")
	stringattr.Get(m.Endpoint, c, "endpoint")
	maps.Copy(c, m.Sender.Values(h))
	return c
}

// Matching

func (m *SESModel) GetName() types.String {
	return m.Name
}

func (m *SESModel) GetID() types.String {
	return m.ID
}

func (m *SESModel) SetID(id types.String) {
	m.ID = id
}
