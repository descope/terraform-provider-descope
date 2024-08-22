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

var HttpStaticIPAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),

	"base_url":                   stringattr.Required(),
	"authentication":             objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator),
	"headers":                    mapattr.StringOptional(),
	"hmac_secret":                stringattr.SecretOptional(),
	"insecure":                   boolattr.Default(false),
	"include_headers_in_context": boolattr.Default(false),
}

// Model

type HttpStaticIPModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	BaseURL                 types.String        `tfsdk:"base_url"`
	Authentication          *HTTPAuthFieldModel `tfsdk:"authentication"`
	Headers                 map[string]string   `tfsdk:"headers"`
	HMACSecret              types.String        `tfsdk:"hmac_secret"`
	Insecure                types.Bool          `tfsdk:"insecure"`
	IncludeHeadersInContext types.Bool          `tfsdk:"include_headers_in_context"`
}

func (m *HttpStaticIPModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "http-static-ip"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *HttpStaticIPModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *HttpStaticIPModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.BaseURL, c, "baseUrl")
	objectattr.Get(m.Authentication, c, "authentication", h)
	c["headers"] = m.Headers
	stringattr.Get(m.HMACSecret, c, "hmacSecret")
	boolattr.Get(m.Insecure, c, "insecure")
	boolattr.Get(m.IncludeHeadersInContext, c, "includeHeadersInContext")
	return c
}

// Matching

func (m *HttpStaticIPModel) GetName() types.String {
	return m.Name
}

func (m *HttpStaticIPModel) GetID() types.String {
	return m.ID
}

func (m *HttpStaticIPModel) SetID(id types.String) {
	m.ID = id
}
