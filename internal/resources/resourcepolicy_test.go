package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseResourcePolicyImportID_returns_composite_identity(t *testing.T) {
	t.Parallel()

	// Given
	importID := "P1/AP1/RS1"

	// When
	identity, err := parseResourcePolicyImportID(importID)

	// Then
	require.NoError(t, err)
	require.Equal(t, resourcePolicyImportIdentity{
		ProjectID:     "P1",
		ApplicationID: "AP1",
		ResourceID:    "RS1",
	}, identity)
	require.Equal(t, "AP1/RS1", resourcePolicyID(identity.ApplicationID, identity.ResourceID))
}

func TestParseResourcePolicyImportID_rejects_invalid_identity(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"P1/AP1",
		"P1/AP1/RS1/extra",
		"P1//RS1",
	}
	for _, importID := range tests {
		t.Run(importID, func(t *testing.T) {
			// When
			_, err := parseResourcePolicyImportID(importID)

			// Then
			require.Error(t, err)
		})
	}
}
