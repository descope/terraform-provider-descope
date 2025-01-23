package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TextServiceValidator = objectattr.NewValidator[TextServiceModel]("must have unique template names and a valid configuration")

var TextServiceAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Optional(TextTemplateAttributes, TextTemplateValidator),
}

type TextServiceModel struct {
	Connector types.String         `tfsdk:"connector"`
	Templates []*TextTemplateModel `tfsdk:"templates"`
}

func (m *TextServiceModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	connector := m.Connector.ValueString()
	if ref := h.Refs.Get(helpers.ConnectorReferenceKey, connector); ref != nil {
		h.Log("Setting textServiceProvider reference to connector '%s'", connector)
		data["textServiceProvider"] = ref.ProviderValue()
	} else {
		h.Error("Unknown connector reference", "No connector named '%s' for text service was defined", connector)
	}
	listattr.Get(m.Templates, data, "textTemplates", h)
	return data
}

func (m *TextServiceModel) SetValues(h *helpers.Handler, data map[string]any) {
	for _, template := range m.Templates {
		name := template.Name.ValueString()
		h.Log("Looking for text template named '%s'", name)
		if id, ok := requireTemplateID(h, data, "textTemplates", name); ok {
			value := types.StringValue(id)
			if !template.ID.Equal(value) {
				h.Log("Setting new ID '%s' for text template named '%s'", id, name)
				template.ID = value
			} else {
				h.Log("Keeping existing ID '%s' for text template named '%s'", id, name)
			}
		}
	}
	if m.Connector.ValueString() == "" {
		stringattr.Set(&m.Connector, data, "textServiceProvider")
	}
	if m.Templates == nil {
		m.Templates = []*TextTemplateModel{}
		listattr.Set(&m.Templates, data, "textTemplates", h)
	}
}

func (m *TextServiceModel) Validate(h *helpers.Handler) {
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
		h.Error("Invalid text service connector", "The connector attribute must be set to Descope or the name of a connector")
	} else if hasActive && connector == helpers.DescopeConnector {
		h.Error("Invalid text service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}
