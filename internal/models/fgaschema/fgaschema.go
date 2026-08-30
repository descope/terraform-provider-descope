package fgaschema

import (
	"context"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages the FGA (fine-grained authorization) schema for a Descope project as a single DSL string. Setting an empty schema, or destroying the resource, clears the project's FGA schema (which also removes any FGA relations that depend on it).",
	Attributes:          FGASchemaAttributes,
}

var FGASchemaAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"schema":     stringattr.Default("", stringvalidator.LengthAtMost(100000), modelPrefixValidator{}),
}

type FGASchemaModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Schema    stringattr.Type `tfsdk:"schema"`
}

func (m *FGASchemaModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Schema, data, "dsl", stringattr.TrimSpaces) // the public API dialect names the schema string "dsl"
	return data
}

func (m *FGASchemaModel) SetValues(h *helpers.Handler, data map[string]any) {
	// the backend stores the DSL as strings.TrimSpace(dsl)+"\n", so a stored value is adopted only when it differs by more
	// than surrounding whitespace, while a state without one (an import) always takes it
	if dsl, ok := data["dsl"].(string); ok {
		unset := m.Schema.IsNull() || m.Schema.IsUnknown()
		if unset || strings.TrimSpace(dsl) != strings.TrimSpace(m.Schema.ValueString()) {
			m.Schema = stringattr.Value(strings.TrimSpace(dsl))
		}
	}
}

func (m *FGASchemaModel) GetID() stringattr.Type {
	return m.ID
}

func (m *FGASchemaModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *FGASchemaModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

type modelPrefixValidator struct{}

func (v modelPrefixValidator) Description(_ context.Context) string {
	return `must start with "model AuthZ" when set`
}

func (v modelPrefixValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v modelPrefixValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if value == "" {
		return
	}
	if !strings.HasPrefix(value, "model AuthZ") {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid FGA Schema", fmt.Sprintf("The %s attribute must start with 'model AuthZ', make sure you're using the schema from the code view in the FGA tab in the Descope console", req.Path)))
	}
}
