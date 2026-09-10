package migrate

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/zclconf/go-cty/cty"
)

// payloadDir holds the JSON documents that generated configuration references with file().
const payloadDir = "payloads"

// attrRef is the destination attribute a value is being translated onto. top is the attribute a Terraform plan
// reports a change on, and name is the full path used to name a sensitive input variable.
type attrRef struct {
	top  string
	name string
	path string
}

func (r attrRef) child(key string) attrRef {
	return attrRef{top: r.top, name: r.name + "_" + key, path: r.path + "." + key}
}

func (r attrRef) index(index int) attrRef {
	return attrRef{top: r.top, name: r.name + "_" + strconv.Itoa(index), path: fmt.Sprintf("%s[%d]", r.path, index)}
}

// emitSingleton adds the destination resource for a legacy object that became a per-project singleton.
func (b *builder) emitSingleton(f family, legacy map[string]any) {
	b.emit(f, legacy, Scope{ProjectID: b.projectID}, "", "")
}

// emitListItem adds the destination resource for one element of a legacy list, taking the backend id from the
// element's own id attribute.
func (b *builder) emitListItem(f family, legacy map[string]any) *Entry {
	return b.emitScoped(f, legacy, Scope{ProjectID: b.projectID})
}

// emitScoped adds the destination resource for one element of a legacy list under an explicit scope.
func (b *builder) emitScoped(f family, legacy map[string]any, scope Scope) *Entry {
	id, _ := legacy[f.idAttr].(string)
	name, _ := legacy[f.nameAttr].(string)
	if id == "" {
		b.blocker(Blocker{
			Kind:       BlockerMissingID,
			LegacyPath: f.path,
			Detail: fmt.Sprintf("%q has no %s in the legacy state, so its %s destination cannot be imported: refresh the state with the source provider, or rerun with -remote to read the id back from the project",
				b.describe(f, name), f.idAttr, f.resourceType),
		})
		return nil
	}
	return b.emit(f, legacy, scope, id, "")
}

// emitKeyed adds the destination resource for one entry of a legacy map, whose key is the entity id.
func (b *builder) emitKeyed(f family, key string, legacy map[string]any) {
	if key == "" {
		b.blocker(Blocker{
			Kind:       BlockerMissingID,
			LegacyPath: f.path,
			Detail:     fmt.Sprintf("a %s entry has an empty key, which is also its backend id", f.resourceType),
		})
		return
	}
	values := map[string]any{}
	for name, value := range legacy {
		values[name] = value
	}
	if f.keyAttr != "" {
		values[f.keyAttr] = key
	}
	b.emit(f, values, Scope{ProjectID: b.projectID}, key, key)
}

// emit translates one legacy object onto its destination resource and records the manifest entry. forEachKey is set
// for a family that is generated as a single for_each resource, in which case it is also the instance key.
func (b *builder) emit(f family, legacy map[string]any, scope Scope, backendID, forEachKey string) *Entry {
	dst, known := b.reg.Lookup(f.resourceType)
	if !known {
		b.blocker(Blocker{
			Kind:       BlockerUnsupportedField,
			LegacyPath: f.path,
			Detail:     fmt.Sprintf("the target provider does not register %s", f.resourceType),
		})
		return nil
	}
	scope.Kind = dst.Scope
	manifestID := backendID
	if dst.Scope == ScopeSingleton {
		manifestID = scope.ProjectID
	}
	importID, err := dst.ImportID(scope, backendID)
	if err != nil {
		b.blocker(Blocker{Kind: BlockerMissingID, LegacyPath: f.path, Detail: err.Error()})
		return nil
	}
	importKey := dst.Type + "\x00" + importID
	if previous, duplicate := b.imports[importKey]; duplicate {
		b.blocker(Blocker{
			Kind:       BlockerDuplicateAddress,
			LegacyPath: f.path,
			Detail:     fmt.Sprintf("%s and %s both map to %s with import id %q", previous, f.path, dst.Type, importID),
		})
		return nil
	}
	b.imports[importKey] = f.path

	existing, matches, unsupportedProviders := b.existing(dst.Type, backendID, scope)
	if len(unsupportedProviders) > 0 {
		b.blocker(Blocker{
			Kind:       BlockerAmbiguousOwnership,
			LegacyPath: f.path,
			Detail:     fmt.Sprintf("%s with backend id %q already has state under an unsupported provider configuration at %s", dst.Type, manifestID, strings.Join(unsupportedProviders, ", ")),
		})
		return nil
	}
	if len(matches) > 0 {
		if len(matches) > 1 {
			b.blocker(Blocker{
				Kind:       BlockerAmbiguousOwnership,
				LegacyPath: f.path,
				Detail:     fmt.Sprintf("%s with backend id %q is already represented by multiple standalone state addresses: %s", dst.Type, manifestID, strings.Join(matches, ", ")),
			})
			return nil
		}
		if previous, claimed := b.existingClaimed[existing.Address]; claimed {
			b.blocker(Blocker{
				Kind:       BlockerAmbiguousOwnership,
				LegacyPath: f.path,
				Address:    existing.Address,
				Detail:     fmt.Sprintf("the already standalone resource is also claimed by %s", previous),
			})
			return nil
		}
		b.existingClaimed[existing.Address] = f.path
		b.manifest.Standalone = append(b.manifest.Standalone, Standalone{
			Address:   existing.Address,
			Type:      existing.Type,
			BackendID: manifestID,
			Detail:    "already has standalone Terraform state and is never generated or imported by this migration",
		})
		return &Entry{
			LegacyPath: f.path,
			Address:    existing.Address,
			Type:       existing.Type,
			Scope:      scope,
			BackendID:  manifestID,
			ImportID:   importID,
			IDSource:   "already standalone state",
		}
	}

	grouped := f.group != "" && forEachKey != "" && groupable(dst, f)
	name, _ := legacy[f.nameAttr].(string)
	var label string
	instance := ""
	if grouped {
		groupKey := dst.Type + "\x00" + f.group
		label = b.groupLabels[groupKey]
		if label == "" {
			groupFamily := f
			groupFamily.labelPrefix = f.group
			label = b.label(dst.Type, groupFamily, "", "")
			b.groupLabels[groupKey] = label
		}
		instance = "[" + strconv.Quote(forEachKey) + "]"
	} else {
		forEachKey = ""
		label = b.label(dst.Type, f, name, backendID)
	}
	address := b.qualify(dst.Type + "." + label + instance)
	if existingOwner, duplicate := b.addresses[address]; duplicate {
		b.blocker(Blocker{
			Kind:       BlockerDuplicateAddress,
			LegacyPath: f.path,
			Address:    address,
			Detail:     fmt.Sprintf("%s is already claimed by %s, so two legacy objects would own the same destination", address, existingOwner),
		})
		return nil
	}
	b.addresses[address] = f.path

	entry := &Entry{
		LegacyPath: f.path,
		Address:    address,
		Type:       dst.Type,
		Label:      label,
		Scope:      scope,
		BackendID:  manifestID,
		ImportID:   importID,
		IDSource:   b.idSource(),
		References: map[string]string{},
	}
	payloadKey := label
	if forEachKey != "" {
		payloadKey = label + "." + sanitizePathSegment(forEachKey)
	}
	attrs := b.translate(dst, f, legacy, scope, entry, payloadKey)
	if len(entry.References) == 0 {
		entry.References = nil
	}
	b.manifest.Entries = append(b.manifest.Entries, *entry)
	b.blocks = append(b.blocks, &resourceBlock{
		Type:       dst.Type,
		Label:      label,
		Attrs:      attrs,
		Entry:      entry,
		Legacy:     f.path,
		ForEachKey: forEachKey,
		KeyAttr:    f.keyAttr,
	})
	return entry
}

// groupable reports whether a map keyed family can be generated as one for_each resource: that holds only when the
// destination has nothing beyond its computed id, its project scope, the collection key and one payload document.
func groupable(dst *Destination, f family) bool {
	if len(f.payload) != 1 || f.keyAttr == "" {
		return false
	}
	expected := map[string]bool{"id": true, "project_id": true, f.keyAttr: true, f.payload[0]: true}
	if len(dst.Schema.Attributes) != len(expected) {
		return false
	}
	for name := range dst.Schema.Attributes {
		if !expected[name] {
			return false
		}
	}
	return true
}

// translate maps a legacy object onto the destination schema. Anything populated that has no destination fails the
// migration closed rather than being dropped.
func (b *builder) translate(dst *Destination, f family, legacy map[string]any, scope Scope, entry *Entry, payloadKey string) map[string]*expr {
	attrs := map[string]*expr{}
	claimed := map[string]string{}

	set := func(name string, value *expr, source string) {
		if value.empty() {
			return
		}
		if previous, taken := claimed[name]; taken {
			b.blocker(Blocker{
				Kind:       BlockerAmbiguousOwnership,
				LegacyPath: f.path,
				Address:    entry.Address,
				Detail:     fmt.Sprintf("attribute %q of %s is claimed by both %s and %s", name, dst.Type, previous, source),
			})
			return
		}
		claimed[name] = source
		attrs[name] = value
	}

	// Scope attributes are literals: an import plan has to be evaluable before anything has been applied.
	if _, ok := dst.Schema.Attributes["project_id"]; ok {
		set("project_id", literalExpr(cty.StringVal(scope.ProjectID)), "the project scope")
	}
	if scope.AppID != "" {
		if _, ok := dst.Schema.Attributes["app_id"]; ok {
			set("app_id", literalExpr(cty.StringVal(scope.AppID)), "the application scope")
		}
	}
	if scope.Method != "" {
		if _, ok := dst.Schema.Attributes["method"]; ok {
			set("method", literalExpr(cty.StringVal(scope.Method)), "the auth method scope")
		}
	}

	for _, name := range sortedKeys(f.fixed) {
		if _, ok := dst.Schema.Attributes[name]; !ok {
			b.blocker(Blocker{
				Kind:       BlockerUnsupportedField,
				LegacyPath: f.path,
				Address:    entry.Address,
				Detail:     fmt.Sprintf("%s has no %q attribute for the value the legacy shape implied", dst.Type, name),
			})
			continue
		}
		set(name, literalExpr(cty.StringVal(f.fixed[name])), "the legacy shape")
	}

	for _, legacyName := range sortedKeys(legacy) {
		value := legacy[legacyName]
		if contains(f.consumed, legacyName) {
			continue
		}
		if f.idAttr != "" && legacyName == f.idAttr && !required(dst.Schema.Attributes[legacyName]) {
			continue
		}
		if b.flattenNested(dst, f, legacyName, value, entry, set) {
			continue
		}
		target := legacyName
		if renamed, ok := f.rename[legacyName]; ok {
			target = renamed
		}
		attribute, exists := dst.Schema.Attributes[target]
		if !exists {
			if populated(value) {
				b.blocker(Blocker{
					Kind:       BlockerUnsupportedField,
					LegacyPath: f.path + "." + legacyName,
					Address:    entry.Address,
					Detail:     fmt.Sprintf("%s has no %q attribute, so the configured value cannot be carried over", dst.Type, target),
				})
			}
			continue
		}
		if skippable(attribute) {
			continue
		}
		ref := attrRef{top: target, name: target, path: target}
		switch {
		case contains(f.forceSecrets, target):
			set(target, b.secret(attribute, ref, value, entry, f.path+"."+legacyName, true), "the conditionally required secret")
		case contains(f.serviceRefs, legacyName):
			set(target, b.serviceRef(dst, target, value, entry, f.path+"."+legacyName), "the legacy delivery service")
		case contains(f.jwtRefs, target):
			set(target, b.jwtRef(value, entry, target, f.path+"."+legacyName), "the legacy JWT template name")
		case target == appRolePermissionIDs:
			set(target, b.permissionIDs(value, scope.AppID, entry, f.path), "the legacy permission names")
		case contains(f.payload, target):
			set(target, b.payload(entry, payloadKey, target, value, f.path), "the legacy payload")
		default:
			set(target, b.attribute(attribute, ref, value, entry, f.path+"."+legacyName), "the legacy state")
		}
	}

	if dst.Type == listResourceType {
		b.translateListData(dst, legacy, entry, payloadKey, f.path, set)
	}
	b.requireAttributes(dst, attrs, entry, f.path)
	return attrs
}

// appRolePermissionIDs is the destination attribute holding app permission ids, which legacy stored as names.
const appRolePermissionIDs = "permission_ids"

// listResourceType is the destination whose single typed data attribute became three typed attributes.
const listResourceType = "descope_list"

// flattenNested carries a legacy nested object into the flat destination attributes the split replaced it with.
func (b *builder) flattenNested(dst *Destination, f family, legacyName string, value any, entry *Entry, set func(string, *expr, string)) bool {
	prefix := legacyName + "."
	handled := false
	for _, source := range sortedKeys(f.flatten) {
		if !strings.HasPrefix(source, prefix) {
			continue
		}
		handled = true
		target := f.flatten[source]
		attribute, exists := dst.Schema.Attributes[target]
		if !exists {
			b.blocker(Blocker{
				Kind:       BlockerUnsupportedField,
				LegacyPath: f.path + "." + source,
				Address:    entry.Address,
				Detail:     fmt.Sprintf("%s has no %q attribute", dst.Type, target),
			})
			continue
		}
		nested, found := nestedValue(value, strings.Split(source, ".")[1:])
		if !found || nested == nil || skippable(attribute) {
			continue
		}
		set(target, b.attribute(attribute, attrRef{top: target, name: target, path: target}, nested, entry, f.path+"."+source), "the legacy "+source)
	}
	if !handled {
		return false
	}
	b.reportUnflattened(dst, f, legacyName, value, entry)
	return true
}

// reportUnflattened fails closed on a leaf of a flattened legacy object that the flattening does not cover.
func (b *builder) reportUnflattened(dst *Destination, f family, legacyName string, value any, entry *Entry) {
	for _, item := range leafPaths(legacyName, value) {
		if _, covered := f.flatten[item.path]; covered || !populated(item.value) {
			continue
		}
		b.blocker(Blocker{
			Kind:       BlockerUnsupportedField,
			LegacyPath: f.path + "." + item.path,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("%s has no destination for the legacy nested value %q", dst.Type, item.path),
		})
	}
}

type leaf struct {
	path  string
	value any
}

func leafPaths(prefix string, value any) []leaf {
	object, ok := asObject(value)
	if !ok {
		return []leaf{{path: prefix, value: value}}
	}
	leaves := []leaf{}
	for _, key := range sortedKeys(object) {
		leaves = append(leaves, leafPaths(prefix+"."+key, object[key])...)
	}
	return leaves
}

func nestedValue(value any, segments []string) (any, bool) {
	current := value
	for _, segment := range segments {
		object, ok := asObject(current)
		if !ok {
			return nil, false
		}
		current, ok = object[segment]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

// translateListData routes the legacy list's typed JSON document onto the typed attribute that replaced it.
func (b *builder) translateListData(dst *Destination, legacy map[string]any, entry *Entry, payloadKey, legacyPath string, set func(string, *expr, string)) {
	listType, _ := legacy["type"].(string)
	document, _ := legacy["data"].(string)
	switch listType {
	case "":
		b.blocker(Blocker{
			Kind:       BlockerInvalidValue,
			LegacyPath: legacyPath,
			Address:    entry.Address,
			Detail:     "the legacy list has no type, so its data cannot be routed to texts, ips or json",
		})
	case "texts", "ips":
		values := []any{}
		if document != "" {
			if err := json.Unmarshal([]byte(document), &values); err != nil {
				b.blocker(Blocker{
					Kind:       BlockerInvalidValue,
					LegacyPath: legacyPath + ".data",
					Address:    entry.Address,
					Detail:     fmt.Sprintf("the legacy %s list data is not a JSON array: %s", listType, err.Error()),
				})
				return
			}
		}
		if len(values) == 0 {
			b.blocker(Blocker{
				Kind:       BlockerInvalidValue,
				LegacyPath: legacyPath + ".data",
				Address:    entry.Address,
				Detail:     fmt.Sprintf("the legacy %s list is empty, but the standalone resource requires a non-empty %s set", listType, listType),
			})
			return
		}
		attribute, exists := dst.Schema.Attributes[listType]
		if !exists {
			b.blocker(Blocker{
				Kind:       BlockerUnsupportedField,
				LegacyPath: legacyPath,
				Address:    entry.Address,
				Detail:     fmt.Sprintf("%s has no %q attribute", dst.Type, listType),
			})
			return
		}
		set(listType, b.attribute(attribute, attrRef{top: listType, name: listType, path: listType}, values, entry, legacyPath+".data"), "the legacy list data")
	case "json":
		if document == "" {
			return
		}
		set("json", b.payload(entry, payloadKey, "json", document, legacyPath), "the legacy list data")
	default:
		b.blocker(Blocker{
			Kind:       BlockerInvalidValue,
			LegacyPath: legacyPath + ".type",
			Address:    entry.Address,
			Detail:     fmt.Sprintf("unknown legacy list type %q", listType),
		})
	}
}

// requireAttributes fails closed when the destination needs a value the legacy state never held.
func (b *builder) requireAttributes(dst *Destination, attrs map[string]*expr, entry *Entry, legacyPath string) {
	for _, name := range sortedKeys(dst.Schema.Attributes) {
		attribute := dst.Schema.Attributes[name]
		if !attribute.IsRequired() {
			continue
		}
		if _, present := attrs[name]; present {
			continue
		}
		kind := BlockerUnsupportedField
		if attribute.IsSensitive() {
			kind = BlockerUnresolvedSecret
		}
		b.blocker(Blocker{
			Kind:       kind,
			LegacyPath: legacyPath,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("%s requires attribute %q and the legacy state holds no value for it", dst.Type, name),
		})
	}
}

// attribute converts one legacy value onto one destination attribute, replacing secrets with variable references so
// that no secret is ever written to a generated file.
func (b *builder) attribute(attribute schema.Attribute, ref attrRef, value any, entry *Entry, legacyPath string) *expr {
	if value == nil {
		return &expr{}
	}
	if attribute.IsSensitive() {
		return b.secret(attribute, ref, value, entry, legacyPath, false)
	}
	if nested, ok := nestedAttributes(attribute); ok {
		return b.nested(nested, ref, value, entry, legacyPath)
	}
	ty, err := ctyTypeFromTF(attribute.GetType().TerraformType(b.ctx))
	if err != nil {
		b.blocker(Blocker{Kind: BlockerInvalidValue, LegacyPath: legacyPath, Address: entry.Address, Detail: err.Error()})
		return &expr{}
	}
	converted, err := valueToCty(value, ty, legacyPath)
	if err != nil {
		b.blocker(Blocker{Kind: BlockerInvalidValue, LegacyPath: legacyPath, Address: entry.Address, Detail: err.Error()})
		return &expr{}
	}
	return literalExpr(converted)
}

// nested walks a destination nested attribute so that a secret buried inside an object still becomes a variable.
func (b *builder) nested(nested nestedSchema, ref attrRef, value any, entry *Entry, legacyPath string) *expr {
	switch nested.kind {
	case nestedSingle:
		object, ok := asObject(value)
		if !ok {
			b.invalidShape(entry, legacyPath, "an object", value)
			return &expr{}
		}
		return b.nestedObject(nested.attributes, ref, object, entry, legacyPath)
	case nestedCollection:
		items, ok := value.([]any)
		if !ok {
			b.invalidShape(entry, legacyPath, "a list", value)
			return &expr{}
		}
		elements := make([]*expr, 0, len(items))
		for i, item := range items {
			object, isObject := asObject(item)
			if !isObject {
				b.invalidShape(entry, fmt.Sprintf("%s[%d]", legacyPath, i), "an object", item)
				continue
			}
			elements = append(elements, b.nestedObject(nested.attributes, ref.index(i), object, entry, fmt.Sprintf("%s[%d]", legacyPath, i)))
		}
		return listExpr(elements)
	case nestedMap:
		entries, ok := asObject(value)
		if !ok {
			b.invalidShape(entry, legacyPath, "an object", value)
			return &expr{}
		}
		attrs := map[string]*expr{}
		for _, key := range sortedKeys(entries) {
			object, isObject := asObject(entries[key])
			if !isObject {
				continue
			}
			attrs[key] = b.nestedObject(nested.attributes, ref.child(key), object, entry, legacyPath+"."+key)
		}
		return objectExpr(attrs)
	}
	return &expr{}
}

func (b *builder) nestedObject(attributes map[string]schema.Attribute, ref attrRef, object map[string]any, entry *Entry, legacyPath string) *expr {
	attrs := map[string]*expr{}
	for _, key := range sortedKeys(object) {
		attribute, exists := attributes[key]
		if !exists {
			if populated(object[key]) {
				b.blocker(Blocker{
					Kind:       BlockerUnsupportedField,
					LegacyPath: legacyPath + "." + key,
					Address:    entry.Address,
					Detail:     fmt.Sprintf("the destination object has no %q attribute", key),
				})
			}
			continue
		}
		if skippable(attribute) {
			continue
		}
		value := b.attribute(attribute, ref.child(key), object[key], entry, legacyPath+"."+key)
		if value.empty() {
			continue
		}
		attrs[key] = value
	}
	return objectExpr(attrs)
}

func (b *builder) invalidShape(entry *Entry, legacyPath, expected string, value any) {
	b.blocker(Blocker{
		Kind:       BlockerInvalidValue,
		LegacyPath: legacyPath,
		Address:    entry.Address,
		Detail:     fmt.Sprintf("%s: expected %s, found %s", legacyPath, expected, jsonKind(value)),
	})
}

// secret records a sensitive value as an input variable, never as a generated value. Backend reads omit write-only
// secrets, so every secret present in legacy state is rehydrated explicitly and listed as the only attribute an
// imported resource may update. A missing optional secret stays omitted; a missing required secret is a blocker.
func (b *builder) secret(attribute schema.Attribute, ref attrRef, value any, entry *Entry, legacyPath string, force bool) *expr {
	present := populated(value)
	if !present && !attribute.IsRequired() && !force {
		return &expr{}
	}
	typeExpression, err := b.secretTypeExpression(attribute)
	if err != nil {
		b.blocker(Blocker{Kind: BlockerInvalidValue, LegacyPath: legacyPath, Address: entry.Address, Detail: err.Error()})
		return &expr{}
	}
	variable := variableName(entry.Address, ref.name)
	owner := entry.Address + "." + ref.path
	if previous, collision := b.variableOwners[variable]; collision && previous != owner {
		b.blocker(Blocker{
			Kind:       BlockerDuplicateAddress,
			LegacyPath: legacyPath,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("sensitive input variable %q is also used by %s", variable, previous),
		})
		return &expr{}
	}
	b.variableOwners[variable] = owner
	b.manifest.Secrets = append(b.manifest.Secrets, SecretRequirement{
		Address:    entry.Address,
		Attribute:  ref.path,
		Variable:   variable,
		Type:       typeExpression,
		LegacyPath: legacyPath,
		Present:    present,
	})
	if !contains(entry.SecretAttributes, ref.path) {
		entry.SecretAttributes = append(entry.SecretAttributes, ref.path)
		sort.Strings(entry.SecretAttributes)
	}
	if !present {
		b.blocker(Blocker{
			Kind:       BlockerUnresolvedSecret,
			LegacyPath: legacyPath,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("%s requires the sensitive attribute %q and the legacy state holds no value for it, so it has to come from the operator's own secret store before adoption", entry.Type, ref.path),
		})
	}
	return rawExpr("var." + variable)
}

func (b *builder) secretTypeExpression(attribute schema.Attribute) (string, error) {
	typeValue, err := ctyTypeFromTF(attribute.GetType().TerraformType(b.ctx))
	if err != nil {
		return "", err
	}
	return hclTypeExpression(typeValue)
}

func hclTypeExpression(value cty.Type) (string, error) {
	switch {
	case value.Equals(cty.String):
		return "string", nil
	case value.Equals(cty.Bool):
		return "bool", nil
	case value.Equals(cty.Number):
		return "number", nil
	case value.IsListType():
		element, err := hclTypeExpression(value.ElementType())
		return "list(" + element + ")", err
	case value.IsSetType():
		element, err := hclTypeExpression(value.ElementType())
		return "set(" + element + ")", err
	case value.IsMapType():
		element, err := hclTypeExpression(value.ElementType())
		return "map(" + element + ")", err
	default:
		return "", fmt.Errorf("sensitive attribute uses unsupported type %s", value.FriendlyName())
	}
}

// serviceRef translates a legacy delivery service object. The legacy connector name becomes the connector id the
// standalone settings resource expects, and whatever else the destination service object still holds - the nested
// message templates of the invitation settings, for instance - is translated in place.
func (b *builder) serviceRef(dst *Destination, target string, value any, entry *Entry, legacyPath string) *expr {
	object, ok := asObject(value)
	if !ok {
		return &expr{}
	}
	nested, hasNested := nestedAttributes(dst.Schema.Attributes[target])
	if !hasNested {
		return &expr{}
	}
	remaining := map[string]any{}
	for key, nestedValue := range object {
		if key == "connector" {
			continue
		}
		if key == "templates" && !serviceKeepsTemplates(nested) {
			continue // the templates became standalone method scoped resources
		}
		remaining[key] = nestedValue
	}
	service := b.nestedObject(nested.attributes, attrRef{top: target, name: target, path: target}, remaining, entry, legacyPath)
	if _, ok := nested.attributes[serviceConnectorID]; !ok {
		return service
	}
	name, _ := object["connector"].(string)
	switch name {
	case "":
		return service
	case helpers.DescopeConnector:
		b.unowned(legacyPath, fmt.Sprintf("the %q delivery service is Descope's built-in provider and has no connector resource", helpers.DescopeConnector))
		service.obj[serviceConnectorID] = literalExpr(cty.StringVal(helpers.DescopeConnector))
		return service
	}
	id, resolved := b.connectorIDs[name]
	if !resolved {
		b.blocker(Blocker{
			Kind:       BlockerUnresolvedRef,
			LegacyPath: legacyPath + ".connector",
			Address:    entry.Address,
			Detail:     fmt.Sprintf("the delivery service references connector %q by name and the standalone resource needs its id, but no connector owned by this descope_project state has that name", name),
		})
		return &expr{}
	}
	entry.References[target+"."+serviceConnectorID] = id
	service.obj[serviceConnectorID] = literalExpr(cty.StringVal(id))
	return service
}

// serviceKeepsTemplates reports whether a destination delivery service object still owns its message templates.
func serviceKeepsTemplates(nested nestedSchema) bool {
	_, ok := nested.attributes["templates"]
	return ok
}

// serviceConnectorID is the delivery service attribute that replaced the legacy connector name.
const serviceConnectorID = "connector_id"

// jwtRef resolves a legacy JWT template name onto the template id the standalone settings resource expects.
func (b *builder) jwtRef(value any, entry *Entry, target, legacyPath string) *expr {
	name, _ := value.(string)
	if name == "" {
		return &expr{}
	}
	id, resolved := b.jwtIDs[name]
	if !resolved {
		b.blocker(Blocker{
			Kind:       BlockerUnresolvedRef,
			LegacyPath: legacyPath,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("%s references JWT template %q by name and the standalone resource needs its id, but no owned JWT template has that name", target, name),
		})
		return &expr{}
	}
	entry.References[target] = id
	return literalExpr(cty.StringVal(id))
}

// permissionIDs resolves the legacy application role permission names onto the app permission ids.
func (b *builder) permissionIDs(value any, appID string, entry *Entry, legacyPath string) *expr {
	names, ok := value.([]any)
	if !ok || len(names) == 0 {
		return &expr{}
	}
	known := b.appPermIDs[appID]
	ids := make([]cty.Value, 0, len(names))
	for _, raw := range names {
		name, isString := raw.(string)
		if !isString || name == "" {
			continue
		}
		id, resolved := known[name]
		if !resolved {
			b.blocker(Blocker{
				Kind:       BlockerUnresolvedRef,
				LegacyPath: legacyPath + ".permissions",
				Address:    entry.Address,
				Detail:     fmt.Sprintf("the application role references permission %q, which is not one of the application's owned permissions, so its id cannot be resolved", name),
			})
			continue
		}
		entry.References[appRolePermissionIDs+"."+name] = id
		ids = append(ids, cty.StringVal(id))
	}
	if len(ids) == 0 {
		return &expr{}
	}
	return literalExpr(cty.SetVal(ids))
}

// payload writes a JSON document to its own file and references it, so reviewers diff the document rather than one
// very long HCL string.
func (b *builder) payload(entry *Entry, payloadKey, attribute string, value any, legacyPath string) *expr {
	document, ok := value.(string)
	if !ok || strings.TrimSpace(document) == "" {
		b.blocker(Blocker{
			Kind:       BlockerMissingPayload,
			LegacyPath: legacyPath + "." + attribute,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("%s needs the %q document and the legacy state holds none", entry.Type, attribute),
		})
		return &expr{}
	}
	pretty := &bytes.Buffer{}
	if err := json.Indent(pretty, []byte(document), "", "  "); err != nil {
		b.blocker(Blocker{
			Kind:       BlockerInvalidValue,
			LegacyPath: legacyPath + "." + attribute,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("the %q document is not valid JSON: %s", attribute, err.Error()),
		})
		return &expr{}
	}
	pretty.WriteString("\n")
	file := path.Join(payloadDir, entry.Type+"."+payloadKey+"."+attribute+".json")
	owner := entry.Address + "." + attribute
	if previous, collision := b.payloadOwners[file]; collision && previous != owner {
		b.blocker(Blocker{
			Kind:       BlockerDuplicateAddress,
			LegacyPath: legacyPath + "." + attribute,
			Address:    entry.Address,
			Detail:     fmt.Sprintf("payload file %q is also used by %s", file, previous),
		})
		return &expr{}
	}
	b.payloadOwners[file] = owner
	b.payloads[file] = pretty.Bytes()
	if !contains(entry.Payloads, file) {
		entry.Payloads = append(entry.Payloads, file)
		sort.Strings(entry.Payloads)
	}
	return rawExpr(`file("${path.module}/` + file + `")`)
}

// label derives a stable, unique HCL identifier for a destination resource.
func (b *builder) label(resourceType string, f family, name, backendID string) string {
	candidates := []string{name, f.labelPrefix, strings.TrimPrefix(resourceType, ProviderTypeName+"_"), backendID}
	candidate := ""
	for _, option := range candidates {
		if candidate = sanitizeLabel(option); candidate != "" {
			break
		}
	}
	if candidate == "" {
		candidate = "entity"
	}
	label := candidate
	for suffix := 2; b.labelTaken(resourceType, label); suffix++ {
		label = candidate + "_" + strconv.Itoa(suffix)
	}
	b.claimLabel(resourceType, label)
	return label
}

func (b *builder) labelTaken(resourceType, label string) bool {
	return b.labels[resourceType][label]
}

func (b *builder) claimLabel(resourceType, label string) {
	if b.labels[resourceType] == nil {
		b.labels[resourceType] = map[string]bool{}
	}
	b.labels[resourceType][label] = true
}

// qualify prefixes a destination address with the module the legacy resource was declared in, because that is where
// the generated configuration goes and therefore how Terraform addresses it.
func (b *builder) qualify(address string) string {
	if b.legacy.ModuleAddress == "" {
		return address
	}
	return b.legacy.ModuleAddress + "." + address
}

func (b *builder) idSource() string {
	if b.manifest.RemoteBackfill {
		return "legacy state, with GET-only backfill enabled"
	}
	return "legacy state"
}

func (b *builder) describe(f family, name string) string {
	if name != "" {
		return name
	}
	return f.path
}

// sanitizeLabel turns a human name into a valid HCL identifier.
func sanitizeLabel(name string) string {
	var out strings.Builder
	pendingSeparator := false
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if pendingSeparator {
				out.WriteByte('_')
				pendingSeparator = false
			}
			out.WriteRune(r)
			continue
		}
		if out.Len() > 0 {
			pendingSeparator = true
		}
	}
	label := out.String()
	if label == "" {
		return ""
	}
	if unicode.IsDigit([]rune(label)[0]) {
		label = "n" + label
	}
	return label
}

// sanitizePathSegment keeps a collection key usable as part of a generated file name.
func sanitizePathSegment(key string) string {
	label := sanitizeLabel(key)
	if label == "" {
		label = "entry"
	}
	return label + "_" + stableSuffix(key)
}

// variableName includes the full resource address and a stable suffix, so equal labels in different resource types
// or modules can never share a credential variable.
func variableName(address, attribute string) string {
	value := address + "\x00" + attribute
	return sanitizeLabel(address+"_"+attribute) + "_" + stableSuffix(value)
}

func stableSuffix(value string) string {
	digest := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", digest[:6])
}

// required reports whether a destination attribute must appear in the configuration.
func required(attribute schema.Attribute) bool {
	return attribute != nil && attribute.IsRequired()
}

// skippable reports whether an attribute is computed only, and therefore must not appear in configuration.
func skippable(attribute schema.Attribute) bool {
	return attribute.IsComputed() && !attribute.IsOptional() && !attribute.IsRequired()
}

type nestedKind int

const (
	nestedSingle nestedKind = iota
	nestedCollection
	nestedMap
)

type nestedSchema struct {
	kind       nestedKind
	attributes map[string]schema.Attribute
}

// nestedAttributes reports the nested schema of a nested attribute, so translation can recurse into it.
func nestedAttributes(attribute schema.Attribute) (nestedSchema, bool) {
	switch typed := attribute.(type) {
	case schema.SingleNestedAttribute:
		return nestedSchema{kind: nestedSingle, attributes: typed.Attributes}, true
	case schema.ListNestedAttribute:
		return nestedSchema{kind: nestedCollection, attributes: typed.NestedObject.Attributes}, true
	case schema.SetNestedAttribute:
		return nestedSchema{kind: nestedCollection, attributes: typed.NestedObject.Attributes}, true
	case schema.MapNestedAttribute:
		return nestedSchema{kind: nestedMap, attributes: typed.NestedObject.Attributes}, true
	}
	return nestedSchema{}, false
}

// existing finds an already standalone resource with the same destination type, backend id and scope.
func (b *builder) existing(resourceType, backendID string, scope Scope) (*ExistingResource, []string, []string) {
	if backendID == "" && scope.Kind == ScopeSingleton {
		backendID = scope.ProjectID
	}
	matches := []*ExistingResource{}
	addresses := []string{}
	unsupportedProviders := []string{}
	for i := range b.legacy.Existing {
		candidate := &b.legacy.Existing[i]
		if candidate.Type != resourceType {
			continue
		}
		id, _ := candidate.Attributes["id"].(string)
		if id == "" && scope.Kind == ScopeSingleton {
			id, _ = candidate.Attributes["project_id"].(string)
		}
		if id != backendID {
			continue
		}
		if scope.Kind != ScopeCompany {
			projectID, _ := candidate.Attributes["project_id"].(string)
			if projectID != scope.ProjectID {
				continue
			}
		}
		if scope.Kind == ScopeApp {
			appID, _ := candidate.Attributes["app_id"].(string)
			if appID != scope.AppID {
				continue
			}
		}
		if scope.Kind == ScopeMethod {
			method, _ := candidate.Attributes["method"].(string)
			if method != scope.Method {
				continue
			}
		}
		if !isDescopeProvider(candidate.ProviderSource) {
			unsupportedProviders = append(unsupportedProviders, candidate.Address)
			continue
		}
		matches = append(matches, candidate)
		addresses = append(addresses, candidate.Address)
	}
	sort.Strings(addresses)
	sort.Strings(unsupportedProviders)
	if len(matches) == 1 {
		return matches[0], addresses, unsupportedProviders
	}
	return nil, addresses, unsupportedProviders
}
