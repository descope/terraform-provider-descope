package policyrule

import (
	"encoding/json"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/stretchr/testify/require"
)

func TestPolicyRuleModel_ToPolicyRule_serializes_nested_values(t *testing.T) {
	// Given
	model := PolicyRuleModel{
		ID:                stringattr.Value("AR1"),
		ProjectID:         stringattr.Value("P1"),
		Version:           intattr.Value(2),
		Name:              stringattr.Value("rule"),
		Description:       stringattr.Value("description"),
		Enabled:           boolattr.Value(true),
		RuleFamily:        stringattr.Value("resource_access"),
		ActionKinds:       strlistattr.Value([]string{"client_access"}),
		Effect:            stringattr.Value("permit"),
		PrincipalType:     stringattr.Value("client"),
		PrincipalSelector: strlistattr.Value([]string{"client-1"}),
		ResourceTargets: listattr.Value([]*ResourceTargetModel{{
			Type:      stringattr.Value("api"),
			AllOfType: boolattr.Value(false),
			IDs:       strlistattr.Value([]string{"api-1"}),
		}}),
		Grants: listattr.Value([]*GrantModel{{
			Scopes:           strlistattr.Value([]string{"read"}),
			AllowedAudiences: strlistattr.Value([]string{"audience"}),
			AllScopes:        boolattr.Value(false),
		}}),
		Conditions: listattr.Value([]*ConditionModel{{
			Key:       stringattr.Value("client.tags"),
			Operator:  stringattr.Value("in"),
			ValueJSON: stringattr.Value(`["trusted"]`),
		}}),
	}

	// When
	rule, diags := model.ToPolicyRule(t.Context())

	// Then
	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, "AR1", rule.ID)
	require.EqualValues(t, 2, rule.Version)
	require.JSONEq(t, `["trusted"]`, string(rule.Conditions[0].Value))
}

func TestPolicyRuleModel_ToPolicyRule_rejects_condition_value_arity_mismatch(t *testing.T) {
	// Given
	model := validPolicyRuleModel()
	model.Conditions = listattr.Value([]*ConditionModel{{
		Key:       stringattr.Value("client.tags"),
		Operator:  stringattr.Value("equal"),
		ValueJSON: stringattr.Value(`["trusted"]`),
	}})

	// When
	_, diags := model.ToPolicyRule(t.Context())

	// Then
	require.True(t, diags.HasError())
	require.Contains(t, diags.Errors()[0].Detail(), "scalar JSON value")
}

func TestPolicyRuleModel_ToPolicyRule_rejects_token_exchange_permit_without_grant(t *testing.T) {
	// Given
	model := validPolicyRuleModel()
	model.RuleFamily = stringattr.Value("token_exchange")
	model.ActionKinds = strlistattr.Value([]string{"exchange_token"})

	// When
	_, diags := model.ToPolicyRule(t.Context())

	// Then
	require.True(t, diags.HasError())
	require.Contains(t, diags.Errors()[0].Detail(), "at least one grant")
}

func TestPolicyRuleModel_SetPolicyRule_maps_server_state(t *testing.T) {
	// Given
	model := PolicyRuleModel{ProjectID: stringattr.Value("P1")}
	rule := &infra.PolicyRule{
		ID:          "AR1",
		Version:     3,
		Name:        "rule",
		Description: "description",
		Enabled:     true,
		Conditions: []infra.PolicyRuleCondition{{
			Key:      "client.tags",
			Operator: "in",
			Value:    json.RawMessage("[ \"trusted\" ]"),
		}},
	}

	// When
	diags := model.SetPolicyRule(t.Context(), rule)

	// Then
	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, "P1", model.ProjectID.ValueString())
	require.Equal(t, "AR1", model.ID.ValueString())
	require.EqualValues(t, 3, model.Version.ValueInt64())
	conditions, conditionDiags := model.Conditions.ToSlice(t.Context())
	require.False(t, conditionDiags.HasError(), conditionDiags.Errors())
	require.Equal(t, `["trusted"]`, conditions[0].ValueJSON.ValueString())
}

func validPolicyRuleModel() PolicyRuleModel {
	return PolicyRuleModel{
		ProjectID:         stringattr.Value("P1"),
		Name:              stringattr.Value("rule"),
		Description:       stringattr.Value(""),
		Enabled:           boolattr.Value(true),
		RuleFamily:        stringattr.Value("resource_access"),
		ActionKinds:       strlistattr.Value([]string{"client_access"}),
		Effect:            stringattr.Value("permit"),
		PrincipalType:     stringattr.Value("client"),
		PrincipalSelector: strlistattr.Empty(),
		ResourceTargets:   listattr.Empty[ResourceTargetModel](),
		Grants:            listattr.Empty[GrantModel](),
		Conditions:        listattr.Empty[ConditionModel](),
	}
}
