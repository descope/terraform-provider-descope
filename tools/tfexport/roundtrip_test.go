package main

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestRoundTrip(t *testing.T) {
	harness.RequireEnv(t)

	scenarios := roundTripScenarios(t)
	if len(scenarios) == 0 {
		t.Fatal("no scenarios loaded")
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			harness.RunScenario(t, scenario)
		})
	}
}
