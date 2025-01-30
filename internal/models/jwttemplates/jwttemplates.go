package jwttemplates

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var JWTTemplatesValidator = objectattr.NewValidator[JWTTemplatesModel]("must have unique template names")

var JWTTemplatesAttributes = map[string]schema.Attribute{
	"user_templates":       listattr.Optional(JWTTemplateAttributes),
	"access_key_templates": listattr.Optional(JWTTemplateAttributes),
}

type JWTTemplatesModel struct {
	UserTemplates      []*JWTTemplateModel `tfsdk:"user_templates"`
	AccessKeyTemplates []*JWTTemplateModel `tfsdk:"access_key_templates"`
}

func (m *JWTTemplatesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.UserTemplates, data, "userTemplates", h)
	listattr.Get(m.AccessKeyTemplates, data, "keyTemplates", h)
	return data
}

func (m *JWTTemplatesModel) SetValues(h *helpers.Handler, data map[string]any) {
	// update templates with their new values
	m.setTemplateValues(h, data, "userTemplates", m.UserTemplates)
	m.setTemplateValues(h, data, "keyTemplates", m.AccessKeyTemplates)
	// we allow setting the templates on import
	if m.UserTemplates == nil && helpers.IsImport(h.Ctx) {
		listattr.Set(&m.UserTemplates, data, "userTemplates", h)
	}
	if m.AccessKeyTemplates == nil && helpers.IsImport(h.Ctx) {
		listattr.Set(&m.AccessKeyTemplates, data, "keyTemplates", h)
	}
}

func (m *JWTTemplatesModel) References(ctx context.Context) helpers.ReferencesMap {
	refs := helpers.ReferencesMap{}
	for _, v := range m.UserTemplates {
		refs.Add(helpers.JWTTemplateReferenceKey, "user", v.ID.ValueString(), v.Name.ValueString())
	}
	for _, v := range m.AccessKeyTemplates {
		refs.Add(helpers.JWTTemplateReferenceKey, "key", v.ID.ValueString(), v.Name.ValueString())
	}
	return refs
}

func (m *JWTTemplatesModel) Validate(h *helpers.Handler) {
	names := map[string]int{}
	for _, v := range m.UserTemplates {
		names[v.Name.ValueString()] += 1
	}
	for _, v := range m.AccessKeyTemplates {
		names[v.Name.ValueString()] += 1
	}
	for k, v := range names {
		if v > 1 {
			h.Error("JWT template names must be unique", "The JWT template name '%s' is used %d times", k, v)
		}
	}
}

func (m *JWTTemplatesModel) setTemplateValues(h *helpers.Handler, data map[string]any, key string, list []*JWTTemplateModel) {
	templates := m.getTemplateIDs(data, key)
	for _, template := range list {
		name := template.Name.ValueString()
		id, found := templates[name]
		if found {
			value := types.StringValue(id)
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
