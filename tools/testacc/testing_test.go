package testacc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequireEnv_returnsConfiguredFixtures_whenAllArePresent(t *testing.T) {
	// Given
	t.Setenv("DESCOPE_TEST_FIXTURE_PROJECT_ID", "P123")
	t.Setenv("DESCOPE_TEST_FIXTURE_RESOURCE_ID", "R123")

	// When
	fixtures := RequireEnv(t, "DESCOPE_TEST_FIXTURE_PROJECT_ID", "DESCOPE_TEST_FIXTURE_RESOURCE_ID")

	// Then
	require.Equal(t, "P123", fixtures["DESCOPE_TEST_FIXTURE_PROJECT_ID"])
	require.Equal(t, "R123", fixtures["DESCOPE_TEST_FIXTURE_RESOURCE_ID"])
}

func TestRequireEnv_stopsOptionalTest_whenFixtureIsMissing(t *testing.T) {
	// Given
	t.Setenv("DESCOPE_TEST_FIXTURE_PROJECT_ID", "")
	continued := false

	// When
	t.Run("missing fixture", func(t *testing.T) {
		RequireEnv(t, "DESCOPE_TEST_FIXTURE_PROJECT_ID")
		continued = true
	})

	// Then
	require.False(t, continued)
}
