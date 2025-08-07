package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EmailServiceValidator = objattr.NewValidator[EmailServiceModel]("must have unique template names and a valid configuration")

var EmailServiceAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Default[EmailTemplateModel](EmailTemplateAttributes, EmailTemplateValidator),
}

type EmailServiceModel struct {
	Connector stringattr.Type                   `tfsdk:"connector"`
	Templates listattr.Type[EmailTemplateModel] `tfsdk:"templates"`
}

func (m *EmailServiceModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	connector := m.Connector.ValueString()
	if ref := h.Refs.Get(helpers.ConnectorReferenceKey, connector); ref != nil {
		h.Log("Setting emailServiceProvider reference to connector '%s'", connector)
		data["emailServiceProvider"] = ref.ProviderValue()
	} else {
		h.Error("Unknown connector reference", "No connector named '%s' for email service was defined", connector)
	}
	listattr.Get(m.Templates, data, "emailTemplates", h)
	return data
}

func (m *EmailServiceModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Connector, data, "emailServiceProvider")

	if m.Templates.IsEmpty() {
		listattr.Set(&m.Templates, data, "emailTemplates", h)
	} else {
		for template := range listattr.MutatingIterator(&m.Templates, h) {
			name := template.Name.ValueString()
			h.Log("Looking for email template named '%s'", name)
			if id, ok := requireTemplateID(h, data, "emailTemplates", name); ok {
				value := stringattr.Value(id)
				if !template.ID.Equal(value) {
					h.Log("Setting new ID '%s' for email template named '%s'", id, name)
					template.ID = value
				} else {
					h.Log("Keeping existing ID '%s' for email template named '%s'", id, name)
				}
			} else if template.ID.ValueString() == "" {
				h.Error("Template not found", "Expected to find email template to match with '%s' template", name)
			}
		}
	}
}

func (m *EmailServiceModel) Validate(h *helpers.Handler) {
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
		h.Error("Invalid email service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}

func (m *EmailServiceModel) UpdateReferences(h *helpers.Handler) {
	replaceConnectorIDWithReference(&m.Connector, h)
}
