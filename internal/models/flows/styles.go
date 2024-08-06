package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var StylesValidator = objectattr.NewValidator[StylesModel]("must be valid JSON data and have all requirements satisfied")

var StylesAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(),
}

type StylesModel struct {
	Data types.String `tfsdk:"data"`
}

func (m *StylesModel) Values(h *helpers.Handler) map[string]any {
	return getStylesData(m.Data, h)
}

func (m *StylesModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all styles values are specified in the configuration
}

func (m *StylesModel) Validate(h *helpers.Handler) {
	data := getStylesData(m.Data, h)
	if data["styles"] == nil {
		h.Error("Invalid styles data", "Expected a JSON object with a styles field")
	}
}

// Computed Mapping

func getStylesData(data types.String, h *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		h.Error("Invalid styles data", "Failed to parse JSON: %s", err.Error())
	}
	return m
}
