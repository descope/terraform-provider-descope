package project

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ProjectAttributes = map[string]schema.Attribute{
	"id":                  stringattr.Identifier(),
	"name":                stringattr.Required(),
	"environment":         stringattr.Optional(stringvalidator.OneOf("", "production")),
	"deletion_protection": boolattr.Tristate(),
	"tags":                strsetattr.Optional(stringvalidator.LengthBetween(1, 50)),
}

var Schema = schema.Schema{
	MarkdownDescription: "Manages a Descope project and its core attributes. The project's configuration is managed with the standalone descope resources that reference it by ID.",
	Attributes:          ProjectAttributes,
}

type ProjectModel struct {
	ID                 stringattr.Type `tfsdk:"id"`
	Name               stringattr.Type `tfsdk:"name"`
	Environment        stringattr.Type `tfsdk:"environment"`
	DeletionProtection boolattr.Type   `tfsdk:"deletion_protection"`
	Tags               strsetattr.Type `tfsdk:"tags"`
}

func (m *ProjectModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Environment, data, "environment")
	strsetattr.Get(m.Tags, data, "tags", h)
	return data
}

func (m *ProjectModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Environment, data, "environment")
	strsetattr.Set(&m.Tags, data, "tags", h)
}

func (m *ProjectModel) GetID() stringattr.Type {
	return m.ID
}

func (m *ProjectModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *ProjectModel) GetProjectID() stringattr.Type {
	return m.ID
}

func (m *ProjectModel) DeletionProtectionDefault(_ context.Context) bool {
	return m.Environment.ValueString() == "production"
}
