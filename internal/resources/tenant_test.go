package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTenantImportID_whenIDIsValid(t *testing.T) {
	// Given
	const importID = "project-id/tenant-id"

	// When
	projectID, tenantID, err := parseTenantImportID(importID)

	// Then
	require.NoError(t, err)
	require.Equal(t, "project-id", projectID)
	require.Equal(t, "tenant-id", tenantID)
}

func TestParseTenantImportID_whenIDIsInvalid(t *testing.T) {
	// Given
	const importID = "tenant-id"

	// When
	_, _, err := parseTenantImportID(importID)

	// Then
	require.EqualError(t, err, "tenant import ID must use the format <project_id>/<tenant_id>")
}
