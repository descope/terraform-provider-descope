package flowtranslations

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages the translations of a flow's texts in a Descope project, the same document produced by exporting a flow's localization in the Descope console.",
	Attributes:          FlowTranslationsAttributes,
}

var FlowTranslationsAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"flow_id":    stringattr.Required(stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace()),
	"languages":  jsonattr.Required(),
}

type FlowTranslationsModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	FlowID    stringattr.Type `tfsdk:"flow_id"`
	Languages jsonattr.Type   `tfsdk:"languages"`
}

func (m *FlowTranslationsModel) Values(h *helpers.Handler) map[string]any {
	overrides := map[string]any{}
	stringattr.Get(m.FlowID, overrides, "flowId")
	jsonattr.Get(m.Languages, overrides, "languages")
	data := map[string]any{"overrides": overrides}
	stringattr.Get(m.FlowID, data, "flowId")
	return data
}

func (m *FlowTranslationsModel) SetValues(h *helpers.Handler, data map[string]any) {
	overrides, _ := data["overrides"].(map[string]any)
	if flowID, _ := overrides["flowId"].(string); flowID != "" {
		stringattr.Set(&m.FlowID, overrides, "flowId", stringattr.SkipIfAlreadySet)
	} else if m.FlowID.IsNull() || m.FlowID.IsUnknown() {
		m.FlowID = m.ID
	}
	languages, _ := overrides["languages"].(map[string]any)
	if languages == nil {
		languages = map[string]any{}
	}
	jsonattr.Set(&m.Languages, map[string]any{"languages": languages}, "languages")
}

func (m *FlowTranslationsModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Languages) || m.Languages.IsNull() {
		return
	}
	var languages map[string]map[string]string
	if err := json.Unmarshal([]byte(m.Languages.ValueString()), &languages); err != nil {
		h.Invalid("The languages attribute must be an object holding the texts of each language as an object of strings by translation key, for example {\"he\": {\"<key>\": \"<text>\"}}")
	}
}

func (m *FlowTranslationsModel) GetID() stringattr.Type        { return m.ID }
func (m *FlowTranslationsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *FlowTranslationsModel) GetProjectID() stringattr.Type { return m.ProjectID }
