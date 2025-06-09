package jwttemplates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var JWTTemplatesValidator = objattr.NewValidator[JWTTemplatesModel]("must have unique template names")

var JWTTemplatesAttributes = map[string]schema.Attribute{
	"user_templates":       listattr.Default[JWTTemplateModel](JWTTemplateAttributes),
	"access_key_templates": listattr.Default[JWTTemplateModel](JWTTemplateAttributes),
}

var JWTTemplatesDefault = &JWTTemplatesModel{
	UserTemplates:      listattr.Empty[JWTTemplateModel](),
	AccessKeyTemplates: listattr.Empty[JWTTemplateModel](),
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
	if m.UserTemplates.IsUnknown() {
		listattr.Set(&m.UserTemplates, data, "userTemplates", h)
	} else {
		m.setTemplateValues(h, data, "userTemplates", &m.UserTemplates)
	}
	if m.AccessKeyTemplates.IsUnknown() {
		listattr.Set(&m.AccessKeyTemplates, data, "keyTemplates", h)
	} else {
		m.setTemplateValues(h, data, "keyTemplates", &m.AccessKeyTemplates)
	}
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

func (m *JWTTemplatesModel) setTemplateValues(h *helpers.Handler, data map[string]any, key string, list *listattr.Type[JWTTemplateModel]) {
	templates := m.getTemplateIDs(data, key)
	for template := range listattr.MutatingIterator(list, h) {
		name := template.Name.ValueString()
		id, found := templates[name]
		if found {
			value := stringattr.Value(id)
			if !template.ID.Equal(value) {
				h.Log("Setting new ID '%s' for %s JWT template named '%s'", id, key, name)
				template.ID = value
			} else {
				h.Log("Keeping existing ID '%s' for %s JWT template named '%s'", id, key, name)
			}
		} else {
			h.Error("JWT template not found", "Expected to find %s JWT template to match with '%s'", key, name)
		}
	}
}

func (m *JWTTemplatesModel) getTemplateIDs(data map[string]any, key string) map[string]string {
	templates := map[string]string{}
	rs, _ := data[key].([]any)
	for _, v := range rs {
		if r, ok := v.(map[string]any); ok {
			id, _ := r["id"].(string)
			name, _ := r["name"].(string)
			templates[name] = id
		}
	}
	return templates
}
