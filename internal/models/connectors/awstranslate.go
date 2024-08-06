package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AWSTranslateAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"access_key_id":     stringattr.Required(),
	"secret_access_key": stringattr.SecretRequired(),
	"session_token":     stringattr.SecretOptional(),
	"region":            stringattr.Required(),
}

// Model

type AWSTranslateModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	SessionToken    types.String `tfsdk:"session_token"`
	Region          types.String `tfsdk:"region"`
}

func (m *AWSTranslateModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "aws-translate"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AWSTranslateModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *AWSTranslateModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyID, c, "accessKeyId")
	stringattr.Get(m.SecretAccessKey, c, "secretAccessKey")
	stringattr.Get(m.SessionToken, c, "sessionToken")
	stringattr.Get(m.Region, c, "region")
	return c
}

// Matching

func (m *AWSTranslateModel) GetName() types.String {
	return m.Name
}

func (m *AWSTranslateModel) GetID() types.String {
	return m.ID
}

func (m *AWSTranslateModel) SetID(id types.String) {
	m.ID = id
}
