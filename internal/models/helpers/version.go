package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const ModelVersion float64 = 1

func EnsureModelVersion(version float64, diagnostics *diag.Diagnostics) {
	if version > ModelVersion {
		diagnostics.AddWarning("Update the Descope terraform provider", "A new version of the Descope terraform provider is available. The Descope servers may still support the provider you have installed, but you should update soon to prevent service interruptions.")
	}
}
