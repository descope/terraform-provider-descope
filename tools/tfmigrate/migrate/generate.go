package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Generated artifact layout. Every path is relative to the output directory and every file is deterministic, so a
// rerun over the same state produces the same bytes and reviews as an empty diff.
const (
	manifestPath    = "manifest.json"
	readmePath      = "README.md"
	blockersPath    = "BLOCKERS.md"
	secretsPath     = "SECRETS.md"
	rollbackPath    = "ROLLBACK.md"
	backupStatePath = "backup/legacy-state.json"
	backupShowPath  = "backup/legacy-show.json"
	detachDir       = "detach"
	adoptDir        = "adopt"
)

// Artifacts is the full set of files a migration run produces, keyed by their path in the output directory.
type Artifacts struct {
	Files map[string][]byte
}

// Generate renders every migration artifact. When the manifest is not ready only the manifest, the blocker report,
// the rollback notes and the state backup are produced: a migration with blockers must never look runnable.
func Generate(m *Manifest, blocks []*resourceBlock, payloads map[string][]byte, stateSnapshot []byte, opts Options) (*Artifacts, error) {
	manifestJSON, err := m.Encode()
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{
		manifestPath:           manifestJSON,
		m.Rollback.StateBackup: stateSnapshot,
		rollbackPath:           []byte(rollbackReport(m)),
	}
	if !m.Ready {
		files[blockersPath] = []byte(blockerReport(m))
		files[readmePath] = []byte(blockedReadme(m))
		return &Artifacts{Files: files}, nil
	}
	files[secretsPath] = []byte(secretReport(m))
	files[readmePath] = []byte(readme(m, opts))
	files[filepath.Join(detachDir, "removed.tf")] = detachConfig(m)
	files[filepath.Join(detachDir, "versions.tf")] = versionsConfig(m, m.SourceProvider)
	files[filepath.Join(adoptDir, "main.tf")] = resourceConfig(m, blocks)
	files[filepath.Join(adoptDir, "imports.tf")] = importConfig(m, blocks)
	files[filepath.Join(adoptDir, "versions.tf")] = versionsConfig(m, m.TargetProvider)
	files[filepath.Join(adoptDir, "variables.tf")] = variablesConfig(m)
	for _, file := range sortedKeys(payloads) {
		files[filepath.Join(adoptDir, file)] = payloads[file]
	}
	if err := validateTerraformFiles(files); err != nil {
		return nil, err
	}
	return &Artifacts{Files: files}, nil
}

// Write builds artifacts in a fresh private directory and atomically publishes that directory at dir. The resolved
// parent must not contain symlinks, and dir must not exist, so state bytes cannot be redirected to another target.
func (a *Artifacts) Write(dir string) error {
	dir = filepath.Clean(dir)
	names := sortedKeys(a.Files)
	for _, name := range names {
		clean := filepath.Clean(name)
		if filepath.IsAbs(clean) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return fmt.Errorf("artifact path %q escapes the output directory", name)
		}
	}
	if _, err := os.Lstat(dir); err == nil {
		return fmt.Errorf("%s already exists: use a new path so a rerun cannot replace or mix reviewed output", dir)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking output directory %s: %w", dir, err)
	}
	parent, err := filepath.Abs(filepath.Dir(dir))
	if err != nil {
		return fmt.Errorf("resolving output parent: %w", err)
	}
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return fmt.Errorf("creating output parent %s: %w", parent, err)
	}
	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return fmt.Errorf("resolving output parent %s: %w", parent, err)
	}
	if resolvedParent != parent {
		return fmt.Errorf("output parent %s contains a symbolic link; choose a direct path", parent)
	}
	temporary, err := os.MkdirTemp(parent, "."+filepath.Base(dir)+".tmp-")
	if err != nil {
		return fmt.Errorf("creating private output directory: %w", err)
	}
	defer func() {
		if temporary != "" {
			_ = os.RemoveAll(temporary)
		}
	}()
	for _, name := range names {
		clean := filepath.Clean(name)
		target := filepath.Join(temporary, clean)
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return fmt.Errorf("creating %s: %w", filepath.Dir(target), err)
		}
		file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return fmt.Errorf("creating %s: %w", target, err)
		}
		if _, err := file.Write(a.Files[name]); err != nil {
			_ = file.Close()
			return fmt.Errorf("writing %s: %w", target, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("closing %s: %w", target, err)
		}
	}
	final := filepath.Join(resolvedParent, filepath.Base(dir))
	if err := os.Rename(temporary, final); err != nil {
		return fmt.Errorf("publishing output directory %s: %w", dir, err)
	}
	temporary = ""
	return nil
}

func validateTerraformFiles(files map[string][]byte) error {
	for _, name := range sortedKeys(files) {
		if filepath.Ext(name) != ".tf" {
			continue
		}
		_, diagnostics := hclparse.NewParser().ParseHCL(files[name], name)
		if diagnostics.HasErrors() {
			return fmt.Errorf("generated invalid Terraform configuration in %s: %s", name, diagnostics.Error())
		}
	}
	return nil
}

// detachConfig renders the removed block that drops the legacy resource from state without deleting anything.
func detachConfig(m *Manifest) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendUnstructuredTokens(comment(
		"Stage 1: detach the legacy resource from Terraform state.",
		"Apply this with the source provider still pinned to "+m.SourceProvider.Version+".",
		"destroy = false keeps every remote object in place: nothing in Descope is deleted.",
		"This file belongs in the module that declares "+m.LegacyLocalAddr+".",
	))
	block := body.AppendNewBlock("removed", nil)
	block.Body().SetAttributeRaw("from", rawExpr(withoutIndex(m.LegacyLocalAddr)).tokens())
	lifecycle := block.Body().AppendNewBlock("lifecycle", nil)
	lifecycle.Body().SetAttributeValue("destroy", cty.False)
	return hclwrite.Format(file.Bytes())
}

// versionsConfig pins Terraform and the provider version for one stage of the migration.
func versionsConfig(m *Manifest, pin ProviderPin) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendUnstructuredTokens(comment(
		"Terraform "+m.RequiredTerraform+" is required end to end. Terraform 1.7 introduced removed blocks and",
		"for_each imports; Terraform 1.8 adds the plan completeness signal used by verify-plan.",
		"The provider version is pinned exactly, so both migration stages run against a known implementation.",
	))
	terraform := body.AppendNewBlock("terraform", nil)
	terraform.Body().SetAttributeValue("required_version", cty.StringVal(m.RequiredTerraform))
	providers := terraform.Body().AppendNewBlock("required_providers", nil)
	providers.Body().SetAttributeValue("descope", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal(pin.Source),
		"version": cty.StringVal(pin.Version),
	}))
	return hclwrite.Format(file.Bytes())
}

// resourceConfig renders the standalone resource configuration for every adopted object.
func resourceConfig(m *Manifest, blocks []*resourceBlock) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendUnstructuredTokens(comment(
		"Stage 2: standalone resource configuration for the objects the legacy resource owned.",
		"Apply this with the provider pinned to "+m.TargetProvider.Version+", together with imports.tf.",
		"Server ids are written as literals so the adoption plan is evaluable before anything has been applied.",
		"Replacing a literal id with a reference to the resource that owns it is a no-op refactor to do afterwards.",
		"This file belongs in the module that declares "+m.LegacyLocalAddr+".",
	))
	body.AppendNewline()
	appendProjectResource(body, m)
	for _, group := range groupBlocks(blocks) {
		body.AppendNewline()
		if !group.forEach {
			block := body.AppendNewBlock("resource", []string{group.resourceType, group.label})
			writeAttributes(block.Body(), group.members[0].Attrs)
			continue
		}
		locals := body.AppendNewBlock("locals", nil)
		locals.Body().SetAttributeValue(group.localName(), group.keyedPayloads())
		body.AppendNewline()
		block := body.AppendNewBlock("resource", []string{group.resourceType, group.label})
		blockBody := block.Body()
		blockBody.SetAttributeRaw("for_each", rawExpr("local."+group.localName()).tokens())
		blockBody.SetAttributeValue("project_id", cty.StringVal(m.ProjectID))
		blockBody.SetAttributeRaw(group.keyAttr, rawExpr("each.key").tokens())
		blockBody.SetAttributeRaw(group.payloadAttr, rawExpr(`file("${path.module}/${each.value}")`).tokens())
	}
	return hclwrite.Format(file.Bytes())
}

// appendProjectResource restores the descope_project instance with only the core attributes that remain after the
// split. A single for_each key and count index zero can be reproduced exactly from state; unsafe shapes are blocked
// during inventory and never reach this renderer.
func appendProjectResource(body *hclwrite.Body, m *Manifest) {
	block := body.AppendNewBlock("resource", []string{ProjectResourceType, legacyResourceLabel(m.LegacyLocalAddr)})
	blockBody := block.Body()
	switch key := m.InstanceKey.(type) {
	case nil:
		writeProjectCore(blockBody, m.ProjectCore)
	case string:
		blockBody.SetAttributeValue("for_each", cty.MapVal(map[string]cty.Value{key: projectCoreValue(m.ProjectCore)}))
		blockBody.SetAttributeRaw("name", rawExpr("each.value.name").tokens())
		blockBody.SetAttributeRaw("environment", rawExpr("each.value.environment").tokens())
		blockBody.SetAttributeRaw("tags", rawExpr("each.value.tags").tokens())
		if m.ProjectCore.DeletionProtection != nil {
			blockBody.SetAttributeRaw("deletion_protection", rawExpr("each.value.deletion_protection").tokens())
		}
	case float64, int, int64, json.Number:
		blockBody.SetAttributeValue("count", cty.NumberIntVal(1))
		writeProjectCore(blockBody, m.ProjectCore)
	}
}

func writeProjectCore(body *hclwrite.Body, core ProjectCore) {
	body.SetAttributeValue("name", cty.StringVal(core.Name))
	body.SetAttributeValue("environment", cty.StringVal(core.Environment))
	body.SetAttributeValue("tags", stringSet(core.Tags))
	if core.DeletionProtection != nil {
		body.SetAttributeValue("deletion_protection", cty.BoolVal(*core.DeletionProtection))
	}
}

func projectCoreValue(core ProjectCore) cty.Value {
	attributes := map[string]cty.Value{
		"name":        cty.StringVal(core.Name),
		"environment": cty.StringVal(core.Environment),
		"tags":        stringSet(core.Tags),
	}
	if core.DeletionProtection != nil {
		attributes["deletion_protection"] = cty.BoolVal(*core.DeletionProtection)
	}
	return cty.ObjectVal(attributes)
}

func stringSet(values []string) cty.Value {
	if len(values) == 0 {
		return cty.SetValEmpty(cty.String)
	}
	items := make([]cty.Value, 0, len(values))
	for _, value := range values {
		items = append(items, cty.StringVal(value))
	}
	return cty.SetVal(items)
}

func legacyResourceLabel(address string) string {
	parts := strings.Split(withoutIndex(address), ".")
	if len(parts) < 2 {
		return "project"
	}
	return parts[len(parts)-1]
}

// importConfig renders one import block per adopted object, and one for_each import block per keyed collection.
// Import ids are literals because Terraform evaluates them before any resource exists.
func importConfig(m *Manifest, blocks []*resourceBlock) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendUnstructuredTokens(comment(
		"Stage 2: native import blocks adopting the existing Descope objects.",
		"Import blocks are evaluated in the root module, so this file belongs in the root module even when the",
		"resources it adopts are declared in a child module.",
		"The legacy project resource is imported back here too, after stage 1 detached it from state.",
	))
	body.AppendNewline()
	project := body.AppendNewBlock("import", nil)
	project.Body().SetAttributeRaw("to", rawExpr(m.LegacyAddress).tokens())
	project.Body().SetAttributeValue("id", cty.StringVal(m.ProjectID))

	for _, group := range groupBlocks(blocks) {
		body.AppendNewline()
		if group.forEach {
			block := body.AppendNewBlock("import", nil)
			block.Body().SetAttributeValue("for_each", group.keyedImportIDs())
			block.Body().SetAttributeRaw("to", rawExpr(group.qualifiedAddress(m)+"[each.key]").tokens())
			block.Body().SetAttributeRaw("id", rawExpr("each.value").tokens())
			continue
		}
		member := group.members[0]
		block := body.AppendNewBlock("import", nil)
		block.Body().SetAttributeRaw("to", rawExpr(member.Entry.Address).tokens())
		block.Body().SetAttributeValue("id", cty.StringVal(member.Entry.ImportID))
	}
	return hclwrite.Format(file.Bytes())
}

// variablesConfig declares the sensitive inputs the adoption needs, without any value.
func variablesConfig(m *Manifest) []byte {
	declared := map[string]SecretRequirement{}
	for _, secret := range m.Secrets {
		if secret.Variable == "" {
			continue
		}
		if _, seen := declared[secret.Variable]; !seen {
			declared[secret.Variable] = secret
		}
	}
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.AppendUnstructuredTokens(comment(
		"Sensitive inputs for the adoption stage. No value is generated: supply each one from the secret store the",
		"team already uses. SECRETS.md records which legacy state attribute each variable corresponds to.",
	))
	for _, name := range sortedKeys(declared) {
		body.AppendNewline()
		block := body.AppendNewBlock("variable", []string{name})
		typeExpression := declared[name].Type
		if typeExpression == "" {
			typeExpression = "string"
		}
		block.Body().SetAttributeRaw("type", rawExpr(typeExpression).tokens())
		block.Body().SetAttributeValue("sensitive", cty.True)
		block.Body().SetAttributeValue("description", cty.StringVal(declared[name].Address+" attribute "+declared[name].Attribute))
	}
	return hclwrite.Format(file.Bytes())
}

// writeAttributes renders one resource body, ordering the addressing attributes first and the rest alphabetically.
func writeAttributes(body *hclwrite.Body, attrs map[string]*expr) {
	leading := []string{"project_id", "app_id", "method", "id", "name"}
	written := map[string]bool{}
	for _, name := range leading {
		if value, ok := attrs[name]; ok {
			setAttribute(body, name, value)
			written[name] = true
		}
	}
	for _, name := range sortedKeys(attrs) {
		if written[name] {
			continue
		}
		setAttribute(body, name, attrs[name])
	}
}

func setAttribute(body *hclwrite.Body, name string, value *expr) {
	if literal, ok := value.literal(); ok {
		body.SetAttributeValue(name, literal)
		return
	}
	body.SetAttributeRaw(name, value.tokens())
}

// blockGroup is either a single resource block or the members of one keyed collection.
type blockGroup struct {
	resourceType string
	label        string
	forEach      bool
	keyAttr      string
	payloadAttr  string
	members      []*resourceBlock
}

func (g *blockGroup) localName() string {
	return sanitizeLabel(g.resourceType + "_" + g.label)
}

// keyedPayloads is the collection key to payload file map that drives the for_each resource and import blocks.
func (g *blockGroup) keyedPayloads() cty.Value {
	entries := map[string]cty.Value{}
	for _, member := range g.members {
		file := ""
		if len(member.Entry.Payloads) > 0 {
			file = member.Entry.Payloads[0]
		}
		entries[member.ForEachKey] = cty.StringVal(file)
	}
	return cty.MapVal(entries)
}

// qualifiedAddress is the group's address as the root module addresses it.
func (g *blockGroup) qualifiedAddress(m *Manifest) string {
	address := g.resourceType + "." + g.label
	if m.ModuleAddress == "" {
		return address
	}
	return m.ModuleAddress + "." + address
}

// keyedImportIDs is self contained in the root import block; it cannot reference the child module local that drives
// the resource's for_each.
func (g *blockGroup) keyedImportIDs() cty.Value {
	entries := map[string]cty.Value{}
	for _, member := range g.members {
		entries[member.ForEachKey] = cty.StringVal(member.Entry.ImportID)
	}
	return cty.MapVal(entries)
}

// groupBlocks orders the generated blocks deterministically, collapsing the members of a keyed collection.
func groupBlocks(blocks []*resourceBlock) []*blockGroup {
	groups := []*blockGroup{}
	index := map[string]*blockGroup{}
	for _, block := range blocks {
		key := block.Type + "." + block.Label
		group, seen := index[key]
		if !seen {
			group = &blockGroup{resourceType: block.Type, label: block.Label, forEach: block.ForEachKey != "", keyAttr: block.KeyAttr}
			if group.forEach {
				group.payloadAttr = payloadAttribute(block)
			}
			index[key] = group
			groups = append(groups, group)
		}
		group.members = append(group.members, block)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		if groups[i].resourceType != groups[j].resourceType {
			return groups[i].resourceType < groups[j].resourceType
		}
		return groups[i].label < groups[j].label
	})
	for _, group := range groups {
		sort.SliceStable(group.members, func(i, j int) bool { return group.members[i].ForEachKey < group.members[j].ForEachKey })
	}
	return groups
}

// payloadAttribute names the single document attribute of a keyed collection member.
func payloadAttribute(block *resourceBlock) string {
	for _, name := range sortedKeys(block.Attrs) {
		if name != "project_id" && name != block.KeyAttr {
			return name
		}
	}
	return ""
}

// withoutIndex strips the final resource instance key while preserving indexed parent module addresses.
func withoutIndex(address string) string {
	_, traversal, err := canonicalTraversal(address)
	if err != nil || len(traversal) == 0 {
		return address
	}
	if _, indexed := traversal[len(traversal)-1].(hcl.TraverseIndex); indexed {
		traversal = traversal[:len(traversal)-1]
	}
	return string(hclwrite.TokensForTraversal(traversal).Bytes())
}

// comment renders leading comment lines as unstructured tokens, which is how hclwrite carries free text.
func comment(lines ...string) hclwrite.Tokens {
	tokens := hclwrite.Tokens{}
	for _, line := range lines {
		tokens = append(tokens, &hclwrite.Token{Type: hclsyntax.TokenComment, Bytes: []byte("# " + line + "\n")})
	}
	return tokens
}
