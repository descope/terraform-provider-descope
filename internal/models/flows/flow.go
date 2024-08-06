package flows

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var FlowValidator = objectattr.NewValidator[FlowModel]("must be valid JSON data and have all requirements satisfied")

var FlowAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(),
}

type FlowModel struct {
	Data types.String `tfsdk:"data"`
}

func (m *FlowModel) Values(h *helpers.Handler) map[string]any {
	return getFlowData(m.Data, h)
}

func (m *FlowModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all styles values are specified in the configuration
}

func (m *FlowModel) Validate(h *helpers.Handler) {
	data := getFlowData(m.Data, h)
	for _, field := range []string{"metadata", "contents", "screens"} {
		if data[field] == nil {
			h.Error("Invalid flow data", "Expected a JSON object with a "+field+" field")
		}
	}
}

// Computed Mapping

func getFlowData(data types.String, h *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		h.Error("Invalid flow data", "Failed to parse JSON: %s", err.Error())
	}
	return m
}
