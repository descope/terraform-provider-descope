package apiresource

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var IDJAGTargetAttributes = map[string]schema.Attribute{
	"authorization_server": stringattr.Required(stringattr.URLValidator, stringvalidator.LengthAtMost(3000)),
	"target_tenant":        stringattr.Default("", stringattr.StandardLenValidator),
	"continuation_trust":   listattr.Default[ContinuationTrustModel](ContinuationTrustAttributes, listvalidator.SizeAtMost(20)),
}

type IDJAGTargetModel struct {
	AuthorizationServer stringattr.Type                       `tfsdk:"authorization_server"`
	TargetTenant        stringattr.Type                       `tfsdk:"target_tenant"`
	ContinuationTrust   listattr.Type[ContinuationTrustModel] `tfsdk:"continuation_trust"`
}

func (m *IDJAGTargetModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.AuthorizationServer, data, "authorizationServer")
	stringattr.Get(m.TargetTenant, data, "targetTenant")
	listattr.Get(m.ContinuationTrust, data, "continuationTrust", h)
	return data
}

func (m *IDJAGTargetModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.AuthorizationServer, data, "authorizationServer")
	stringattr.Set(&m.TargetTenant, data, "targetTenant")
	listattr.Set(&m.ContinuationTrust, data, "continuationTrust", h)
}

var ContinuationTrustAttributes = map[string]schema.Attribute{
	"issuer":          stringattr.Required(stringattr.URLValidator, stringvalidator.LengthAtMost(3000)),
	"jwks_uri":        stringattr.Required(stringattr.URLValidator, stringvalidator.LengthAtMost(3000)),
	"sign_algorithm":  stringattr.Required(stringvalidator.OneOf("RS256", "ES256", "ES384", "ES512")),
	"actor_client_id": stringattr.Required(stringattr.StandardLenValidator),
	"enabled":         boolattr.Default(true),
}

type ContinuationTrustModel struct {
	Issuer        stringattr.Type `tfsdk:"issuer"`
	JWKSURI       stringattr.Type `tfsdk:"jwks_uri"`
	SignAlgorithm stringattr.Type `tfsdk:"sign_algorithm"`
	ActorClientID stringattr.Type `tfsdk:"actor_client_id"`
	Enabled       boolattr.Type   `tfsdk:"enabled"`
}

func (m *ContinuationTrustModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Issuer, data, "issuer")
	stringattr.Get(m.JWKSURI, data, "jwksUri")
	stringattr.Get(m.SignAlgorithm, data, "signAlgorithm")
	stringattr.Get(m.ActorClientID, data, "actorClientId")
	boolattr.Get(m.Enabled, data, "enabled")
	return data
}

func (m *ContinuationTrustModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Issuer, data, "issuer")
	stringattr.Set(&m.JWKSURI, data, "jwksUri")
	stringattr.Set(&m.SignAlgorithm, data, "signAlgorithm")
	stringattr.Set(&m.ActorClientID, data, "actorClientId")
	boolattr.SetDefault(&m.Enabled, data, "enabled", false)
}
