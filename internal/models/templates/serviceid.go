package templates

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Email Service

var EmailServiceIDValidator = objattr.NewValidator[EmailServiceIDModel]("must have unique template names and a valid configuration")

var EmailServiceIDAttributes = map[string]schema.Attribute{
	"connector_id": stringattr.Default(helpers.DescopeConnector),
	"templates":    listattr.Default[EmailTemplateModel](EmailTemplateAttributes, EmailTemplateValidator),
}

type EmailServiceIDModel struct {
	ConnectorID stringattr.Type                   `tfsdk:"connector_id"`
	Templates   listattr.Type[EmailTemplateModel] `tfsdk:"templates"`
}

func (m *EmailServiceIDModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ConnectorID, data, "emailServiceProvider")
	listattr.Get(m.Templates, data, "emailTemplates", h)
	return data
}

func (m *EmailServiceIDModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ConnectorID, data, "emailServiceProvider")
	SetServiceConnectorID(&m.ConnectorID)

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

func (m *EmailServiceIDModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.ConnectorID, m.Templates) {
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

	// the connector_id default isn't applied yet during config validation, so an absent value counts as the Descope sentinel too
	if connectorID := m.ConnectorID.ValueString(); hasActive && (connectorID == "" || connectorID == helpers.DescopeConnector) {
		h.Error("Invalid email service connector", "The connector_id attribute must not be set to Descope if any template is marked as active")
	}
}

// Provider references come back from the server as `type:id` values, and an empty string means the built-in Descope delivery service.
func SetServiceConnectorID(s *stringattr.Type) {
	value := s.ValueString()
	if value == "" {
		*s = stringattr.Value(helpers.DescopeConnector)
	} else if _, id, found := strings.Cut(value, ":"); found {
		*s = stringattr.Value(id)
	}
}
