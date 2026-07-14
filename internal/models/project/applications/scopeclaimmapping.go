package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ScopeClaimMappingAttributes is the federated (SSO/OIDC) variant of the scope claim mapping. Unlike
// the inbound third-party app variant it intentionally omits `mandatory`: mandatory scopes are only
// enforced for inbound apps, never for federated applications.
var ScopeClaimMappingAttributes = map[string]schema.Attribute{
	"scope":               stringattr.Required(),
	"claims":              strmapattr.Default(),
	"description":         stringattr.Default(""),
	"use_project_mapping": boolattr.Default(false),
}

type ScopeClaimMappingModel struct {
	Scope             stringattr.Type `tfsdk:"scope"`
	Claims            strmapattr.Type `tfsdk:"claims"`
	Description       stringattr.Type `tfsdk:"description"`
	UseProjectMapping boolattr.Type   `tfsdk:"use_project_mapping"`
}

func (m *ScopeClaimMappingModel) Values(h *helpers.Handler) map[string]any {
	// A scope inheriting the project-wide mapping (use_project_mapping = true) must not define its own
	// claims — the project-wide mapping's claims are used instead, so per-app claims would be ignored.
	if m.UseProjectMapping.ValueBool() {
		if claims := helpers.Require(m.Claims.ToMap(h.Ctx)); len(claims) > 0 {
			h.Error("Invalid scope claim mapping",
				"Scope %q sets use_project_mapping = true and therefore cannot define its own claims; the project-wide scope claim mapping is used instead.",
				m.Scope.ValueString())
		}
	}
	data := map[string]any{}
	stringattr.Get(m.Scope, data, "scope")
	strmapattr.Get(m.Claims, data, "claims", h)
	stringattr.Get(m.Description, data, "description")
	boolattr.Get(m.UseProjectMapping, data, "useProjectMapping")
	return data
}

func (m *ScopeClaimMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Scope, data, "scope")
	strmapattr.Set(&m.Claims, data, "claims", h)
	stringattr.Set(&m.Description, data, "description")
	boolattr.Set(&m.UseProjectMapping, data, "useProjectMapping")
}
