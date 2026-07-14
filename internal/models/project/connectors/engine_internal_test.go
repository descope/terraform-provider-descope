package connectors

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
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
