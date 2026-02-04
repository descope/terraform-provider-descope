package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const NoProjectID = ""

type Response struct {
	Entity string         `json:"entity"`
	ID     string         `json:"id"`
	Data   map[string]any `json:"data"`
}

type Client struct {
	version       string
	managementKey string
	baseURL       string

	apiClients map[string]*api.Client
	lock       sync.Mutex
}

func NewClient(version, managementKey, baseURL string) *Client {
	return &Client{
		version:       version,
		managementKey: managementKey,
		baseURL:       baseURL,
		apiClients:    map[string]*api.Client{},
	}
}

func (c *Client) Create(ctx context.Context, projectID, entity string, data map[string]any) (*Response, error) {
	httpBody := map[string]any{
		"entity": entity,
		"data":   data,
	}

	tflog.Info(ctx, "Starting CREATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, "/v1/mgmt/infra", httpBody, nil, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished CREATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Read(ctx context.Context, projectID, entity, entityID string) (*Response, error) {
	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting READ request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, "/v1/mgmt/infra", &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished READ request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Update(ctx context.Context, projectID, entity, entityID string, data map[string]any) (*Response, error) {
	httpBody := map[string]any{
		"entity": entity,
		"id":     entityID,
		"data":   data,
	}

	tflog.Info(ctx, "Starting UPDATE request", map[string]any{"body": debugRequest(httpBody)})
	httpRes, err := c.getAPIClient(projectID).DoPutRequest(ctx, "/v1/mgmt/infra", httpBody, nil, c.managementKey)
	if err != nil {
		return nil, err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return nil, err
	}

	tflog.Info(ctx, "Finished UPDATE request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return res, nil
}

func (c *Client) Delete(ctx context.Context, projectID, entity, entityID string) error {
	httpQuery := map[string]string{
		"entity": entity,
		"id":     entityID,
	}

	tflog.Info(ctx, "Starting DELETE request", map[string]any{"query": debugRequest(httpQuery)})
	httpRes, err := c.getAPIClient(projectID).DoDeleteRequest(ctx, "/v1/mgmt/infra", &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey)
	if err != nil {
		return err
	}

	res := &Response{}
	if err := json.Unmarshal([]byte(httpRes.BodyStr), res); err != nil {
		return err
	}

	tflog.Info(ctx, "Finished DELETE request")
	return nil
}

func (c *Client) getAPIClient(projectID string) *api.Client {
	c.lock.Lock()
	defer c.lock.Unlock()

	apiClient, ok := c.apiClients[projectID]
	if !ok {
		apiClient = makeAPIClient(c.version, projectID, c.baseURL)
		c.apiClients[projectID] = apiClient
	}

	return apiClient
}

func makeAPIClient(version, projectID, baseURL string) *api.Client {
	headers := map[string]string{
		"user-agent": makeUserAgent(version),
	}

	params := api.ClientParams{
		ProjectID:            projectID,
		BaseURL:              baseURL,
		CustomDefaultHeaders: headers,
	}

	return api.NewClient(params)
}

func makeUserAgent(version string) string {
	if v := os.Getenv("DESCOPE_USER_AGENT"); v != "" {
		return v
	}
	return fmt.Sprintf("terraform-provider-descope/%s", version)
}
