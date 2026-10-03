package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanCommandExitCodes(t *testing.T) {
	state := filepath.Join("migrate", "testdata", "states", "legacy-root.tfstate.json")
	args := func(out, source string) []string {
		return []string{
			"plan",
			"-state", state,
			"-address", "descope_project.main",
			"-source-provider-version", source,
			"-target-provider-version", "2.0.0",
			"-out", out,
		}
	}

	t.Run("ready", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		out := filepath.Join(t.TempDir(), "migration")
		assert.Equal(t, exitOK, run(args(out, "1.8.3"), stdout, stderr))
		assert.Contains(t, stdout.String(), "ready")
		assert.Empty(t, stderr.String())
		require.FileExists(t, filepath.Join(out, "manifest.json"))
		require.FileExists(t, filepath.Join(out, "adopt", "imports.tf"))
	})

	t.Run("invalid provider constraint", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		assert.Equal(t, exitUsage, run(args(filepath.Join(t.TempDir(), "migration"), ">= 1.8"), stdout, stderr))
		assert.Contains(t, stderr.String(), "must be an exact semantic version")
	})

	t.Run("blockers", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		out := filepath.Join(t.TempDir(), "migration")
		blocked := args(out, "1.8.3")
		blocked[2] = filepath.Join("migrate", "testdata", "states", "legacy-missing-id.tfstate.json")
		assert.Equal(t, exitBlocked, run(blocked, stdout, stderr))
		assert.Contains(t, stdout.String(), "blocked")
		assert.Contains(t, stderr.String(), "blocker")
		require.FileExists(t, filepath.Join(out, "BLOCKERS.md"))
		assert.NoFileExists(t, filepath.Join(out, "adopt", "imports.tf"))
	})
}
