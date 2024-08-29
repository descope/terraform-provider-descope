package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SumoLogicValidator = objectattr.NewValidator[SumoLogicModel]("must have a valid configuration")

var SumoLogicAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"http_source_url":          stringattr.SecretRequired(),
	"audit_enabled":            boolattr.Default(true),
	"audit_filters":            stringattr.Default(""),
	"troubleshoot_log_enabled": boolattr.Default(false),
}

// Model

type SumoLogicModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	HTTPSourceURL          types.String `tfsdk:"http_source_url"`
	AuditEnabled           types.Bool   `tfsdk:"audit_enabled"`
	AuditFilters           types.String `tfsdk:"audit_filters"`
	TroubleshootLogEnabled types.Bool   `tfsdk:"troubleshoot_log_enabled"`
}

func (m *SumoLogicModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "sumologic"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SumoLogicModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

func (m *SumoLogicModel) Validate(h *helpers.Handler) {
	if !m.AuditFilters.IsNull() && !m.AuditEnabled.IsNull() && !m.AuditEnabled.ValueBool() {
		h.Error("Invalid connector configuration", "The audit_filters field cannot be used when audit_enabled is set to false")
	}
}

// Configuration

func (m *SumoLogicModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.HTTPSourceURL, c, "httpSourceUrl")
	boolattr.Get(m.AuditEnabled, c, "auditEnabled")
	stringattr.Get(m.AuditFilters, c, "auditFilters")
	boolattr.Get(m.TroubleshootLogEnabled, c, "troubleshootLogEnabled")
	return c
}

// Matching

func (m *SumoLogicModel) GetName() types.String {
	return m.Name
}

func (m *SumoLogicModel) GetID() types.String {
	return m.ID
}

func (m *SumoLogicModel) SetID(id types.String) {
	m.ID = id
}
