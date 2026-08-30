package texttemplate

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single text message template for an authentication method in a Descope project. The active template is selected on the method's settings resource, e.g. the `text_template_id` attribute of `descope_magiclink_settings`.",
	Attributes:          TextTemplateAttributes,
}

var TextTemplateAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"method":     stringattr.Required(stringvalidator.OneOf("magiclink", "otp"), stringplanmodifier.RequiresReplace()),
	"name":       stringattr.Required(),
	"body":       stringattr.Required(),
}

type TextTemplateModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Method    stringattr.Type `tfsdk:"method"`
	Name      stringattr.Type `tfsdk:"name"`
	Body      stringattr.Type `tfsdk:"body"`
}

func (m *TextTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Method, data, "method")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Body, data, "body")
	return data
}

func (m *TextTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Method, data, "method")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Body, data, "body", stringattr.SkipIfAlreadySet) // the backend canonicalizes template macros so return value might be different
}

func (m *TextTemplateModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Name) {
		return
	}
	if m.Name.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid text template", "Cannot use 'System' as the name of a template")
	}
}

func (m *TextTemplateModel) GetID() stringattr.Type {
	return m.ID
}

func (m *TextTemplateModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *TextTemplateModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

func (m *TextTemplateModel) GetMethod() stringattr.Type {
	return m.Method
}
