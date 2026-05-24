package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SCIMAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"disabled":         boolattr.Default(false),
	"federated_app_id": stringattr.Required(),
	"base_url":         stringattr.Required(),
	"authentication":   objattr.Default(HTTPAuthFieldDefault, HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":          strmapattr.Default(),
	"hmac_secret":      stringattr.SecretOptional(),
	"insecure":         boolattr.Default(false),
}

// Model

type SCIMModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	Disabled       boolattr.Type                    `tfsdk:"disabled"`
	FederatedAppID stringattr.Type                  `tfsdk:"federated_app_id"`
	BaseURL        stringattr.Type                  `tfsdk:"base_url"`
	Authentication objattr.Type[HTTPAuthFieldModel] `tfsdk:"authentication"`
	Headers        strmapattr.Type                  `tfsdk:"headers"`
	HMACSecret     stringattr.Type                  `tfsdk:"hmac_secret"`
	Insecure       boolattr.Type                    `tfsdk:"insecure"`
}

func (m *SCIMModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "scim"
	boolattr.Get(m.Disabled, data, "disabled")
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SCIMModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	boolattr.Set(&m.Disabled, data, "disabled")
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *SCIMModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.FederatedAppID, c, "federatedAppId")
	stringattr.Get(m.BaseURL, c, "baseUrl")
	objattr.Get(m.Authentication, c, "authentication", h)
	getHeaders(m.Headers, c, "headers", h)
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	return c
}

func (m *SCIMModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.FederatedAppID, c, "federatedAppId")
	stringattr.Set(&m.BaseURL, c, "baseUrl")
	objattr.Set(&m.Authentication, c, "authentication", h)
	setHeaders(&m.Headers, c, "headers", h)
	stringattr.Nil(&m.HMACSecret)
	boolattr.Set(&m.Insecure, c, "insecure")
}

// Matching

func (m *SCIMModel) GetName() stringattr.Type {
	return m.Name
}

func (m *SCIMModel) GetID() stringattr.Type {
	return m.ID
}

func (m *SCIMModel) SetID(id stringattr.Type) {
	m.ID = id
}
