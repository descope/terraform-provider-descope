package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/floatattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GenericSmsGatewayAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"post_url":       stringattr.Required(),
	"sender":         stringattr.Default(""),
	"timeout":        floatattr.Default(0),
	"authentication": objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":        mapattr.StringOptional(),
	"hmac_secret":    stringattr.SecretOptional(),
	"insecure":       boolattr.Default(false),
	"use_static_ips": boolattr.Default(false),
}

// Model

type GenericSmsGatewayModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	PostUrl        types.String        `tfsdk:"post_url"`
	Sender         types.String        `tfsdk:"sender"`
	Timeout        types.Float64       `tfsdk:"timeout"`
	Authentication *HTTPAuthFieldModel `tfsdk:"authentication"`
	Headers        map[string]string   `tfsdk:"headers"`
	HMACSecret     types.String        `tfsdk:"hmac_secret"`
	Insecure       types.Bool          `tfsdk:"insecure"`
	UseStaticIPs   types.Bool          `tfsdk:"use_static_ips"`
}

func (m *GenericSmsGatewayModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "generic-sms-gateway"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *GenericSmsGatewayModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *GenericSmsGatewayModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.PostUrl, c, "postUrl")
	stringattr.Get(m.Sender, c, "sender")
	floatattr.Get(m.Timeout, c, "timeout")
	objectattr.Get(m.Authentication, c, "authentication", h)
	c["headers"] = m.Headers
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	return c
}

// Matching

func (m *GenericSmsGatewayModel) GetName() types.String {
	return m.Name
}

func (m *GenericSmsGatewayModel) GetID() types.String {
	return m.ID
}

func (m *GenericSmsGatewayModel) SetID(id types.String) {
	m.ID = id
}
