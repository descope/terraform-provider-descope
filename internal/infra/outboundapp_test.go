package infra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/stretchr/testify/require"
)

func Test_CreateOutboundApp_posts_typed_request(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, outboundAppCreatePath, r.URL.Path)
		var request descope.CreateOutboundAppRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "example", request.Name)
		require.Equal(t, "secret", request.ClientSecret)
		_, err := w.Write([]byte(`{"app":{"id":"app-1","name":"example"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	app, err := client.CreateOutboundApp(context.Background(), "project-1", &descope.CreateOutboundAppRequest{
		OutboundApp:  descope.OutboundApp{Name: "example"},
		ClientSecret: "secret",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "app-1", app.ID)
}

func Test_LoadOutboundApp_gets_stable_ID_path(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, outboundAppLoadPath+"/app-1", r.URL.Path)
		_, err := w.Write([]byte(`{"app":{"id":"app-1","name":"example"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	app, err := client.LoadOutboundApp(context.Background(), "project-1", "app-1")

	// Then
	require.NoError(t, err)
	require.Equal(t, "example", app.Name)
}

func Test_UpdateOutboundApp_omits_unconfigured_secret(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, outboundAppUpdatePath, r.URL.Path)
		var request struct {
			App struct {
				ID           string  `json:"id"`
				ClientSecret *string `json:"clientSecret"`
			} `json:"app"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "app-1", request.App.ID)
		require.Nil(t, request.App.ClientSecret)
		_, err := w.Write([]byte(`{"app":{"id":"app-1","name":"renamed"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	app, err := client.UpdateOutboundApp(context.Background(), "project-1", OutboundAppUpdateRequest{
		App: &descope.OutboundApp{ID: "app-1", Name: "renamed"},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "renamed", app.Name)
}

func Test_UpdateOutboundApp_sends_configured_secret(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			App struct {
				ClientSecret string `json:"clientSecret"`
			} `json:"app"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "replacement", request.App.ClientSecret)
		_, err := w.Write([]byte(`{"app":{"id":"app-1","name":"example"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)
	secret := "replacement"

	// When
	_, err := client.UpdateOutboundApp(context.Background(), "project-1", OutboundAppUpdateRequest{
		App:          &descope.OutboundApp{ID: "app-1", Name: "example"},
		ClientSecret: &secret,
	})

	// Then
	require.NoError(t, err)
}

func Test_DeleteOutboundApp_posts_ID(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, outboundAppDeletePath, r.URL.Path)
		var request map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "app-1", request["id"])
		_, err := w.Write([]byte(`{}`))
		require.NoError(t, err)
	}))
	defer server.Close()
	client := NewClient("test", "management-key", server.URL)

	// When
	err := client.DeleteOutboundApp(context.Background(), "project-1", "app-1")

	// Then
	require.NoError(t, err)
}
