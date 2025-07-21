package attributes

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/iancoleman/strcase"
)

// Validators

var attributeTypeValidator = stringvalidator.OneOf("string", "number", "boolean", "singleselect", "multiselect", "date")

// Base

type AttributeModel struct {
	ID            stringattr.Type `tfsdk:"id"`
	Name          stringattr.Type `tfsdk:"name"`
	Type          stringattr.Type `tfsdk:"type"`
	SelectOptions strsetattr.Type `tfsdk:"select_options"`
}

func (m *AttributeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "name")
	stringattr.Get(m.Name, data, "displayName")
	stringattr.Get(m.Type, data, "type")
	getOptions(m.SelectOptions, data, "options", h)
	return data
}

func (m *AttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "name")
	stringattr.Set(&m.Name, data, "displayName")
	stringattr.Set(&m.Type, data, "type")
	setOptions(&m.SelectOptions, data, "options", h)
}

func (m *AttributeModel) Modify(h *helpers.Handler) {
	if v := m.Name.ValueString(); v != "" && m.ID.IsUnknown() {
		id := strcase.ToLowerCamel(v)
		if len(id) > 20 {
			id = id[:20]
		}
		m.ID = stringattr.Value(id)
	}
}

// Matching

func (m *AttributeModel) GetName() stringattr.Type {
	return m.Name
}

func (m *AttributeModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AttributeModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Options

func getOptions(s strsetattr.Type, data map[string]any, key string, h *helpers.Handler) {
	options := []map[string]any{}
	for option := range strsetattr.Iterator(s, h) {
		options = append(options, map[string]any{"label": option, "value": option})
	}
	data[key] = options
}

func setOptions(s *strsetattr.Type, data map[string]any, key string, _ *helpers.Handler) {
	result := []string{}
	if vs, ok := data[key].([]any); ok {
		for _, v := range vs {
			if os, ok := v.(map[string]any); ok {
				if option, ok := os["label"].(string); ok {
					result = append(result, option)
				}
			}
		}
	}
	*s = strsetattr.Value(result)
}
