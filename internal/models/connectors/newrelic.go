package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var NewRelicValidator = objectattr.NewValidator[NewRelicModel]("must have a valid configuration")

var NewRelicAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"api_key":                  stringattr.SecretRequired(),
	"data_center":              stringattr.Default(""),
	"audit_enabled":            boolattr.Default(true),
	"audit_filters":            strlistattr.Optional(),
	"troubleshoot_log_enabled": boolattr.Default(false),
	"override_logs_prefix":     boolattr.Default(false),
	"logs_prefix":              stringattr.Default("descope."),
}

// Model

type NewRelicModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	APIKey                 types.String   `tfsdk:"api_key"`
	DataCenter             types.String   `tfsdk:"data_center"`
	AuditEnabled           types.Bool     `tfsdk:"audit_enabled"`
	AuditFilters           []types.String `tfsdk:"audit_filters"`
	TroubleshootLogEnabled types.Bool     `tfsdk:"troubleshoot_log_enabled"`
	OverrideLogsPrefix     types.Bool     `tfsdk:"override_logs_prefix"`
	LogsPrefix             types.String   `tfsdk:"logs_prefix"`
}

func (m *NewRelicModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "newrelic"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *NewRelicModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

func (m *NewRelicModel) Validate(h *helpers.Handler) {
	if len(m.AuditFilters) > 0 && !m.AuditEnabled.IsNull() && !m.AuditEnabled.ValueBool() {
		h.Error("Invalid connector configuration", "The audit_filters field cannot be used when audit_enabled is set to false")
	}
	if !m.LogsPrefix.IsNull() && !m.OverrideLogsPrefix.ValueBool() {
		h.Error("Invalid connector configuration", "The logs_prefix field cannot be used unless override_logs_prefix is set to true")
	}
}

func (m *NewRelicModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.APIKey, c, "apiKey")
	stringattr.Get(m.DataCenter, c, "dataCenter")
	boolattr.Get(m.AuditEnabled, c, "auditEnabled")

	// Convert list of types.String to a standard Go slice of strings
	var auditFilters []string
	for _, filter := range m.AuditFilters {
		if !filter.IsNull() && !filter.IsUnknown() {
			auditFilters = append(auditFilters, filter.ValueString())
		}
	}
	c["auditFilters"] = auditFilters

	boolattr.Get(m.TroubleshootLogEnabled, c, "troubleshootLogEnabled")
	boolattr.Get(m.OverrideLogsPrefix, c, "overrideLogsPrefix")
	stringattr.Get(m.LogsPrefix, c, "logsPrefix")
	return c
}

// Matching

func (m *NewRelicModel) GetName() types.String {
	return m.Name
}

func (m *NewRelicModel) GetID() types.String {
	return m.ID
}

func (m *NewRelicModel) SetID(id types.String) {
	m.ID = id
}
