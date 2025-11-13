package jwttemplates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var JWTTemplatesValidator = objattr.NewValidator[JWTTemplatesModel]("must have unique template names")

var JWTTemplatesAttributes = map[string]schema.Attribute{
	"user_templates":       listattr.Default[JWTTemplateModel](JWTTemplateAttributes),
	"access_key_templates": listattr.Default[JWTTemplateModel](JWTTemplateAttributes),
}

type JWTTemplatesModel struct {
	UserTemplates      listattr.Type[JWTTemplateModel] `tfsdk:"user_templates"`
	AccessKeyTemplates listattr.Type[JWTTemplateModel] `tfsdk:"access_key_templates"`
}

func (m *JWTTemplatesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.UserTemplates, data, "userTemplates", h)
	listattr.Get(m.AccessKeyTemplates, data, "keyTemplates", h)
	return data
}

func (m *JWTTemplatesModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.UserTemplates, data, "userTemplates", "name", h)
	listattr.SetMatchingNames(&m.AccessKeyTemplates, data, "keyTemplates", "name", h)
}

func (m *JWTTemplatesModel) CollectReferences(h *helpers.Handler) {
	for v := range listattr.Iterator(m.UserTemplates, h) {
		h.Refs.Add(helpers.JWTTemplateReferenceKey, "user", v.ID.ValueString(), v.Name.ValueString())
	}
	for v := range listattr.Iterator(m.AccessKeyTemplates, h) {
		h.Refs.Add(helpers.JWTTemplateReferenceKey, "key", v.ID.ValueString(), v.Name.ValueString())
	}
}

func (m *JWTTemplatesModel) Validate(h *helpers.Handler) {
	names := map[string]int{}
	for v := range listattr.Iterator(m.UserTemplates, h) {
		names[v.Name.ValueString()] += 1
	}
	for v := range listattr.Iterator(m.AccessKeyTemplates, h) {
		names[v.Name.ValueString()] += 1
	}
	for k, v := range names {
		if v > 1 {
			h.Conflict("The JWT template name '%s' is used %d times", k, v)
		}
	}
}
