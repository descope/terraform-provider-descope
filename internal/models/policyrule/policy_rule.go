package policyrule

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ResourceTargetAttributes = map[string]schema.Attribute{
	"type":        stringattr.Required(stringvalidator.OneOf("api", "mcp", "outbound_app", "any")),
	"all_of_type": boolattr.Default(false),
	"ids":         strlistattr.Default(stringattr.NonEmptyValidator),
}

var GrantAttributes = map[string]schema.Attribute{
	"scopes":            strlistattr.Default(stringattr.NonEmptyValidator),
	"allowed_audiences": strlistattr.Default(stringattr.NonEmptyValidator),
	"all_scopes":        boolattr.Default(false),
}

var ConditionAttributes = map[string]schema.Attribute{
	"key":        stringattr.Required(),
	"operator":   stringattr.Required(stringvalidator.OneOf("equal", "notEqual", "contains", "notContains", "in", "notIn")),
	"value_json": stringattr.Required(jsonValueValidator{}, compactJSONPlanModifier{}),
}

var PolicyRuleAttributes = map[string]schema.Attribute{
	"id":                 stringattr.Identifier(),
	"project_id":         stringattr.Required(stringplanmodifier.RequiresReplace()),
	"version":            intattr.Generated(),
	"name":               stringattr.Required(),
	"description":        stringattr.Default(""),
	"enabled":            boolattr.Required(),
	"rule_family":        stringattr.Required(stringvalidator.OneOf("resource_access", "outbound_access", "token_exchange")),
	"action_kinds":       strlistattr.Required(listvalidator.SizeAtLeast(1), stringvalidator.OneOf("user_access", "client_access", "exchange_token", "fetch_outbound_token")),
	"effect":             stringattr.Required(stringvalidator.OneOf("permit", "forbid")),
	"principal_type":     stringattr.Required(stringvalidator.OneOf("any", "user", "client")),
	"principal_selector": strlistattr.Default(stringattr.NonEmptyValidator),
	"resource_targets":   listattr.Default[ResourceTargetModel](ResourceTargetAttributes),
	"grants":             listattr.Default[GrantModel](GrantAttributes),
	"conditions":         listattr.Default[ConditionModel](ConditionAttributes),
	"cedar_text":         stringattr.Generated(),
	"created_time":       intattr.Generated(),
	"modified_time":      intattr.Generated(),
}

var Schema = schema.Schema{
	MarkdownDescription: "Manages a declarative Descope policy rule through the Management Service policy-rule CRUD API.",
	Attributes:          PolicyRuleAttributes,
}

type jsonValueValidator struct{}

func (jsonValueValidator) Description(context.Context) string {
	return "must be a valid JSON value"
}

func (v jsonValueValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (jsonValueValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if !json.Valid([]byte(req.ConfigValue.ValueString())) {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path,
			"Invalid JSON Value",
			fmt.Sprintf("Attribute %s must contain a valid JSON value", req.Path),
		))
	}
}

type compactJSONPlanModifier struct{}

func (compactJSONPlanModifier) Description(context.Context) string {
	return "normalizes equivalent JSON representations before apply"
}

func (m compactJSONPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (compactJSONPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	compact, err := compactJSON(json.RawMessage(req.PlanValue.ValueString()))
	if err == nil {
		resp.PlanValue = types.StringValue(compact)
	}
}
