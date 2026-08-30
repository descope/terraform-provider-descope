package voicetemplate

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single voice call message template for an authentication method in a Descope project. The active template is selected on the method's settings resource, e.g. the `voice_template_id` attribute of `descope_otp_settings`.",
	Attributes:          VoiceTemplateAttributes,
}

var VoiceTemplateAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"method":     stringattr.Required(stringvalidator.OneOf("otp"), stringplanmodifier.RequiresReplace()),
	"name":       stringattr.Required(),
	"body":       stringattr.Required(),
}

type VoiceTemplateModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Method    stringattr.Type `tfsdk:"method"`
	Name      stringattr.Type `tfsdk:"name"`
	Body      stringattr.Type `tfsdk:"body"`
}

func (m *VoiceTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Method, data, "method")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Body, data, "body")
	return data
}

func (m *VoiceTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Method, data, "method")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Body, data, "body", stringattr.SkipIfAlreadySet) // the backend canonicalizes template macros so return value might be different
}

func (m *VoiceTemplateModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Name) {
		return
	}
	if m.Name.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid voice template", "Cannot use 'System' as the name of a template")
	}
}

func (m *VoiceTemplateModel) GetID() stringattr.Type {
	return m.ID
}

func (m *VoiceTemplateModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *VoiceTemplateModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

func (m *VoiceTemplateModel) GetMethod() stringattr.Type {
	return m.Method
}
