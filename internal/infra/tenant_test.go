package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/stretchr/testify/require"
)

func TestCreateTenant_whenRequestSucceeds(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/v1/mgmt/tenant/create", req.URL.Path)
		require.Equal(t, "Bearer project-id:management-key", req.Header.Get("Authorization"))

		var body TenantCreateRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "Tenant name", body.Name)
		require.Equal(t, []string{"example.com"}, body.SelfProvisioningDomains)

		_, err := w.Write([]byte(`{"id":"tenant-id"}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	id, err := client.CreateTenant(context.Background(), "project-id", TenantCreateRequest{
		Name:                    "Tenant name",
		SelfProvisioningDomains: []string{"example.com"},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "tenant-id", id)
}

func TestReadTenant_whenRequestSucceeds(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/v1/mgmt/tenant", req.URL.Path)
		require.Equal(t, "tenant-id", req.URL.Query().Get("id"))

		_, err := w.Write([]byte(`{"id":"tenant-id","name":"Tenant name","selfProvisioningDomains":["example.com"],"enforceSSO":true}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	tenant, err := client.ReadTenant(context.Background(), "project-id", "tenant-id")

	// Then
	require.NoError(t, err)
	require.Equal(t, "tenant-id", tenant.ID)
	require.Equal(t, "Tenant name", tenant.Name)
	require.Equal(t, []string{"example.com"}, tenant.SelfProvisioningDomains)
	require.True(t, tenant.EnforceSSO)
}

func TestUpdateTenant_whenRequestSucceeds(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/v1/mgmt/tenant/update", req.URL.Path)

		var body TenantUpdateRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "tenant-id", body.ID)
		require.Equal(t, "Updated tenant", body.Name)
		require.JSONEq(t, `{"department":"engineering"}`, string(body.CustomAttributes))

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	err := client.UpdateTenant(context.Background(), "project-id", TenantUpdateRequest{
		ID:               "tenant-id",
		Name:             "Updated tenant",
		CustomAttributes: json.RawMessage(`{"department":"engineering"}`),
	})

	// Then
	require.NoError(t, err)
}

func TestDeleteTenant_whenRequestSucceeds(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/v1/mgmt/tenant/delete", req.URL.Path)

		var body TenantDeleteRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "tenant-id", body.ID)
		require.False(t, body.Cascade)

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	err := client.DeleteTenant(context.Background(), "project-id", "tenant-id")

	// Then
	require.NoError(t, err)
}

func TestIsTenantNotFound_whenErrorIsWrapped(t *testing.T) {
	// Given
	err := fmt.Errorf("wrapped: %w", &descope.Error{Code: tenantNotFoundErrorCode})

	// When
	notFound := IsTenantNotFound(err)

	// Then
	require.True(t, notFound)
}
