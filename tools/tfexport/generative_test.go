package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestGenerativeCoverage(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	seed := int64(7)
	if value := os.Getenv("DESCOPE_TFEXPORT_FUZZ_SEED"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			t.Fatalf("invalid DESCOPE_TFEXPORT_FUZZ_SEED: %s", err)
		}
		seed = parsed
	}

	generator := newConfigGenerator(t, seed)
	const generations = 4
	steps := make([]harness.Step, 0, generations)
	for generation := 0; generation < generations; generation++ {
		steps = append(steps, harness.Step{Config: generator.generation(generation)})
	}

	var singletons []string
	for _, resourceType := range generator.types {
		if generator.loaded[resourceType].ExportSingleton() {
			singletons = append(singletons, resourceType)
		}
	}

	harness.RunScenario(t, harness.Scenario{
		Name:          fmt.Sprintf("generative-%d", seed),
		ExpectMissing: singletons,
		Steps:         steps,
	})
	generator.report(t)
}

var generativeSkippedResources = map[string]string{
	"descope_project":           "the container the generated resources live in, created by the seed preamble",
	"descope_flow":              "payload-based, covered by flow scenarios",
	"descope_widget":            "payload-based, covered by the widget scenario",
	"descope_styles":            "payload-based, covered by flow scenarios",
	"descope_app_role":          "app-scoped wiring, covered by the apps scenario",
	"descope_app_permission":    "app-scoped wiring, covered by the apps scenario",
	"descope_access_key":        "not exportable by design",
	"descope_management_key":    "company-level key material",
	"descope_descoper":          "company-level, not exportable",
	"descope_engine":            "export deferred, so it would report as a discovery gap",
	"descope_oauth_settings":    "only has feature toggles, covered by polarity scenarios",
	"descope_session_migration": "vendor-locked attribute coupling, covered by settings scenarios",
	"descope_fga_schema":        "DSL payload, covered by the portal-lists-fga scenario",
	"descope_wsfed_app":         "covered by the wsfed scenario; groups_mapping coupling is nontrivial",
}

var generativeSkippedAttributes = map[string]string{
	"descope_list.texts":                                  "the kind of list is pinned to texts in every generation",
	"descope_list.ips":                                    "exactly one kind may be set, and texts is the generated one",
	"descope_list.json":                                   "exactly one kind may be set, and texts is the generated one",
	"descope_email_template.use_plain_text_body":          "requires plain_text_body pairing, exercised by templates scenarios",
	"descope_oidc_app.client_id":                          "create-only pairing with client_secret",
	"descope_oidc_app.client_secret":                      "create-only pairing with client_id",
	"descope_saml_app.dynamic_configuration":              "exactly one of dynamic/manual, manual is generated",
	"descope_oauth_provider.apple_key_generator":          "apple-specific key material object",
	"descope_oauth_provider.native_apple_key_generator":   "apple-specific key material object",
	"descope_oauth_provider.native_client_id":             "native pairing",
	"descope_oauth_provider.native_client_secret":         "native pairing",
	"descope_sso_settings.sso_suite_settings":             "nested suite block with its own coupling",
	"descope_inbound_app.permissions_scopes":              "nested scope objects, exercised by the inbound-app scenario",
	"descope_inbound_app.connections_scopes":              "nested scope objects, exercised by the inbound-app scenario",
	"descope_inbound_app.session_settings":                "nested settings object, exercised by the inbound-app scenario",
	"descope_inbound_app.scope_claim_mapping":             "nested mapping objects, exercised by the inbound-app scenario",
	"descope_password_settings.expiration_weeks":          "requires expiration=true pairing, set explicitly",
	"descope_password_settings.lock_attempts":             "requires lock=true pairing, set explicitly",
	"descope_password_settings.reuse_amount":              "requires reuse=true pairing, set explicitly",
	"descope_session_settings.inactivity_time":            "requires enable_inactivity=true pairing",
	"descope_user_attribute.select_options":               "requires select types, handled by pairing",
	"descope_tenant_attribute.select_options":             "requires select types, handled by pairing",
	"descope_access_key_attribute.select_options":         "requires select types, handled by pairing",
	"descope_project_settings.custom_domain":              "requires a real CNAME-verified domain",
	"descope_project_settings.test_users_static_otp":      "coupled pair, handled by pairing",
	"descope_project_settings.test_users_verifier_regexp": "coupled pair, handled by pairing",
	"descope_oauth_provider.claim_mapping":                "claims must map to existing descope attributes",
	"descope_admin_portal.enabled":                        "requires widgets, handled by pairing",
}

var attributeValueOverrides = map[string]string{
	"method":                            `"magiclink"`,
	"descope_text_template.method":      `"otp"`,
	"descope_voice_template.method":     `"otp"`,
	"descope_jwt_template.type":         `"user"|"key"`,
	"descope_user_attribute.type":       `"string"`,
	"descope_tenant_attribute.type":     `"number"`,
	"descope_access_key_attribute.type": `"string"`,
	"auth_type":                         `"legacyJWTSecret"`,
	"default_signature_algorithm":       `""|"sha256"`,
	"subject_name_id_type":              `""|"email"`,
	"subject_name_id_format":            `""`,
	"default_audience":                  `""|"projectId"`,
	"trusted_apps_audience":             `""|"projectId"`,
	"app_type":                          `"oauth"|"apikey"`,
	"access_type":                       `""|"offline"`,
	"prompt":                            `["none"]|["consent"]`,
	"empty_claim_policy":                `"none"|"delete"`,
	"auth_schema":                       `"default"|"tenantOnly"`,
	"status":                            `"active"`,
	"environment":                       `""`,
	"invite_expiration":                 `"2 weeks"|"3 weeks"`,
	"descope_jwt_template.template":     `jsonencode({ cov = 1 })|jsonencode({ cov = 2 })`,
	"descope_role.permissions":          `[descope_permission.cov_permission.name]`,
	"expiration_weeks":                  `13|26`,
	"deletion_protection":               `false`,
	"refresh_token_expiration":          `"3 weeks"|"4 weeks"`,
	"issuer_type":                       `"legacy"|"inbound"`,
	"enforce_strength":                  `"weak"|"strong"`,
	"prompts":                           `["none"]|["login", "consent"]`,
	"client_type":                       `""|"confidential"`,
	"min_length":                        `10|12`,
	"test_users_static_otp":             `"123456"|"654321"`,
	"android_fingerprints":              `["00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF"]`,
	"issuer":                            `"https://cov.example.com/issuer"|"https://cov.example.com/issuer2"`,
	"refresh_token_response_method":     `"cookies"|"response_body"`,
	"session_token_response_method":     `"response_body"|"cookies"`,
}

var nestedValueOverrides = map[string]string{
	"widget_authorization":                   `{ view_permissions = [descope_permission.cov_permission.name] }|{ edit_permissions = [descope_permission.cov_permission.name] }`,
	"descope_tenant_attribute.authorization": `{ view_permissions = [descope_permission.cov_permission.name] }`,
	"mandatory_user_attributes":              `[{ id = "email" }]|[{ id = "email" }, { id = "givenName" }]`,
	"attributes_scopes":                      `[{ name = "profile", description = "cov scope" }]`,
	"attribute_mapping":                      `[{ name = "cov-attr", value = "user.email" }]|[{ name = "cov-attr", value = "user.name" }]`,
}

var namePatternValues = []struct {
	fragment string
	values   string
}{
	{"endpoint", `"https://cov.example.com/endpoint"|"https://cov.example.com/endpoint2"`},
	{"cookie_policy", `"strict"|"lax"`},
	{"response_method", `"onBody"`},
	{"redirect_url", `"https://cov.example.com/redirect"|"https://cov.example.com/redirect2"`},
	{"logo", `"https://cov.example.com/logo.png"`},
	{"_url", `"https://cov.example.com/page"|"https://cov.example.com/page2"`},
	{"_urls", `["https://cov.example.com/cb"]|["https://cov.example.com/cb", "https://cov.example.com/cb2"]`},
	{"_domains", `["cov.example.com"]|["cov.example.com", "cov2.example.com"]`},
	{"domain", `"cov.example.com"|"cov2.example.com"`},
	{"email", `"cov@example.com"|"cov2@example.com"`},
	{"expiration", `"7 minutes"|"9 minutes"`},
	{"_duration", `"7 minutes"|"9 minutes"`},
	{"_time", `"7 minutes"|"9 minutes"`},
}

type configGenerator struct {
	random   *rand.Rand
	types    []string
	loaded   map[string]resources.ExportableResource
	covered  map[string]bool // "type.attribute" pairs exercised with generated values
	settable map[string]bool // all generable "type.attribute" pairs
}

func newConfigGenerator(t *testing.T, seed int64) *configGenerator {
	t.Helper()
	loaded := registry.Load(context.Background())
	generator := &configGenerator{
		random:   rand.New(rand.NewSource(seed)),
		loaded:   loaded,
		covered:  map[string]bool{},
		settable: map[string]bool{},
	}
	for resourceType, exportable := range loaded {
		if _, skipped := generativeSkippedResources[resourceType]; skipped {
			continue
		}
		if strings.HasSuffix(resourceType, "_connector") {
			continue // the connector sweep covers every connector type
		}
		generator.types = append(generator.types, resourceType)
		for name, attribute := range exportable.ExportSchema().Attributes {
			if generator.generable(resourceType, name, attribute) {
				generator.settable[resourceType+"."+name] = true
			}
		}
	}
	sort.Strings(generator.types)
	return generator
}

func (g *configGenerator) generable(resourceType, name string, attribute schema.Attribute) bool {
	if name == "project_id" {
		return false
	}
	if attribute.IsComputed() && !attribute.IsOptional() && !attribute.IsRequired() {
		return false
	}
	if _, skipped := generativeSkippedAttributes[resourceType+"."+name]; skipped {
		return false
	}
	if strings.HasSuffix(name, "_template_id") || strings.HasSuffix(name, "_connector_id") || name == "access_key_jwt_template" {
		return false // references to other entities, covered by the scenario catalog
	}
	return true
}

func (g *configGenerator) generation(generation int) string {
	var blocks []string
	for _, resourceType := range g.types {
		blocks = append(blocks, g.resourceBlock(resourceType, generation))
	}
	return strings.Join(blocks, "\n")
}

func (g *configGenerator) resourceBlock(resourceType string, generation int) string {
	exportable := g.loaded[resourceType]
	sc := exportable.ExportSchema()

	names := make([]string, 0, len(sc.Attributes))
	for name := range sc.Attributes {
		names = append(names, name)
	}
	sort.Strings(names)

	pairedLines, paired := g.pairings(resourceType, generation)

	var lines []string
	lines = append(lines, "  project_id = descope_project.test.id")
	for _, name := range names {
		attribute := sc.Attributes[name]
		if !g.generable(resourceType, name, attribute) || paired[name] {
			continue
		}
		include := attribute.IsRequired() || name == "name" || name == "deletion_protection"
		if !include {
			key := resourceType + "." + name
			if !g.covered[key] {
				include = true
			} else {
				include = g.random.Intn(3) == 0
			}
		}
		if !include {
			continue
		}
		value, ok := g.value(resourceType, name, attribute, generation)
		if !ok {
			continue
		}
		g.covered[resourceType+"."+name] = true
		lines = append(lines, fmt.Sprintf("  %s = %s", name, value))
	}
	lines = append(lines, pairedLines...)

	label := strings.TrimPrefix(resourceType, "descope_")
	return fmt.Sprintf("resource %q %q {\n%s\n}\n", resourceType, "cov_"+label, strings.Join(lines, "\n"))
}

func (g *configGenerator) pairings(resourceType string, generation int) ([]string, map[string]bool) {
	paired := map[string]bool{}
	var lines []string
	emit := func(line string, names ...string) {
		lines = append(lines, line)
		for _, name := range names {
			g.covered[resourceType+"."+name] = true
			paired[name] = true
		}
	}
	switch resourceType {
	case "descope_inbound_app":
		// a public client has no secret and the provider rejects that pairing, and DPoP needs a client type, so it is pinned rather than rotated
		emit(`  client_type = "confidential"`, "client_type")
	case "descope_email_template":
		// required unless use_plain_text_body is enabled, so pinned in every generation
		emit(fmt.Sprintf(`  html_body = "cov body %d"`, generation), "html_body")
	case "descope_list":
		// exactly one kind may be set, so the generator always writes the texts one
		emit(fmt.Sprintf(`  texts = ["cov-%d-a", "cov-%d-b"]`, generation, generation), "texts")
	case "descope_password_settings":
		emit(fmt.Sprintf("  expiration_weeks = %d", 13+generation), "expiration_weeks")
		emit(fmt.Sprintf("  lock_attempts = %d", 5+generation), "lock_attempts")
		emit(fmt.Sprintf("  reuse_amount = %d", 4+generation), "reuse_amount")
	case "descope_session_settings":
		// session_token_expiration is pinned in every generation: sparse updates trip a provider inconsistency
		emit(fmt.Sprintf(`  inactivity_time = "%d minutes"`, 12+generation), "inactivity_time")
		emit(fmt.Sprintf(`  session_token_expiration = "%d minutes"`, 8+generation), "session_token_expiration")
	case "descope_oauth_provider":
		// custom providers require these on every write, not only when uncovered
		emit(`  allowed_grant_types = ["authorization_code"]`, "allowed_grant_types")
		emit(`  user_info_endpoint  = "https://cov.example.com/userinfo"`, "user_info_endpoint")
		emit(fmt.Sprintf(`  client_id = "cov-client-%d"`, generation), "client_id")
		emit(`  client_secret = "not-a-real-secret"`, "client_secret")
		emit(`  authorization_endpoint = "https://cov.example.com/authorize"`, "authorization_endpoint")
		emit(`  token_endpoint = "https://cov.example.com/token"`, "token_endpoint")
	case "descope_project_settings":
		emit(fmt.Sprintf(`  test_users_static_otp = "12345%d"`, generation), "test_users_static_otp")
		emit(`  test_users_verifier_regexp = "^cov-.*$"`, "test_users_verifier_regexp")
	case "descope_admin_portal":
		emit("  enabled = true", "enabled")
		emit(`  widgets = [{ widget_id = "user-management", type = "user-management" }]`, "widgets")
	case "descope_saml_app":
		emit(fmt.Sprintf(`  manual_configuration = {
    acs_url   = "https://cov.example.com/acs%d"
    entity_id = "cov-entity"
  }`, generation), "manual_configuration")
	}
	return lines, paired
}

func (g *configGenerator) value(resourceType, name string, attribute schema.Attribute, generation int) (string, bool) {
	pick := func(options string) string {
		parts := strings.Split(options, "|")
		return parts[generation%len(parts)]
	}
	if override, ok := nestedValueOverrides[resourceType+"."+name]; ok {
		return pick(override), true
	}
	if override, ok := nestedValueOverrides[name]; ok {
		return pick(override), true
	}
	if name == "id" {
		// settable ids are machine names, stable per resource so generations update in place, and short: some validators cap them at 20
		initials := ""
		for _, word := range strings.Split(strings.TrimPrefix(resourceType, "descope_"), "_") {
			initials += word[:1]
		}
		return fmt.Sprintf("%q", "cov"+initials), true
	}
	if override, ok := attributeValueOverrides[resourceType+"."+name]; ok {
		return pick(override), true
	}
	if override, ok := attributeValueOverrides[name]; ok {
		return pick(override), true
	}

	switch attribute.(type) {
	case schema.StringAttribute:
		for _, pattern := range namePatternValues {
			if strings.Contains(name, pattern.fragment) {
				return pick(pattern.values), true
			}
		}
		if attribute.IsSensitive() {
			return fmt.Sprintf("%q", "not-a-real-secret"), true
		}
		return fmt.Sprintf(`"cov %s %d"`, strings.ReplaceAll(name, "_", " "), generation), true
	case schema.BoolAttribute:
		if generation%2 == 0 {
			return "true", true
		}
		return "false", true
	case schema.Int64Attribute:
		return strconv.Itoa(3 + generation), true
	case schema.Float64Attribute:
		return strconv.Itoa(25 + generation), true
	case schema.ListAttribute, schema.SetAttribute:
		for _, pattern := range namePatternValues {
			if strings.Contains(name, pattern.fragment) && strings.HasPrefix(pick(pattern.values), "[") {
				return pick(pattern.values), true
			}
		}
		return fmt.Sprintf(`["cov-%s-%d"]`, strings.ReplaceAll(name, "_", "-"), generation), true
	case schema.MapAttribute:
		return fmt.Sprintf(`{ "cov-key" = "cov-value-%d" }`, generation), true
	}
	return "", false // nested objects and collections are exercised by the scenario catalog
}

func (g *configGenerator) report(t *testing.T) {
	t.Helper()
	perType := map[string][2]int{}
	var missing []string
	for key := range g.settable {
		resourceType, _, _ := strings.Cut(key, ".")
		counts := perType[resourceType]
		counts[1]++
		if g.covered[key] {
			counts[0]++
		} else {
			missing = append(missing, key)
		}
		perType[resourceType] = counts
	}
	total, covered := 0, 0
	types := make([]string, 0, len(perType))
	for resourceType := range perType {
		types = append(types, resourceType)
	}
	sort.Strings(types)
	for _, resourceType := range types {
		counts := perType[resourceType]
		covered += counts[0]
		total += counts[1]
		t.Logf("coverage: %s %d/%d attributes", resourceType, counts[0], counts[1])
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Logf("coverage: not exercised: %s", strings.Join(missing, ", "))
	}
	t.Logf("coverage: TOTAL %d/%d settable attributes exercised with generated values", covered, total)
}
