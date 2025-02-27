package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AuditWebhookAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"base_url":       stringattr.Required(),
	"authentication": objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":        mapattr.StringOptional(),
	"hmac_secret":    stringattr.SecretOptional(),
	"insecure":       boolattr.Default(false),
	"audit_filters":  listattr.Optional(AuditFilterFieldAttributes),
}

// Model

type AuditWebhookModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	BaseURL        types.String             `tfsdk:"base_url"`
	Authentication *HTTPAuthFieldModel      `tfsdk:"authentication"`
	Headers        map[string]string        `tfsdk:"headers"`
	HMACSecret     types.String             `tfsdk:"hmac_secret"`
	Insecure       types.Bool               `tfsdk:"insecure"`
	AuditFilters   []*AuditFilterFieldModel `tfsdk:"audit_filters"`
}

func (m *AuditWebhookModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "audit-webhook"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AuditWebhookModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.BaseURL, c, "baseUrl")
		objectattr.Set(&m.Authentication, c, "authentication", h)
		if vs, ok := c["headers"].(map[string]any); ok {
			for k, v := range vs {
				if s, ok := v.(string); ok {
					m.Headers[k] = s
				}
			}
		}
		stringattr.Set(&m.HMACSecret, c, "hmacSecret")
		boolattr.Set(&m.Insecure, c, "insecure")
		listattr.Set(&m.AuditFilters, c, "auditFilters", h)
	}
}

// Configuration

func (m *AuditWebhookModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.BaseURL, c, "baseUrl")
	objectattr.Get(m.Authentication, c, "authentication", h)
	getHeaders(m.Headers, c, "headers")
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	listattr.Get(m.AuditFilters, c, "auditFilters", h)
	return c
}

// Matching

func (m *AuditWebhookModel) GetName() types.String {
	return m.Name
}

func (m *AuditWebhookModel) GetID() types.String {
	return m.ID
}

func (m *AuditWebhookModel) SetID(id types.String) {
	m.ID = id
}
