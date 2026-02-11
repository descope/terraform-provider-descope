package inboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ApplicationScopeAttributes = map[string]schema.Attribute{
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"optional":    boolattr.Default(false),
	"values":      strlistattr.Default(),
}

type ApplicationScopeModel struct {
	Name        stringattr.Type  `tfsdk:"name"`
	Description stringattr.Type  `tfsdk:"description"`
	Optional    boolattr.Type    `tfsdk:"optional"`
	ScopeValues strlistattr.Type `tfsdk:"values"`
}

func (m *ApplicationScopeModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	boolattr.Get(m.Optional, data, "optional")
	strlistattr.Get(m.ScopeValues, data, "values", h)
	return data
}

func (m *ApplicationScopeModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	boolattr.Set(&m.Optional, data, "optional")
	strlistattr.Set(&m.ScopeValues, data, "values", h)
}
