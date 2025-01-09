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

var HTTPAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"base_url":                   stringattr.Required(),
	"authentication":             objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":                    mapattr.StringOptional(),
	"hmac_secret":                stringattr.SecretOptional(),
	"insecure":                   boolattr.Default(false),
	"include_headers_in_context": boolattr.Default(false),
	"use_static_ips":             boolattr.Default(false),
}

// Model

type HTTPModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	BaseURL                 types.String        `tfsdk:"base_url"`
	Authentication          *HTTPAuthFieldModel `tfsdk:"authentication"`
	Headers                 map[string]string   `tfsdk:"headers"`
	HMACSecret              types.String        `tfsdk:"hmac_secret"`
	Insecure                types.Bool          `tfsdk:"insecure"`
	IncludeHeadersInContext types.Bool          `tfsdk:"include_headers_in_context"`
	UseStaticIPs            types.Bool          `tfsdk:"use_static_ips"`
}

func (m *HTTPModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "http"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *HTTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *HTTPModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.BaseURL, c, "baseUrl")
	objectattr.Get(m.Authentication, c, "authentication", h)
	c["headers"] = m.Headers
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	boolattr.Get(m.IncludeHeadersInContext, c, "includeHeadersInContext")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	return c
}

// Matching

func (m *HTTPModel) GetName() types.String {
	return m.Name
}

func (m *HTTPModel) GetID() types.String {
	return m.ID
}

func (m *HTTPModel) SetID(id types.String) {
	m.ID = id
}
