package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Common values

func connectorValues(id, name, description stringattr.Type, h *helpers.Handler) map[string]any {
	data := map[string]any{}

	stringattr.Get(name, data, "name")
	stringattr.Get(description, data, "description")

	// use the name as a lookup key to set the connector reference or existing id
	connectorName := name.ValueString()
	if ref := h.Refs.Get(helpers.ConnectorReferenceKey, connectorName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for connector '%s' to: %s", connectorName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown connector reference", "No connector named '%s' was defined", connectorName)
		data["id"] = id.ValueString()
	}

	return data
}

func setConnectorValues(id, name, description *stringattr.Type, data map[string]any, _ *helpers.Handler) {
	stringattr.Set(id, data, "id")
	stringattr.Set(name, data, "name")
	stringattr.Set(description, data, "description")
}

// Engine assignment

// executorTypeEngine marks a connector as running on a specific engine rather than locally.
// The backend keeps the assignment in the top-level executorType/executorId fields of the
// connector and defaults to a local executor when they are unset, so an empty engine id
// simply leaves the connector running locally.
const executorTypeEngine = "engine"

// setConnectorEngine writes the engine assignment onto the connector object as the top-level
// executor fields the backend expects. An empty engine id emits nothing, so the connector
// keeps running locally.
func setConnectorEngine(data map[string]any, engineID stringattr.Type) {
	if id := engineID.ValueString(); id != "" {
		data["executorType"] = executorTypeEngine
		data["executorId"] = id
	}
}

// getConnectorEngine reads the top-level executor fields back into the engine_id attribute,
// leaving it empty for connectors that run locally.
func getConnectorEngine(data map[string]any, engineID *stringattr.Type) {
	if data["executorType"] == executorTypeEngine {
		stringattr.Set(engineID, data, "executorId")
	} else {
		*engineID = stringattr.Value("")
	}
}

// Connector Utils

func addConnectorReferences[T any, M helpers.NamedModel[T]](h *helpers.Handler, key string, connectors listattr.Type[T]) {
	for v := range listattr.Iterator(connectors, h) {
		var connector M = v
		h.Refs.Add(helpers.ConnectorReferenceKey, key, connector.GetID().ValueString(), connector.GetName().ValueString())
	}
}

func addConnectorNames[T any, M helpers.NamedModel[T]](h *helpers.Handler, names map[string]int, connectors listattr.Type[T]) {
	for v := range listattr.Iterator(connectors, h) {
		var connector M = v
		if v := connector.GetName().ValueString(); v != "" {
			names[v] += 1
		}
	}
}

// Sender Field

var SenderFieldAttributes = map[string]schema.Attribute{
	"email": stringattr.Required(),
	"name":  stringattr.Default(""),
}

type SenderFieldModel struct {
	Email stringattr.Type `tfsdk:"email"`
	Name  stringattr.Type `tfsdk:"name"`
}

func (m *SenderFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Email, data, "fromEmail")
	stringattr.Get(m.Name, data, "fromName")
	return data
}

func (m *SenderFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Email, data, "fromEmail")
	stringattr.Set(&m.Name, data, "fromName")
}

// Server Field

var ServerFieldAttributes = map[string]schema.Attribute{
	"host": stringattr.Required(),
	"port": intattr.Default(25),
}

type ServerFieldModel struct {
	Host stringattr.Type `tfsdk:"host"`
	Port intattr.Type    `tfsdk:"port"`
}

func (m *ServerFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Host, data, "host")
	intattr.Get(m.Port, data, "port")
	return data
}

func (m *ServerFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Host, data, "host")
	intattr.Set(&m.Port, data, "port")
}

// Audit Filter Field

var AuditFilterFieldAttributes = map[string]schema.Attribute{
	"key":      stringattr.Required(stringvalidator.OneOf("actions", "tenants")),
	"operator": stringattr.Required(stringvalidator.OneOf("includes", "excludes")),
	"values":   strlistattr.Required(listvalidator.SizeAtLeast(1)),
}

type AuditFilterFieldModel struct {
	Key      stringattr.Type  `tfsdk:"key"`
	Operator stringattr.Type  `tfsdk:"operator"`
	Vals     strlistattr.Type `tfsdk:"values"`
}

func (m *AuditFilterFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Key, data, "key")
	stringattr.Get(m.Operator, data, "operator")
	strlistattr.Get(m.Vals, data, "values", h)
	return data
}

func (m *AuditFilterFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Key, data, "key")
	stringattr.Set(&m.Operator, data, "operator")
	strlistattr.Set(&m.Vals, data, "values", h)
}

// HTTP Headers

func getHeaders(s strmapattr.Type, data map[string]any, key string, h *helpers.Handler) { // nolint:unparam
	headers := []any{}
	for k, v := range strmapattr.Iterator(s, h) {
		headers = append(headers, map[string]any{"key": k, "value": v})
	}
	data[key] = headers
}

func setHeaders(s *strmapattr.Type, data map[string]any, key string, _ *helpers.Handler) { // nolint:unparam
	headers := map[string]string{}
	if v, ok := data[key].([]any); ok {
		for i := range v {
			if m, ok := v[i].(map[string]any); ok {
				key, _ := m["key"].(string)
				value, _ := m["value"].(string)
				headers[key] = value
			}
		}
	}
	*s = strmapattr.Value(headers)
}

// HTTP Auth Field

var HTTPAuthFieldValidator = objattr.NewValidator[HTTPAuthFieldModel]("must specify exactly one authentication method")

var HTTPAuthFieldAttributes = map[string]schema.Attribute{
	"bearer_token":              stringattr.SecretOptional(),
	"basic":                     objattr.Default[HTTPAuthBasicFieldModel](nil, HTTPAuthBasicFieldAttributes),
	"api_key":                   objattr.Default[HTTPAuthAPIKeyFieldModel](nil, HTTPAuthAPIKeyFieldAttributes),
	"oauth2_client_credentials": objattr.Default[HTTPAuthOAuth2ClientCredentialsFieldModel](nil, HTTPAuthOAuth2ClientCredentialsFieldAttributes),
}

var HTTPAuthFieldDefault = &HTTPAuthFieldModel{
	BearerToken:             stringattr.Value(""),
	Basic:                   objattr.Value[HTTPAuthBasicFieldModel](nil),
	ApiKey:                  objattr.Value[HTTPAuthAPIKeyFieldModel](nil),
	OAuth2ClientCredentials: objattr.Value[HTTPAuthOAuth2ClientCredentialsFieldModel](nil),
}

type HTTPAuthFieldModel struct {
	BearerToken             stringattr.Type                                         `tfsdk:"bearer_token"`
	Basic                   objattr.Type[HTTPAuthBasicFieldModel]                   `tfsdk:"basic"`
	ApiKey                  objattr.Type[HTTPAuthAPIKeyFieldModel]                  `tfsdk:"api_key"`
	OAuth2ClientCredentials objattr.Type[HTTPAuthOAuth2ClientCredentialsFieldModel] `tfsdk:"oauth2_client_credentials"`
}

func (m *HTTPAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["method"] = "none"
	if v := m.BearerToken.ValueString(); v != "" {
		data["method"] = "bearerToken"
		data["bearerToken"] = v
	}
	if m.Basic.IsSet() {
		data["method"] = "basic"
		objattr.Get(m.Basic, data, "basic", h)
	}
	if m.ApiKey.IsSet() {
		data["method"] = "apiKey"
		objattr.Get(m.ApiKey, data, "apiKey", h)
	}
	if m.OAuth2ClientCredentials.IsSet() {
		data["method"] = "oauth2ClientCredentials"
		objattr.Get(m.OAuth2ClientCredentials, data, "oauth2ClientCredentials", h)
	}
	return data
}

func (m *HTTPAuthFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Nil(&m.BearerToken)
	if data["method"] == "basic" {
		objattr.Set(&m.Basic, data, "basic", h)
	} else {
		objattr.Nil(&m.Basic)
	}
	if data["method"] == "apiKey" {
		objattr.Set(&m.ApiKey, data, "apiKey", h)
	} else {
		objattr.Nil(&m.ApiKey)
	}
	if data["method"] == "oauth2ClientCredentials" {
		objattr.Set(&m.OAuth2ClientCredentials, data, "oauth2ClientCredentials", h)
	} else {
		objattr.Nil(&m.OAuth2ClientCredentials)
	}
}

func (m *HTTPAuthFieldModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.BearerToken) {
		return // skip validation if there are unknown values
	}

	count := 0
	if m.BearerToken.ValueString() != "" {
		count += 1
	}
	if m.Basic.IsSet() {
		count += 1
	}
	if m.ApiKey.IsSet() {
		count += 1
	}
	if m.OAuth2ClientCredentials.IsSet() {
		count += 1
	}

	if count > 1 {
		h.Invalid("Cannot specify more than one connector authentication method")
	}
}

// HTTP Auth Basic Field

var HTTPAuthBasicFieldAttributes = map[string]schema.Attribute{
	"username": stringattr.Required(),
	"password": stringattr.SecretRequired(),
}

type HTTPAuthBasicFieldModel struct {
	Username stringattr.Type `tfsdk:"username"`
	Password stringattr.Type `tfsdk:"password"`
}

func (m *HTTPAuthBasicFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Username, data, "username")
	stringattr.Get(m.Password, data, "password")
	return data
}

func (m *HTTPAuthBasicFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Username, data, "username")
	stringattr.Nil(&m.Password)
}

// HTTP Auth APIKey Field

var HTTPAuthAPIKeyFieldAttributes = map[string]schema.Attribute{
	"key":   stringattr.Required(),
	"token": stringattr.SecretRequired(),
}

type HTTPAuthAPIKeyFieldModel struct {
	Key   stringattr.Type `tfsdk:"key"`
	Token stringattr.Type `tfsdk:"token"`
}

func (m *HTTPAuthAPIKeyFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Key, data, "key")
	stringattr.Get(m.Token, data, "token")
	return data
}

func (m *HTTPAuthAPIKeyFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Key, data, "key")
	stringattr.Nil(&m.Token)
}

// HTTP Auth OAuth2 Client Credentials Field

var HTTPAuthOAuth2ClientCredentialsFieldAttributes = map[string]schema.Attribute{
	"client_id":             stringattr.Required(),
	"client_secret":         stringattr.SecretRequired(),
	"auth_url":              stringattr.Required(),
	"auth_style":            stringattr.Default("header", stringvalidator.OneOf("header", "params")),
	"scopes":                stringattr.Default(""),
	"token_request_headers": strmapattr.Default(),
}

type HTTPAuthOAuth2ClientCredentialsFieldModel struct {
	ClientID            stringattr.Type `tfsdk:"client_id"`
	ClientSecret        stringattr.Type `tfsdk:"client_secret"`
	AuthURL             stringattr.Type `tfsdk:"auth_url"`
	AuthStyle           stringattr.Type `tfsdk:"auth_style"`
	Scopes              stringattr.Type `tfsdk:"scopes"`
	TokenRequestHeaders strmapattr.Type `tfsdk:"token_request_headers"`
}

func (m *HTTPAuthOAuth2ClientCredentialsFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ClientID, data, "clientId")
	stringattr.Get(m.ClientSecret, data, "clientSecret")
	stringattr.Get(m.AuthURL, data, "authUrl")
	stringattr.Get(m.AuthStyle, data, "authStyle")
	stringattr.Get(m.Scopes, data, "scopes")
	getHeaders(m.TokenRequestHeaders, data, "tokenRequestHeaders", h)
	return data
}

func (m *HTTPAuthOAuth2ClientCredentialsFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ClientID, data, "clientId")
	stringattr.Nil(&m.ClientSecret)
	stringattr.Set(&m.AuthURL, data, "authUrl")
	stringattr.Set(&m.AuthStyle, data, "authStyle")
	stringattr.Set(&m.Scopes, data, "scopes")
	setHeaders(&m.TokenRequestHeaders, data, "tokenRequestHeaders", h)
}
