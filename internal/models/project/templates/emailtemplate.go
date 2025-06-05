package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EmailTemplateValidator = objattr.NewValidator[EmailTemplateModel]("must have a valid name and contain at least one body attribute set")

var EmailTemplateAttributes = map[string]schema.Attribute{
	"active":              boolattr.Default(false),
	"id":                  stringattr.Identifier(),
	"name":                stringattr.Required(),
	"subject":             stringattr.Required(),
	"html_body":           stringattr.Default(""),
	"plain_text_body":     stringattr.Default(""),
	"use_plain_text_body": boolattr.Default(false),
}

type EmailTemplateModel struct {
	Active           boolattr.Type   `tfsdk:"active"`
	ID               stringattr.Type `tfsdk:"id"`
	Name             stringattr.Type `tfsdk:"name"`
	Subject          stringattr.Type `tfsdk:"subject"`
	HTMLBody         stringattr.Type `tfsdk:"html_body"`
	PlainTextBody    stringattr.Type `tfsdk:"plain_text_body"`
	UsePlainTextBody boolattr.Type   `tfsdk:"use_plain_text_body"`
}

func (m *EmailTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	boolattr.Get(m.Active, data, "active")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Subject, data, "subject")
	stringattr.Get(m.HTMLBody, data, "body")
	stringattr.Get(m.PlainTextBody, data, "bodyPlainText")
	boolattr.Get(m.UsePlainTextBody, data, "useBodyPlainText")
	return data
}

func (m *EmailTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	boolattr.Set(&m.Active, data, "active")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Subject, data, "subject")
	stringattr.Set(&m.HTMLBody, data, "body")
	stringattr.Set(&m.PlainTextBody, data, "bodyPlainText")
	boolattr.Set(&m.UsePlainTextBody, data, "useBodyPlainText")
}

func (m *EmailTemplateModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Name, m.UsePlainTextBody, m.PlainTextBody, m.HTMLBody) {
		return // skip validation if there are unknown values
	}
	if m.Name.ValueString() == helpers.DescopeTemplate || m.ID.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid email template", "Cannot use 'System' as the name or id of a template")
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
