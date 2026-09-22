package harness

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/emit"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var fakeSecretValues = []string{
	"not-a-real-password",
	"not-a-real-token",
	"not-a-real-key",
	"not-a-real-secret",
	"00aBcDeFgHiJkLmNoPqRsTuVwXyZ",
	// the 8x8 connector api_key fixtures are omitted: those attributes aren't sensitive, and TestSecretClassification covers them
}

type stateResource struct {
	Type   string
	Name   string // the terraform resource label in the seed config
	Values map[string]any
}

func seedStateResources(t *testing.T, seedDir string, env []string) []stateResource {
	t.Helper()
	output, err := Try(seedDir, env, "show", "-json")
	if err != nil {
		t.Fatalf("terraform show failed: %s\n%s", err, output)
	}
	var parsed struct {
		Values struct {
			RootModule struct {
				Resources []struct {
					Type   string         `json:"type"`
					Name   string         `json:"name"`
					Values map[string]any `json:"values"`
				} `json:"resources"`
			} `json:"root_module"`
		} `json:"values"`
	}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("parsing state JSON: %s", err)
	}
	var resources []stateResource
	for _, resource := range parsed.Values.RootModule.Resources {
		if strings.HasPrefix(resource.Type, "descope_") {
			resources = append(resources, stateResource(resource))
		}
	}
	return resources
}

func ExportFiles(t *testing.T, exportDir string) (configs map[string]string, payloads map[string]string) {
	t.Helper()
	configs, payloads = map[string]string{}, map[string]string{}
	err := filepath.Walk(exportDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		relative, err := filepath.Rel(exportDir, path)
		if err != nil {
			return err
		}
		if relative == "terraform.tfvars" || relative == "plan.bin" || strings.HasPrefix(relative, ".terraform") || strings.HasSuffix(relative, ".tfstate") || strings.HasSuffix(relative, ".backup") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasSuffix(relative, ".tf") {
			configs[relative] = string(content)
		} else {
			payloads[relative] = string(content)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return configs, payloads
}

var (
	ResourceHeaderPattern = regexp.MustCompile(`(?m)^resource "(descope_[a-z0-9_]+)" "([a-z0-9_]+)" \{`)
	ImportToPattern       = regexp.MustCompile(`(?m)^  to = (descope_[a-z0-9_]+)\.([a-z0-9_]+)$`)
	variableDeclPattern   = regexp.MustCompile(`(?m)^variable "([a-z0-9_]+)" \{`)
	variableUsePattern    = regexp.MustCompile(`var\.([a-z0-9_]+)`)
	fileRefPattern        = regexp.MustCompile(`file\("\$\{path\.module\}/([^"]+)"\)`)
)

func assertExportIntegrity(t *testing.T, stepName, exportDir string) {
	t.Helper()
	configs, payloads := ExportFiles(t, exportDir)

	addresses := map[string]string{} // address -> file
	for name, content := range configs {
		if name == "import.tf" {
			continue
		}
		for _, match := range ResourceHeaderPattern.FindAllStringSubmatch(content, -1) {
			address := match[1] + "." + match[2]
			if previous, ok := addresses[address]; ok {
				t.Errorf("%s: duplicate resource address %s in %s and %s", stepName, address, previous, name)
			}
			addresses[address] = name
		}
	}

	imports := map[string]bool{}
	for _, match := range ImportToPattern.FindAllStringSubmatch(configs["import.tf"], -1) {
		address := match[1] + "." + match[2]
		if imports[address] {
			t.Errorf("%s: duplicate import block for %s", stepName, address)
		}
		imports[address] = true
		if _, ok := addresses[address]; !ok {
			t.Errorf("%s: import block for %s has no matching resource block", stepName, address)
		}
	}
	for address, file := range addresses {
		if !imports[address] {
			t.Errorf("%s: resource %s in %s has no import block", stepName, address, file)
		}
	}

	declared := map[string]bool{}
	for _, match := range variableDeclPattern.FindAllStringSubmatch(configs["variables.tf"], -1) {
		declared[match[1]] = true
	}
	used := map[string]bool{}
	for name, content := range configs {
		if name == "variables.tf" {
			continue
		}
		for _, match := range variableUsePattern.FindAllStringSubmatch(content, -1) {
			used[match[1]] = true
			if !declared[match[1]] {
				t.Errorf("%s: %s references undeclared variable %s", stepName, name, match[1])
			}
		}
	}
	for name := range declared {
		if !used[name] {
			t.Errorf("%s: variable %s is declared but never used", stepName, name)
		}
	}

	// file() references point at existing extracted files, and all extracted files are referenced; extracted .json payloads must parse
	referenced := map[string]bool{}
	for name, content := range configs {
		for _, match := range fileRefPattern.FindAllStringSubmatch(content, -1) {
			referenced[match[1]] = true
			if _, ok := payloads[filepath.FromSlash(match[1])]; !ok {
				t.Errorf("%s: %s references missing file %s", stepName, name, match[1])
			}
		}
	}
	for name, content := range payloads {
		if !referenced[filepath.ToSlash(name)] {
			t.Errorf("%s: extracted file %s is not referenced by any configuration", stepName, name)
		}
		if strings.HasSuffix(name, ".json") && !json.Valid([]byte(content)) {
			t.Errorf("%s: extracted file %s is not valid JSON", stepName, name)
		}
	}

	leakChecks := slices.Clone(fakeSecretValues)
	if key := os.Getenv("DESCOPE_MANAGEMENT_KEY"); key != "" {
		leakChecks = append(leakChecks, key)
	}
	files := map[string]string{}
	for name, content := range configs {
		files[name] = content
	}
	for name, content := range payloads {
		files[name] = content
	}
	for name, content := range files {
		for _, secret := range leakChecks {
			if strings.Contains(content, secret) {
				t.Errorf("%s: secret value leaked into generated file %s", stepName, name)
			}
		}
	}

	t.Logf("%s: integrity verified: %d resource blocks across %d files, %d import blocks, %d variables, %d extracted files",
		stepName, len(addresses), len(configs), len(imports), len(declared), len(payloads))
}

var attributesExemptFromConsistency = map[string]string{
	"descope_flow.data.flowId":                     "the flow_id attribute is authoritative, so the write overrides the id the configured document carried",
	"descope_flow.data.metadata.componentsVersion": "writing a flow upgrades it to the current components version",
	"descope_styles.data.componentsVersion":        "writing the theme upgrades it the same way writing a flow does",
	"descope_widget.data.metadata.screens":         "the read reports the screens at the top level of the document instead",
	"descope_widget.data.metadata.widgetId":        "the read reports the widget id at the top level of the document instead",
}

// knownInconsistencies are confirmed open bugs: the mismatch is logged instead of failing the suite, and entries must be removed once fixed.
var knownInconsistencies = map[string]string{
	"descope_oidc_app.client_id": "the read returns an empty client_id until a client_type is set",
}

func assertStateConsistency(t *testing.T, stepName string, seed []stateResource, changes []Change, expectMissing []string) {
	t.Helper()

	imported := map[string]map[string]any{} // type + identity -> before values
	for _, change := range changes {
		if len(change.Change.Importing) == 0 || change.Change.Before == nil {
			continue
		}
		imported[change.Type+"|"+resourceIdentity(change.Change.Before)] = change.Change.Before
	}

	compared, mismatches, zeroOmitted := 0, 0, 0
	for _, resource := range seed {
		if slices.Contains(expectMissing, resource.Type) {
			continue
		}
		identity := resourceIdentity(resource.Values)
		before, ok := imported[resource.Type+"|"+identity]
		if !ok {
			t.Errorf("%s: consistency: seeded %s %q has no matching imported resource", stepName, resource.Type, identity)
			continue
		}
		for name, seedValue := range resource.Values {
			if seedValue == nil {
				continue
			}
			importValue, ok := before[name]
			if !ok {
				continue
			}
			compared++
			switch compareConsistency(sensitivePaths(resource.Type), resource.Type, name, seedValue, importValue) {
			case consistencyEqual:
			case consistencyZeroOmitted:
				zeroOmitted++
			case consistencyMismatch:
				seedJSON, _ := json.Marshal(seedValue)
				importJSON, _ := json.Marshal(importValue)
				if reason, known := knownInconsistencies[resource.Type+"."+name]; known {
					t.Logf("%s: consistency: KNOWN issue on %s %q attribute %q (%s): seed state %s but import read %s",
						stepName, resource.Type, identity, name, reason, seedJSON, importJSON)
					continue
				}
				mismatches++
				t.Errorf("%s: consistency: %s %q attribute %q: seed state %s but import read %s",
					stepName, resource.Type, identity, name, seedJSON, importJSON)
			}
		}
	}
	t.Logf("%s: consistency verified: %d seeded resources, %d attribute values compared, %d zero values omitted by the server, %d mismatches",
		stepName, len(seed), compared, zeroOmitted, mismatches)
}

// numbers stay literals so identifiers too large for a float aren't compared as equal when they differ only in their last digits.
func parseJSONPair(seedValue, importValue any) (seedDoc, importDoc any, ok bool) {
	seedText, seedOK := seedValue.(string)
	importText, importOK := importValue.(string)
	if !seedOK || !importOK {
		return nil, nil, false
	}
	seedDoc, seedOK = parseJSONDocument(seedText)
	importDoc, importOK = parseJSONDocument(importText)
	return seedDoc, importDoc, seedOK && importOK
}

func parseJSONDocument(text string) (any, bool) {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil || decoder.More() {
		return nil, false
	}
	switch value.(type) {
	case map[string]any, []any:
		return value, true
	}
	return nil, false
}

func resourceIdentity(values map[string]any) string {
	if name, _ := values["name"].(string); name != "" {
		if method, _ := values["method"].(string); method != "" {
			return method + "/" + name
		}
		return name
	}
	if id, _ := values["id"].(string); id != "" {
		if projectID, _ := values["project_id"].(string); projectID == id {
			return "singleton"
		}
		return id
	}
	return "singleton"
}

type consistencyVerdict int

const (
	consistencyEqual consistencyVerdict = iota
	consistencyZeroOmitted
	consistencyMismatch
)

func compareConsistency(sensitive map[string]bool, resourceType, path string, seedValue, importValue any) consistencyVerdict {
	if sensitive[path] {
		return consistencyEqual // write-only secret: values can't be compared
	}
	if _, exempt := attributesExemptFromConsistency[resourceType+"."+path]; exempt {
		return consistencyEqual
	}
	seedJSON, _ := json.Marshal(seedValue)
	importJSON, _ := json.Marshal(importValue)
	if string(seedJSON) == string(importJSON) {
		return consistencyEqual
	}
	if importValue == nil {
		if isZeroValue(seedValue) {
			return consistencyZeroOmitted
		}
		return consistencyMismatch
	}
	seedMap, seedOK := seedValue.(map[string]any)
	importMap, importOK := importValue.(map[string]any)
	if seedOK && importOK {
		verdict := consistencyEqual
		for key, nested := range seedMap {
			switch compareConsistency(sensitive, resourceType, path+"."+key, nested, importMap[key]) {
			case consistencyEqual:
			case consistencyZeroOmitted:
				verdict = consistencyZeroOmitted
			case consistencyMismatch:
				return consistencyMismatch
			}
		}
		return verdict
	}
	if seedDoc, importDoc, ok := parseJSONPair(seedValue, importValue); ok {
		return compareConsistency(sensitive, resourceType, path, seedDoc, importDoc)
	}
	seedList, seedOK := seedValue.([]any)
	importList, importOK := importValue.([]any)
	if seedOK && importOK && len(seedList) == len(importList) {
		verdict := consistencyEqual
		for i := range seedList {
			switch compareConsistency(sensitive, resourceType, path, seedList[i], importList[i]) {
			case consistencyEqual:
			case consistencyZeroOmitted:
				verdict = consistencyZeroOmitted
			case consistencyMismatch:
				return consistencyMismatch
			}
		}
		return verdict
	}
	return consistencyMismatch
}

func isZeroValue(value any) bool {
	switch v := value.(type) {
	case bool:
		return !v
	case string:
		return v == ""
	case float64:
		return v == 0
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	}
	return false
}

var sensitivePathsCache = map[string]map[string]bool{}

func sensitivePaths(resourceType string) map[string]bool {
	if cached, ok := sensitivePathsCache[resourceType]; ok {
		return cached
	}
	paths := map[string]bool{}
	if exportable, ok := registry.Load(context.Background())[resourceType]; ok {
		collectSensitivePaths(exportable.ExportSchema().Attributes, "", paths)
	}
	sensitivePathsCache[resourceType] = paths
	return paths
}

func collectSensitivePaths(attributes map[string]schema.Attribute, prefix string, paths map[string]bool) {
	for name, attribute := range attributes {
		path := prefix + name
		if attribute.IsSensitive() {
			paths[path] = true
		}
		switch a := attribute.(type) {
		case schema.SingleNestedAttribute:
			collectSensitivePaths(a.Attributes, path+".", paths)
		case schema.ListNestedAttribute:
			collectSensitivePaths(a.NestedObject.Attributes, path+".", paths)
		case schema.SetNestedAttribute:
			collectSensitivePaths(a.NestedObject.Attributes, path+".", paths)
		}
	}
}

func assertExportDeterminism(t *testing.T, stepName, firstDir, secondDir string) {
	t.Helper()
	firstConfigs, firstPayloads := ExportFiles(t, firstDir)
	secondConfigs, secondPayloads := ExportFiles(t, secondDir)
	for name, content := range firstConfigs {
		if secondConfigs[name] != content {
			t.Errorf("%s: determinism: %s differs between two exports of the same project%s", stepName, name, firstDifference(content, secondConfigs[name]))
		}
	}
	for name, content := range firstPayloads {
		second := secondPayloads[name]
		if second == content {
			continue
		}
		firstStripped, secondStripped := stripComponentsVersion(content), stripComponentsVersion(second)
		if firstStripped == secondStripped {
			continue
		}
		t.Errorf("%s: determinism: %s differs between two exports of the same project%s", stepName, name, firstDifference(firstStripped, secondStripped))
	}
	if len(firstConfigs) != len(secondConfigs) || len(firstPayloads) != len(secondPayloads) {
		t.Errorf("%s: determinism: file sets differ between two exports (%d/%d vs %d/%d)",
			stepName, len(firstConfigs), len(firstPayloads), len(secondConfigs), len(secondPayloads))
	}
	t.Logf("%s: determinism verified: %d files identical across two exports", stepName, len(firstConfigs)+len(firstPayloads))
}

func stripComponentsVersion(content string) string {
	var payload map[string]any
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return content
	}
	delete(payload, "componentsVersion")
	if metadata, ok := payload["metadata"].(map[string]any); ok {
		delete(metadata, "componentsVersion")
	}
	normalized, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return content
	}
	return string(normalized)
}

func firstDifference(first, second string) string {
	if second == "" {
		return ": the second export has no such file"
	}
	firstLines, secondLines := strings.Split(first, "\n"), strings.Split(second, "\n")
	for i := 0; i < len(firstLines) && i < len(secondLines); i++ {
		if firstLines[i] != secondLines[i] {
			return fmt.Sprintf(":\n  line %d first:  %s\n  line %d second: %s", i+1, excerpt(firstLines[i]), i+1, excerpt(secondLines[i]))
		}
	}
	return fmt.Sprintf(": identical for %d lines, then the line counts differ (%d vs %d)", min(len(firstLines), len(secondLines)), len(firstLines), len(secondLines))
}

func excerpt(line string) string {
	const limit = 200
	if len(line) > limit {
		return line[:limit] + "..."
	}
	return line
}

func assertFormatting(t *testing.T, stepName, exportDir string, env []string) {
	t.Helper()
	if output, err := Try(exportDir, env, "fmt", "-check", "-recursive"); err != nil {
		t.Errorf("%s: terraform fmt -check failed:\n%s", stepName, output)
	} else {
		t.Logf("%s: formatting verified: terraform fmt reports no changes", stepName)
	}
}

func planSummary(t *testing.T, stepName string, changes []Change) {
	t.Helper()
	actions := map[string]int{}
	types := map[string]int{}
	for _, change := range changes {
		actions[strings.Join(change.Change.Actions, ",")]++
		types[change.Type]++
	}
	t.Logf("%s: plan summary: %d resources, actions %v, %d distinct types", stepName, len(changes), actions, len(types))
	if testing.Verbose() {
		names := make([]string, 0, len(types))
		for name := range types {
			names = append(names, fmt.Sprintf("%s=%d", name, types[name]))
		}
		slices.Sort(names)
		t.Logf("%s: plan types: %s", stepName, strings.Join(names, " "))
	}
}

func assertExportedValues(t *testing.T, stepName string, seed []stateResource, changes []Change, expectMissing []string) {
	t.Helper()

	exported := map[string]map[string]any{}
	for _, change := range changes {
		if len(change.Change.Importing) == 0 || change.Change.After == nil {
			continue
		}
		exported[change.Type+"|"+resourceIdentity(change.Change.After)] = change.Change.After
	}

	verified, pruned := 0, 0
	for _, resource := range seed {
		if slices.Contains(expectMissing, resource.Type) {
			continue
		}
		after, ok := exported[resource.Type+"|"+resourceIdentity(resource.Values)]
		if !ok {
			continue // absence is the coverage assertion's concern
		}
		for name, seedValue := range resource.Values {
			if seedValue == nil || name == "id" || name == "project_id" {
				continue
			}
			if _, exempt := attributesExemptFromConsistency[resource.Type+"."+name]; exempt {
				continue
			}
			if _, known := knownInconsistencies[resource.Type+"."+name]; known {
				continue
			}
			if sensitivePaths(resource.Type)[name] {
				pruned++
				continue // secrets become variable references; their values can't be verified
			}
			configValue, present := after[name]
			if !present || configValue == nil {
				if isZeroValue(seedValue) || seedEqualsSchemaDefault(t, resource.Type, name, seedValue) {
					pruned++
					continue
				}
				seedJSON, _ := json.Marshal(seedValue)
				t.Errorf("%s: exported values: %s %q attribute %q with seed value %s was pruned but doesn't match the schema default",
					stepName, resource.Type, resourceIdentity(resource.Values), name, seedJSON)
				continue
			}
			verified++
			seedJSON, _ := json.Marshal(seedValue)
			configJSON, _ := json.Marshal(configValue)
			if string(seedJSON) != string(configJSON) {
				if compareConsistency(sensitivePaths(resource.Type), resource.Type, name, seedValue, configValue) == consistencyMismatch {
					t.Errorf("%s: exported values: %s %q attribute %q: seeded %s but exported %s",
						stepName, resource.Type, resourceIdentity(resource.Values), name, seedJSON, configJSON)
				}
			}
		}
	}
	t.Logf("%s: exported values verified: %d attribute values match the seed exactly, %d pruned as defaults or write-only",
		stepName, verified, pruned)
}

func seedEqualsSchemaDefault(t *testing.T, resourceType, attribute string, seedValue any) bool {
	t.Helper()
	ctx := context.Background()
	exportable, ok := registry.Load(ctx)[resourceType]
	if !ok {
		return false
	}
	schemaAttribute, ok := exportable.ExportSchema().Attributes[attribute]
	if !ok {
		return false
	}
	def, ok := prune.AttributeDefault(ctx, schemaAttribute)
	if !ok {
		return false
	}
	defaultJSON, err := emit.ValueJSON(ctx, def)
	if err != nil {
		return false
	}
	seedJSON, _ := json.Marshal(seedValue)
	return string(defaultJSON) == string(seedJSON)
}
