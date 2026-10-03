package migrate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadStateFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fixture       string
		address       string
		module        string
		local         string
		index         any
		peers         []string
		standaloneNum int
	}{
		{
			name:          "root module",
			fixture:       "legacy-root.tfstate.json",
			address:       "descope_project.main",
			local:         "descope_project.main",
			standaloneNum: 2,
		},
		{
			name:    "nested count",
			fixture: "legacy-nested-count.tfstate.json",
			address: "module.identity.descope_project.main[0]",
			module:  "module.identity",
			local:   "descope_project.main[0]",
			index:   json.Number("0"),
		},
		{
			name:    "nested for_each",
			fixture: "legacy-nested-foreach.tfstate.json",
			address: `module.tenant.descope_project.main["prod"]`,
			module:  `module.tenant`,
			local:   `descope_project.main["prod"]`,
			index:   "prod",
		},
		{
			name:    "terraform show json",
			fixture: "legacy-nested-show.json",
			address: `module.tenant.descope_project.main["prod"]`,
			module:  `module.tenant`,
			local:   `descope_project.main["prod"]`,
			index:   "prod",
		},
		{
			name:    "count peers",
			fixture: "legacy-count-peers.tfstate.json",
			address: "descope_project.main[0]",
			local:   "descope_project.main[0]",
			index:   json.Number("0"),
			peers:   []string{"descope_project.main[1]"},
		},
		{
			name:    "module instance peers",
			fixture: "legacy-module-peers.tfstate.json",
			address: `module.tenant["acme"].descope_project.main`,
			module:  `module.tenant["acme"]`,
			local:   "descope_project.main",
			peers:   []string{`module.tenant["beta"].descope_project.main`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(filepath.Join("testdata", "states", tc.fixture))
			require.NoError(t, err)
			state, err := LoadStateFile(data)
			require.NoError(t, err)
			resource, err := state.Find(tc.address)
			require.NoError(t, err)
			assert.Equal(t, tc.address, resource.Address)
			assert.Equal(t, tc.module, resource.ModuleAddress)
			assert.Equal(t, tc.local, resource.LocalAddress)
			assert.Equal(t, tc.index, resource.IndexKey)
			assert.Equal(t, tc.peers, resource.Peers)
			assert.Len(t, resource.Existing, tc.standaloneNum)
		})
	}
}

func TestFindLegacyAddress(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join("testdata", "states", "legacy-nested-show.json"))
	require.NoError(t, err)
	state, err := LoadStateFile(data)
	require.NoError(t, err)

	resource, err := state.Find(`descope_project.main["prod"]`)
	require.NoError(t, err)
	assert.Equal(t, `module.tenant.descope_project.main["prod"]`, resource.Address)
}

func TestFindDistinguishesCountAndNumericForEachKeys(t *testing.T) {
	state, err := LoadStateFile([]byte(`{
  "version": 4,
  "terraform_version": "1.9.5",
  "resources": [{
    "mode": "managed",
    "type": "descope_project",
    "provider": "provider[\"registry.terraform.io/descope/descope\"]",
    "name": "main",
    "instances": [
      {"index_key": 0, "schema_version": 0, "attributes": {"id": "count"}},
      {"index_key": "0", "schema_version": 0, "attributes": {"id": "for-each"}}
    ]
  }]
}`))
	require.NoError(t, err)

	count, err := state.Find("descope_project.main[0]")
	require.NoError(t, err)
	forEach, err := state.Find(`descope_project.main["0"]`)
	require.NoError(t, err)
	assert.Equal(t, "count", count.Attributes["id"])
	assert.Equal(t, "for-each", forEach.Attributes["id"])
}

func TestResourceAddressCanonicalization(t *testing.T) {
	_, _, _, _, err := resourceAddresses("module.invalid\nresource.hack", ProjectResourceType, "main", nil)
	assert.Error(t, err)

	address, _, _, _, err := resourceAddresses("", ProjectResourceType, "main", "prod\nresource.hack")
	require.NoError(t, err)
	assert.NotContains(t, address, "\n")
	assert.Contains(t, address, `\n`)
}

func TestLoadStateRejectsWrongProvider(t *testing.T) {
	_, err := LoadStateFile([]byte(`{
  "version": 4,
  "terraform_version": "1.9.5",
  "resources": [{
    "mode": "managed",
    "type": "descope_project",
    "name": "main",
    "provider": "provider[\"registry.terraform.io/example/descope\"]",
    "instances": [{"schema_version": 0, "attributes": {"id": "P2project"}}]
  }]
}`))
	assert.ErrorContains(t, err, officialProviderSource)
}

func TestDescopeProviderSource(t *testing.T) {
	assert.True(t, isDescopeProvider(officialProviderSource))
	assert.True(t, isDescopeProvider(`provider["registry.terraform.io/descope/descope"]`))
	assert.True(t, isDescopeProvider(`module.identity.provider["registry.terraform.io/descope/descope"]`))
	assert.False(t, isDescopeProvider(`provider["registry.terraform.io/descope/descope"].eu`))
	assert.False(t, isDescopeProvider(`provider["registry.terraform.io/example/descope"]`))
}
