package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SESValidator = objattr.NewValidator[SESModel]("must have a valid configuration")

var SESAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"auth_type":     stringattr.Default("credentials", stringvalidator.OneOf("credentials", "assumeRole")),
	"access_key_id": stringattr.SecretOptional(),
	"secret":        stringattr.SecretOptional(),
	"role_arn":      stringattr.Default(""),
	"external_id":   stringattr.Default(""),
	"region":        stringattr.Required(),
	"endpoint":      stringattr.Default(""),
	"sender":        objattr.Required[SenderFieldModel](SenderFieldAttributes),
}

// Model

type SESModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	AuthType    stringattr.Type                `tfsdk:"auth_type"`
	AccessKeyId stringattr.Type                `tfsdk:"access_key_id"`
	Secret      stringattr.Type                `tfsdk:"secret"`
	RoleARN     stringattr.Type                `tfsdk:"role_arn"`
	ExternalID  stringattr.Type                `tfsdk:"external_id"`
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

func (m *SESModel) Validate(h *helpers.Handler) {
	// checking for "assumeRole" value even in the "credentials" cehcks to take into account Null value during validation
	if m.AccessKeyId.ValueString() != "" && m.AuthType.ValueString() != "" && m.AuthType.ValueString() == "assumeRole" {
		h.Conflict("The access_key_id field can only be used when auth_type is set to 'credentials'")
	}
	if m.AccessKeyId.ValueString() == "" && !m.AccessKeyId.IsUnknown() && m.AuthType.ValueString() != "assumeRole" {
		h.Conflict("The access_key_id field is required when auth_type is set to 'credentials'")
	}
	if m.Secret.ValueString() != "" && m.AuthType.ValueString() != "" && m.AuthType.ValueString() == "assumeRole" {
		h.Conflict("The secret field can only be used when auth_type is set to 'credentials'")
	}
	if m.Secret.ValueString() == "" && !m.Secret.IsUnknown() && m.AuthType.ValueString() != "assumeRole" {
		h.Conflict("The secret field is required when auth_type is set to 'credentials'")
	}
	if m.RoleARN.ValueString() != "" && m.AuthType.ValueString() != "assumeRole" {
		h.Conflict("The role_arn field can only be used when auth_type is set to 'assumeRole'")
	}
	if m.RoleARN.ValueString() == "" && !m.RoleARN.IsUnknown() && m.AuthType.ValueString() == "assumeRole" {
		h.Conflict("The role_arn field is required when auth_type is set to 'assumeRole'")
	}
	if m.ExternalID.ValueString() != "" && m.AuthType.ValueString() != "assumeRole" {
		h.Conflict("The external_id field can only be used when auth_type is set to 'assumeRole'")
	}
	if m.ExternalID.ValueString() == "" && !m.ExternalID.IsUnknown() && m.AuthType.ValueString() == "assumeRole" {
		h.Conflict("The external_id field is required when auth_type is set to 'assumeRole'")
	}
}

// Configuration

func (m *SESModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AuthType, c, "authType")
	stringattr.Get(m.AccessKeyId, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secretAccessKey")
	stringattr.Get(m.RoleARN, c, "roleArn")
	stringattr.Get(m.ExternalID, c, "externalId")
	stringattr.Get(m.Region, c, "region")
	stringattr.Get(m.Endpoint, c, "endpoint")
	objattr.Get(m.Sender, c, helpers.RootKey, h)
	return c
}

func (m *SESModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.AuthType, c, "authType")
	stringattr.Nil(&m.AccessKeyId)
	stringattr.Nil(&m.Secret)
	stringattr.Set(&m.RoleARN, c, "roleArn")
	stringattr.Set(&m.ExternalID, c, "externalId")
	stringattr.Set(&m.Region, c, "region")
	stringattr.Set(&m.Endpoint, c, "endpoint")
	if !m.Sender.IsSet() {
		m.Sender = objattr.Value(&SenderFieldModel{}) // XXX switch this together with the Set / Opt refactor
	}
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
