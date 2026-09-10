package migrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// ManifestFormatVersion is bumped when the manifest layout changes in a way consumers must notice.
const ManifestFormatVersion = "1"

// RequiredTerraformVersion is the minimum safe end-to-end migration release. Terraform 1.7 introduced `removed`
// and `for_each` import blocks; 1.8 adds the plan completeness signal verify-plan requires to fail closed.
const RequiredTerraformVersion = ">= 1.8.0"

// ToolName identifies the producer in generated artifacts.
const ToolName = "tfmigrate"

// Blocker kinds. A manifest with any blocker is never ready, and no migration artifacts are written for it.
const (
	BlockerMissingID          = "missing_id"
	BlockerUnsupportedField   = "unsupported_legacy_field"
	BlockerDuplicateAddress   = "duplicate_destination"
	BlockerAmbiguousOwnership = "ambiguous_ownership"
	BlockerUnresolvedRef      = "unresolved_reference"
	BlockerMissingPayload     = "missing_payload"
	BlockerUnresolvedSecret   = "unresolved_secret"
	BlockerUnknownAttribute   = "unknown_legacy_attribute"
	BlockerInvalidValue       = "invalid_legacy_value"
)

// Scope carries the addressing values a destination resource needs in its import id and configuration.
type Scope struct {
	Kind      ScopeKind `json:"kind"`
	ProjectID string    `json:"project_id,omitempty"`
	AppID     string    `json:"app_id,omitempty"`
	Method    string    `json:"method,omitempty"`
}

// ProviderPin records an exact provider version, so the two migration stages run against known implementations.
type ProviderPin struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

// Entry maps one formerly nested object to exactly one destination resource instance.
type Entry struct {
	LegacyPath string   `json:"legacy_path"`
	Address    string   `json:"address"`
	Type       string   `json:"type"`
	Label      string   `json:"label"`
	Scope      Scope    `json:"scope"`
	BackendID  string   `json:"backend_id"`
	ImportID   string   `json:"import_id"`
	IDSource   string   `json:"id_source"`
	Payloads   []string `json:"payloads,omitempty"`
	// SecretAttributes are required sensitive attributes that the backend never returns, so an adoption plan writes
	// them back in place. Verify accepts an update touching only these attributes and rejects anything else.
	SecretAttributes []string `json:"secret_attributes,omitempty"`
	// References records resolved cross resource ids, so a reviewer can replace literals with resource references.
	References map[string]string `json:"references,omitempty"`
}

// Blocker is a reason the migration is not safe to run. Blockers are fatal, never advisory.
type Blocker struct {
	Kind       string `json:"kind"`
	LegacyPath string `json:"legacy_path,omitempty"`
	Address    string `json:"address,omitempty"`
	Detail     string `json:"detail"`
}

type SecretRequirement struct {
	Address    string `json:"address"`
	Attribute  string `json:"attribute"`
	Variable   string `json:"variable"`
	Type       string `json:"type"`
	LegacyPath string `json:"legacy_path"`
	Present    bool   `json:"present_in_legacy_state"`
}

// Retained is a legacy attribute that stays on descope_project after the split.
type Retained struct {
	Attribute string `json:"attribute"`
	Detail    string `json:"detail"`
}

// Unowned is an object present in the legacy state that the migration deliberately does not adopt, because the
// backend or the console owns it rather than Terraform.
type Unowned struct {
	LegacyPath string `json:"legacy_path"`
	Detail     string `json:"detail"`
}

// Standalone is an already managed destination that the migration leaves alone rather than importing twice.
type Standalone struct {
	Address   string `json:"address"`
	Type      string `json:"type"`
	BackendID string `json:"backend_id"`
	Detail    string `json:"detail"`
}

// ProjectCore is the part of descope_project that remains on that resource after the split.
type ProjectCore struct {
	Name               string   `json:"name"`
	Environment        string   `json:"environment"`
	Tags               []string `json:"tags"`
	DeletionProtection *bool    `json:"deletion_protection,omitempty"`
}

// Rollback records what is needed to restore the pre-migration state. The utility never performs a rollback.
type Rollback struct {
	StateBackup     string   `json:"state_backup"`
	Restorable      bool     `json:"restorable"`
	LegacyAddress   string   `json:"legacy_address"`
	SourceProvider  string   `json:"source_provider_version"`
	Instructions    []string `json:"instructions"`
	BackendMutation string   `json:"backend_mutation"`
}

// Manifest is the machine readable ownership map for one legacy descope_project instance.
type Manifest struct {
	Tool              string              `json:"tool"`
	FormatVersion     string              `json:"format_version"`
	Ready             bool                `json:"ready"`
	ProjectID         string              `json:"project_id"`
	LegacyAddress     string              `json:"legacy_address"`
	LegacyLocalAddr   string              `json:"legacy_local_address"`
	InputFormat       string              `json:"input_format"`
	ModuleAddress     string              `json:"module_address"`
	SourceProvider    ProviderPin         `json:"source_provider"`
	ProjectCore       ProjectCore         `json:"project_core"`
	InstanceKey       any                 `json:"instance_key,omitempty"`
	TargetProvider    ProviderPin         `json:"target_provider"`
	RequiredTerraform string              `json:"required_terraform"`
	RemoteBackfill    bool                `json:"remote_backfill"`
	OwnershipPolicy   string              `json:"ownership_policy"`
	Entries           []Entry             `json:"entries"`
	Blockers          []Blocker           `json:"blockers"`
	Secrets           []SecretRequirement `json:"secrets"`
	Retained          []Retained          `json:"retained"`
	Unowned           []Unowned           `json:"unowned"`
	Standalone        []Standalone        `json:"already_standalone"`
	Rollback          Rollback            `json:"rollback"`
}

// Summary is the one line result a migration run prints.
func (m *Manifest) Summary() string {
	state := "ready"
	if !m.Ready {
		state = "blocked"
	}
	return fmt.Sprintf("%s: %s, %d destination(s), %d already standalone, %d blocker(s), %d secret requirement(s), %d object(s) not adopted",
		m.LegacyAddress, state, len(m.Entries), len(m.Standalone), len(m.Blockers), len(m.Secrets), len(m.Unowned))
}

// sortAll orders every collection so reruns over the same input produce byte identical output.
func (m *Manifest) sortAll() {
	sort.SliceStable(m.Entries, func(i, j int) bool {
		if m.Entries[i].Address != m.Entries[j].Address {
			return m.Entries[i].Address < m.Entries[j].Address
		}
		return m.Entries[i].LegacyPath < m.Entries[j].LegacyPath
	})
	sort.SliceStable(m.Blockers, func(i, j int) bool {
		a, b := m.Blockers[i], m.Blockers[j]
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.LegacyPath != b.LegacyPath {
			return a.LegacyPath < b.LegacyPath
		}
		if a.Address != b.Address {
			return a.Address < b.Address
		}
		return a.Detail < b.Detail
	})
	sort.SliceStable(m.Secrets, func(i, j int) bool {
		if m.Secrets[i].Address != m.Secrets[j].Address {
			return m.Secrets[i].Address < m.Secrets[j].Address
		}
		return m.Secrets[i].Attribute < m.Secrets[j].Attribute
	})
	sort.SliceStable(m.Retained, func(i, j int) bool { return m.Retained[i].Attribute < m.Retained[j].Attribute })
	sort.SliceStable(m.Standalone, func(i, j int) bool { return m.Standalone[i].Address < m.Standalone[j].Address })
	sort.SliceStable(m.Unowned, func(i, j int) bool { return m.Unowned[i].LegacyPath < m.Unowned[j].LegacyPath })
}

// Encode renders the manifest as stable, indented JSON with a trailing newline.
func (m *Manifest) Encode() ([]byte, error) {
	m.sortAll()
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding manifest: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeManifest reads a manifest written by Encode.
func DecodeManifest(data []byte) (*Manifest, error) {
	m := &Manifest{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(m); err != nil {
		return nil, fmt.Errorf("decoding manifest: %w", err)
	}
	if m.Tool != ToolName {
		return nil, fmt.Errorf("manifest was not produced by %s", ToolName)
	}
	if m.FormatVersion != ManifestFormatVersion {
		return nil, fmt.Errorf("manifest format version %q is not supported, expected %q", m.FormatVersion, ManifestFormatVersion)
	}
	if m.SourceProvider.Source != ProviderSource || m.TargetProvider.Source != ProviderSource {
		return nil, fmt.Errorf("manifest provider source must be %q", ProviderSource)
	}
	return m, nil
}
