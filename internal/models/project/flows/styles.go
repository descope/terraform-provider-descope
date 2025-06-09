package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var StylesValidator = objattr.NewValidator[StylesModel]("must be valid JSON data and have all requirements satisfied")

var StylesAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(),
}

type StylesModel struct {
	Data stringattr.Type `tfsdk:"data"`
}

func (m *StylesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["data"] = getStylesData(m.Data, h)
	return data
}

func (m *StylesModel) SetValues(h *helpers.Handler, data map[string]any) {
	if m.Data.IsUnknown() || m.Data.IsNull() {
		if styleMap, ok := data["data"].(map[string]any); ok {
			b, err := json.Marshal(styleMap)
			if err != nil {
				h.Error("Invalid style data", "Failed to parse JSON: %s", err.Error())
				return
			}
			m.Data = stringattr.Value(string(b))
		}
	}
}

func (m *StylesModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Data) {
		return // skip validation if there are unknown values
	}
	data := getStylesData(m.Data, h)
	if data["styles"] == nil {
		h.Error("Invalid styles data", "Expected a JSON object with a styles field")
	}
}

// Computed Mapping

func getStylesData(data stringattr.Type, h *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		h.Error("Invalid styles data", "Failed to parse JSON: %s", err.Error())
	}
	return m
}
