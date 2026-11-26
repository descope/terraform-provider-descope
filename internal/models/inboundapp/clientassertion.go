package inboundapp

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ClientAssertionValidator = objattr.NewValidator[ClientAssertionModel]("must carry the key source's field")

var ClientAssertionAttributes = map[string]schema.Attribute{
	"key_source":     stringattr.Required(stringvalidator.OneOf("jwks", "publicKey")),
	"jwks_uri":       stringattr.Default(""),
	"public_key":     stringattr.Default(""),
	"sign_algorithm": stringattr.Default(""),
}

type ClientAssertionModel struct {
	KeySource     stringattr.Type `tfsdk:"key_source"`
	JwksUri       stringattr.Type `tfsdk:"jwks_uri"`
	PublicKey     stringattr.Type `tfsdk:"public_key"`
	SignAlgorithm stringattr.Type `tfsdk:"sign_algorithm"`
}

func (m *ClientAssertionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.KeySource, data, "keySource")
	stringattr.Get(m.JwksUri, data, "jwksUri")
	stringattr.Get(m.PublicKey, data, "publicKey")
	stringattr.Get(m.SignAlgorithm, data, "signAlgorithm")
	return data
}

func (m *ClientAssertionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.KeySource, data, "keySource")
	setTrimmed(&m.JwksUri, data, "jwksUri")
	setTrimmed(&m.PublicKey, data, "publicKey")
	setTrimmed(&m.SignAlgorithm, data, "signAlgorithm")
}

func (m *ClientAssertionModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.KeySource, m.JwksUri, m.PublicKey) {
		return
	}
	switch m.KeySource.ValueString() {
	case "jwks":
		if m.JwksUri.ValueString() == "" {
			h.Missing("The jwks_uri field is required when key_source is jwks")
		}
		if m.PublicKey.ValueString() != "" {
			h.Conflict("The public_key field cannot be used when key_source is jwks")
		}
	case "publicKey":
		if m.PublicKey.ValueString() == "" {
			h.Missing("The public_key field is required when key_source is publicKey")
		}
		if m.JwksUri.ValueString() != "" {
			h.Conflict("The jwks_uri field cannot be used when key_source is publicKey")
		}
	}
}

// The backend stores entry fields trimmed and omits empty ones, so a planned value that only differs by surrounding
// whitespace is kept and a missing field is empty, whatever the state holds.
func setTrimmed(s *stringattr.Type, data map[string]any, key string) {
	planned := *s
	v, _ := data[key].(string)
	*s = stringattr.Value(v)
	if !planned.IsNull() && !planned.IsUnknown() && strings.TrimSpace(planned.ValueString()) == s.ValueString() {
		*s = planned
	}
}
