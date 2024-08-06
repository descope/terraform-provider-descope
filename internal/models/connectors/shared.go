package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
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
		h.Error("Unknown connector reference", "No connector named '"+connectorName+"' was defined")
		data["id"] = id.ValueString()
	}

	return data
}

// Connector Identifiers

func SetConnectorIDs[T any, M helpers.MatchableModel[T]](h *helpers.Handler, data map[string]any, key string, connectors []M) {
	for _, connector := range connectors {
		n := connector.GetName().ValueString()
		h.Log("Looking for " + key + " connector named '" + n + "'")
		if connectorID, ok := findConnectorID(h, data, key, n); ok {
			value := types.StringValue(connectorID)
			if !connector.GetID().Equal(value) {
				h.Log("Setting new ID '" + connectorID + "' for " + key + " connector named '" + n + "'")
				connector.SetID(value)
			} else {
				h.Log("Keeping existing ID '" + connectorID + "' for " + key + " connector named '" + n + "'")
			}
		}

	}
}

func findConnectorID(h *helpers.Handler, data map[string]any, key string, name string) (string, bool) {
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

	h.Error("Connector not found", "Expected to find connector of type '"+key+"' to match with '"+name+"' connector")
	return "", false
}

// Connector Utils

func addConnectorReferences[T any, M helpers.MatchableModel[T]](refs helpers.ReferencesMap, key string, connectors []M) {
	for _, connector := range connectors {
		refs.Add(helpers.ConnectorReferenceKey, key, connector.GetID().ValueString(), connector.GetName().ValueString())
	}
}

func addConnectorNames[T any, M helpers.MatchableModel[T]](names map[string]int, connectors []M) {
	for _, connector := range connectors {
		names[connector.GetName().ValueString()] += 1
	}
}

// HTTP Auth Field

var HTTPAuthFieldValidator = objectattr.NewValidator[HTTPAuthFieldModel]("must specify exactly one authentication method")

var HTTPAuthFieldAttributes = map[string]schema.Attribute{
	"bearer_token": stringattr.SecretOptional(),
	"basic": objectattr.Optional(map[string]schema.Attribute{
		"username": stringattr.Required(),
		"password": stringattr.SecretRequired(),
	}),
	"api_key": objectattr.Optional(map[string]schema.Attribute{
		"key":   stringattr.Required(),
		"token": stringattr.SecretRequired(),
	}),
}

type HTTPAuthFieldModel struct {
	BearerToken types.String              `tfsdk:"bearer_token"`
	Basic       *HTTPAuthFieldBasicModel  `tfsdk:"basic"`
	ApiKey      *HTTPAuthFieldApiKeyModel `tfsdk:"api_key"`
}

type HTTPAuthFieldBasicModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type HTTPAuthFieldApiKeyModel struct {
	Key   types.String `tfsdk:"key"`
	Token types.String `tfsdk:"token"`
}

func (m *HTTPAuthFieldModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	if v := m.BearerToken.ValueString(); v != "" {
		data["method"] = "bearerToken"
		data["bearerToken"] = v
	}
	if v := m.Basic; v != nil {
		data["method"] = "basic"
		data["basic"] = map[string]any{
			"username": v.Username.ValueString(),
			"password": v.Password.ValueString(),
		}
	}
	if v := m.ApiKey; v != nil {
		data["method"] = "apiKey"
		data["apiKey"] = map[string]any{
			"key":   v.Key.ValueString(),
			"token": v.Token.ValueString(),
		}
	}
	return data
}

func (m *HTTPAuthFieldModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all auth values are specified in the configuration
}

func (m *HTTPAuthFieldModel) Validate(h *helpers.Handler) {
	count := 0
	if m.BearerToken.ValueString() != "" {
		count += 1
	}
	if m.Basic != nil {
		count += 1
	}
	if m.ApiKey != nil {
		count += 1
	}

	if count == 0 {
		h.Error("Invalid HTTP authentication type", "An HTTP authentication method is required")
	} else if count > 1 {
		h.Error("Invalid HTTP authentication type", "Only one HTTP authentication method is allowed")
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
