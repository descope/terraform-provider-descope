package project

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ScopeClaimMappingAttributes is the project-wide variant of the scope claim mapping. Unlike the app
// variants it has neither `use_project_mapping` (it *is* the project mapping) nor `mandatory`
// (mandatory is an app-level concern).
var ScopeClaimMappingAttributes = map[string]schema.Attribute{
	"scope":       stringattr.Required(),
	"claims":      strmapattr.Default(),
	"description": stringattr.Default(""),
}

type ScopeClaimMappingModel struct {
	Scope       stringattr.Type `tfsdk:"scope"`
	Claims      strmapattr.Type `tfsdk:"claims"`
	Description stringattr.Type `tfsdk:"description"`
}

func (m *ScopeClaimMappingModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Scope, data, "scope")
	strmapattr.Get(m.Claims, data, "claims", h)
	stringattr.Get(m.Description, data, "description")
	return data
}

func (m *ScopeClaimMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Scope, data, "scope")
	strmapattr.Set(&m.Claims, data, "claims", h)
	stringattr.Set(&m.Description, data, "description")
}
