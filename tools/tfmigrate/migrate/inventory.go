package migrate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/helpers"
)

// legacyOwnershipAttributes are the descope_project attributes that held configuration the split moved into
// standalone resources. Every one of them is claimed by the walker below.
var legacyOwnershipAttributes = []string{
	"admin_portal",
	"applications",
	"attributes",
	"authentication",
	"authorization",
	"connectors",
	"flows",
	"invite_settings",
	"jwt_templates",
	"lists",
	"project_settings",
	"styles",
	"widgets",
}

// builder walks one legacy descope_project instance and produces the ownership manifest plus the configuration
// blocks and payload files needed to adopt every owned object.
type builder struct {
	ctx       context.Context
	reg       *Registry
	legacy    *LegacyResource
	projectID string
	manifest  *Manifest

	labels         map[string]map[string]bool
	addresses      map[string]string
	payloads       map[string][]byte
	imports        map[string]string
	variableOwners map[string]string
	payloadOwners  map[string]string
	blocks         []*resourceBlock
	groupLabels    map[string]string

	// identity tables, filled by the families that own them before the families that reference them.
	connectorIDs    map[string]string
	existingClaimed map[string]string
	jwtIDs          map[string]string
	appPermIDs      map[string]map[string]string
	activeIDs       map[string]string
}

// resourceBlock is one generated standalone resource configuration. ForEachKey is set for a member of a keyed
// collection that is generated as a single for_each resource.
type resourceBlock struct {
	Type       string
	Label      string
	Attrs      map[string]*expr
	Entry      *Entry
	Legacy     string
	ForEachKey string
	KeyAttr    string
}

// Build produces the manifest and the generated configuration for one legacy resource instance.
func Build(ctx context.Context, reg *Registry, legacy *LegacyResource, opts Options) (*Manifest, []*resourceBlock, map[string][]byte, error) {
	if opts.ProviderSource != ProviderSource {
		return nil, nil, nil, fmt.Errorf("provider source must be %q", ProviderSource)
	}
	if !isDescopeProvider(legacy.ProviderSource) {
		return nil, nil, nil, fmt.Errorf("legacy resource provider must be %q", ProviderSource)
	}
	projectID, _ := legacy.Attributes["id"].(string)
	if projectID == "" {
		return nil, nil, nil, fmt.Errorf("resource %s has no project id in its state, the migration cannot address any destination without it", legacy.Address)
	}
	b := &builder{
		ctx:       ctx,
		reg:       reg,
		legacy:    legacy,
		projectID: projectID,
		manifest: &Manifest{
			Tool:              ToolName,
			FormatVersion:     ManifestFormatVersion,
			ProjectID:         projectID,
			LegacyAddress:     legacy.Address,
			LegacyLocalAddr:   legacy.LocalAddress,
			InputFormat:       legacy.SourceFormat,
			ModuleAddress:     legacy.ModuleAddress,
			SourceProvider:    ProviderPin{Source: opts.ProviderSource, Version: opts.SourceProviderVersion},
			TargetProvider:    ProviderPin{Source: opts.ProviderSource, Version: opts.TargetProviderVersion},
			RequiredTerraform: RequiredTerraformVersion,
			ProjectCore: ProjectCore{
				Name:               stringValue(legacy.Attributes["name"]),
				Environment:        stringValue(legacy.Attributes["environment"]),
				Tags:               stringValues(legacy.Attributes["tags"]),
				DeletionProtection: boolPointer(legacy.Attributes["deletion_protection"]),
			},
			InstanceKey:     legacy.IndexKey,
			RemoteBackfill:  opts.Remote,
			OwnershipPolicy: "Only objects represented in the selected legacy descope_project state are adopted. Already standalone state is left alone, built-in Descope objects are recorded as unowned, and remote-only console-managed objects are never generated.",
			Retained:        retainedAttributes,
		},
		labels:          map[string]map[string]bool{},
		addresses:       map[string]string{},
		payloads:        map[string][]byte{},
		imports:         map[string]string{},
		connectorIDs:    map[string]string{},
		variableOwners:  map[string]string{},
		payloadOwners:   map[string]string{},
		existingClaimed: map[string]string{},
		jwtIDs:          map[string]string{},
		appPermIDs:      map[string]map[string]string{},
		activeIDs:       map[string]string{},
		groupLabels:     map[string]string{},
	}
	for _, existing := range legacy.Existing {
		b.addresses[existing.Address] = "existing Terraform state"
		if existing.ModuleAddress == legacy.ModuleAddress && existing.Name != "" {
			b.claimLabel(existing.Type, existing.Name)
		}
	}
	if len(legacy.Peers) > 0 {
		b.blocker(Blocker{
			Kind:       BlockerAmbiguousOwnership,
			LegacyPath: legacy.Address,
			Detail: fmt.Sprintf("%s has sibling instances %s; Terraform removed blocks cannot name one count or for_each instance, so detaching this address would forget the siblings too: migrate all instances in one reviewed change or first separate their resource declarations",
				legacy.Address, strings.Join(legacy.Peers, ", ")),
		})
	}
	if legacy.ModuleIndexed {
		b.blocker(Blocker{
			Kind:       BlockerAmbiguousOwnership,
			LegacyPath: legacy.Address,
			Detail:     "the declaring module uses count or for_each; generated child-module configuration would apply to every module instance, including instances absent from this state snapshot: first move the selected project into an unindexed module or root resource declaration",
		})
	}
	if unsafeCountIndex(legacy.IndexKey) {
		b.blocker(Blocker{
			Kind:       BlockerAmbiguousOwnership,
			LegacyPath: legacy.Address,
			Detail:     "the selected count instance is not index 0, and recreating its resource declaration from state would also create lower indexes: preserve the original count expression and migrate every instance together",
		})
	}
	b.walk()
	backupPath := backupStatePath
	restorable := legacy.SourceFormat == formatRawState
	instructions := []string{"Stop writers and obtain explicit human approval for rollback."}
	if restorable {
		instructions = append(instructions,
			"Verify the workspace, lineage and serial in "+backupPath+" before any human runs `terraform state push`.",
			"A human may restore the state backup from the workspace it was pulled from; the utility and agent never run `terraform state push`.")
	} else {
		backupPath = backupShowPath
		instructions = append(instructions,
			"This input is Terraform show JSON and cannot be restored with `terraform state push`.",
			"Before detach, capture a separate raw backup with `(umask 077; set -o noclobber; terraform state pull > legacy.tfstate.json)` and verify its workspace and lineage.")
	}
	instructions = append(instructions,
		"Re-pin the descope provider to "+opts.SourceProviderVersion+" and run `terraform init -upgrade`.",
		"Restore or rotate any write-only secret that the adoption plan rehydrated; restoring Terraform state does not restore a remote credential.",
		"Remove the generated detach and adopt configuration from the module, then run `terraform plan` and confirm it reports no changes.")
	b.manifest.Rollback = Rollback{
		StateBackup:     backupPath,
		Restorable:      restorable,
		LegacyAddress:   legacy.Address,
		SourceProvider:  opts.SourceProviderVersion,
		BackendMutation: "no create or delete; adoption can rehydrate state-held write-only secrets with in-place updates explicitly enumerated in the manifest",
		Instructions:    instructions,
	}
	b.manifest.Ready = len(b.manifest.Blockers) == 0
	return b.manifest, b.blocks, b.payloads, nil
}

// Options are the inputs a migration run needs beyond the state snapshot itself.
type Options struct {
	ProviderSource        string
	SourceProviderVersion string
	TargetProviderVersion string
	Remote                bool
}

// walk claims every legacy attribute exactly once. Families that own an identity other families reference are
// visited first, so a reference always resolves against an already recorded id.
func (b *builder) walk() {
	b.walkConnectors()
	b.walkJWTTemplates()
	b.walkApplications()
	b.walkAuthentication()
	b.walkSimpleFamilies()
	b.walkAuthorization()
	b.walkAttributes()
	b.walkProjectSettings()
	b.checkUnclaimedAttributes()
}

// checkUnclaimedAttributes fails closed on a populated legacy attribute this version does not know about, which is
// what a state written by a newer legacy provider would look like.
func (b *builder) checkUnclaimedAttributes() {
	known := map[string]bool{}
	for _, name := range legacyOwnershipAttributes {
		known[name] = true
	}
	for _, retained := range retainedAttributes {
		known[retained.Attribute] = true
	}
	for _, name := range sortedKeys(b.legacy.Attributes) {
		if known[name] || !populated(b.legacy.Attributes[name]) {
			continue
		}
		b.blocker(Blocker{
			Kind:       BlockerUnknownAttribute,
			LegacyPath: name,
			Detail:     fmt.Sprintf("legacy attribute %q is populated but this version of the migration does not know where it belongs, upgrade the migration utility to a version that covers the provider that wrote this state", name),
		})
	}
}

func (b *builder) walkConnectors() {
	connectors, ok := b.object("connectors")
	if !ok {
		return
	}
	families := b.reg.ConnectorFamilies()
	for _, name := range sortedKeys(connectors) {
		resourceType, known := families[name]
		if !known {
			if populated(connectors[name]) {
				b.blocker(Blocker{
					Kind:       BlockerUnsupportedField,
					LegacyPath: "connectors." + name,
					Detail:     fmt.Sprintf("connector family %q has no standalone resource in the target provider", name),
				})
			}
			continue
		}
		f := connectorFamily(name, resourceType)
		for index, item := range b.list("connectors." + name) {
			itemFamily := indexedFamily(f, index)
			connectorName, _ := item[f.nameAttr].(string)
			if connectorName == helpers.DescopeConnector {
				b.unowned(itemFamily.path, fmt.Sprintf("a connector named %q shadows Descope's built-in delivery service and is not adopted", helpers.DescopeConnector))
				continue
			}
			entry := b.emitListItem(itemFamily, item)
			if entry == nil || connectorName == "" {
				continue
			}
			if previous, taken := b.connectorIDs[connectorName]; taken && previous != entry.BackendID {
				b.blocker(Blocker{
					Kind:       BlockerAmbiguousOwnership,
					LegacyPath: itemFamily.path,
					Detail:     fmt.Sprintf("two connectors are named %q, so a delivery service reference to that name cannot be resolved to one id", connectorName),
				})
				continue
			}
			b.connectorIDs[connectorName] = entry.BackendID
		}
	}
}

func (b *builder) walkJWTTemplates() {
	if _, ok := b.object("jwt_templates"); !ok {
		return
	}
	b.claimObjectKeys("jwt_templates", []string{"user_templates", "access_key_templates"})
	for _, f := range jwtTemplateKinds {
		for index, item := range b.list(f.path) {
			itemFamily := indexedFamily(f, index)
			entry := b.emitListItem(itemFamily, item)
			if entry == nil {
				continue
			}
			if name, ok := item[f.nameAttr].(string); ok && name != "" {
				if previous, duplicate := b.jwtIDs[name]; duplicate && previous != entry.BackendID {
					b.blocker(Blocker{Kind: BlockerAmbiguousOwnership, LegacyPath: itemFamily.path, Detail: fmt.Sprintf("JWT template name %q resolves to both %q and %q", name, previous, entry.BackendID)})
					continue
				}
				b.jwtIDs[name] = entry.BackendID
			}
		}
	}
}

func (b *builder) walkApplications() {
	if _, ok := b.object("applications"); !ok {
		return
	}
	b.claimObjectKeys("applications", []string{"oidc_applications", "saml_applications", "wsfed_applications"})
	for _, f := range applicationFamilies {
		for index, item := range b.list(f.path) {
			appFamily := indexedFamily(f, index)
			entry := b.emitListItem(appFamily, item)
			if entry == nil {
				continue
			}
			appID := entry.BackendID
			scope := Scope{Kind: ScopeApp, ProjectID: b.projectID, AppID: appID}
			b.appPermIDs[appID] = map[string]string{}
			for childIndex, permission := range asObjects(item["permissions"]) {
				child := indexedFamily(appPermissionFamily, childIndex)
				child.path = fmt.Sprintf("%s.permissions[%d]", appFamily.path, childIndex)
				permEntry := b.emitScoped(child, permission, scope)
				if permEntry == nil {
					continue
				}
				if name, ok := permission["name"].(string); ok && name != "" {
					if previous, duplicate := b.appPermIDs[appID][name]; duplicate && previous != permEntry.BackendID {
						b.blocker(Blocker{Kind: BlockerAmbiguousOwnership, LegacyPath: child.path, Detail: fmt.Sprintf("application permission name %q resolves to both %q and %q", name, previous, permEntry.BackendID)})
						continue
					}
					b.appPermIDs[appID][name] = permEntry.BackendID
				}
			}
			for childIndex, role := range asObjects(item["roles"]) {
				child := indexedFamily(appRoleFamily, childIndex)
				child.path = fmt.Sprintf("%s.roles[%d]", appFamily.path, childIndex)
				b.emitScoped(child, role, scope)
			}
		}
	}
}

func (b *builder) walkAuthentication() {
	if _, ok := b.object("authentication"); !ok {
		return
	}
	claimed := []string{"oauth"}
	for _, method := range authMethods {
		claimed = append(claimed, method.legacy)
	}
	b.claimObjectKeys("authentication", claimed)

	for _, method := range authMethods {
		settings, present := b.object("authentication." + method.legacy)
		if !present {
			continue
		}
		fixed := map[string]string{}
		for _, service := range method.services() {
			serviceObject, hasService := b.object("authentication." + method.legacy + "." + service.legacy)
			if !hasService {
				continue
			}
			b.emitServiceTemplates(method, service, "authentication."+method.legacy+"."+service.legacy, serviceObject)
			if activeID, found := b.activeIDs[method.method+"/"+service.legacy]; found {
				fixed[service.templateIDAttr] = activeID
			}
		}
		f := family{
			path:         "authentication." + method.legacy,
			resourceType: method.resourceType,
			shape:        shapeObject,
			fixed:        fixed,
			serviceRefs:  method.serviceRefs(),
		}
		b.emitSingleton(f, settings)
	}
	b.walkOAuth()
}

// services returns the delivery services this legacy auth method could configure.
func (m authMethod) services() []serviceKind {
	services := []serviceKind{}
	if m.email {
		services = append(services, emailService)
	}
	if m.text {
		services = append(services, textService)
	}
	if m.voice {
		services = append(services, voiceService)
	}
	return services
}

// serviceRefs names the legacy delivery service objects whose connector name must become a connector id.
func (m authMethod) serviceRefs() []string {
	refs := []string{}
	for _, service := range m.services() {
		refs = append(refs, service.legacy)
	}
	return refs
}

// emitServiceTemplates turns the message templates nested in a legacy delivery service into standalone method
// scoped template resources, and records which one the settings resource must select.
func (b *builder) emitServiceTemplates(method authMethod, service serviceKind, path string, serviceObject map[string]any) {
	scope := Scope{Kind: ScopeMethod, ProjectID: b.projectID, Method: method.method}
	active := []string{}
	for index, template := range asObjects(serviceObject["templates"]) {
		name, _ := template["name"].(string)
		itemPath := fmt.Sprintf("%s.templates[%d]", path, index)
		if name == helpers.DescopeTemplate {
			b.unowned(itemPath, fmt.Sprintf("the %q template is the built-in Descope template and is not managed by Terraform", helpers.DescopeTemplate))
			continue
		}
		f := family{
			path:         itemPath,
			resourceType: service.resourceType,
			shape:        shapeList,
			idAttr:       "id",
			nameAttr:     "name",
			consumed:     []string{"active"},
			labelPrefix:  method.method + "_" + strings.TrimSuffix(service.legacy, "_service") + "_template",
		}
		entry := b.emitScoped(f, template, scope)
		if entry == nil {
			continue
		}
		if activeFlag, _ := template["active"].(bool); activeFlag {
			active = append(active, entry.BackendID)
		}
	}
	switch len(active) {
	case 0:
	case 1:
		b.activeIDs[method.method+"/"+service.legacy] = active[0]
	default:
		sort.Strings(active)
		b.blocker(Blocker{
			Kind:       BlockerAmbiguousOwnership,
			LegacyPath: path + ".templates",
			Detail: fmt.Sprintf("%d templates are marked active, so the %s attribute of %s cannot be resolved: ids %s",
				len(active), service.templateIDAttr, method.resourceType, strings.Join(active, ", ")),
		})
	}
}

func (b *builder) walkOAuth() {
	oauth, ok := b.object("authentication.oauth")
	if !ok {
		return
	}
	system, _ := b.object("authentication.oauth.system")
	for _, name := range oauthSystemProviders {
		provider, present := asObject(system[name])
		if !present {
			continue
		}
		f := oauthProviderFamily
		f.path = "authentication.oauth.system." + name
		f.forceSecrets = oauthRequiredSecrets(name, provider, true)
		f.consumed = append(f.consumed, oauthSystemReservedFields...)
		for _, field := range oauthSystemReservedFields {
			if populated(provider[field]) {
				b.blocker(Blocker{
					Kind:       BlockerUnsupportedField,
					LegacyPath: f.path + "." + field,
					Detail:     fmt.Sprintf("%s is reserved on the %s system provider and cannot be carried into the standalone resource", field, name),
				})
			}
		}
		b.emitKeyed(f, name, provider)
	}
	for _, name := range sortedKeys(system) {
		if !contains(oauthSystemProviders, name) && populated(system[name]) {
			b.blocker(Blocker{
				Kind:       BlockerUnsupportedField,
				LegacyPath: "authentication.oauth.system." + name,
				Detail:     fmt.Sprintf("%q is not a known system OAuth provider", name),
			})
		}
	}
	custom, _ := b.object("authentication.oauth.custom")
	for _, name := range sortedKeys(custom) {
		provider, present := asObject(custom[name])
		if !present {
			continue
		}
		f := oauthProviderFamily
		f.path = "authentication.oauth.custom." + name
		f.forceSecrets = oauthRequiredSecrets(name, provider, false)
		b.emitKeyed(f, name, provider)
	}
	b.emitSingleton(oauthSettingsFamily, oauth)
}

var oauthSystemReservedFields = []string{
	"description",
	"logo",
	"issuer",
	"authorization_endpoint",
	"token_endpoint",
	"user_info_endpoint",
	"jwks_endpoint",
	"use_client_assertion",
	"claim_mapping",
}

func oauthRequiredSecrets(name string, provider map[string]any, system bool) []string {
	required := []string{}
	clientID, _ := provider["client_id"].(string)
	useAssertion, _ := provider["use_client_assertion"].(bool)
	if system {
		if clientID != "" && (name != "apple" || !populated(provider["apple_key_generator"])) {
			required = append(required, "client_secret")
		}
		if name == "apple" && stringValue(provider["native_client_id"]) != "" && !populated(provider["native_apple_key_generator"]) {
			required = append(required, "native_client_secret")
		}
		return required
	}
	if !useAssertion {
		required = append(required, "client_secret")
	}
	return required
}

func (b *builder) walkSimpleFamilies() {
	for _, f := range simpleFamilies {
		if strings.HasPrefix(f.path, "authorization.") || strings.HasPrefix(f.path, "attributes.") {
			continue // dedicated walkers also claim their parent object's other keys
		}
		switch f.shape {
		case shapeObject:
			object, ok := b.object(f.path)
			if !ok {
				continue
			}
			b.emitSingleton(f, object)
		case shapeList:
			for index, item := range b.list(f.path) {
				b.emitListItem(indexedFamily(f, index), item)
			}
		case shapeMap:
			entries, ok := b.object(f.path)
			if !ok {
				continue
			}
			for _, key := range sortedKeys(entries) {
				item, present := asObject(entries[key])
				if !present {
					continue
				}
				b.emitKeyed(keyedFamily(f, key), key, item)
			}
		}
	}
}

func (b *builder) walkAuthorization() {
	authorization, ok := b.object("authorization")
	if !ok {
		return
	}
	b.claimObjectKeys("authorization", []string{"roles", "permissions", "fga"})
	for _, f := range simpleFamilies {
		if f.path != "authorization.roles" && f.path != "authorization.permissions" {
			continue
		}
		for index, item := range b.list(f.path) {
			b.emitListItem(indexedFamily(f, index), item)
		}
	}
	if schemaText, _ := authorization["fga"].(string); strings.TrimSpace(schemaText) != "" {
		b.emitSingleton(fgaFamily, map[string]any{"fga": strings.TrimSpace(schemaText)})
	}
}

func (b *builder) walkAttributes() {
	if _, ok := b.object("attributes"); !ok {
		return
	}
	b.claimObjectKeys("attributes", []string{"tenant", "user", "access_key"})
	for _, f := range simpleFamilies {
		if !strings.HasPrefix(f.path, "attributes.") {
			continue
		}
		for index, item := range b.list(f.path) {
			b.emitListItem(indexedFamily(f, index), item)
		}
	}
}

func (b *builder) walkProjectSettings() {
	settings, ok := b.object("project_settings")
	if !ok {
		return
	}
	if migration, present := b.object("project_settings.session_migration"); present {
		b.emitSingleton(sessionMigrationFamily, migration)
	}
	claimed := map[string]bool{"session_migration": true}
	for _, resourceType := range projectSettingsSplit {
		dst, known := b.reg.Lookup(resourceType)
		if !known {
			b.blocker(Blocker{
				Kind:       BlockerUnsupportedField,
				LegacyPath: "project_settings",
				Detail:     fmt.Sprintf("the target provider does not register %s", resourceType),
			})
			continue
		}
		subset := map[string]any{}
		for name, value := range settings {
			if _, belongs := dst.Schema.Attributes[name]; belongs {
				subset[name] = value
				claimed[name] = true
			}
		}
		f := family{
			path:         "project_settings",
			resourceType: resourceType,
			shape:        shapeObject,
			jwtRefs:      jwtTemplateRefs(dst),
		}
		b.emitSingleton(f, subset)
	}
	for _, name := range sortedKeys(settings) {
		if claimed[name] || !populated(settings[name]) {
			continue
		}
		b.blocker(Blocker{
			Kind:       BlockerUnsupportedField,
			LegacyPath: "project_settings." + name,
			Detail:     fmt.Sprintf("legacy setting %q has no destination in %s", name, strings.Join(projectSettingsSplit, " or ")),
		})
	}
}

// jwtTemplateRefs names the destination attributes that hold a JWT template id, so a legacy template name is
// resolved to the id the standalone resource expects.
func jwtTemplateRefs(dst *Destination) []string {
	refs := []string{}
	for _, name := range []string{"user_jwt_template", "access_key_jwt_template"} {
		if _, ok := dst.Schema.Attributes[name]; ok {
			refs = append(refs, name)
		}
	}
	return refs
}

// claimObjectKeys reports any key of a legacy container object that this version does not claim.
func (b *builder) claimObjectKeys(path string, claimed []string) {
	object, ok := b.object(path)
	if !ok {
		return
	}
	for _, name := range sortedKeys(object) {
		if contains(claimed, name) || !populated(object[name]) {
			continue
		}
		b.blocker(Blocker{
			Kind:       BlockerUnsupportedField,
			LegacyPath: path + "." + name,
			Detail:     fmt.Sprintf("legacy attribute %q is populated but has no destination resource", path+"."+name),
		})
	}
}

// object reads a nested legacy object, reporting whether it holds anything at all.
func (b *builder) object(path string) (map[string]any, bool) {
	value := b.valueAt(path)
	object, ok := asObject(value)
	if !ok || len(object) == 0 {
		return nil, false
	}
	return object, true
}

// list reads a nested legacy list of objects.
func (b *builder) list(path string) []map[string]any {
	return asObjects(b.valueAt(path))
}

func (b *builder) valueAt(path string) any {
	var current any = b.legacy.Attributes
	for _, segment := range strings.Split(path, ".") {
		object, ok := asObject(current)
		if !ok {
			return nil
		}
		current = object[segment]
	}
	return current
}

func (b *builder) blocker(blocker Blocker) {
	b.manifest.Blockers = append(b.manifest.Blockers, blocker)
}

func (b *builder) unowned(path, detail string) {
	b.manifest.Unowned = append(b.manifest.Unowned, Unowned{LegacyPath: path, Detail: detail})
}

func asObject(value any) (map[string]any, bool) {
	object, ok := value.(map[string]any)
	return object, ok
}

func asObjects(value any) []map[string]any {
	items, ok := value.([]any)
	if !ok {
		return nil
	}
	objects := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if object, ok := asObject(item); ok {
			objects = append(objects, object)
		}
	}
	return objects
}

func indexedFamily(source family, index int) family {
	source.path = fmt.Sprintf("%s[%d]", source.path, index)
	return source
}

func keyedFamily(source family, key string) family {
	source.path += "[" + strconv.Quote(key) + "]"
	return source
}

// populated reports whether a legacy value carries configuration, as opposed to being absent or empty.
func populated(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case string:
		return typed != ""
	case bool:
		return typed
	case float64:
		return typed != 0
	case json.Number:
		value, err := typed.Float64()
		return err != nil || value != 0
	case []any:
		return len(typed) > 0
	case map[string]any:
		for _, nested := range typed {
			if populated(nested) {
				return true
			}
		}
		return false
	default:
		return true
	}
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}

func stringValue(value any) string {
	text, _ := value.(string)
	return text
}

func boolPointer(value any) *bool {
	flag, ok := value.(bool)
	if !ok {
		return nil
	}
	return &flag
}

func stringValues(value any) []string {
	items, _ := value.([]any)
	values := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok {
			values = append(values, text)
		}
	}
	sort.Strings(values)
	return values
}

func unsafeCountIndex(value any) bool {
	switch index := value.(type) {
	case float64:
		return index != 0
	case int:
		return index != 0
	case int64:
		return index != 0
	case json.Number:
		parsed, err := index.Int64()
		return err != nil || parsed != 0
	default:
		return false
	}
}
