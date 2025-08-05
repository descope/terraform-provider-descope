package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

func getEmailValues(h *helpers.Handler, connector stringattr.Type, connectorKey string, templates listattr.Type[EmailTemplateModel], templatesKey string) map[string]any {
	data := map[string]any{}
	conn := connector.ValueString()
	if ref := h.Refs.Get(helpers.ConnectorReferenceKey, conn); ref != nil {
		h.Log("Setting %s reference to connector '%s'", connectorKey, conn)
		data[connectorKey] = ref.ProviderValue()
	} else {
		h.Error("Unknown connector reference", "No connector named '%s' for email service was defined", conn)
	}
	listattr.Get(templates, data, templatesKey, h)
	return data
}

func setEmailValues(h *helpers.Handler, data map[string]any, connector *stringattr.Type, connectorKey string, templates *listattr.Type[EmailTemplateModel], templatesKey string) {
	stringattr.Set(connector, data, connectorKey)

	if templates.IsEmpty() {
		listattr.Set(templates, data, templatesKey, h)
	} else {
		for template := range listattr.MutatingIterator(templates, h) {
			name := template.Name.ValueString()
			h.Log("Looking for email template named '%s'", name)
			if id, ok := requireTemplateID(h, data, templatesKey, name); ok {
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

func validateEmailValues(h *helpers.Handler, connector stringattr.Type, templates listattr.Type[EmailTemplateModel]) {
	if helpers.HasUnknownValues(connector, templates) {
		return
	}

	hasActive := false
	names := map[string]int{}
	for v := range listattr.Iterator(templates, h) {
		hasActive = hasActive || v.Active.ValueBool()
		names[v.Name.ValueString()] += 1
	}

	for k, v := range names {
		if v > 1 {
			h.Error("Template names must be unique", "The template name '%s' is used %d times", k, v)
		}
	}

	conn := connector.ValueString()
	if hasActive && conn == helpers.DescopeConnector {
		h.Error("Invalid email service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}
