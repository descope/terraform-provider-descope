package jwttemplates

import (
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var JWTTemplateAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"name":               stringattr.Required(),
	"description":        stringattr.Default(""),
	"auth_schema":        stringattr.Default("default", stringvalidator.OneOf("default", "tenantOnly", "none")),
	"conformance_issuer": boolattr.Default(false),
	"template":           stringattr.Required(),
	"type":               stringattr.Required(stringvalidator.OneOf("key", "user")),
}

type JWTTemplateModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	AuthSchema        types.String `tfsdk:"auth_schema"`
	ConformanceIssuer types.Bool   `tfsdk:"conformance_issuer"`
	Template          types.String `tfsdk:"template"`
	Type              types.String `tfsdk:"type"`
}

func (m *JWTTemplateModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	stringattr.Get(m.AuthSchema, data, "authSchema")
	boolattr.Get(m.ConformanceIssuer, data, "conformanceIssuer")
	stringattr.Get(m.Type, data, "type")

	// convert template JSON string to map
	template := map[string]any{}
	if err := json.Unmarshal([]byte(m.Template.ValueString()), &template); err != nil {
		h.Error("Unexpected template structure", "Could not deserialize template attribute, json expected")
	}
	data["template"] = template

	// use the name as a lookup key to set the JWT template reference or existing id
	templateName := m.Name.ValueString()
	if ref := h.Refs.Get(helpers.JWTTemplateReferenceKey, templateName); ref != nil {
		refValue := ref.ReferenceValue()
		h.Log("Updating reference for JWT template '%s' to: %s", templateName, refValue)
		data["id"] = refValue
	} else {
		h.Error("Unknown JWT template reference", "No JWT template named '"+templateName+"' was defined")
	}

	return data
}

func (m *JWTTemplateModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all JWT template values are specified in the configuration
}
