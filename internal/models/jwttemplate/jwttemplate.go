package jwttemplate

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single JWT template in a Descope project. A JWT template customizes the claims in the JWTs that are generated for user sessions or access keys, and can be referenced from the project's token settings.",
	Attributes:          JWTTemplateAttributes,
}

var JWTTemplateAttributes = map[string]schema.Attribute{
	"id":                       stringattr.Identifier(),
	"project_id":               stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":                     stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description":              stringattr.Default("", stringvalidator.LengthAtMost(500)),
	"type":                     stringattr.Required(stringvalidator.OneOf("user", "key")),
	"issuer_type":              stringattr.Default("legacy", stringvalidator.OneOf("legacy", "inbound", "federated")),
	"auth_schema":              stringattr.Default("default", stringvalidator.OneOf("default", "tenantOnly", "none")),
	"empty_claim_policy":       stringattr.Default("none", stringvalidator.OneOf("none", "nil", "delete")),
	"auto_tenant_claim":        boolattr.Default(false),
	"conformance_issuer":       boolattr.Default(false),
	"enforce_issuer":           boolattr.Default(false),
	"exclude_permission_claim": boolattr.Default(false),
	"override_subject_claim":   boolattr.Default(false),
	"add_jti_claim":            boolattr.Default(false),
	"template":                 jsonattr.Required(),
}

type JWTTemplateModel struct {
	ID                     stringattr.Type `tfsdk:"id"`
	ProjectID              stringattr.Type `tfsdk:"project_id"`
	Name                   stringattr.Type `tfsdk:"name"`
	Description            stringattr.Type `tfsdk:"description"`
	Type                   stringattr.Type `tfsdk:"type"`
	IssuerType             stringattr.Type `tfsdk:"issuer_type"`
	AuthSchema             stringattr.Type `tfsdk:"auth_schema"`
	EmptyClaimPolicy       stringattr.Type `tfsdk:"empty_claim_policy"`
	AutoDCT                boolattr.Type   `tfsdk:"auto_tenant_claim"`
	ConformanceIssuer      boolattr.Type   `tfsdk:"conformance_issuer"`
	EnforceIssuer          boolattr.Type   `tfsdk:"enforce_issuer"`
	ExcludePermissionClaim boolattr.Type   `tfsdk:"exclude_permission_claim"`
	OverrideSubjectClaim   boolattr.Type   `tfsdk:"override_subject_claim"`
	AddJtiClaim            boolattr.Type   `tfsdk:"add_jti_claim"`
	Template               jsonattr.Type   `tfsdk:"template"`
}

func (m *JWTTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.Type, data, "type")
	stringattr.Get(m.IssuerType, data, "issuerType")
	stringattr.Get(m.AuthSchema, data, "authSchema")
	stringattr.Get(m.EmptyClaimPolicy, data, "emptyClaimPolicy")
	boolattr.Get(m.AutoDCT, data, "autoDCT")
	boolattr.Get(m.ConformanceIssuer, data, "conformanceIssuer")
	boolattr.Get(m.EnforceIssuer, data, "enforceIssuer")
	boolattr.Get(m.ExcludePermissionClaim, data, "excludePermissions")
	boolattr.Get(m.OverrideSubjectClaim, data, "overrideSubject")
	boolattr.Get(m.AddJtiClaim, data, "addJti")
	jsonattr.Get(m.Template, data, "template")
	return data
}

func (m *JWTTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.Type, data, "type")
	stringattr.SetDefault(&m.IssuerType, data, "issuerType", "legacy") // enum fields come back as empty strings from templates that were created without them, e.g. when importing
	stringattr.SetDefault(&m.AuthSchema, data, "authSchema", "default")
	stringattr.SetDefault(&m.EmptyClaimPolicy, data, "emptyClaimPolicy", "none")
	boolattr.Set(&m.AutoDCT, data, "autoDCT")
	boolattr.Set(&m.ConformanceIssuer, data, "conformanceIssuer")
	boolattr.Set(&m.EnforceIssuer, data, "enforceIssuer")
	boolattr.Set(&m.ExcludePermissionClaim, data, "excludePermissions")
	boolattr.Set(&m.OverrideSubjectClaim, data, "overrideSubject")
	boolattr.Set(&m.AddJtiClaim, data, "addJti")
	jsonattr.Set(&m.Template, data, "template")
}

func (m *JWTTemplateModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Template, m.IssuerType, m.ConformanceIssuer) {
		return
	}
	if strings.TrimSpace(m.Template.ValueString()) == "null" {
		h.Invalid("The 'template' attribute must be a JSON object")
	}
	if m.IssuerType.ValueString() == "federated" && m.ConformanceIssuer.ValueBool() {
		h.Invalid("The 'conformance_issuer' attribute cannot be enabled when 'issuer_type' is 'federated'")
	}
}

func (m *JWTTemplateModel) GetID() stringattr.Type {
	return m.ID
}

func (m *JWTTemplateModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *JWTTemplateModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
