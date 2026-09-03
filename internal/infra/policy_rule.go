package infra

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	policyRuleCreatePath = "/v1/mgmt/policies/rule/create"
	policyRuleLoadPath   = "/v1/mgmt/policies/rule/load"
	policyRuleUpdatePath = "/v1/mgmt/policies/rule/update"
	policyRuleDeletePath = "/v1/mgmt/policies/rule/delete"
)

type PolicyRule struct {
	ID                string                     `json:"id,omitempty"`
	Version           int64                      `json:"version,string,omitempty"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Enabled           bool                       `json:"enabled"`
	RuleFamily        string                     `json:"ruleFamily"`
	ActionKinds       []string                   `json:"actionKinds"`
	Effect            string                     `json:"effect"`
	PrincipalType     string                     `json:"principalType"`
	PrincipalSelector []string                   `json:"principalSelector"`
	ResourceTargets   []PolicyRuleResourceTarget `json:"resourceTargets"`
	Grants            []PolicyRuleGrant          `json:"grants"`
	Conditions        []PolicyRuleCondition      `json:"conditions"`
	CedarText         string                     `json:"cedarText,omitempty"`
	CreatedTime       int64                      `json:"createdTime,omitempty"`
	ModifiedTime      int64                      `json:"modifiedTime,omitempty"`
}

type PolicyRuleResourceTarget struct {
	Type      string   `json:"type"`
	AllOfType bool     `json:"allOfType"`
	IDs       []string `json:"ids"`
}

type PolicyRuleGrant struct {
	Scopes           []string `json:"scopes"`
	AllowedAudiences []string `json:"allowedAudiences"`
	AllScopes        bool     `json:"allScopes"`
}

type PolicyRuleCondition struct {
	Key      string          `json:"key"`
	Operator string          `json:"operator"`
	Value    json.RawMessage `json:"value"`
}

type createPolicyRuleRequest struct {
	PolicyRule PolicyRule `json:"policyRule"`
}

type updatePolicyRuleRequest struct {
	PolicyRule PolicyRule `json:"policyRule"`
}

type policyRuleIDRequest struct {
	ID string `json:"id"`
}

type policyRuleResponse struct {
	PolicyRule PolicyRule `json:"policyRule"`
}

func (c *Client) CreatePolicyRule(ctx context.Context, projectID string, rule PolicyRule) (*PolicyRule, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, policyRuleCreatePath, createPolicyRuleRequest{PolicyRule: rule}, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("create policy rule: %w", err)
	}

	return decodePolicyRuleResponse(httpRes.BodyStr, "create")
}

func (c *Client) LoadPolicyRule(ctx context.Context, projectID, id string) (*PolicyRule, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, policyRuleLoadPath, policyRuleIDRequest{ID: id}, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("load policy rule %s: %w", id, err)
	}

	return decodePolicyRuleResponse(httpRes.BodyStr, "load")
}

func (c *Client) UpdatePolicyRule(ctx context.Context, projectID string, rule PolicyRule) (*PolicyRule, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, policyRuleUpdatePath, updatePolicyRuleRequest{PolicyRule: rule}, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("update policy rule %s: %w", rule.ID, err)
	}

	return decodePolicyRuleResponse(httpRes.BodyStr, "update")
}

func (c *Client) DeletePolicyRule(ctx context.Context, projectID, id string) error {
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, policyRuleDeletePath, policyRuleIDRequest{ID: id}, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("delete policy rule %s: %w", id, err)
	}
	return nil
}

func decodePolicyRuleResponse(body, operation string) (*PolicyRule, error) {
	var response policyRuleResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, fmt.Errorf("decode %s policy rule response: %w", operation, err)
	}
	if response.PolicyRule.ID == "" {
		return nil, fmt.Errorf("decode %s policy rule response: missing policy rule id", operation)
	}
	return &response.PolicyRule, nil
}
