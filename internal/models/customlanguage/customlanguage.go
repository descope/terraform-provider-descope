package customlanguage

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var CustomLanguageAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	// The locale code (language + optional region) is immutable; changing it replaces the resource.
	"language": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"region":   stringattr.Optional(stringplanmodifier.RequiresReplace()),
	"name":     stringattr.Required(),
}

var Schema = schema.Schema{
	Attributes: CustomLanguageAttributes,
}

type CustomLanguageModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Language  stringattr.Type `tfsdk:"language"`
	Region    stringattr.Type `tfsdk:"region"`
	Name      stringattr.Type `tfsdk:"name"`
}

func (m *CustomLanguageModel) Values(_ *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Language, data, "language")
	stringattr.Get(m.Region, data, "region")
	stringattr.Get(m.Name, data, "name")
	return data
}

func (m *CustomLanguageModel) SetValues(_ *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Language, data, "language")
	stringattr.Set(&m.Region, data, "region")
	stringattr.Set(&m.Name, data, "name")
}

func (m *CustomLanguageModel) GetID() stringattr.Type {
	return m.ID
}

func (m *CustomLanguageModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *CustomLanguageModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
