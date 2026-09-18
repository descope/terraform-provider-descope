package migrate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixtureProjectReader struct {
	project map[string]any
	reads   int
}

func (r *fixtureProjectReader) ReadProject(_ context.Context) (map[string]any, error) {
	r.reads++
	return r.project, nil
}

func TestGETOnlyBackfill(t *testing.T) {
	snapshot := readStateFixture(t, "legacy-missing-id.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	reader := &fixtureProjectReader{project: map[string]any{
		"authorization": map[string]any{
			"roles": []any{map[string]any{"id": "R-admin", "name": "Admin"}},
		},
	}}

	filled, err := Backfill(context.Background(), legacy, reader)
	require.NoError(t, err)
	assert.Equal(t, 1, reader.reads)
	assert.Equal(t, []string{`authorization.roles "Admin"`}, filled)
	authorization, ok := legacy.Attributes["authorization"].(map[string]any)
	require.True(t, ok)
	roles, ok := authorization["roles"].([]any)
	require.True(t, ok)
	role, ok := roles[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "R-admin", role["id"])
}

func TestProjectReaderRejectsUnsafeLogs(t *testing.T) {
	t.Setenv(managementKeyEnv, "test-key")
	t.Setenv(unsafeLogsEnv, "1")
	reader, err := NewProjectReader("P2project", "")
	assert.Nil(t, reader)
	assert.ErrorContains(t, err, unsafeLogsEnv)
}

func TestGETOnlyBackfillRejectsAmbiguousNames(t *testing.T) {
	snapshot := readStateFixture(t, "legacy-missing-id.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	reader := &fixtureProjectReader{project: map[string]any{
		"authorization": map[string]any{
			"roles": []any{
				map[string]any{"id": "R-one", "name": "Admin"},
				map[string]any{"id": "R-two", "name": "Admin"},
			},
		},
	}}
	filled, err := Backfill(context.Background(), legacy, reader)
	assert.Nil(t, filled)
	assert.ErrorContains(t, err, "more than one backend id")
}

func TestGETOnlyBackfillConnectorID(t *testing.T) {
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	connectors, ok := legacy.Attributes["connectors"].(map[string]any)
	require.True(t, ok)
	smtp, ok := connectors["smtp"].([]any)
	require.True(t, ok)
	connector, ok := smtp[0].(map[string]any)
	require.True(t, ok)
	connector["id"] = ""

	reader := &fixtureProjectReader{project: map[string]any{
		"connectors": map[string]any{
			"smtp": []any{map[string]any{"id": "C-smtp", "name": "acme-smtp"}},
		},
	}}
	filled, err := Backfill(context.Background(), legacy, reader)
	require.NoError(t, err)
	assert.Contains(t, filled, `connectors.smtp "acme-smtp"`)
	assert.Equal(t, "C-smtp", connector["id"])
}

func TestValidateRemoteBaseURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{url: ""},
		{url: "https://api.descope.com"},
		{url: "https://api.example.test/path"},
		{url: "http://api.example.test", wantErr: true},
		{url: "https://127.0.0.1", wantErr: true},
		{url: "https://localhost", wantErr: true},
		{url: "https://api.example.test:8443", wantErr: true},
		{url: "https://user:secret@api.example.test", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			err := validateRemoteBaseURL(tc.url)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
