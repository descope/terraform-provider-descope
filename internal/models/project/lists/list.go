package lists

import (
	"encoding/json"
	"net"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ListsModifier = listattr.NewModifierMatchingNames[ListModel]("maintains list identifiers between plan changes")

var ListValidator = objattr.NewValidator[ListModel]("must have a valid data value")

var ListAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description": stringattr.Default("", stringattr.StandardLenValidator),
	"type":        stringattr.Required(stringvalidator.OneOf("texts", "ips", "json")),
	"data":        stringattr.Required(),
}

type ListModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Type        stringattr.Type `tfsdk:"type"`
	Data        stringattr.Type `tfsdk:"data"`
}

func (m *ListModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Type, data, "type")

	var v any
	if err := json.Unmarshal([]byte(m.Data.ValueString()), &v); err != nil {
		panic("Invalid template data after validation: " + err.Error())
	}
	data["data"] = v

	return data
}

func (m *ListModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Type, data, "type")

	// always update the data value from the API response if present
	if v, ok := data["data"]; ok {
		if b, err := json.Marshal(v); err == nil {
			m.Data = stringattr.Value(string(b))
		}
		// in case of marshaling fails, keep existing value to avoid breaking state
	} else if m.Data.ValueString() == "" {
		// set default if data is not in response and not already set
		m.Data = stringattr.Value("{}")
	}
}

func (m *ListModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Type, m.Data) {
		return // skip validation if there are unknown values
	}

	var v any
	_ = json.Unmarshal([]byte(m.Data.ValueString()), &v)

	switch m.Type.ValueString() {
	case "texts":
		valid := true
		arr, ok := v.([]any)
		if !ok {
			valid = false
		}
		for _, item := range arr {
			if _, ok := item.(string); !ok {
				valid = false
			}
		}
		if !valid {
			h.Invalid("The 'data' attribute must be a JSON array of strings for list type 'texts'")
		}
	case "ips":
		valid := true
		arr, ok := v.([]any)
		if !ok {
			valid = false
		}
		for _, item := range arr {
			if s, ok := item.(string); !ok || !isPermittedIPValid(s) {
				valid = false
			}
		}
		if !valid {
			h.Invalid("The 'data' attribute must be a JSON array of IP strings for list type 'ips'")
		}
	case "json":
		if _, ok := v.(map[string]any); !ok {
			h.Invalid("The 'data' attribute must be a JSON object for list type 'json'")
		}
	}
}

func isPermittedIPValid(ipOrCIDR string) bool {
	if _, _, err := net.ParseCIDR(ipOrCIDR); err == nil {
		return true // It's a valid CIDR range
	}
	if net.ParseIP(ipOrCIDR) != nil {
		return true // It's a valid IP address
	}
	return false
}

// Matching

func (m *ListModel) GetName() stringattr.Type {
	return m.Name
}

func (m *ListModel) GetID() stringattr.Type {
	return m.ID
}

func (m *ListModel) SetID(id stringattr.Type) {
	m.ID = id
}
