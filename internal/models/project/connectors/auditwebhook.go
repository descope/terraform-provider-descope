// Code generated by terragen. DO NOT EDIT.

package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var AuditWebhookAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"base_url":       stringattr.Required(),
	"authentication": objattr.Default(HTTPAuthFieldDefault, HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":        strmapattr.Default(),
	"hmac_secret":    stringattr.SecretOptional(),
	"insecure":       boolattr.Default(false),
	"audit_filters":  listattr.Default[AuditFilterFieldModel](AuditFilterFieldAttributes),
}

// Model

type AuditWebhookModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	BaseURL        stringattr.Type                      `tfsdk:"base_url"`
	Authentication objattr.Type[HTTPAuthFieldModel]     `tfsdk:"authentication"`
	Headers        strmapattr.Type                      `tfsdk:"headers"`
	HMACSecret     stringattr.Type                      `tfsdk:"hmac_secret"`
	Insecure       boolattr.Type                        `tfsdk:"insecure"`
	AuditFilters   listattr.Type[AuditFilterFieldModel] `tfsdk:"audit_filters"`
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
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *AuditWebhookModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.BaseURL, c, "baseUrl")
	objattr.Get(m.Authentication, c, "authentication", h)
	getHeaders(m.Headers, c, "headers", h)
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	listattr.Get(m.AuditFilters, c, "auditFilters", h)
	return c
}

func (m *AuditWebhookModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.BaseURL, c, "baseUrl")
	objattr.Set(&m.Authentication, c, "authentication", h)
	setHeaders(&m.Headers, c, "headers", h)
	stringattr.Nil(&m.HMACSecret)
	boolattr.Set(&m.Insecure, c, "insecure")
	listattr.Set(&m.AuditFilters, c, "auditFilters", h)
}

// Matching

func (m *AuditWebhookModel) GetName() stringattr.Type {
	return m.Name
}

func (m *AuditWebhookModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AuditWebhookModel) SetID(id stringattr.Type) {
	m.ID = id
}
