package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TextServiceValidator = objattr.NewValidator[TextServiceModel]("must have unique template names and a valid configuration")

var TextServiceAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Default[TextTemplateModel](TextTemplateAttributes, TextTemplateValidator),
}

type TextServiceModel struct {
	Connector stringattr.Type                  `tfsdk:"connector"`
	Templates listattr.Type[TextTemplateModel] `tfsdk:"templates"`
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
	stringattr.Set(&m.Connector, data, "textServiceProvider")

	if m.Templates.IsEmpty() {
		listattr.Set(&m.Templates, data, "textTemplates", h)
	} else {
		for template := range listattr.MutatingIterator(&m.Templates, h) {
			name := template.Name.ValueString()
			h.Log("Looking for text template named '%s'", name)
			if id, ok := requireTemplateID(h, data, "textTemplates", name); ok {
				value := stringattr.Value(id)
				if !template.ID.Equal(value) {
					h.Log("Setting new ID '%s' for text template named '%s'", id, name)
					template.ID = value
				} else {
					h.Log("Keeping existing ID '%s' for text template named '%s'", id, name)
				}
			} else if template.ID.ValueString() == "" {
				h.Error("Template not found", "Expected to find text template to match with '%s' template", name)
			}
		}
	}
}

func (m *TextServiceModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Connector, m.Templates) {
		return
	}

	hasActive := false
	names := map[string]int{}
	for v := range listattr.Iterator(m.Templates, h) {
		hasActive = hasActive || v.Active.ValueBool()
		names[v.Name.ValueString()] += 1
	}

	for k, v := range names {
		if v > 1 {
			h.Error("Template names must be unique", "The template name '%s' is used %d times", k, v)
		}
	}

	connector := m.Connector.ValueString()
	if hasActive && connector == helpers.DescopeConnector {
		h.Error("Invalid text service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}

func (m *TextServiceModel) UpdateReferences(h *helpers.Handler) {
	replaceConnectorIDWithReference(&m.Connector, h)
}
