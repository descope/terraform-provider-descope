package flow

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single flow in a Descope project via its exported JSON representation, the same format produced by exporting a flow in the Descope console.",
	Attributes:          FlowAttributes,
}

var FlowAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"flow_id":    stringattr.Required(stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace()),
	"data":       jsonattr.Required("metadata", "contents"),
}

type FlowModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	FlowID    stringattr.Type `tfsdk:"flow_id"`
	Data      jsonattr.Type   `tfsdk:"data"`
}

func (m *FlowModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	jsonattr.Get(m.Data, data, helpers.RootKey)
	flowID := m.FlowID.ValueString()
	if valueID, _ := data["flowId"].(string); valueID != "" && valueID != flowID {
		h.Warn("Possible flow mismatch", "The flow data specifies a different flowId '%s'. You can update the flow data to use the same flowId or ignore this warning to use the '%s' flowId.", valueID, flowID)
	}
	data["flowId"] = flowID
	return data
}

func (m *FlowModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.FlowID, data, "flowId", stringattr.SkipIfAlreadySet)
	jsonattr.Set(&m.Data, data, helpers.RootKey, jsonattr.SkipIfAlreadySet)
}

func (m *FlowModel) GetID() stringattr.Type        { return m.ID }
func (m *FlowModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *FlowModel) GetProjectID() stringattr.Type { return m.ProjectID }
