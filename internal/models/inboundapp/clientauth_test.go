package inboundapp_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const testPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA89j2dNDIGQJ6/kvqrBaG
PlAbx36ofVw53VVPRw4XuwUDjkKBopjIh4so0jsMH+94y162xsHss85QKPtxUJGj
zGKrLWZTpR1Tc5jPaqBRbOq3hT7RhU/sOD6UheGjq2SptqAFWzzJLvh0ti/zCt0D
KNtg/sS1UEKqlhHMAl6la5JsupegB7yfNV0oxHQzs2kKuwFKlpo6eT0y2fIM+ZmB
/RWYIA8OvM7/kKf6TPuNQ6ZboKSxC0bZBfCG+9feG4urZnZrtiSwyODmnbgsg5jQ
/8r1e77CoDSDsABFUg5TiuIV+NNf2VDdmfTZ13aovMJJJDIvxj9F4bxw6tNzztUM
gQIDAQAB
-----END PUBLIC KEY-----`

func validate(t *testing.T, validatable interface{ Validate(*helpers.Handler) }) diag.Diagnostics {
	t.Helper()
	diags := diag.Diagnostics{}
	validatable.Validate(helpers.NewHandler(context.Background(), &diags))
	return diags
}

func TestInboundAppClientAuthMethodValidation(t *testing.T) {
	jwks := &inboundapp.ClientAssertionModel{KeySource: types.StringValue("jwks"), JwksUri: types.StringValue("https://example.com/jwks")}
	spiffe := &inboundapp.WorkloadClientAssertionModel{FederationType: types.StringValue("spiffe"), TrustDomain: types.StringValue("prod.example.com"), JwksUri: types.StringValue("https://example.com/bundle"), AllowedSubjects: strlistattr.Value([]string{"spiffe://prod.example.com/ns/api"})}
	testCases := []struct {
		name       string
		method     string
		assertions []*inboundapp.ClientAssertionModel
		workload   []*inboundapp.WorkloadClientAssertionModel
		valid      bool
	}{
		{"empty method", "", nil, nil, true},
		{"client secret", "clientSecret", nil, nil, true},
		{"client secret keeps stored entries", "clientSecret", []*inboundapp.ClientAssertionModel{jwks}, []*inboundapp.WorkloadClientAssertionModel{spiffe}, true},
		{"private key jwt with an entry", "privateKeyJwt", []*inboundapp.ClientAssertionModel{jwks}, nil, true},
		{"private key jwt without entries", "privateKeyJwt", nil, nil, false},
		{"private key jwt with only workload entries", "privateKeyJwt", nil, []*inboundapp.WorkloadClientAssertionModel{spiffe}, false},
		{"workload identity with an entry", "workloadIdentity", nil, []*inboundapp.WorkloadClientAssertionModel{spiffe}, true},
		{"workload identity without entries", "workloadIdentity", nil, nil, false},
		{"workload identity with only key entries", "workloadIdentity", []*inboundapp.ClientAssertionModel{jwks}, nil, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, "confidential"))
			m.ClientAuthMethod = types.StringValue(tc.method)
			m.ClientAssertions = listattr.Value(tc.assertions)
			m.WorkloadClientAssertions = listattr.Value(tc.workload)
			if diags := validate(t, m); diags.HasError() == tc.valid {
				t.Errorf("expected valid=%t, got %v", tc.valid, diags)
			}
		})
	}
	t.Run("unknown entries are not validated", func(t *testing.T) {
		m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, "confidential"))
		m.ClientAuthMethod = types.StringValue("privateKeyJwt")
		m.ClientAssertions = listattr.Type[inboundapp.ClientAssertionModel]{}
		m.ClientAssertions.ListValue = types.ListUnknown(m.ClientAssertions.ElementType(context.Background()))
		if diags := validate(t, m); diags.HasError() {
			t.Errorf("expected no error, got %v", diags)
		}
	})
}

func TestClientAssertionValidation(t *testing.T) {
	str := types.StringValue
	testCases := []struct {
		name  string
		model inboundapp.ClientAssertionModel
		valid bool
	}{
		{"jwks with uri", inboundapp.ClientAssertionModel{KeySource: str("jwks"), JwksUri: str("https://example.com/jwks")}, true},
		{"jwks without uri", inboundapp.ClientAssertionModel{KeySource: str("jwks")}, false},
		{"jwks with a public key", inboundapp.ClientAssertionModel{KeySource: str("jwks"), JwksUri: str("https://example.com/jwks"), PublicKey: str(testPublicKeyPEM)}, false},
		{"public key with pem", inboundapp.ClientAssertionModel{KeySource: str("publicKey"), PublicKey: str(testPublicKeyPEM), SignAlgorithm: str("RS256")}, true},
		{"public key without pem", inboundapp.ClientAssertionModel{KeySource: str("publicKey")}, false},
		{"public key with a jwks uri", inboundapp.ClientAssertionModel{KeySource: str("publicKey"), PublicKey: str(testPublicKeyPEM), JwksUri: str("https://example.com/jwks")}, false},
		{"unknown uri skips validation", inboundapp.ClientAssertionModel{KeySource: str("jwks"), JwksUri: types.StringUnknown()}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if diags := validate(t, &tc.model); diags.HasError() == tc.valid {
				t.Errorf("expected valid=%t, got %v", tc.valid, diags)
			}
		})
	}
}

func TestWorkloadClientAssertionValidation(t *testing.T) {
	str := types.StringValue
	subjects := strlistattr.Value([]string{"spiffe://prod.example.com/ns/api"})
	testCases := []struct {
		name  string
		model inboundapp.WorkloadClientAssertionModel
		valid bool
	}{
		{"spiffe with bundle endpoint", inboundapp.WorkloadClientAssertionModel{FederationType: str("spiffe"), TrustDomain: str("prod.example.com"), JwksUri: str("https://example.com/bundle"), AllowedSubjects: subjects}, true},
		{"spiffe with discovery provider", inboundapp.WorkloadClientAssertionModel{FederationType: str("spiffe"), TrustDomain: str("prod.example.com"), Issuer: str("https://example.com/spire"), AllowedSubjects: subjects}, true},
		{"spiffe without trust domain", inboundapp.WorkloadClientAssertionModel{FederationType: str("spiffe"), JwksUri: str("https://example.com/bundle"), AllowedSubjects: subjects}, false},
		{"spiffe without key source", inboundapp.WorkloadClientAssertionModel{FederationType: str("spiffe"), TrustDomain: str("prod.example.com"), AllowedSubjects: subjects}, false},
		{"oidc with issuer", inboundapp.WorkloadClientAssertionModel{FederationType: str("oidc"), Issuer: str("https://example.com/issuer"), AllowedSubjects: subjects}, true},
		{"oidc with jwks override", inboundapp.WorkloadClientAssertionModel{FederationType: str("oidc"), Issuer: str("https://example.com/issuer"), JwksUri: str("https://example.com/keys"), AllowedSubjects: subjects}, true},
		{"oidc without issuer", inboundapp.WorkloadClientAssertionModel{FederationType: str("oidc"), AllowedSubjects: subjects}, false},
		{"oidc with a trust domain", inboundapp.WorkloadClientAssertionModel{FederationType: str("oidc"), Issuer: str("https://example.com/issuer"), TrustDomain: str("prod.example.com"), AllowedSubjects: subjects}, false},
		{"unknown issuer skips validation", inboundapp.WorkloadClientAssertionModel{FederationType: str("oidc"), Issuer: types.StringUnknown(), AllowedSubjects: subjects}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if diags := validate(t, &tc.model); diags.HasError() == tc.valid {
				t.Errorf("expected valid=%t, got %v", tc.valid, diags)
			}
		})
	}
}

func TestClientAssertionSetValuesKeepsUntrimmedPlan(t *testing.T) {
	h := helpers.NewHandler(context.Background(), &diag.Diagnostics{})
	m := &inboundapp.ClientAssertionModel{KeySource: types.StringValue("publicKey"), PublicKey: types.StringValue(testPublicKeyPEM + "\n"), SignAlgorithm: types.StringValue("")}
	m.SetValues(h, map[string]any{"keySource": "publicKey", "publicKey": testPublicKeyPEM, "signAlgorithm": "RS256"})
	if m.PublicKey.ValueString() != testPublicKeyPEM+"\n" {
		t.Errorf("expected the planned public key to be kept, got %q", m.PublicKey.ValueString())
	}
	if m.SignAlgorithm.ValueString() != "RS256" {
		t.Errorf("expected a changed value to be adopted, got %q", m.SignAlgorithm.ValueString())
	}
	fresh := &inboundapp.ClientAssertionModel{}
	fresh.SetValues(h, map[string]any{"keySource": "jwks", "jwksUri": "https://example.com/jwks"})
	if fresh.JwksUri.ValueString() != "https://example.com/jwks" || fresh.PublicKey.ValueString() != "" {
		t.Errorf("expected a fresh model to adopt the stored values, got %v", fresh)
	}
	if fresh.PublicKey.IsNull() || fresh.SignAlgorithm.IsNull() {
		t.Errorf("expected fields the backend omits to be empty rather than null on import, got %v", fresh)
	}
	stale := &inboundapp.ClientAssertionModel{KeySource: types.StringValue("jwks"), JwksUri: types.StringValue("https://example.com/jwks"), SignAlgorithm: types.StringValue("RS256")}
	stale.SetValues(h, map[string]any{"keySource": "jwks", "jwksUri": "https://example.com/jwks"})
	if !stale.SignAlgorithm.Equal(types.StringValue("")) {
		t.Errorf("expected a field the backend omits to read back as empty over the state value, got %s", stale.SignAlgorithm)
	}
}

func TestWorkloadClientAssertionValues(t *testing.T) {
	h := helpers.NewHandler(context.Background(), &diag.Diagnostics{})
	m := &inboundapp.WorkloadClientAssertionModel{FederationType: types.StringValue("spiffe"), TrustDomain: types.StringValue("prod.example.com"), JwksUri: types.StringValue("https://example.com/bundle"), AllowedSubjects: strlistattr.Value([]string{"spiffe://prod.example.com/ns/api", "spiffe://prod.example.com/ns/batch/*"})}
	data := m.Values(h)
	if data["federationType"] != "spiffe" || data["trustDomain"] != "prod.example.com" || data["jwksUri"] != "https://example.com/bundle" {
		t.Errorf("unexpected values: %v", data)
	}
	subjects, ok := data["allowedSubjects"].([]string)
	if !ok || len(subjects) != 2 || subjects[1] != "spiffe://prod.example.com/ns/batch/*" {
		t.Errorf("unexpected allowed subjects: %v", data["allowedSubjects"])
	}
	data["allowedSubjects"] = []any{"spiffe://prod.example.com/ns/api", "spiffe://prod.example.com/ns/batch/*"}
	parsed := &inboundapp.WorkloadClientAssertionModel{}
	parsed.SetValues(h, data)
	if !parsed.AllowedSubjects.Equal(m.AllowedSubjects) || !parsed.TrustDomain.Equal(m.TrustDomain) {
		t.Errorf("round trip changed the model: %v", parsed)
	}
}

func TestInboundAppClientAuth(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	testacc.Run(t,
		// privateKeyJwt needs at least one client assertion
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				client_auth_method = "privateKeyJwt"
			`),
			ExpectError: regexp.MustCompile(`client_assertions`),
		},
		// A key source must carry its own field
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				client_auth_method = "privateKeyJwt"
				client_assertions = [{ key_source = "jwks" }]
			`),
			ExpectError: regexp.MustCompile(`jwks_uri`),
		},
		// Create with private_key_jwt keys from a JWKS URL and a PEM public key
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				client_auth_method = "privateKeyJwt"
				client_assertions = [
					{
						key_source = "jwks"
						jwks_uri = "https://example.com/jwks"
						sign_algorithm = "RS256"
					},
					{
						key_source = "publicKey"
						public_key = <<-EOT
`+testPublicKeyPEM+`
						EOT
					},
				]
			`),
			Check: a.Check(map[string]any{
				"client_auth_method":                 "privateKeyJwt",
				"client_assertions.#":                "2",
				"client_assertions.0.key_source":     "jwks",
				"client_assertions.0.jwks_uri":       "https://example.com/jwks",
				"client_assertions.0.sign_algorithm": "RS256",
				"client_assertions.1.key_source":     "publicKey",
				"client_assertions.1.public_key":     testacc.AttributeHasPrefix("-----BEGIN PUBLIC KEY-----"),
				"workload_client_assertions.#":       "0",
			}),
		},
		// Switch to workload identity with a SPIFFE trust domain and an OIDC issuer
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				client_auth_method = "workloadIdentity"
				workload_client_assertions = [
					{
						federation_type = "spiffe"
						trust_domain = "prod.example.com"
						jwks_uri = "https://example.com/bundle"
						allowed_subjects = ["spiffe://prod.example.com/ns/api", "spiffe://prod.example.com/ns/batch/*"]
					},
					{
						federation_type = "oidc"
						issuer = "https://example.com/issuer"
						allowed_subjects = ["repo:acme/.*:ref:refs/heads/main"]
						sign_algorithm = "ES256"
					},
				]
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(a.Path(), plancheck.ResourceActionUpdate)},
			},
			Check: a.Check(map[string]any{
				"client_auth_method":                              "workloadIdentity",
				"client_assertions.#":                             "0",
				"workload_client_assertions.#":                    "2",
				"workload_client_assertions.0.federation_type":    "spiffe",
				"workload_client_assertions.0.trust_domain":       "prod.example.com",
				"workload_client_assertions.0.jwks_uri":           "https://example.com/bundle",
				"workload_client_assertions.0.allowed_subjects.#": "2",
				"workload_client_assertions.0.allowed_subjects.1": "spiffe://prod.example.com/ns/batch/*",
				"workload_client_assertions.1.federation_type":    "oidc",
				"workload_client_assertions.1.issuer":             "https://example.com/issuer",
				"workload_client_assertions.1.allowed_subjects.0": "repo:acme/.*:ref:refs/heads/main",
				"workload_client_assertions.1.sign_algorithm":     "ES256",
			}),
		},
		// Import keeps every entry
		resource.TestStep{
			ResourceName:            a.Path(),
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
			ImportStateVerifyIgnore: []string{"client_secret", "deletion_protection"},
		},
		// Back to the client secret, which clears the entries
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				client_auth_method = "clientSecret"
			`),
			Check: a.Check(map[string]any{
				"client_auth_method":           "clientSecret",
				"client_assertions.#":          "0",
				"workload_client_assertions.#": "0",
			}),
		},
	)
}
