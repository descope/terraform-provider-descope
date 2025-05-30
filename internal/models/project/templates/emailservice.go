package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EmailServiceValidator = objectattr.NewValidator[EmailServiceModel]("must have unique template names and a valid configuration")

var EmailServiceAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Optional(EmailTemplateAttributes, EmailTemplateValidator),
}

type EmailServiceModel struct {
	Connector types.String          `tfsdk:"connector"`
	Templates []*EmailTemplateModel `tfsdk:"templates"`
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
	// update known templates with their new values
	for _, template := range m.Templates {
		name := template.Name.ValueString()
		h.Log("Looking for email template named '%s'", name)
		if id, ok := requireTemplateID(h, data, "emailTemplates", name); ok {
			value := types.StringValue(id)
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
	// we allow to set templates on import
	if m.Templates == nil && helpers.IsImport(h.Ctx) {
		listattr.Set(&m.Templates, data, "emailTemplates", h)
	}
}

func (m *EmailServiceModel) Validate(h *helpers.Handler) {
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
	if hasActive && connector == helpers.DescopeConnector {
		h.Error("Invalid email service connector", "The connector attribute must not be set to Descope if any template is marked as active")
	}
}

func (m *EmailServiceModel) SetReferences(h *helpers.Handler) {
	replaceConnectorIDWithReference(&m.Connector, h)
}
