package migrate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// LegacyResource is the located legacy descope_project instance and the addresses derived from it.
type LegacyResource struct {
	// Address is the full instance address, module qualified and including any instance key.
	Address string
	// ModuleAddress is the module the resource is declared in, empty for the root module.
	ModuleAddress string
	// ModuleConfigAddress omits module instance keys, grouping repeated instances of the same module declaration.
	ModuleConfigAddress string
	// ModuleIndexed is true when the declaring module itself uses count or for_each.
	ModuleIndexed bool
	// LocalAddress is the instance address inside its own module, which is what `removed` blocks and generated
	// resource configuration in that module must use.
	LocalAddress string
	// Type and Name are the resource type and configuration label.
	Type string
	Name string
	// IndexKey is the instance key of a count or for_each resource, nil for a single instance.
	IndexKey any
	// Attributes is the instance's flat attribute tree as stored by Terraform.
	Attributes map[string]any
	// SchemaVersion is the schema version the state was written with, when the source reports it.
	SchemaVersion int64
	// SourceFormat records which input shape the resource was read from, for the migration report.
	SourceFormat   string
	ProviderSource string
	// Existing contains managed resources that already have their own state and must not be adopted again.
	Existing []ExistingResource
	// Peers are other instances of the same descope_project configuration. A removed block cannot name just one
	// count or for_each instance, so their presence blocks a single-instance migration.
	Peers []string
}

// ExistingResource is an already standalone managed resource from the same state snapshot.
type ExistingResource struct {
	Address        string
	ModuleAddress  string
	Type           string
	Name           string
	Attributes     map[string]any
	ProviderSource string
}

// StateFile is a legacy Terraform state or plan-format snapshot holding descope_project instances.
type StateFile struct {
	TerraformVersion string
	Format           string
	instances        []LegacyResource
	existing         []ExistingResource
}

const (
	formatRawState = "terraform.tfstate"
	formatShowJSON = "terraform show -json"
)

const officialProviderSource = "registry.terraform.io/descope/descope"

func isDescopeProvider(source string) bool {
	if source == officialProviderSource || source == ProviderSource {
		return true
	}
	provider := `provider["` + officialProviderSource + `"]`
	if source == provider {
		return true
	}
	prefix, found := strings.CutSuffix(source, "."+provider)
	if !found || prefix == "" {
		return false
	}
	canonical, _, err := canonicalTraversal(prefix)
	return err == nil && canonical == prefix && strings.HasPrefix(prefix, "module.")
}

// rawState is the subset of the version 4 state file layout the migration reads.
type rawState struct {
	Version          int              `json:"version"`
	TerraformVersion string           `json:"terraform_version"`
	Resources        []rawStateModule `json:"resources"`
}

type rawStateModule struct {
	Module    string             `json:"module"`
	Mode      string             `json:"mode"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Provider  string             `json:"provider"`
	Instances []rawStateInstance `json:"instances"`
}
type rawStateInstance struct {
	IndexKey      any            `json:"index_key"`
	SchemaVersion int64          `json:"schema_version"`
	Attributes    map[string]any `json:"attributes"`
	Deposed       string         `json:"deposed"`
}

// showState is the subset of the `terraform show -json` state layout the migration reads.
type showState struct {
	FormatVersion    string     `json:"format_version"`
	TerraformVersion string     `json:"terraform_version"`
	Values           *showValue `json:"values"`
}

type showValue struct {
	RootModule *showModule `json:"root_module"`
}

type showModule struct {
	Address      string          `json:"address"`
	Resources    []showResource  `json:"resources"`
	ChildModules []*showModule   `json:"child_modules"`
	Outputs      json.RawMessage `json:"outputs,omitempty"`
}

type showResource struct {
	Address       string         `json:"address"`
	Mode          string         `json:"mode"`
	Type          string         `json:"type"`
	Name          string         `json:"name"`
	Index         any            `json:"index"`
	SchemaVersion int64          `json:"schema_version"`
	ProviderName  string         `json:"provider_name"`
	Values        map[string]any `json:"values"`
}

// LoadStateFile reads either a raw Terraform state file or the JSON produced by `terraform show -json`, and collects
// every managed descope_project instance in the root module and in all child modules.

func decodeStateJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(target)
}
func LoadStateFile(data []byte) (*StateFile, error) {
	probe := struct {
		Version       json.RawMessage `json:"version"`
		FormatVersion string          `json:"format_version"`
	}{}
	if err := decodeStateJSON(data, &probe); err != nil {
		return nil, fmt.Errorf("input is not JSON: %w", err)
	}
	switch {
	case probe.FormatVersion != "":
		return loadShowState(data)
	case len(probe.Version) > 0:
		return loadRawState(data)
	}
	return nil, fmt.Errorf("input has neither a state `version` nor a `format_version` field, so it is not a Terraform state or `terraform show -json` snapshot")
}

func resourceAddresses(moduleAddress, resourceType, name string, index any) (full, module, moduleConfig, local string, err error) {
	local, _, err = canonicalTraversal(resourceType + "." + name + indexSuffix(index))
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid resource address: %w", err)
	}
	if err := validateResourceTraversal(resourceType, name, local); err != nil {
		return "", "", "", "", err
	}
	if moduleAddress == "" {
		return local, "", "", local, nil
	}
	module, traversal, err := canonicalTraversal(moduleAddress)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid module address: %w", err)
	}
	if err := validateModuleTraversal(module, traversal); err != nil {
		return "", "", "", "", err
	}
	withoutIndexes := make(hcl.Traversal, 0, len(traversal))
	for _, step := range traversal {
		if _, isIndex := step.(hcl.TraverseIndex); !isIndex {
			withoutIndexes = append(withoutIndexes, step)
		}
	}
	moduleConfig = string(hclwrite.TokensForTraversal(withoutIndexes).Bytes())
	full, _, err = canonicalTraversal(module + "." + local)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid module resource address: %w", err)
	}
	return full, module, moduleConfig, local, nil
}

func canonicalTraversal(address string) (string, hcl.Traversal, error) {
	traversal, diagnostics := hclsyntax.ParseTraversalAbs([]byte(address), "state address", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return "", nil, fmt.Errorf("%q: %s", address, diagnostics.Error())
	}
	return string(hclwrite.TokensForTraversal(traversal).Bytes()), traversal, nil
}

func validateResourceTraversal(resourceType, name, address string) error {
	_, traversal, err := canonicalTraversal(address)
	if err != nil {
		return err
	}
	if len(traversal) < 2 || len(traversal) > 3 {
		return fmt.Errorf("resource address %q has an invalid traversal", address)
	}
	root, rootOK := traversal[0].(hcl.TraverseRoot)
	attribute, attributeOK := traversal[1].(hcl.TraverseAttr)
	if !rootOK || !attributeOK || root.Name != resourceType || attribute.Name != name {
		return fmt.Errorf("resource address %q does not match %s.%s", address, resourceType, name)
	}
	if len(traversal) == 3 {
		if _, indexOK := traversal[2].(hcl.TraverseIndex); !indexOK {
			return fmt.Errorf("resource address %q has an invalid instance key", address)
		}
	}
	return nil
}

func validateModuleTraversal(address string, traversal hcl.Traversal) error {
	if len(traversal) < 2 {
		return fmt.Errorf("module address %q is incomplete", address)
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok || root.Name != "module" {
		return fmt.Errorf("module address %q does not begin with module", address)
	}
	position := 1
	for {
		if position >= len(traversal) {
			return fmt.Errorf("module address %q is missing a module name", address)
		}
		if _, ok := traversal[position].(hcl.TraverseAttr); !ok {
			return fmt.Errorf("module address %q has an invalid module name", address)
		}
		position++
		if position < len(traversal) {
			if _, ok := traversal[position].(hcl.TraverseIndex); ok {
				position++
			}
		}
		if position == len(traversal) {
			return nil
		}
		next, ok := traversal[position].(hcl.TraverseAttr)
		if !ok || next.Name != "module" {
			return fmt.Errorf("module address %q has an invalid segment", address)
		}
		position++
	}
}

func traversalHasIndex(address string) bool {
	if address == "" {
		return false
	}
	_, traversal, err := canonicalTraversal(address)
	if err != nil {
		return true
	}
	for _, step := range traversal {
		if _, ok := step.(hcl.TraverseIndex); ok {
			return true
		}
	}
	return false
}

func loadRawState(data []byte) (*StateFile, error) {
	state := &rawState{}
	if err := decodeStateJSON(data, state); err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}
	if state.Version != 4 {
		return nil, fmt.Errorf("state file version %d is not supported, expected version 4", state.Version)
	}
	file := &StateFile{TerraformVersion: state.TerraformVersion, Format: formatRawState}
	for _, res := range state.Resources {
		if res.Mode != "managed" {
			continue
		}
		for _, instance := range res.Instances {
			if instance.Deposed != "" {
				continue
			}
			address, module, moduleConfig, local, err := resourceAddresses(res.Module, res.Type, res.Name, instance.IndexKey)
			if err != nil {
				return nil, err
			}
			if res.Type != ProjectResourceType {
				file.existing = append(file.existing, ExistingResource{
					Address: address, ModuleAddress: module, Type: res.Type, Name: res.Name,
					Attributes: instance.Attributes, ProviderSource: res.Provider,
				})
				continue
			}
			if !isDescopeProvider(res.Provider) {
				return nil, fmt.Errorf("legacy resource %s is associated with provider %q, expected %s", address, res.Provider, officialProviderSource)
			}
			file.instances = append(file.instances, LegacyResource{
				Address:             address,
				ModuleAddress:       module,
				ModuleConfigAddress: moduleConfig,
				ModuleIndexed:       traversalHasIndex(module),
				LocalAddress:        local,
				Type:                res.Type,
				Name:                res.Name,
				IndexKey:            instance.IndexKey,
				Attributes:          instance.Attributes,
				SchemaVersion:       instance.SchemaVersion,
				SourceFormat:        formatRawState,
				ProviderSource:      res.Provider,
			})
		}
	}
	file.finalize()
	return file, nil
}

func loadShowState(data []byte) (*StateFile, error) {
	state := &showState{}
	if err := decodeStateJSON(data, state); err != nil {
		return nil, fmt.Errorf("reading `terraform show -json` snapshot: %w", err)
	}
	major, _, _ := strings.Cut(state.FormatVersion, ".")
	if major != "1" {
		return nil, fmt.Errorf("terraform show JSON format version %q is not supported, expected major version 1", state.FormatVersion)
	}
	file := &StateFile{TerraformVersion: state.TerraformVersion, Format: formatShowJSON}
	if state.Values == nil || state.Values.RootModule == nil {
		return file, nil
	}
	if err := collectShowModule(state.Values.RootModule, file); err != nil {
		return nil, err
	}
	file.finalize()
	return file, nil
}

func collectShowModule(module *showModule, file *StateFile) error {
	for _, res := range module.Resources {
		if res.Mode != "managed" {
			continue
		}
		address, moduleAddress, moduleConfig, local, err := resourceAddresses(module.Address, res.Type, res.Name, res.Index)
		if err != nil {
			return err
		}
		if res.Address != "" {
			provided, _, parseErr := canonicalTraversal(res.Address)
			if parseErr != nil {
				return fmt.Errorf("invalid resource address in Terraform show output: %w", parseErr)
			}
			if provided != address {
				return fmt.Errorf("resource address %q does not match its module, type, name and index (%q)", provided, address)
			}
		}
		if res.Type != ProjectResourceType {
			file.existing = append(file.existing, ExistingResource{
				Address: address, ModuleAddress: moduleAddress, Type: res.Type, Name: res.Name,
				Attributes: res.Values, ProviderSource: res.ProviderName,
			})
			continue
		}
		if !isDescopeProvider(res.ProviderName) {
			return fmt.Errorf("legacy resource %s is associated with provider %q, expected %s", address, res.ProviderName, officialProviderSource)
		}
		file.instances = append(file.instances, LegacyResource{
			Address:             address,
			ModuleAddress:       moduleAddress,
			ModuleConfigAddress: moduleConfig,
			ModuleIndexed:       traversalHasIndex(moduleAddress),
			LocalAddress:        local,
			Type:                res.Type,
			Name:                res.Name,
			IndexKey:            res.Index,
			Attributes:          res.Values,
			SchemaVersion:       res.SchemaVersion,
			SourceFormat:        formatShowJSON,
			ProviderSource:      res.ProviderName,
		})
	}
	for _, child := range module.ChildModules {
		if child != nil {
			if err := collectShowModule(child, file); err != nil {
				return err
			}
		}
	}
	return nil
}

// finalize attaches the already standalone resources and sibling project instances to each migration candidate.
func (f *StateFile) finalize() {
	for i := range f.instances {
		instance := &f.instances[i]
		instance.Existing = append([]ExistingResource(nil), f.existing...)
		for j := range f.instances {
			peer := &f.instances[j]
			if i == j || peer.ModuleConfigAddress != instance.ModuleConfigAddress || peer.Name != instance.Name {
				continue
			}
			instance.Peers = append(instance.Peers, peer.Address)
		}
		sort.Strings(instance.Peers)
	}
}

// Addresses lists every legacy instance address found, sorted, for error messages and discovery.
func (f *StateFile) Addresses() []string {
	addresses := make([]string, 0, len(f.instances))
	for _, instance := range f.instances {
		addresses = append(addresses, instance.Address)
	}
	sort.Strings(addresses)
	return addresses
}

// Find locates one legacy instance by address. The address may be given with or without a module prefix as long as it
// is unambiguous, and an exact match always wins over a suffix match.
func (f *StateFile) Find(address string) (*LegacyResource, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("no legacy resource address was given")
	}
	for i := range f.instances {
		if f.instances[i].Address == address {
			return &f.instances[i], nil
		}
	}
	var matches []*LegacyResource
	for i := range f.instances {
		if f.instances[i].LocalAddress == address {
			matches = append(matches, &f.instances[i])
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		if len(f.instances) == 0 {
			return nil, fmt.Errorf("no managed %s resources found in the %s input", ProjectResourceType, f.Format)
		}
		return nil, fmt.Errorf("resource %s not found in the %s input, available addresses: %s", address, f.Format, strings.Join(f.Addresses(), ", "))
	default:
		found := make([]string, 0, len(matches))
		for _, match := range matches {
			found = append(found, match.Address)
		}
		sort.Strings(found)
		return nil, fmt.Errorf("address %s is ambiguous, it matches %s: pass the module qualified address", address, strings.Join(found, ", "))
	}
}

// indexSuffix spells an instance key the way Terraform addresses it: numeric keys bare, string keys quoted.
func indexSuffix(key any) string {
	switch value := key.(type) {
	case nil:
		return ""
	case string:
		return "[" + strconv.Quote(value) + "]"
	case float64:
		if value == float64(int64(value)) {
			return "[" + strconv.FormatInt(int64(value), 10) + "]"
		}
		return "[" + strconv.FormatFloat(value, 'f', -1, 64) + "]"
	case int:
		return "[" + strconv.Itoa(value) + "]"
	case int64:
		return "[" + strconv.FormatInt(value, 10) + "]"
	case json.Number:
		return "[" + value.String() + "]"
	default:
		return "[" + fmt.Sprintf("%v", value) + "]"
	}
}
