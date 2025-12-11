package lists

import (
	"encoding/json"
	"net"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ListValidator = objattr.NewValidator[ListModel]("must have a valid data value")

var ListAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
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

	// We do not currently update the data value if it's already set because it might be different after apply
	if m.Data.ValueString() == "" {
		value := "{}"
		if v, ok := data["data"]; ok {
			if b, err := json.Marshal(v); err == nil {
				value = string(b)
			}
		}
		m.Data = stringattr.Value(value)
	}
}

func (m *ListModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Type, m.Data) {
		return // skip validation if there are unknown values
	}

	var v any
	json.Unmarshal([]byte(m.Data.ValueString()), &v)

	switch m.Type.ValueString() {
	case "texts", "ips":
		if _, ok := v.([]any); ok {
			for _, item := range v.([]any) {
				if s, ok := item.(string); !ok {
					h.Invalid("The 'data' attribute must be a JSON array of strings for list types 'texts' and 'ips'")
				} else if m.Type.ValueString() == "ips" && !isPermittedIPValid(s) {
					h.Invalid("The 'data' attribute must be a JSON array of valid IP strings for list type 'ips'")
				}
			}
		} else {
			h.Invalid("The 'data' attribute must be a JSON array of strings for list types 'texts' and 'ips'")
		}
	case "json":
		if _, ok := v.(map[string]any); !ok {
			h.Invalid("The 'data' attribute must be a JSON object for list type 'json'")
		}
	default:
		panic("Unhandled list type in validation: " + m.Type.ValueString())
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
