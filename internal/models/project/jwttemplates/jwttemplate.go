package jwttemplates

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var JWTTemplateAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"name":               stringattr.Required(),
	"description":        stringattr.Default(""),
	"auth_schema":        stringattr.Default("default", stringvalidator.OneOf("default", "tenantOnly", "none")),
	"empty_claim_policy": stringattr.Default("none", stringvalidator.OneOf("none", "nil", "delete")),
	"conformance_issuer": boolattr.Default(false),
	"enforce_issuer":     boolattr.Default(false),
	"template":           stringattr.Required(),
}

type JWTTemplateModel struct {
	ID                stringattr.Type `tfsdk:"id"`
	Name              stringattr.Type `tfsdk:"name"`
	Description       stringattr.Type `tfsdk:"description"`
	AuthSchema        stringattr.Type `tfsdk:"auth_schema"`
	EmptyClaimPolicy  stringattr.Type `tfsdk:"empty_claim_policy"`
	ConformanceIssuer boolattr.Type   `tfsdk:"conformance_issuer"`
	EnforceIssuer     boolattr.Type   `tfsdk:"enforce_issuer"`
	Template          stringattr.Type `tfsdk:"template"`
}

func (m *JWTTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.AuthSchema, data, "authSchema")
	stringattr.Get(m.EmptyClaimPolicy, data, "emptyClaimPolicy")
	boolattr.Get(m.ConformanceIssuer, data, "conformanceIssuer")
	boolattr.Get(m.EnforceIssuer, data, "enforceIssuer")

	// convert template JSON string to map
	template := map[string]any{}
	if err := json.Unmarshal([]byte(m.Template.ValueString()), &template); err != nil {
		h.Invalid("Expected a valid JSON string for the template attribute")
	}
	data["template"] = template

	// use the name as a lookup key to set the JWT template reference or existing id
	templateName := m.Name.ValueString()
	if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, templateName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for JWT template '%s' to: %s", templateName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown JWT template reference", "No JWT template named '%s' was defined", templateName)
	}

	return data
}

func (m *JWTTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	stringattr.Set(&m.AuthSchema, data, "authSchema")
	stringattr.Set(&m.EmptyClaimPolicy, data, "emptyClaimPolicy")
	boolattr.Set(&m.ConformanceIssuer, data, "conformanceIssuer")
	boolattr.Set(&m.EnforceIssuer, data, "enforceIssuer")

	if m.Template.ValueString() == "" {
		template := "{}"
		if t, ok := data["template"].(map[string]any); ok {
			if b, err := json.Marshal(t); err == nil { // XXX json.MarshalIndent(t, "", "  ") ?
				template = string(b)
			}
		}
		m.Template = stringattr.Value(template)
	}
}
