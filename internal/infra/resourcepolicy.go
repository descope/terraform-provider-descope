package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	createResourcePolicyPath = "/v1/mgmt/resourcepolicy/create"
	loadResourcePolicyPath   = "/v1/mgmt/resourcepolicy/app/load"
	updateResourcePolicyPath = "/v1/mgmt/resourcepolicy/update"
	deleteResourcePolicyPath = "/v1/mgmt/resourcepolicy/delete"
)

var ErrResourcePolicyNotFound = errors.New("resource policy not found")

type ResourcePolicyIdentity struct {
	ApplicationID string `json:"thirdPartyApplicationId"`
	ResourceID    string `json:"resourceId"`
}

type ResourcePolicy struct {
	ApplicationID      string   `json:"thirdPartyApplicationId"`
	ResourceID         string   `json:"resourceId"`
	UserAccessScopes   []string `json:"userAccessScopes"`
	ClientAccessScopes []string `json:"clientAccessScopes"`
	AllUserScopes      bool     `json:"allUserScopes"`
	AllClientScopes    bool     `json:"allClientScopes"`
}

type resourcePolicyResponse struct {
	ResourcePolicy ResourcePolicy `json:"resourcePolicy"`
}

type resourcePoliciesResponse struct {
	ResourcePolicies []ResourcePolicy `json:"resourcePolicies"`
}

func (c *Client) CreateResourcePolicy(ctx context.Context, projectID string, policy ResourcePolicy) (*ResourcePolicy, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, createResourcePolicyPath, policy, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("create resource policy: %w", err)
	}

	var response resourcePolicyResponse
	if err := json.Unmarshal([]byte(httpRes.BodyStr), &response); err != nil {
		return nil, fmt.Errorf("decode created resource policy: %w", err)
	}
	return &response.ResourcePolicy, nil
}

func (c *Client) ReadResourcePolicy(ctx context.Context, projectID string, identity ResourcePolicyIdentity) (*ResourcePolicy, error) {
	request := struct {
		ApplicationID string `json:"thirdPartyApplicationId"`
	}{ApplicationID: identity.ApplicationID}
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, loadResourcePolicyPath, request, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("load resource policies: %w", err)
	}

	var response resourcePoliciesResponse
	if err := json.Unmarshal([]byte(httpRes.BodyStr), &response); err != nil {
		return nil, fmt.Errorf("decode resource policies: %w", err)
	}
	for i := range response.ResourcePolicies {
		if response.ResourcePolicies[i].ResourceID == identity.ResourceID {
			return &response.ResourcePolicies[i], nil
		}
	}
	return nil, fmt.Errorf("%w: application %q resource %q", ErrResourcePolicyNotFound, identity.ApplicationID, identity.ResourceID)
}

func (c *Client) UpdateResourcePolicy(ctx context.Context, projectID string, policy ResourcePolicy) (*ResourcePolicy, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, updateResourcePolicyPath, policy, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("update resource policy: %w", err)
	}

	var response resourcePolicyResponse
	if err := json.Unmarshal([]byte(httpRes.BodyStr), &response); err != nil {
		return nil, fmt.Errorf("decode updated resource policy: %w", err)
	}
	return &response.ResourcePolicy, nil
}

func (c *Client) DeleteResourcePolicy(ctx context.Context, projectID string, identity ResourcePolicyIdentity) error {
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, deleteResourcePolicyPath, identity, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("delete resource policy: %w", err)
	}
	return nil
}
