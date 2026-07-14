package inboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ScopeClaimMappingAttributes = map[string]schema.Attribute{
	"scope":               stringattr.Required(),
	"claims":              strmapattr.Default(),
	"description":         stringattr.Default(""),
	"use_project_mapping": boolattr.Default(false),
	"mandatory":           boolattr.Default(false),
}

type ScopeClaimMappingModel struct {
	Scope             stringattr.Type `tfsdk:"scope"`
	Claims            strmapattr.Type `tfsdk:"claims"`
	Description       stringattr.Type `tfsdk:"description"`
	UseProjectMapping boolattr.Type   `tfsdk:"use_project_mapping"`
	Mandatory         boolattr.Type   `tfsdk:"mandatory"`
}

func (m *ScopeClaimMappingModel) Values(h *helpers.Handler) map[string]any {
	validateScopeClaimMapping(h, m.Scope, m.Claims, m.UseProjectMapping)
	data := map[string]any{}
	stringattr.Get(m.Scope, data, "scope")
	strmapattr.Get(m.Claims, data, "claims", h)
	stringattr.Get(m.Description, data, "description")
	boolattr.Get(m.UseProjectMapping, data, "useProjectMapping")
	boolattr.Get(m.Mandatory, data, "mandatory")
	return data
}

// validateScopeClaimMapping enforces that a scope inheriting the project-wide mapping
// (use_project_mapping = true) does not also define its own claims — the project-wide mapping's
// claims are used instead, so per-app claims would be silently ignored. Shared by the inbound and
// federated scope-claim-mapping models.
func validateScopeClaimMapping(h *helpers.Handler, scope stringattr.Type, claims strmapattr.Type, useProjectMapping boolattr.Type) {
	if !useProjectMapping.ValueBool() {
		return
	}
	if claimsMap := helpers.Require(claims.ToMap(h.Ctx)); len(claimsMap) > 0 {
		h.Error("Invalid scope claim mapping",
			"Scope %q sets use_project_mapping = true and therefore cannot define its own claims; the project-wide scope claim mapping is used instead.",
			scope.ValueString())
	}
}

func (m *ScopeClaimMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Scope, data, "scope")
	strmapattr.Set(&m.Claims, data, "claims", h)
	stringattr.Set(&m.Description, data, "description")
	boolattr.Set(&m.UseProjectMapping, data, "useProjectMapping")
	boolattr.Set(&m.Mandatory, data, "mandatory")
}
