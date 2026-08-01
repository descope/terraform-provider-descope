package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SNSValidator = objattr.NewValidator[SNSModel]("must have a valid configuration")

var SNSAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"auth_type":          stringattr.Default("credentials", stringvalidator.OneOf("credentials", "assumeRole")),
	"access_key_id":      stringattr.SecretOptional(),
	"secret":             stringattr.SecretOptional(),
	"role_arn":           stringattr.Default(""),
	"external_id":        stringattr.Default(""),
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

	AuthType          stringattr.Type `tfsdk:"auth_type"`
	AccessKeyId       stringattr.Type `tfsdk:"access_key_id"`
	Secret            stringattr.Type `tfsdk:"secret"`
	RoleARN           stringattr.Type `tfsdk:"role_arn"`
	ExternalID        stringattr.Type `tfsdk:"external_id"`
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

	isCredentials := m.AuthType.ValueString() == "credentials" || m.AuthType.ValueString() == ""
	isAssumeRole := m.AuthType.ValueString() == "assumeRole"

	if isCredentials && m.AccessKeyId.ValueString() == "" && !m.AccessKeyId.IsUnknown() {
		h.Conflict("The access_key_id field is required when auth_type is set to 'credentials'")
	}
	if isCredentials && m.Secret.ValueString() == "" && !m.Secret.IsUnknown() {
		h.Conflict("The secret field is required when auth_type is set to 'credentials'")
	}
	if !isCredentials && m.AccessKeyId.ValueString() != "" {
		h.Conflict("The access_key_id field can only be used when auth_type is set to 'credentials'")
	}
	if !isCredentials && m.Secret.ValueString() != "" {
		h.Conflict("The secret field can only be used when auth_type is set to 'credentials'")
	}

	if isAssumeRole && m.RoleARN.ValueString() == "" && !m.RoleARN.IsUnknown() {
		h.Conflict("The role_arn field is required when auth_type is set to 'assumeRole'")
	}
	if isAssumeRole && m.ExternalID.ValueString() == "" && !m.ExternalID.IsUnknown() {
		h.Conflict("The external_id field is required when auth_type is set to 'assumeRole'")
	}
	if !isAssumeRole && m.RoleARN.ValueString() != "" {
		h.Conflict("The role_arn field can only be used when auth_type is set to 'assumeRole'")
	}
	if !isAssumeRole && m.ExternalID.ValueString() != "" {
		h.Conflict("The external_id field can only be used when auth_type is set to 'assumeRole'")
	}
}

// Configuration

func (m *SNSModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AuthType, c, "authType")
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.RoleARN, c, "roleArn")
	stringattr.Get(m.ExternalID, c, "externalId")
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
	stringattr.Set(&m.AuthType, c, "authType")
	stringattr.Nil(&m.AccessKeyId)
	stringattr.Nil(&m.Secret)
	stringattr.Set(&m.RoleARN, c, "roleArn")
	stringattr.Set(&m.ExternalID, c, "externalId")
	stringattr.Set(&m.Region, c, "awsSNSRegion")
	stringattr.Set(&m.Endpoint, c, "awsEndpoint")
	if m.OrganizationNumber.ValueString() == "" { // Don't overwrite when deprecated field is set
		stringattr.Set(&m.OriginationNumber, c, "originationNumber")
	}
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
