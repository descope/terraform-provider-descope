package jwttemplates

import (
	"context"
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

// JWT Template Attributes

var JWTTemplateAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"name":               stringattr.Required(),
	"description":        stringattr.Default(""),
	"auth_schema":        stringattr.Default("default", stringvalidator.OneOf("default", "tenantOnly", "none")),
	"conformance_issuer": boolattr.Default(false),
	"template":           stringattr.Required(),
	"type":               stringattr.Required(stringvalidator.OneOf("key", "user")),
}

type JWTTemplateModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	AuthSchema        types.String `tfsdk:"auth_schema"`
	ConformanceIssuer types.Bool   `tfsdk:"conformance_issuer"`
	Template          types.String `tfsdk:"template"`
	Type              types.String `tfsdk:"type"`
}

func (m *JWTTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.AuthSchema, data, "authSchema")
	boolattr.Get(m.ConformanceIssuer, data, "conformanceIssuer")
	stringattr.Get(m.Type, data, "type")

	// convert template JSON string to map
	template := map[string]any{}
	if err := json.Unmarshal([]byte(m.Template.ValueString()), &template); err != nil {
		h.Error("Unexpected template structure", "Could not deserialize template attribute, json expected")
	}
	data["template"] = template

	// use the name as a lookup key to set the JWT template reference or existing id
	templateName := m.Name.ValueString()
	if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, templateName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for JWT template '%s' to: %s", templateName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown JWT template reference", "No JWT template named '"+templateName+"' was defined")
	}

	return data
}

func (m *JWTTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all JWT template values are specified in the configuration
}
