package testacc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAlias_usesLocalPrefix_whenPrefixIsUnset(t *testing.T) {
	// Given
	t.Setenv("DESCOPE_TESTACC_PREFIX", "")

	// When
	alias := GenerateAlias(t)

	// Then
	require.True(t, strings.HasPrefix(alias, "testacc-local-GenerateAlias_usesLocalPrefix_whenPrefixIsUnset-"))
}

func TestGenerateAlias_usesConfiguredRunPrefix_whenPrefixIsSet(t *testing.T) {
	// Given
	t.Setenv("DESCOPE_TESTACC_PREFIX", "testacc-gh-29485327177")

	// When
	alias := GenerateAlias(t)

	// Then
	require.True(t, strings.HasPrefix(alias, "testacc-gh-29485327177-GenerateAlias_usesConfiguredRunPrefix_whenPrefixIsSet-"))
}
