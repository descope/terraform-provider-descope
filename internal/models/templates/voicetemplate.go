package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var VoiceTemplateValidator = objectattr.NewValidator[VoiceTemplateModel]("must have a valid name")

var VoiceTemplateAttributes = map[string]schema.Attribute{
	"active": boolattr.Default(false),
	"id":     stringattr.Identifier(),
	"name":   stringattr.Required(),
	"body":   stringattr.Required(),
}

type VoiceTemplateModel struct {
	Active types.Bool   `tfsdk:"active"`
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Body   types.String `tfsdk:"body"`
}

func (m *VoiceTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	boolattr.Get(m.Active, data, "active")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Body, data, "body")
	return data
}

func (m *VoiceTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all template values are specified in the configuration
}

func (m *VoiceTemplateModel) Validate(h *helpers.Handler) {
	if m.Name.ValueString() == helpers.DescopeTemplate || m.ID.ValueString() == helpers.DescopeTemplate {
		h.Error("Invalid voice template", "Cannot use 'System' as the name or id of a template.")
	}
}
