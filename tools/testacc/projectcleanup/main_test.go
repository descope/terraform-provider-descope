package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_selectsOnlyProjectsWithCurrentRunPrefix(t *testing.T) {
	// Given
	input := strings.NewReader(`{
  "ok": true,
  "projects": [
    {"id": "P-current-1", "name": "testacc-gh-123-TestProject-a"},
    {"id": "P-other", "name": "testacc-gh-456-TestProject-b"},
    {"id": "P-current-2", "name": "testacc-gh-123-TestSettings-c"},
    {"id": "P-local", "name": "testacc-local-TestProject-d"}
  ]
}`)
	var output bytes.Buffer

	// When
	err := run([]string{"testacc-gh-123-"}, input, &output)

	// Then
	require.NoError(t, err)
	require.Equal(t, "P-current-1\nP-current-2\n", output.String())
}

func TestRun_rejectsMissingPrefix(t *testing.T) {
	// Given
	input := strings.NewReader(`{"ok":true,"projects":[]}`)
	var output bytes.Buffer

	// When
	err := run(nil, input, &output)

	// Then
	require.EqualError(t, err, "usage: projectcleanup <project-name-prefix>")
}
