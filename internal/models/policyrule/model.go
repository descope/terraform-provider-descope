package policyrule

import (
	"context"
	"encoding/json"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/listtype"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/types/valuelisttype"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PolicyRuleModel struct {
	ID                stringattr.Type                    `tfsdk:"id"`
	ProjectID         stringattr.Type                    `tfsdk:"project_id"`
	Version           intattr.Type                       `tfsdk:"version"`
	Name              stringattr.Type                    `tfsdk:"name"`
	Description       stringattr.Type                    `tfsdk:"description"`
	Enabled           boolattr.Type                      `tfsdk:"enabled"`
	RuleFamily        stringattr.Type                    `tfsdk:"rule_family"`
	ActionKinds       strlistattr.Type                   `tfsdk:"action_kinds"`
	Effect            stringattr.Type                    `tfsdk:"effect"`
	PrincipalType     stringattr.Type                    `tfsdk:"principal_type"`
	PrincipalSelector strlistattr.Type                   `tfsdk:"principal_selector"`
	ResourceTargets   listattr.Type[ResourceTargetModel] `tfsdk:"resource_targets"`
	Grants            listattr.Type[GrantModel]          `tfsdk:"grants"`
	Conditions        listattr.Type[ConditionModel]      `tfsdk:"conditions"`
	CedarText         stringattr.Type                    `tfsdk:"cedar_text"`
	CreatedTime       intattr.Type                       `tfsdk:"created_time"`
	ModifiedTime      intattr.Type                       `tfsdk:"modified_time"`
}

type ResourceTargetModel struct {
	Type      stringattr.Type  `tfsdk:"type"`
	AllOfType boolattr.Type    `tfsdk:"all_of_type"`
	IDs       strlistattr.Type `tfsdk:"ids"`
}

type GrantModel struct {
	Scopes           strlistattr.Type `tfsdk:"scopes"`
	AllowedAudiences strlistattr.Type `tfsdk:"allowed_audiences"`
	AllScopes        boolattr.Type    `tfsdk:"all_scopes"`
}

type ConditionModel struct {
	Key       stringattr.Type `tfsdk:"key"`
	Operator  stringattr.Type `tfsdk:"operator"`
	ValueJSON stringattr.Type `tfsdk:"value_json"`
}

func (m PolicyRuleModel) ToPolicyRule(ctx context.Context) (infra.PolicyRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	rule := infra.PolicyRule{
		ID:                m.ID.ValueString(),
		Version:           m.Version.ValueInt64(),
		Name:              m.Name.ValueString(),
		Description:       m.Description.ValueString(),
		Enabled:           m.Enabled.ValueBool(),
		RuleFamily:        m.RuleFamily.ValueString(),
		ActionKinds:       stringList(ctx, m.ActionKinds, &diags),
		Effect:            m.Effect.ValueString(),
		PrincipalType:     m.PrincipalType.ValueString(),
		PrincipalSelector: stringList(ctx, m.PrincipalSelector, &diags),
		ResourceTargets:   []infra.PolicyRuleResourceTarget{},
		Grants:            []infra.PolicyRuleGrant{},
		Conditions:        []infra.PolicyRuleCondition{},
	}

	targets, targetDiags := m.ResourceTargets.ToSlice(ctx)
	diags.Append(targetDiags...)
	for _, target := range targets {
		rule.ResourceTargets = append(rule.ResourceTargets, infra.PolicyRuleResourceTarget{
			Type:      target.Type.ValueString(),
			AllOfType: target.AllOfType.ValueBool(),
			IDs:       stringList(ctx, target.IDs, &diags),
		})
	}

	grants, grantDiags := m.Grants.ToSlice(ctx)
	diags.Append(grantDiags...)
	for _, grant := range grants {
		value := infra.PolicyRuleGrant{
			Scopes:           stringList(ctx, grant.Scopes, &diags),
			AllowedAudiences: stringList(ctx, grant.AllowedAudiences, &diags),
			AllScopes:        grant.AllScopes.ValueBool(),
		}
		if len(value.Scopes) == 0 && len(value.AllowedAudiences) == 0 && !value.AllScopes {
			diags.AddAttributeError(path.Root("grants"), "Invalid Policy Rule Grant", "Each grant must define scopes, allowed audiences, or all_scopes")
		}
		rule.Grants = append(rule.Grants, value)
	}

	conditions, conditionDiags := m.Conditions.ToSlice(ctx)
	diags.Append(conditionDiags...)
	for i, condition := range conditions {
		value := json.RawMessage(condition.ValueJSON.ValueString())
		if err := validateConditionValue(condition.Operator.ValueString(), value); err != nil {
			diags.AddAttributeError(path.Root("conditions").AtListIndex(i).AtName("value_json"), "Invalid Condition Value", err.Error())
		}
		rule.Conditions = append(rule.Conditions, infra.PolicyRuleCondition{
			Key:      condition.Key.ValueString(),
			Operator: condition.Operator.ValueString(),
			Value:    value,
		})
	}

	if rule.RuleFamily == "token_exchange" && rule.Effect == "permit" && len(rule.Grants) == 0 {
		diags.AddAttributeError(path.Root("grants"), "Missing Policy Rule Grant", "A token_exchange permit rule requires at least one grant")
	}
	return rule, diags
}

func (m *PolicyRuleModel) SetPolicyRule(ctx context.Context, rule *infra.PolicyRule) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID = stringattr.Value(rule.ID)
	m.Version = intattr.Value(rule.Version)
	m.Name = stringattr.Value(rule.Name)
	m.Description = stringattr.Value(rule.Description)
	m.Enabled = boolattr.Value(rule.Enabled)
	m.RuleFamily = stringattr.Value(rule.RuleFamily)
	m.ActionKinds = stringListValue(ctx, rule.ActionKinds, &diags)
	m.Effect = stringattr.Value(rule.Effect)
	m.PrincipalType = stringattr.Value(rule.PrincipalType)
	m.PrincipalSelector = stringListValue(ctx, rule.PrincipalSelector, &diags)
	m.CedarText = stringattr.Value(rule.CedarText)
	m.CreatedTime = intattr.Value(rule.CreatedTime)
	m.ModifiedTime = intattr.Value(rule.ModifiedTime)

	targets := make([]*ResourceTargetModel, 0, len(rule.ResourceTargets))
	for _, target := range rule.ResourceTargets {
		targets = append(targets, &ResourceTargetModel{
			Type:      stringattr.Value(target.Type),
			AllOfType: boolattr.Value(target.AllOfType),
			IDs:       stringListValue(ctx, target.IDs, &diags),
		})
	}
	m.ResourceTargets = nestedListValue(ctx, targets, &diags)

	grants := make([]*GrantModel, 0, len(rule.Grants))
	for _, grant := range rule.Grants {
		grants = append(grants, &GrantModel{
			Scopes:           stringListValue(ctx, grant.Scopes, &diags),
			AllowedAudiences: stringListValue(ctx, grant.AllowedAudiences, &diags),
			AllScopes:        boolattr.Value(grant.AllScopes),
		})
	}
	m.Grants = nestedListValue(ctx, grants, &diags)

	conditions := make([]*ConditionModel, 0, len(rule.Conditions))
	for i, condition := range rule.Conditions {
		value, err := compactJSON(condition.Value)
		if err != nil {
			diags.AddAttributeError(path.Root("conditions").AtListIndex(i).AtName("value_json"), "Invalid Policy Rule Response", err.Error())
			continue
		}
		conditions = append(conditions, &ConditionModel{
			Key:       stringattr.Value(condition.Key),
			Operator:  stringattr.Value(condition.Operator),
			ValueJSON: stringattr.Value(value),
		})
	}
	m.Conditions = nestedListValue(ctx, conditions, &diags)
	return diags
}

func stringList(ctx context.Context, value strlistattr.Type, diags *diag.Diagnostics) []string {
	elements, elementDiags := value.ToSlice(ctx)
	diags.Append(elementDiags...)
	result := make([]string, 0, len(elements))
	for _, element := range elements {
		result = append(result, element.ValueString())
	}
	return result
}

func stringListValue(ctx context.Context, values []string, diags *diag.Diagnostics) strlistattr.Type {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}
	result, valueDiags := valuelisttype.NewValue[types.String](ctx, elements)
	diags.Append(valueDiags...)
	return result
}

func nestedListValue[T any](ctx context.Context, values []*T, diags *diag.Diagnostics) listattr.Type[T] {
	result, valueDiags := listtype.NewValue(ctx, values)
	diags.Append(valueDiags...)
	return result
}
