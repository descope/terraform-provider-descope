package lists

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ListAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description": stringattr.Optional(stringattr.StandardLenValidator),
	"type":        stringattr.Optional(stringvalidator.OneOf("texts", "ips", "json")),
	// current internal implementation allows string of array OR json object.
	"data": schema.DynamicAttribute{
		Required: true,
	},
}

type ListModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Type        stringattr.Type `tfsdk:"type"`
	Data        types.Dynamic   `tfsdk:"data"`
}

func (m *ListModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Type, data, "type")

	stringattr.Get(m.ID, data, "id")

	if !m.Data.IsNull() && !m.Data.IsUnknown() {
		underlying := m.Data.UnderlyingValue()
		if underlying != nil {
			if str, ok := underlying.(types.String); ok {
				// JSON string - unmarshal it
				var jsonValue any
				if err := json.Unmarshal([]byte(str.ValueString()), &jsonValue); err == nil {
					data["data"] = jsonValue
				}
			} else if list, ok := underlying.(types.List); ok {
				// texts/ips array - convert to []string
				stringArray := make([]string, 0, len(list.Elements()))
				for _, elem := range list.Elements() {
					if str, ok := elem.(types.String); ok {
						stringArray = append(stringArray, str.ValueString())
					}
				}
				data["data"] = stringArray
			}
		}
	}

	return data
}

func (m *ListModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Type, data, "type")

	// if already set, do not update
	if !m.Data.IsNull() && !m.Data.IsUnknown() {
		return
	}

	// Extract the data field from the response
	listType, _ := data["type"].(string)
	if dataField, exists := data["data"]; exists {
		if listType == "texts" || listType == "ips" {
			// texts/ips types: extract array of strings
			var stringArray []string
			if arr, ok := dataField.([]string); ok {
				stringArray = arr
			} else if arr, ok := dataField.([]any); ok {
				stringArray = make([]string, 0, len(arr))
				for _, v := range arr {
					if str, ok := v.(string); ok {
						stringArray = append(stringArray, str)
					}
				}
			}
			if stringArray != nil {
				elems := make([]attr.Value, len(stringArray))
				for i, s := range stringArray {
					elems[i] = types.StringValue(s)
				}
				m.Data = types.DynamicValue(types.ListValueMust(types.StringType, elems))
			}
		} else {
			// JSON type: marshal to JSON string
			if jsonBytes, err := json.Marshal(dataField); err == nil {
				m.Data = types.DynamicValue(types.StringValue(string(jsonBytes)))
			}
		}
	}
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
