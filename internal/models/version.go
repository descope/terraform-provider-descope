package models

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const ModelVersion float64 = 1

func ensureModelVersion(version float64, diagnostics *diag.Diagnostics) {
	if version > ModelVersion {
		diagnostics.AddWarning("Update the Descope terraform provider", fmt.Sprintf("A new version of the Descope terraform provider is available that uses model version %d. The Descope servers still support the model version %d used by the provider you have installed, but you should update soon to prevent service interruptions.", int(version), int(ModelVersion)))
	}
}
