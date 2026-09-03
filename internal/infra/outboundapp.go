package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/descope/go-sdk/descope"
)

const (
	outboundAppCreatePath = "/v1/mgmt/outbound/app/create"
	outboundAppUpdatePath = "/v1/mgmt/outbound/app/update"
	outboundAppDeletePath = "/v1/mgmt/outbound/app/delete"
	outboundAppLoadPath   = "/v1/mgmt/outbound/app"
)

type outboundAppResponse struct {
	App *descope.OutboundApp `json:"app"`
}

type outboundAppUpdate struct {
	descope.OutboundApp
	ClientSecret *string `json:"clientSecret,omitempty"`
}

type outboundAppUpdateRequest struct {
	App outboundAppUpdate `json:"app"`
}

type OutboundAppUpdateRequest struct {
	App          *descope.OutboundApp
	ClientSecret *string
}

func (c *Client) CreateOutboundApp(ctx context.Context, projectID string, request *descope.CreateOutboundAppRequest) (*descope.OutboundApp, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundAppCreatePath, request, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("create outbound app: %w", err)
	}
	return unmarshalOutboundApp(httpRes.BodyStr)
}

func (c *Client) LoadOutboundApp(ctx context.Context, projectID, id string) (*descope.OutboundApp, error) {
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, path.Join(outboundAppLoadPath, id), nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("load outbound app %s: %w", id, err)
	}
	return unmarshalOutboundApp(httpRes.BodyStr)
}

func (c *Client) UpdateOutboundApp(ctx context.Context, projectID string, request OutboundAppUpdateRequest) (*descope.OutboundApp, error) {
	body := outboundAppUpdateRequest{App: outboundAppUpdate{OutboundApp: *request.App, ClientSecret: request.ClientSecret}}
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundAppUpdatePath, body, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("update outbound app %s: %w", request.App.ID, err)
	}
	return unmarshalOutboundApp(httpRes.BodyStr)
}

func (c *Client) DeleteOutboundApp(ctx context.Context, projectID, id string) error {
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundAppDeletePath, map[string]string{"id": id}, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("delete outbound app %s: %w", id, err)
	}
	return nil
}

func unmarshalOutboundApp(body string) (*descope.OutboundApp, error) {
	var response outboundAppResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, fmt.Errorf("decode outbound app response: %w", err)
	}
	if response.App == nil {
		return nil, fmt.Errorf("decode outbound app response: missing app")
	}
	return response.App, nil
}
