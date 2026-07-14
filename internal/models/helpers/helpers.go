package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Used as a sentinel value when the JSON values for an object are at the root of the map.
const RootKey string = ""

// Require is a helper function that panics if the provided diagnostics contain errors.
func Require[T any](v T, diags diag.Diagnostics) T {
	if errs := diags.Errors(); len(errs) > 0 {
		panic(fmt.Sprintf("%s: %s", errs[0].Summary(), errs[0].Detail()))
	}
	return v
}

// ValidateScopeClaimMapping enforces that a scope inheriting the project-wide mapping
// (use_project_mapping = true) does not also define its own claims — the project-wide mapping's
// claims are used instead, so per-app claims would be silently ignored. Shared by the inbound and
// federated scope-claim-mapping models so the rule and message stay identical across both.
func ValidateScopeClaimMapping(h *Handler, scope string, useProjectMapping bool, claimCount int) {
	if useProjectMapping && claimCount > 0 {
		h.Error("Invalid scope claim mapping",
			"Scope %q sets use_project_mapping = true and therefore cannot define its own claims; the project-wide scope claim mapping is used instead.",
			scope)
	}
}

// Checks if any of the provided values are in an Unknown state.
func HasUnknownValues(values ...any) bool {
	for _, v := range values {
		switch v := v.(type) {
		case interface{ IsUnknown() bool }:
			if v.IsUnknown() {
				return true
			}
		case interface{ Elements() map[string]attr.Value }:
			for _, elem := range v.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		case interface{ Elements() []attr.Value }:
			for _, elem := range v.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		case interface{ Attributes() map[string]attr.Value }:
			for _, elem := range v.Attributes() {
				if elem.IsUnknown() {
					return true
				}
			}
		}
	}
	return false
}
