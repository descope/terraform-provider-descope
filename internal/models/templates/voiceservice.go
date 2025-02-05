package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var VoiceServiceValidator = objectattr.NewValidator[VoiceServiceModel]("must have unique template names and a valid configuration")

var VoiceServiceAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Optional(VoiceTemplateAttributes, VoiceTemplateValidator),
}

type VoiceServiceModel struct {
	Connector types.String          `tfsdk:"connector"`
	Templates []*VoiceTemplateModel `tfsdk:"templates"`
}

func (m *VoiceServiceModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	connector := m.Connector.ValueString()
	if ref := h.Refs.Get(helpers.ConnectorReferenceKey, connector); ref != nil {
		h.Log("Setting voiceServiceProvider reference to connector '%s'", connector)
		data["voiceServiceProvider"] = ref.ProviderValue()
	} else {
		h.Error("Unknown connector reference", "No connector named '%s' for voice service was defined", connector)
	}
	listattr.Get(m.Templates, data, "voiceTemplates", h)
	return data
}

func (m *VoiceServiceModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Connector, data, "voiceServiceProvider")
	// update known templates with their new values
	for _, template := range m.Templates {
		name := template.Name.ValueString()
		h.Log("Looking for voice template named '%s'", name)
		if id, ok := requireTemplateID(h, data, "voiceTemplates", name); ok {
			value := types.StringValue(id)
			if !template.ID.Equal(value) {
				h.Log("Setting new ID '%s' for voice template named '%s'", id, name)
				template.ID = value
			} else {
				h.Log("Keeping existing ID '%s' for voice template named '%s'", id, name)
			}
		}
	}
	// we allow to set templates on import
	if m.Templates == nil && helpers.IsImport(h.Ctx) {
		listattr.Set(&m.Templates, data, "voiceTemplates", h)
	}
}

func (m *VoiceServiceModel) Validate(h *helpers.Handler) {
	hasActive := false
	names := map[string]int{}
	for _, v := range m.Templates {
		hasActive = hasActive || v.Active.ValueBool()
		names[v.Name.ValueString()] += 1
	}

	for k, v := range names {
		if v > 1 {
			h.Error("Template names must be unique", "The template name '%s' is used %d times", k, v)
		}
	}

	connector := m.Connector.ValueString()
	if connector == "" {
		h.Error("Invalid voice service connector", "The connector attribute must be set to Descope or the name of a connector")
	} else if hasActive && connector == helpers.DescopeConnector {
		h.Error("Invalid voice service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}

func (m *VoiceServiceModel) SetReferences(h *helpers.Handler) {
	replaceConnectorIDWithReference(&m.Connector, h)
}
