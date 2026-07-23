package infra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOutboundSCIMConfiguration_posts_typed_management_request(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/mgmt/outbound/scim/create", r.URL.Path)
		var request OutboundSCIMWriteRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "app-1", request.AppID)
		require.Equal(t, "https://scim.example.com", request.Configuration["host"])
		_, err := w.Write([]byte(`{"configuration":{"appId":"app-1","configuration":{"host":"https://scim.example.com"},"enabled":true,"version":"3"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	configuration, err := client.CreateOutboundSCIMConfiguration(context.Background(), "project-1", OutboundSCIMWriteRequest{
		AppID:         "app-1",
		Configuration: map[string]any{"host": "https://scim.example.com"},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "app-1", configuration.AppID)
	require.Equal(t, int64(3), configuration.Version)
}

func TestUpdateOutboundSCIMConfiguration_posts_version(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/mgmt/outbound/scim/update", r.URL.Path)
		var request OutboundSCIMWriteRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, int64(8), request.Version)
		_, err := w.Write([]byte(`{"configuration":{"appId":"app-1","configuration":{},"version":"9"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	configuration, err := client.UpdateOutboundSCIMConfiguration(context.Background(), "project-1", OutboundSCIMWriteRequest{
		AppID:         "app-1",
		Configuration: map[string]any{},
		Version:       8,
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, int64(9), configuration.Version)
}

func TestLoadOutboundSCIMConfiguration_gets_app_path(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/mgmt/outbound/scim/app-1", r.URL.Path)
		_, err := w.Write([]byte(`{"configuration":{"appId":"app-1","configuration":{},"enabled":true,"lastExportTime":10,"lastProcessingTime":11,"failures":2,"version":"4"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	configuration, err := client.LoadOutboundSCIMConfiguration(context.Background(), "project-1", "app-1")

	// Then
	require.NoError(t, err)
	require.True(t, configuration.Enabled)
	require.Equal(t, int64(10), configuration.LastExportTime)
	require.Equal(t, int64(11), configuration.LastProcessingTime)
	require.Equal(t, int64(2), configuration.Failures)
}

func TestSetOutboundSCIMEnabled_posts_enabled_state(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/mgmt/outbound/scim/enabled/set", r.URL.Path)
		var request OutboundSCIMEnabledRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.False(t, request.Enabled)
		_, err := w.Write([]byte(`{"configuration":{"appId":"app-1","configuration":{},"enabled":false,"version":"4"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	configuration, err := client.SetOutboundSCIMEnabled(context.Background(), "project-1", OutboundSCIMEnabledRequest{AppID: "app-1", Enabled: false})

	// Then
	require.NoError(t, err)
	require.False(t, configuration.Enabled)
}

func TestDeleteOutboundSCIMConfiguration_posts_app_id(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/mgmt/outbound/scim/delete", r.URL.Path)
		var request OutboundSCIMDeleteRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "app-1", request.AppID)
		_, err := w.Write([]byte(`{}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	err := client.DeleteOutboundSCIMConfiguration(context.Background(), "project-1", "app-1")

	// Then
	require.NoError(t, err)
}
