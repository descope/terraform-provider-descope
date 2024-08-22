package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AWSS3Validator = objectattr.NewValidator[AWSS3Model]("must have a valid configuration")

var AWSS3Attributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"access_key_id":            stringattr.SecretRequired(),
	"secret_access_key":        stringattr.SecretRequired(),
	"region":                   stringattr.Required(),
	"bucket":                   stringattr.Required(),
	"audit_enabled":            boolattr.Default(true),
	"audit_filters":            stringattr.Default(""),
	"troubleshoot_log_enabled": boolattr.Default(false),
}

// Model

type AWSS3Model struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyID            types.String `tfsdk:"access_key_id"`
	SecretAccessKey        types.String `tfsdk:"secret_access_key"`
	Region                 types.String `tfsdk:"region"`
	Bucket                 types.String `tfsdk:"bucket"`
	AuditEnabled           types.Bool   `tfsdk:"audit_enabled"`
	AuditFilters           types.String `tfsdk:"audit_filters"`
	TroubleshootLogEnabled types.Bool   `tfsdk:"troubleshoot_log_enabled"`
}

func (m *AWSS3Model) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "aws-s3"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AWSS3Model) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

func (m *AWSS3Model) Validate(h *helpers.Handler) {
	if !m.AuditFilters.IsNull() && !m.AuditEnabled.IsNull() && !m.AuditEnabled.ValueBool() {
		h.Error("Invalid connector configuration", "The audit_filters field cannot be used when audit_enabled is set to false")
	}
}

// Configuration

func (m *AWSS3Model) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyID, c, "accessKeyId")
	stringattr.Get(m.SecretAccessKey, c, "secretAccessKey")
	stringattr.Get(m.Region, c, "region")
	stringattr.Get(m.Bucket, c, "bucket")
	boolattr.Get(m.AuditEnabled, c, "auditEnabled")
	stringattr.Get(m.AuditFilters, c, "auditFilters")
	boolattr.Get(m.TroubleshootLogEnabled, c, "troubleshootLogEnabled")
	return c
}

// Matching

func (m *AWSS3Model) GetName() types.String {
	return m.Name
}

func (m *AWSS3Model) GetID() types.String {
	return m.ID
}

func (m *AWSS3Model) SetID(id types.String) {
	m.ID = id
}
