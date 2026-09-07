package testacc

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/stretchr/testify/require"
)

// Returns the JSON fixture at path as a compacted quoted HCL string, applying old/new replacement pairs first and escaping
// the template interpolation sequences that are meaningful inside quoted HCL strings.
func FixtureJSON(t *testing.T, path string, replacements ...string) string {
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	s := string(b)
	for i := 0; i+1 < len(replacements); i += 2 {
		s = strings.ReplaceAll(s, replacements[i], replacements[i+1])
	}

	var v any
	require.NoError(t, json.Unmarshal([]byte(s), &v))
	compact, err := json.Marshal(v)
	require.NoError(t, err)
	q := strconv.Quote(string(compact))
	q = strings.ReplaceAll(q, "${", "$${")
	q = strings.ReplaceAll(q, "%{", "%%{")
	return q
}

// Removes or modifies an entity through the management API, bypassing Terraform, to simulate out-of-band changes.
func OutOfBandPost(t *testing.T, projectID, path string, body map[string]any) {
	base := os.Getenv("DESCOPE_BASE_URL")
	require.NotEmpty(t, base, "The DESCOPE_BASE_URL environment variable must be set for out-of-band requests")
	client := infra.NewClient("testacc", os.Getenv("DESCOPE_MANAGEMENT_KEY"), base)
	err := client.Post(context.Background(), projectID, path, body)
	require.NoError(t, err, "Out-of-band POST request to %s failed", path)
}

// Used to skip specific tests in environments where the required license is expected to not be available usually.
func RequireLicense(t *testing.T, license string) {
	if os.Getenv("DESCOPE_TESTACC_LICENSED") == "" {
		t.Skipf("Skipping test that requires the %s license (set DESCOPE_TESTACC_LICENSED to run)", license)
	}
}
