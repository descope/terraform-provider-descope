package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TextTemplateValidator = objattr.NewValidator[TextTemplateModel]("must have a valid name")

var TextTemplateAttributes = map[string]schema.Attribute{
	"active": boolattr.Default(false),
	"id":     stringattr.Identifier(),
	"name":   stringattr.Required(),
	"body":   stringattr.Required(),
}

type TextTemplateModel struct {
	Active boolattr.Type   `tfsdk:"active"`
	ID     stringattr.Type `tfsdk:"id"`
	Name   stringattr.Type `tfsdk:"name"`
	Body   stringattr.Type `tfsdk:"body"`
}

func (m *TextTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	boolattr.Get(m.Active, data, "active")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Body, data, "body")
	return data
}

func (m *TextTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	boolattr.Set(&m.Active, data, "active")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Body, data, "body")
}

func (m *TextTemplateModel) Validate(h *helpers.Handler) {
	if m.Name.ValueString() == helpers.DescopeTemplate || m.ID.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid text template", "Cannot use 'System' as the name or id of a template")
	}
}
