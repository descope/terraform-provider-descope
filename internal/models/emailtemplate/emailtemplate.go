package emailtemplate

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single email message template for an authentication method in a Descope project. The active template is selected on the method's settings resource, e.g. the `email_template_id` attribute of `descope_magiclink_settings`.",
	Attributes:          EmailTemplateAttributes,
}

var EmailTemplateAttributes = map[string]schema.Attribute{
	"id":                  stringattr.Identifier(),
	"project_id":          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"method":              stringattr.Required(stringvalidator.OneOf("enchantedlink", "magiclink", "otp", "password", "sso"), stringplanmodifier.RequiresReplace()),
	"name":                stringattr.Required(),
	"subject":             stringattr.Required(),
	"html_body":           stringattr.Default(""),
	"plain_text_body":     stringattr.Default(""),
	"use_plain_text_body": boolattr.Default(false),
}

type EmailTemplateModel struct {
	ID               stringattr.Type `tfsdk:"id"`
	ProjectID        stringattr.Type `tfsdk:"project_id"`
	Method           stringattr.Type `tfsdk:"method"`
	Name             stringattr.Type `tfsdk:"name"`
	Subject          stringattr.Type `tfsdk:"subject"`
	HTMLBody         stringattr.Type `tfsdk:"html_body"`
	PlainTextBody    stringattr.Type `tfsdk:"plain_text_body"`
	UsePlainTextBody boolattr.Type   `tfsdk:"use_plain_text_body"`
}

func (m *EmailTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Method, data, "method")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Subject, data, "subject")
	stringattr.Get(m.HTMLBody, data, "body")
	stringattr.Get(m.PlainTextBody, data, "bodyPlainText")
	boolattr.Get(m.UsePlainTextBody, data, "useBodyPlainText")
	return data
}

func (m *EmailTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Method, data, "method")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Subject, data, "subject", stringattr.SkipIfAlreadySet) // the backend canonicalizes template macros so return value might be different
	stringattr.Set(&m.HTMLBody, data, "body", stringattr.SkipIfAlreadySet)
	stringattr.Set(&m.PlainTextBody, data, "bodyPlainText", stringattr.SkipIfAlreadySet)
	boolattr.Set(&m.UsePlainTextBody, data, "useBodyPlainText")
}

func (m *EmailTemplateModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Name, m.UsePlainTextBody, m.PlainTextBody, m.HTMLBody) {
		return
	}
	if m.Name.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid email template", "Cannot use 'System' as the name of a template")
	}
	if m.UsePlainTextBody.ValueBool() {
		if m.PlainTextBody.ValueString() == "" {
			h.Missing("The plain_text_body attribute is required when use_plain_text_body is enabled")
		}
	} else {
		if m.HTMLBody.ValueString() == "" {
			h.Missing("The html_body attribute is required unless use_plain_text_body is enabled")
		}
	}
}

func (m *EmailTemplateModel) GetID() stringattr.Type {
	return m.ID
}

func (m *EmailTemplateModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *EmailTemplateModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

func (m *EmailTemplateModel) GetMethod() stringattr.Type {
	return m.Method
}
