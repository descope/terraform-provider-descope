package jwttemplates

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var JWTTemplatesAttributes = map[string]schema.Attribute{
	"templates": listattr.Optional(JWTTemplateAttributes),
}

type JWTTemplatesModel struct {
	Templates []*JWTTemplateModel `tfsdk:"templates"`
}

func (m *JWTTemplatesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Templates, data, "templates", h)
	return data
}

func (m *JWTTemplatesModel) SetValues(h *helpers.Handler, data map[string]any) {
	templates := m.getTemplateIDs(data)
	for _, template := range m.Templates {
		name := template.Name.ValueString()
		id, found := templates[name]
		if found {
			value := types.StringValue(id)
			if !template.ID.Equal(value) {
				h.Log("Setting new ID '" + id + "' for JWT template named '" + name + "'")
				template.ID = value
			} else {
				h.Log("Keeping existing ID '" + id + "' for JWT template named '" + name + "'")
			}
		} else {
			h.Error("JWT template not found", "Expected to find JWT template to match with '"+name+"'")
		}
	}
}

func (m *JWTTemplatesModel) References(ctx context.Context) helpers.ReferencesMap {
	refs := helpers.ReferencesMap{}
	for _, v := range m.Templates {
		refs.Add(helpers.JWTTemplateReferenceKey, "", v.ID.ValueString(), v.Name.ValueString())
	}
	return refs
}

func (m *JWTTemplatesModel) getTemplateIDs(data map[string]any) (templates map[string]string) {
	templates = map[string]string{}

	rs, _ := data["templates"].([]any)
	for _, v := range rs {
		if r, ok := v.(map[string]any); ok {
			id, _ := r["id"].(string)
			name, _ := r["name"].(string)
			templates[name] = id
		}
	}

	return
}
