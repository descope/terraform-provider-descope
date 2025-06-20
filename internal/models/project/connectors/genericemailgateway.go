// Code generated by terragen. DO NOT EDIT.

package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var GenericEmailGatewayAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"post_url":       stringattr.Required(),
	"sender":         stringattr.Default(""),
	"authentication": objattr.Default(HTTPAuthFieldDefault, HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":        strmapattr.Default(),
	"hmac_secret":    stringattr.SecretOptional(),
	"insecure":       boolattr.Default(false),
	"use_static_ips": boolattr.Default(false),
}

// Model

type GenericEmailGatewayModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`

	PostURL        stringattr.Type                  `tfsdk:"post_url"`
	Sender         stringattr.Type                  `tfsdk:"sender"`
	Authentication objattr.Type[HTTPAuthFieldModel] `tfsdk:"authentication"`
	Headers        strmapattr.Type                  `tfsdk:"headers"`
	HMACSecret     stringattr.Type                  `tfsdk:"hmac_secret"`
	Insecure       boolattr.Type                    `tfsdk:"insecure"`
	UseStaticIPs   boolattr.Type                    `tfsdk:"use_static_ips"`
}

func (m *GenericEmailGatewayModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "generic-email-gateway"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *GenericEmailGatewayModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		m.SetConfigurationValues(c, h)
	}
}

// Configuration

func (m *GenericEmailGatewayModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.PostURL, c, "postUrl")
	stringattr.Get(m.Sender, c, "sender")
	objattr.Get(m.Authentication, c, "authentication", h)
	getHeaders(m.Headers, c, "headers", h)
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	return c
}

func (m *GenericEmailGatewayModel) SetConfigurationValues(c map[string]any, h *helpers.Handler) {
	stringattr.Set(&m.PostURL, c, "postUrl")
	stringattr.Set(&m.Sender, c, "sender")
	objattr.Set(&m.Authentication, c, "authentication", h)
	setHeaders(&m.Headers, c, "headers", h)
	stringattr.Nil(&m.HMACSecret)
	boolattr.Set(&m.Insecure, c, "insecure")
	boolattr.Set(&m.UseStaticIPs, c, "useStaticIps")
}

// Matching

func (m *GenericEmailGatewayModel) GetName() stringattr.Type {
	return m.Name
}

func (m *GenericEmailGatewayModel) GetID() stringattr.Type {
	return m.ID
}

func (m *GenericEmailGatewayModel) SetID(id stringattr.Type) {
	m.ID = id
}
