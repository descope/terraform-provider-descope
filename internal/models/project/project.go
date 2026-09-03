package project

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
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

type ProjectModel struct {
	ID                 stringattr.Type `tfsdk:"id"`
	Name               stringattr.Type `tfsdk:"name"`
	Environment        stringattr.Type `tfsdk:"environment"`
	DeletionProtection boolattr.Type   `tfsdk:"deletion_protection"`
	Tags               strsetattr.Type `tfsdk:"tags"`
}

func (m *ProjectModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["version"] = helpers.ModelVersion
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Environment, data, "environment")
	strsetattr.Get(m.Tags, data, "tags", h)
	return data
}

func (m *ProjectModel) SetValues(h *helpers.Handler, data map[string]any) {
	if v, ok := data["version"].(float64); ok {
		helpers.EnsureModelVersion(v, h.Diagnostics)
	}

	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Environment, data, "environment")
	strsetattr.Set(&m.Tags, data, "tags", h)
}

func (m *ProjectModel) DeletionProtectionDefault(_ context.Context) bool {
	return m.Environment.ValueString() == "production"
}
