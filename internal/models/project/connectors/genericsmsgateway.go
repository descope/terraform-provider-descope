package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GenericSMSGatewayAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"post_url":       stringattr.Required(),
	"sender":         stringattr.Default(""),
	"authentication": objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":        mapattr.StringOptional(),
	"hmac_secret":    stringattr.SecretOptional(),
	"insecure":       boolattr.Default(false),
	"use_static_ips": boolattr.Default(false),
}

// Model

type GenericSMSGatewayModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	PostURL        types.String        `tfsdk:"post_url"`
	Sender         types.String        `tfsdk:"sender"`
	Authentication *HTTPAuthFieldModel `tfsdk:"authentication"`
	Headers        map[string]string   `tfsdk:"headers"`
	HMACSecret     types.String        `tfsdk:"hmac_secret"`
	Insecure       types.Bool          `tfsdk:"insecure"`
	UseStaticIPs   types.Bool          `tfsdk:"use_static_ips"`
}

func (m *GenericSMSGatewayModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "generic-sms-gateway"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *GenericSMSGatewayModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.PostURL, c, "postUrl")
		stringattr.Set(&m.Sender, c, "sender")
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
		boolattr.Set(&m.UseStaticIPs, c, "useStaticIps")
	}
}

// Configuration

func (m *GenericSMSGatewayModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.PostURL, c, "postUrl")
	stringattr.Get(m.Sender, c, "sender")
	objectattr.Get(m.Authentication, c, "authentication", h)
	getHeaders(m.Headers, c, "headers")
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	return c
}

// Matching

func (m *GenericSMSGatewayModel) GetName() types.String {
	return m.Name
}

func (m *GenericSMSGatewayModel) GetID() types.String {
	return m.ID
}

func (m *GenericSMSGatewayModel) SetID(id types.String) {
	m.ID = id
}
