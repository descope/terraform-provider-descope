package migrate

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationLifecycle(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)

	opts := fixtureOptions()
	manifest, blocks, payloads, err := Build(ctx, registry, legacy, opts)
	require.NoError(t, err)
	require.True(t, manifest.Ready, "blockers: %#v", manifest.Blockers)
	assert.Len(t, manifest.Entries, 36)
	assert.Len(t, manifest.Standalone, 1)
	assert.Equal(t, "descope_permission.already", manifest.Standalone[0].Address)
	assert.Equal(t, "PM-manage", manifest.Standalone[0].BackendID)

	// Every nested object has one and only one destination/import tuple. The already standalone permission is
	// counted in the standalone inventory instead of generated imports.
	mapped := map[string]string{}
	for _, entry := range manifest.Entries {
		key := entry.Type + "\x00" + entry.ImportID
		if previous, duplicate := mapped[key]; duplicate {
			t.Fatalf("duplicate ownership of %s: %s and %s", key, previous, entry.LegacyPath)
		}
		assert.NotEmpty(t, entry.BackendID, entry.Address)
		mapped[key] = entry.LegacyPath
		destination, ok := registry.Lookup(entry.Type)
		require.True(t, ok)
		expected, err := destination.ImportID(entry.Scope, entry.BackendID)
		require.NoError(t, err)
		assert.Equal(t, expected, entry.ImportID)
	}
	assert.Len(t, mapped, 36)

	counts := map[string]int{}
	for _, entry := range manifest.Entries {
		counts[entry.Type]++
	}
	assert.Equal(t, 2, counts["descope_flow"])
	assert.Equal(t, 1, counts["descope_widget"])
	assert.Equal(t, 1, counts["descope_session_migration"])
	assert.Equal(t, 3, counts["descope_list"])
	assert.Equal(t, 2, counts["descope_email_template"])
	assert.Equal(t, 1, counts["descope_text_template"])
	assert.Equal(t, 1, counts["descope_voice_template"])
	assert.Equal(t, 1, counts["descope_app_role"])
	assert.Equal(t, 1, counts["descope_app_permission"])
	assert.Equal(t, 3, counts["descope_smtp_connector"]+counts["descope_sendgrid_connector"]+counts["descope_twilio_core_connector"])

	artifacts, err := Generate(manifest, blocks, payloads, snapshot, opts)
	require.NoError(t, err)
	for name, contents := range artifacts.Files {
		if filepath.Ext(name) != ".tf" {
			continue
		}
		_, diags := hclparse.NewParser().ParseHCL(contents, name)
		assert.False(t, diags.HasErrors(), "%s: %s", name, diags.Error())
	}
	combined := bytes.Join(artifactValuesExcept(artifacts.Files, backupStatePath), nil)
	for _, secret := range []string{"smtp-secret", "sg-secret", "twilio-token", "generated-secret", "migration-secret", "shh"} {
		assert.NotContains(t, string(combined), secret)
	}
	mainHCL := string(artifacts.Files[filepath.Join(adoptDir, "main.tf")])
	assert.Contains(t, mainHCL, `resource "descope_project" "main"`)
	assert.Contains(t, mainHCL, `local.descope_flow_flows`)
	assert.Contains(t, mainHCL, `var.descope_oauth_provider_oauth_provider_client_secret_`)
	assert.NotContains(t, mainHCL, "authorization_endpoint")
	assert.Contains(t, string(artifacts.Files[filepath.Join(adoptDir, "imports.tf")]), `for_each = {`)
	assert.Contains(t, string(artifacts.Files[filepath.Join(detachDir, "removed.tf")]), `destroy = false`)
	variables := map[string]bool{}
	for _, secret := range manifest.Secrets {
		assert.NotEmpty(t, secret.Variable, secret.Address+" "+secret.Attribute)
		assert.False(t, variables[secret.Variable], "duplicate variable "+secret.Variable)
		variables[secret.Variable] = true
	}
	changes := []map[string]any{{
		"address":       manifest.LegacyAddress,
		"mode":          "managed",
		"provider_name": officialProviderSource,
		"change":        map[string]any{"actions": []string{"no-op"}, "importing": map[string]any{"id": manifest.ProjectID}},
	}}
	for _, entry := range manifest.Entries {
		changes = append(changes, map[string]any{
			"address":       entry.Address,
			"mode":          "managed",
			"provider_name": officialProviderSource,
			"change":        map[string]any{"actions": []string{"no-op"}, "importing": map[string]any{"id": entry.ImportID}},
		})
	}
	plan, err := json.Marshal(map[string]any{"format_version": "1.2", "complete": true, "resource_changes": changes})
	require.NoError(t, err)
	verification, err := VerifyPlan(plan, manifest)
	require.NoError(t, err)
	assert.True(t, verification.OK, verification.Report())
}

func TestMigrationBlockers(t *testing.T) {
	ctx := context.Background()
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)

	tests := []struct {
		name    string
		fixture string
		kind    string
	}{
		{name: "missing server id", fixture: "legacy-missing-id.tfstate.json", kind: BlockerMissingID},
		{name: "single indexed module", fixture: "legacy-indexed-module-single.tfstate.json", kind: BlockerAmbiguousOwnership},
		{name: "duplicate backend mapping", fixture: "legacy-duplicate-mapping.tfstate.json", kind: BlockerDuplicateAddress},
		{name: "missing required secret", fixture: "legacy-secret-blocker.tfstate.json", kind: BlockerUnresolvedSecret},
		{name: "count peers", fixture: "legacy-count-peers.tfstate.json", kind: BlockerAmbiguousOwnership},
		{name: "module instance peers", fixture: "legacy-module-peers.tfstate.json", kind: BlockerAmbiguousOwnership},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			snapshot := readStateFixture(t, tc.fixture)
			state, err := LoadStateFile(snapshot)
			require.NoError(t, err)
			address := state.Addresses()[0]
			legacy, err := state.Find(address)
			require.NoError(t, err)
			manifest, _, _, err := Build(ctx, registry, legacy, fixtureOptions())
			require.NoError(t, err)
			assert.False(t, manifest.Ready)
			assert.Contains(t, blockerKinds(manifest), tc.kind)
		})
	}
}

func TestIndexedProjectGeneration(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		address  string
		mainWant string
	}{
		{
			name:     "count",
			fixture:  "legacy-nested-count.tfstate.json",
			address:  "module.identity.descope_project.main[0]",
			mainWant: "count       = 1",
		},
		{
			name:     "for_each",
			fixture:  "legacy-nested-foreach.tfstate.json",
			address:  `module.tenant.descope_project.main["prod"]`,
			mainWant: "for_each = {",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			snapshot := readStateFixture(t, tc.fixture)
			legacy := findLegacyFixture(t, snapshot, tc.address)
			registry, err := NewRegistry(ctx)
			require.NoError(t, err)
			manifest, blocks, payloads, err := Build(ctx, registry, legacy, fixtureOptions())
			require.NoError(t, err)
			require.True(t, manifest.Ready, "blockers: %#v", manifest.Blockers)
			artifacts, err := Generate(manifest, blocks, payloads, snapshot, fixtureOptions())
			require.NoError(t, err)
			main := string(artifacts.Files[filepath.Join(adoptDir, "main.tf")])
			imports := string(artifacts.Files[filepath.Join(adoptDir, "imports.tf")])
			removed := string(artifacts.Files[filepath.Join(detachDir, "removed.tf")])
			assert.Contains(t, main, tc.mainWant)
			assert.Contains(t, imports, "to = "+tc.address)
			assert.Contains(t, removed, "from = descope_project.main")
		})
	}
}

func TestShowSnapshotIsNotStateBackup(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-nested-show.json")
	legacy := findLegacyFixture(t, snapshot, `module.tenant.descope_project.main["prod"]`)
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, blocks, payloads, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	assert.False(t, manifest.Rollback.Restorable)
	assert.Equal(t, backupShowPath, manifest.Rollback.StateBackup)
	artifacts, err := Generate(manifest, blocks, payloads, snapshot, fixtureOptions())
	require.NoError(t, err)
	assert.Contains(t, artifacts.Files, backupShowPath)
	assert.NotContains(t, artifacts.Files, backupStatePath)
	assert.Contains(t, string(artifacts.Files[rollbackPath]), "cannot be restored")
}

func TestEmptyTypedListIsBlocked(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	lists, ok := legacy.Attributes["lists"].([]any)
	require.True(t, ok)
	first, ok := lists[0].(map[string]any)
	require.True(t, ok)
	first["data"] = "[]"
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, _, _, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	assert.False(t, manifest.Ready)
	assert.Contains(t, blockerKinds(manifest), BlockerInvalidValue)
}

func TestExistingStandaloneAddressIsReserved(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	legacy.Existing = append(legacy.Existing, ExistingResource{
		Address:        "descope_role.admin",
		Type:           "descope_role",
		Name:           "admin",
		ProviderSource: `provider["registry.terraform.io/descope/descope"]`,
		Attributes: map[string]any{
			"id":         "R-existing",
			"project_id": legacy.Attributes["id"],
		},
	})
	legacy.Existing = append(legacy.Existing, ExistingResource{
		Address:        "descope_flow.flows[\"existing\"]",
		Type:           "descope_flow",
		Name:           "flows",
		ProviderSource: `provider["registry.terraform.io/descope/descope"]`,
		Attributes: map[string]any{
			"id":         "existing",
			"project_id": legacy.Attributes["id"],
		},
	})
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, _, _, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	require.True(t, manifest.Ready, "blockers: %#v", manifest.Blockers)
	roleFound := false
	flowFound := false
	for _, entry := range manifest.Entries {
		if entry.Type == "descope_role" {
			assert.Equal(t, "descope_role.admin_2", entry.Address)
			roleFound = true
		}
		if entry.Type == "descope_flow" {
			assert.Contains(t, entry.Address, "descope_flow.flows_2[")
			flowFound = true
		}
	}
	assert.True(t, roleFound, "descope_role entry not found")
	assert.True(t, flowFound, "descope_flow entry not found")
}

func TestMatchingForkProviderStateIsBlocked(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	legacy.Existing = append(legacy.Existing, ExistingResource{
		Address:        "descope_role.fork",
		Type:           "descope_role",
		Name:           "fork",
		ProviderSource: `provider["registry.terraform.io/example/descope"]`,
		Attributes: map[string]any{
			"id":         "R-admin",
			"project_id": legacy.Attributes["id"],
		},
	})
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, _, _, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	assert.False(t, manifest.Ready)
	assert.Contains(t, blockerKinds(manifest), BlockerAmbiguousOwnership)
}

func TestGeneratedNamesAreCollisionSafe(t *testing.T) {
	assert.NotEqual(t,
		variableName("descope_abuseipdb_connector.prod_api", "api_key"),
		variableName("descope_amplitude_connector.prod_api", "api_key"),
	)
	assert.NotEqual(t, sanitizePathSegment("sign-in"), sanitizePathSegment("sign_in"))
}

func TestConnectorMappingCoverage(t *testing.T) {
	registry, err := NewRegistry(context.Background())
	require.NoError(t, err)
	families := registry.ConnectorFamilies()
	assert.Len(t, families, 73)
	assert.Equal(t, "descope_smtp_connector", families["smtp"])
	assert.Equal(t, "descope_twilio_core_connector", families["twilio_core"])
	assert.Equal(t, "descope_google_cloud_translation_connector", families["google_cloud_translation"])

	for name, renames := range connectorRenames {
		destination, ok := registry.Lookup(families[name])
		require.True(t, ok, name)
		for _, target := range renames {
			_, exists := destination.Schema.Attributes[target]
			assert.True(t, exists, "%s target %s", name, target)
		}
	}
	for name, flattened := range connectorFlattening {
		destination, ok := registry.Lookup(families[name])
		require.True(t, ok, name)
		for _, target := range flattened {
			_, exists := destination.Schema.Attributes[target]
			assert.True(t, exists, "%s target %s", name, target)
		}
	}
}

func TestGeneratedArtifactsGolden(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, blocks, payloads, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	artifacts, err := Generate(manifest, blocks, payloads, snapshot, fixtureOptions())
	require.NoError(t, err)

	goldens := map[string]string{
		manifestPath:                          "manifest.golden.json",
		filepath.Join(adoptDir, "main.tf"):    "main.golden.tf",
		filepath.Join(adoptDir, "imports.tf"): "imports.golden.tf",
	}
	for generated, golden := range goldens {
		expectedPath := filepath.Join("testdata", "golden", golden)
		if os.Getenv("UPDATE_GOLDEN") == "1" {
			require.NoError(t, os.MkdirAll(filepath.Dir(expectedPath), 0o755))
			require.NoError(t, os.WriteFile(expectedPath, artifacts.Files[generated], 0o600))
		}
		expected, err := os.ReadFile(expectedPath)
		require.NoError(t, err)
		assert.Equal(t, string(expected), string(artifacts.Files[generated]), generated)
	}

	// A second run over the same source is byte-identical.
	secondLegacy := findLegacyFixture(t, snapshot, "descope_project.main")
	secondManifest, secondBlocks, secondPayloads, err := Build(ctx, registry, secondLegacy, fixtureOptions())
	require.NoError(t, err)
	second, err := Generate(secondManifest, secondBlocks, secondPayloads, snapshot, fixtureOptions())
	require.NoError(t, err)
	assert.Equal(t, artifacts.Files, second.Files)
}

func fixtureOptions() Options {
	return Options{
		ProviderSource:        ProviderSource,
		SourceProviderVersion: "1.8.3",
		TargetProviderVersion: "2.0.0",
	}
}

func readStateFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "states", name))
	require.NoError(t, err)
	return data
}

func findLegacyFixture(t *testing.T, snapshot []byte, address string) *LegacyResource {
	t.Helper()
	state, err := LoadStateFile(snapshot)
	require.NoError(t, err)
	legacy, err := state.Find(address)
	require.NoError(t, err)
	return legacy
}

func blockerKinds(manifest *Manifest) []string {
	kinds := make([]string, 0, len(manifest.Blockers))
	for _, blocker := range manifest.Blockers {
		kinds = append(kinds, blocker.Kind)
	}
	sort.Strings(kinds)
	return kinds
}

func artifactValuesExcept(files map[string][]byte, excluded string) [][]byte {
	values := [][]byte{}
	for name, value := range files {
		if name == excluded {
			continue
		}
		values = append(values, value)
	}
	return values
}

func TestSensitiveStateWarning(t *testing.T) {
	ctx := context.Background()
	snapshot := readStateFixture(t, "legacy-root.tfstate.json")
	legacy := findLegacyFixture(t, snapshot, "descope_project.main")
	registry, err := NewRegistry(ctx)
	require.NoError(t, err)
	manifest, blocks, payloads, err := Build(ctx, registry, legacy, fixtureOptions())
	require.NoError(t, err)
	artifacts, err := Generate(manifest, blocks, payloads, snapshot, fixtureOptions())
	require.NoError(t, err)

	assert.Contains(t, strings.ToLower(string(artifacts.Files[readmePath])), "sensitive")
	assert.Contains(t, strings.ToLower(string(artifacts.Files[secretsPath])), "cleartext")
	assert.Contains(t, strings.ToLower(string(artifacts.Files[rollbackPath])), "credential")
}

func TestArtifactsRejectUnsafeOutput(t *testing.T) {
	artifacts := &Artifacts{Files: map[string][]byte{"manifest.json": []byte("{}\n")}}
	nonempty := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(nonempty, "reviewed.txt"), []byte("keep"), 0o600))
	assert.ErrorContains(t, artifacts.Write(nonempty), "already exists")

	target := t.TempDir()
	link := filepath.Join(t.TempDir(), "output")
	require.NoError(t, os.Symlink(target, link))
	assert.ErrorContains(t, artifacts.Write(link), "already exists")
	assert.NoFileExists(t, filepath.Join(target, "manifest.json"))

	realParent := t.TempDir()
	linkedParent := filepath.Join(t.TempDir(), "linked-parent")
	require.NoError(t, os.Symlink(realParent, linkedParent))
	assert.ErrorContains(t, artifacts.Write(filepath.Join(linkedParent, "output")), "contains a symbolic link")

	escapeParent := t.TempDir()
	escapeRoot := filepath.Join(escapeParent, "output")
	escaping := &Artifacts{Files: map[string][]byte{"../outside": []byte("secret")}}
	assert.ErrorContains(t, escaping.Write(escapeRoot), "escapes")
	assert.NoFileExists(t, filepath.Join(escapeParent, "outside"))

	fresh := filepath.Join(t.TempDir(), "fresh")
	require.NoError(t, artifacts.Write(fresh))
	rootInfo, err := os.Stat(fresh)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), rootInfo.Mode().Perm())
	fileInfo, err := os.Stat(filepath.Join(fresh, "manifest.json"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), fileInfo.Mode().Perm())
}

func TestChildModuleSecretGuidance(t *testing.T) {
	manifest := &Manifest{
		ModuleAddress: "module.identity",
		Rollback:      Rollback{StateBackup: backupStatePath},
		Secrets: []SecretRequirement{{
			Address:   "module.identity.descope_smtp_connector.mail",
			Attribute: "password",
			Variable:  "mail_password",
			Present:   true,
		}},
	}
	report := secretReport(manifest)
	assert.Contains(t, report, "`TF_VAR_*` only sets root-module variables")
	assert.Contains(t, report, "pass each one through every existing module call")
}
