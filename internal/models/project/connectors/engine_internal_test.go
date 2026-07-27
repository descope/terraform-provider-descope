package connectors

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
)

// TestConnectorEngine verifies that the engine_id attribute round-trips through the
// top-level executor fields the backend expects for connectors assigned to an engine.
func TestConnectorEngine(t *testing.T) {
	// assigning an engine writes the top-level executor fields
	data := map[string]any{}
	setConnectorEngine(data, stringattr.Value("CIEngine123"))
	assert.Equal(t, "engine", data["executorType"])
	assert.Equal(t, "CIEngine123", data["executorId"])

	// reading them back recovers the engine id
	var engineID stringattr.Type
	getConnectorEngine(data, &engineID)
	assert.Equal(t, "CIEngine123", engineID.ValueString())

	// an empty engine id leaves the connector local: no executor fields are emitted
	local := map[string]any{}
	setConnectorEngine(local, stringattr.Value(""))
	_, hasType := local["executorType"]
	_, hasID := local["executorId"]
	assert.False(t, hasType)
	assert.False(t, hasID)

	// reading a locally executed connector yields an empty engine id
	var localEngineID stringattr.Type
	getConnectorEngine(map[string]any{"executorType": "local"}, &localEngineID)
	assert.Equal(t, "", localEngineID.ValueString())

	// a connector object with no executor fields also yields an empty engine id
	var missingEngineID stringattr.Type
	getConnectorEngine(map[string]any{}, &missingEngineID)
	assert.Equal(t, "", missingEngineID.ValueString())
}

// TestHIBPConnectorEngine verifies that the generated HIBP model wires the engine_id
// attribute through the top-level executor fields, like the HTTP connector does.
func TestHIBPConnectorEngine(t *testing.T) {
	diags := diag.Diagnostics{}
	h := helpers.NewHandler(context.Background(), &diags)
	h.Refs.Add(helpers.ConnectorReferenceKey, "hibp", "CIHIBPExample", "My HIBP Connector")

	m := &HIBPModel{
		ID:       stringattr.Value("CIHIBPExample"),
		Name:     stringattr.Value("My HIBP Connector"),
		EngineID: stringattr.Value("CIEngineExample"),
	}
	data := m.Values(h)
	assert.False(t, diags.HasError())
	assert.Equal(t, "engine", data["executorType"])
	assert.Equal(t, "CIEngineExample", data["executorId"])

	var round HIBPModel
	round.SetValues(h, data)
	assert.False(t, diags.HasError())
	assert.Equal(t, "CIEngineExample", round.EngineID.ValueString())
}
