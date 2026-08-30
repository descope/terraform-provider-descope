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
	ID   string         `json:"id"`
	Data map[string]any `json:"data"`
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
	if _, err := c.getAPIClient(projectID).DoDeleteRequest(ctx, "/v1/mgmt/infra", &api.HTTPRequest{QueryParams: httpQuery}, c.managementKey); err != nil {
		return err
	}

	tflog.Info(ctx, "Finished DELETE request")
	return nil
}

// Post/PostData/Get/Del are thin helpers for dedicated per-resource endpoints. Unlike the entity
// methods above they take a full path and don't wrap the body in an {entity,id,data} envelope —
// the request/response bodies are the resource's JSON directly.

func (c *Client) Post(ctx context.Context, projectID, path string, body map[string]any) error {
	tflog.Info(ctx, "Starting POST request", map[string]any{"path": path, "body": debugRequest(body)})
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, path, body, nil, c.managementKey)
	return err
}

func (c *Client) PostData(ctx context.Context, projectID, path string, body map[string]any) (map[string]any, error) {
	tflog.Info(ctx, "Starting POST request", map[string]any{"path": path, "body": debugRequest(body)})
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, path, body, nil, c.managementKey)
	if err != nil {
		return nil, err
	}
	data := map[string]any{}
	if httpRes.BodyStr != "" {
		if err := json.Unmarshal([]byte(httpRes.BodyStr), &data); err != nil {
			return nil, err
		}
	}
	tflog.Info(ctx, "Finished POST request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return data, nil
}

func (c *Client) Get(ctx context.Context, projectID, path string, query map[string]string) (map[string]any, error) {
	var req *api.HTTPRequest
	if query != nil {
		req = &api.HTTPRequest{QueryParams: query}
	}
	tflog.Info(ctx, "Starting GET request", map[string]any{"path": path, "query": debugRequest(query)})
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, path, req, c.managementKey)
	if err != nil {
		return nil, err
	}
	data := map[string]any{}
	if httpRes.BodyStr != "" {
		if err := json.Unmarshal([]byte(httpRes.BodyStr), &data); err != nil {
			return nil, err
		}
	}
	tflog.Info(ctx, "Finished GET request", map[string]any{"response": debugResponse(httpRes.BodyStr)})
	return data, nil
}

func (c *Client) Del(ctx context.Context, projectID, path string, query map[string]string) error {
	var req *api.HTTPRequest
	if query != nil {
		req = &api.HTTPRequest{QueryParams: query}
	}
	tflog.Info(ctx, "Starting DELETE request", map[string]any{"path": path, "query": debugRequest(query)})
	_, err := c.getAPIClient(projectID).DoDeleteRequest(ctx, path, req, c.managementKey)
	return err
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
