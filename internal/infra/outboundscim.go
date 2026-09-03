package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const outboundSCIMPath = "/v1/mgmt/outbound/scim"

type OutboundSCIMConfiguration struct {
	AppID              string         `json:"appId"`
	Configuration      map[string]any `json:"configuration"`
	Enabled            bool           `json:"enabled"`
	LastExportTime     int64          `json:"lastExportTime"`
	LastProcessingTime int64          `json:"lastProcessingTime"`
	Failures           int64          `json:"failures"`
	Version            int64          `json:"-"`
}

type OutboundSCIMWriteRequest struct {
	AppID         string         `json:"appId"`
	Configuration map[string]any `json:"configuration"`
	Version       int64          `json:"version,omitempty"`
}

type OutboundSCIMEnabledRequest struct {
	AppID   string `json:"appId"`
	Enabled bool   `json:"enabled"`
}

type OutboundSCIMDeleteRequest struct {
	AppID string `json:"appId"`
}

func (c *Client) CreateOutboundSCIMConfiguration(ctx context.Context, projectID string, request OutboundSCIMWriteRequest) (*OutboundSCIMConfiguration, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundSCIMPath+"/create", request, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("create outbound SCIM configuration: %w", err)
	}
	return decodeOutboundSCIMResponse(httpRes.BodyStr)
}

func (c *Client) UpdateOutboundSCIMConfiguration(ctx context.Context, projectID string, request OutboundSCIMWriteRequest) (*OutboundSCIMConfiguration, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundSCIMPath+"/update", request, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("update outbound SCIM configuration: %w", err)
	}
	return decodeOutboundSCIMResponse(httpRes.BodyStr)
}

func (c *Client) LoadOutboundSCIMConfiguration(ctx context.Context, projectID, appID string) (*OutboundSCIMConfiguration, error) {
	httpRes, err := c.getAPIClient(projectID).DoGetRequest(ctx, outboundSCIMPath+"/"+url.PathEscape(appID), nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("load outbound SCIM configuration: %w", err)
	}
	return decodeOutboundSCIMResponse(httpRes.BodyStr)
}

func (c *Client) SetOutboundSCIMEnabled(ctx context.Context, projectID string, request OutboundSCIMEnabledRequest) (*OutboundSCIMConfiguration, error) {
	httpRes, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundSCIMPath+"/enabled/set", request, nil, c.managementKey)
	if err != nil {
		return nil, fmt.Errorf("set outbound SCIM enabled state: %w", err)
	}
	return decodeOutboundSCIMResponse(httpRes.BodyStr)
}

func (c *Client) DeleteOutboundSCIMConfiguration(ctx context.Context, projectID, appID string) error {
	_, err := c.getAPIClient(projectID).DoPostRequest(ctx, outboundSCIMPath+"/delete", OutboundSCIMDeleteRequest{AppID: appID}, nil, c.managementKey)
	if err != nil {
		return fmt.Errorf("delete outbound SCIM configuration: %w", err)
	}
	return nil
}

func decodeOutboundSCIMResponse(body string) (*OutboundSCIMConfiguration, error) {
	var response struct {
		Configuration json.RawMessage `json:"configuration"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, fmt.Errorf("decode outbound SCIM response: %w", err)
	}
	if len(response.Configuration) == 0 || string(response.Configuration) == "null" {
		return nil, fmt.Errorf("decode outbound SCIM response: missing configuration")
	}

	var wire struct {
		AppID              string          `json:"appId"`
		Configuration      map[string]any  `json:"configuration"`
		Enabled            bool            `json:"enabled"`
		LastExportTime     int64           `json:"lastExportTime"`
		LastProcessingTime int64           `json:"lastProcessingTime"`
		Failures           int64           `json:"failures"`
		Version            json.RawMessage `json:"version"`
	}
	if err := json.Unmarshal(response.Configuration, &wire); err != nil {
		return nil, fmt.Errorf("decode outbound SCIM configuration: %w", err)
	}
	version, err := decodeProtoInt64(wire.Version)
	if err != nil {
		return nil, fmt.Errorf("decode outbound SCIM version: %w", err)
	}
	return &OutboundSCIMConfiguration{
		AppID:              wire.AppID,
		Configuration:      wire.Configuration,
		Enabled:            wire.Enabled,
		LastExportTime:     wire.LastExportTime,
		LastProcessingTime: wire.LastProcessingTime,
		Failures:           wire.Failures,
		Version:            version,
	}, nil
}

func decodeProtoInt64(raw json.RawMessage) (int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, nil
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strconv.ParseInt(text, 10, 64)
	}
	var number int64
	if err := json.Unmarshal(raw, &number); err != nil {
		return 0, err
	}
	return number, nil
}
