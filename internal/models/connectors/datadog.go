package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DatadogValidator = objectattr.NewValidator[DatadogModel]("must have a valid configuration")

var DatadogAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_key":                  stringattr.SecretRequired(),
	"site":                     stringattr.Default(""),
	"audit_enabled":            boolattr.Default(true),
	"audit_filters":            listattr.Optional(AuditFilterFieldAttributes),
	"troubleshoot_log_enabled": boolattr.Default(false),
}

// Model

type DatadogModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey                 types.String             `tfsdk:"api_key"`
	Site                   types.String             `tfsdk:"site"`
	AuditEnabled           types.Bool               `tfsdk:"audit_enabled"`
	AuditFilters           []*AuditFilterFieldModel `tfsdk:"audit_filters"`
	TroubleshootLogEnabled types.Bool               `tfsdk:"troubleshoot_log_enabled"`
}

func (m *DatadogModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "datadog"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *DatadogModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.APIKey, c, "apiKey")
		stringattr.Set(&m.Site, c, "site")
		boolattr.Set(&m.AuditEnabled, c, "auditEnabled")
		listattr.Set(&m.AuditFilters, c, "auditFilters", h)
		boolattr.Set(&m.TroubleshootLogEnabled, c, "troubleshootLogEnabled")
	}
}

func (m *DatadogModel) Validate(h *helpers.Handler) {
	if len(m.AuditFilters) != 0 && !m.AuditEnabled.IsNull() && !m.AuditEnabled.ValueBool() {
		h.Error("Invalid connector configuration", "The audit_filters field cannot be used when audit_enabled is set to false")
	}
}

// Configuration

func (m *DatadogModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.Site, c, "site")
	boolattr.Get(m.AuditEnabled, c, "auditEnabled")
	listattr.Get(m.AuditFilters, c, "auditFilters", h)
	boolattr.Get(m.TroubleshootLogEnabled, c, "troubleshootLogEnabled")
	return c
}

// Matching

func (m *DatadogModel) GetName() types.String {
	return m.Name
}

func (m *DatadogModel) GetID() types.String {
	return m.ID
}

func (m *DatadogModel) SetID(id types.String) {
	m.ID = id
}
