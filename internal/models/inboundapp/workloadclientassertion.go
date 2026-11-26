package inboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var WorkloadClientAssertionValidator = objattr.NewValidator[WorkloadClientAssertionModel]("must carry the federation type's fields")

var WorkloadClientAssertionAttributes = map[string]schema.Attribute{
	"federation_type":  stringattr.Required(stringvalidator.OneOf("spiffe", "oidc")),
	"trust_domain":     stringattr.Default(""),
	"issuer":           stringattr.Default(""),
	"jwks_uri":         stringattr.Default(""),
	"allowed_subjects": strlistattr.Required(listvalidator.SizeAtLeast(1)),
	"sign_algorithm":   stringattr.Default(""),
}

type WorkloadClientAssertionModel struct {
	FederationType  stringattr.Type  `tfsdk:"federation_type"`
	TrustDomain     stringattr.Type  `tfsdk:"trust_domain"`
	Issuer          stringattr.Type  `tfsdk:"issuer"`
	JwksUri         stringattr.Type  `tfsdk:"jwks_uri"`
	AllowedSubjects strlistattr.Type `tfsdk:"allowed_subjects"`
	SignAlgorithm   stringattr.Type  `tfsdk:"sign_algorithm"`
}

func (m *WorkloadClientAssertionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.FederationType, data, "federationType")
	stringattr.Get(m.TrustDomain, data, "trustDomain")
	stringattr.Get(m.Issuer, data, "issuer")
	stringattr.Get(m.JwksUri, data, "jwksUri")
	strlistattr.Get(m.AllowedSubjects, data, "allowedSubjects", h)
	stringattr.Get(m.SignAlgorithm, data, "signAlgorithm")
	return data
}

func (m *WorkloadClientAssertionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.FederationType, data, "federationType")
	setTrimmed(&m.TrustDomain, data, "trustDomain")
	setTrimmed(&m.Issuer, data, "issuer")
	setTrimmed(&m.JwksUri, data, "jwksUri")
	strlistattr.Set(&m.AllowedSubjects, data, "allowedSubjects", h)
	setTrimmed(&m.SignAlgorithm, data, "signAlgorithm")
}

func (m *WorkloadClientAssertionModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.FederationType, m.TrustDomain, m.Issuer, m.JwksUri) {
		return
	}
	switch m.FederationType.ValueString() {
	case "spiffe":
		if m.TrustDomain.ValueString() == "" {
			h.Missing("The trust_domain field is required when federation_type is spiffe")
		}
		if m.Issuer.ValueString() == "" && m.JwksUri.ValueString() == "" {
			h.Missing("A spiffe entry requires the jwks_uri of its bundle endpoint or the issuer of its OIDC Discovery Provider")
		}
	case "oidc":
		if m.Issuer.ValueString() == "" {
			h.Missing("The issuer field is required when federation_type is oidc")
		}
		if m.TrustDomain.ValueString() != "" {
			h.Conflict("The trust_domain field cannot be used when federation_type is oidc")
		}
	}
}
