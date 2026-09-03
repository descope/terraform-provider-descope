package infra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tenantNotFoundErrorCode = "E112202"

type Tenant struct {
	ID                      string          `json:"id"`
	Name                    string          `json:"name"`
	SelfProvisioningDomains []string        `json:"selfProvisioningDomains"`
	CustomAttributes        json.RawMessage `json:"customAttributes"`
	AuthType                string          `json:"authType"`
	Disabled                bool            `json:"disabled"`
	EnforceSSO              bool            `json:"enforceSSO"`
	EnforceSSOExclusions    []string        `json:"enforceSSOExclusions"`
	FederatedApplicationIDs []string        `json:"federatedAppIds"`
	Parent                  string          `json:"parent"`
	RoleInheritance         string          `json:"roleInheritance"`
	IDJagSettings           json.RawMessage `json:"idJagSettings"`
	IDJagEnabled            bool            `json:"idJagEnabled"`
}

type TenantCreateRequest struct {
	ID                      string   `json:"id,omitempty"`
	Name                    string   `json:"name"`
	SelfProvisioningDomains []string `json:"selfProvisioningDomains"`
	Disabled                bool     `json:"disabled"`
	EnforceSSO              bool     `json:"enforceSSO"`
	EnforceSSOExclusions    []string `json:"enforceSSOExclusions"`
	FederatedApplicationIDs []string `json:"federatedAppIds"`
	Parent                  string   `json:"parent,omitempty"`
	RoleInheritance         string   `json:"roleInheritance,omitempty"`
}

type TenantUpdateRequest struct {
	ID                      string          `json:"id"`
	Name                    string          `json:"name"`
	SelfProvisioningDomains []string        `json:"selfProvisioningDomains"`
	CustomAttributes        json.RawMessage `json:"customAttributes,omitempty"`
	AuthType                string          `json:"authType,omitempty"`
	Disabled                bool            `json:"disabled"`
	EnforceSSO              bool            `json:"enforceSSO"`
	EnforceSSOExclusions    []string        `json:"enforceSSOExclusions"`
	FederatedApplicationIDs []string        `json:"federatedAppIds"`
	RoleInheritance         string          `json:"roleInheritance,omitempty"`
	IDJagSettings           json.RawMessage `json:"idJagSettings,omitempty"`
	IDJagEnabled            bool            `json:"idJagEnabled"`
}

type TenantDeleteRequest struct {
	ID      string `json:"id"`
	Cascade bool   `json:"cascade"`
}

func (c *Client) CreateTenant(ctx context.Context, projectID string, req TenantCreateRequest) (string, error) {
	tflog.Info(ctx, "Starting tenant CREATE request", map[string]any{"body": debugRequest(req)})
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, "/v1/mgmt/tenant/create", req, nil, c.managementKey)
	if err != nil {
		return "", fmt.Errorf("create tenant: %w", err)
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), &res); err != nil {
		return "", fmt.Errorf("decode create tenant response: %w", err)
	}

	tflog.Info(ctx, "Finished tenant CREATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res.ID, nil
}

func (c *Client) ReadTenant(ctx context.Context, projectID, tenantID string) (*Tenant, error) {
	tflog.Info(ctx, "Starting tenant READ request", map[string]any{"id": tenantID})
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, "/v1/mgmt/tenant", &api.HTTPRequest{
		QueryParams: map[string]string{"id": tenantID},
	}, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("read tenant: %w", err)
	}

	var tenant Tenant
	if err := json.Unmarshal([]byte(httpRes.BodyStr), &tenant); err != nil {
		return nil, fmt.Errorf("decode tenant response: %w", err)
	}

	tflog.Info(ctx, "Finished tenant READ request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return &tenant, nil
}

func (c *Client) UpdateTenant(ctx context.Context, projectID string, req TenantUpdateRequest) error {
	tflog.Info(ctx, "Starting tenant UPDATE request", map[string]any{"body": debugRequest(req)})
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, "/v1/mgmt/tenant/update", req, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("update tenant: %w", err)
	}

	tflog.Info(ctx, "Finished tenant UPDATE request")
	return nil
}

func (c *Client) DeleteTenant(ctx context.Context, projectID, tenantID string) error {
	req := TenantDeleteRequest{ID: tenantID}
	tflog.Info(ctx, "Starting tenant DELETE request", map[string]any{"body": debugRequest(req)})
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, "/v1/mgmt/tenant/delete", req, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("delete tenant: %w", err)
	}

	tflog.Info(ctx, "Finished tenant DELETE request")
	return nil
}

func IsTenantNotFound(err error) bool {
	var descopeErr *descope.Error
	return errors.As(err, &descopeErr) && descopeErr.Code == tenantNotFoundErrorCode
}
