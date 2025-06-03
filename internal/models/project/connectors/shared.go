package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strmapattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Common values

func connectorValues(id, name, description types.String, h *helpers.Handler) map[string]any {
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

func setConnectorValues(id, name, description *types.String, data map[string]any, _ *helpers.Handler) {
	stringattr.Set(id, data, "id")
	stringattr.Set(name, data, "name")
	stringattr.Set(description, data, "description")
}

// Connector Identifiers

func setConnectorsValues[T any, M helpers.MatchableModel[T]](h *helpers.Handler, data map[string]any, key string, connectors *listattr.Type[T]) {
	listattr.Set[T, M](connectors, data, key, h)

	// if connectors.IsEmpty() { // TODO should be safe to always set the values
	// 	listattr.Set2[T, M](connectors, data, key, h)
	// } else {
	// 	for v := range listattr.MutatingIterator(connectors, h) {
	// 		var connector M = v
	// 		name := connector.GetName().ValueString()
	// 		h.Log("Looking for %s connector named '%s'", key, name)
	// 		if connectorID, ok := findConnectorID(data, key, name, h); ok {
	// 			value := types.StringValue(connectorID)
	// 			if !connector.GetID().Equal(value) {
	// 				h.Log("Setting new ID '%s' for %s connector named '%s'", connectorID, key, name)
	// 				connector.SetID(value)
	// 			} else {
	// 				h.Log("Keeping existing ID '%s' for %s connector named '%s'", connectorID, key, name)
	// 			}
	// 		}
	// 	}
	// }
}

func findConnectorID(data map[string]any, key string, name string, h *helpers.Handler) (string, bool) {
	list, ok := data[key].([]any)
	if !ok {
		return "", false
	}

	for _, v := range list {
		if connector, ok := v.(map[string]any); ok {
			if n, ok := connector["name"].(string); ok && name == n {
				if id, ok := connector["id"].(string); ok {
					return id, true
				}
			}
			if common, ok := connector["common"].(map[string]any); ok {
				if n, ok := common["displayName"].(string); ok && name == n {
					if id, ok := common["id"].(string); ok {
						return id, true
					}
				}
			}
		}
	}

	h.Error("Connector not found", "Expected to find connector of type '%s' to match with '%s' connector", key, name)
	return "", false
}

// Connector Utils

func addConnectorReferences[T any, M helpers.MatchableModel[T]](h *helpers.Handler, key string, connectors listattr.Type[T]) {
	for v := range listattr.Iterator(connectors, h) {
		var connector M = v
		h.Refs.Add(helpers.ConnectorReferenceKey, key, connector.GetID().ValueString(), connector.GetName().ValueString())
	}
}

func addConnectorNames[T any, M helpers.MatchableModel[T]](h *helpers.Handler, names map[string]int, connectors listattr.Type[T]) {
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
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
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
	Host types.String `tfsdk:"host"`
	Port types.Int64  `tfsdk:"port"`
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
	Key      types.String     `tfsdk:"key"`
	Operator types.String     `tfsdk:"operator"`
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

func getHeaders(s strmapattr.Type, data map[string]any, key string, h *helpers.Handler) {
	headers := []any{}
	for k, v := range strmapattr.Iterator(s, h) {
		headers = append(headers, map[string]any{"key": k, "value": v})
	}
	data[key] = headers
}

func setHeaders(s *strmapattr.Type, data map[string]any, key string, _ *helpers.Handler) {
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
	"bearer_token": stringattr.SecretOptional(),
	"basic":        objattr.Default[HTTPAuthBasicFieldModel](nil, HTTPAuthBasicFieldAttributes),
	"api_key":      objattr.Default[HTTPAuthAPIKeyFieldModel](nil, HTTPAuthAPIKeyFieldAttributes),
}

var HTTPAuthFieldDefault = &HTTPAuthFieldModel{
	BearerToken: stringattr.Value(""),
	Basic:       objattr.Value[HTTPAuthBasicFieldModel](nil),
	ApiKey:      objattr.Value[HTTPAuthAPIKeyFieldModel](nil),
}

type HTTPAuthFieldModel struct {
	BearerToken types.String                           `tfsdk:"bearer_token"`
	Basic       objattr.Type[HTTPAuthBasicFieldModel]  `tfsdk:"basic"`
	ApiKey      objattr.Type[HTTPAuthAPIKeyFieldModel] `tfsdk:"api_key"`
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
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
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
	Key   types.String `tfsdk:"key"`
	Token types.String `tfsdk:"token"`
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
