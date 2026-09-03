package infra

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_ResourcePolicyCRUD_uses_management_routes(t *testing.T) {
	t.Parallel()

	// Given
	type capturedRequest struct {
		Method string
		Path   string
		Body   map[string]any
	}

	var (
		mu       sync.Mutex
		requests []capturedRequest
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		requests = append(requests, capturedRequest{Method: r.Method, Path: r.URL.Path, Body: body})
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/mgmt/resourcepolicy/create", "/v1/mgmt/resourcepolicy/update":
			_, _ = w.Write([]byte(`{"resourcePolicy":{"thirdPartyApplicationId":"AP1","resourceId":"RS1","userAccessScopes":["read"],"clientAccessScopes":["write"],"allUserScopes":false,"allClientScopes":true}}`))
		case "/v1/mgmt/resourcepolicy/app/load":
			_, _ = w.Write([]byte(`{"resourcePolicies":[{"thirdPartyApplicationId":"AP1","resourceId":"RS2"},{"thirdPartyApplicationId":"AP1","resourceId":"RS1","userAccessScopes":["read"],"clientAccessScopes":["write"],"allUserScopes":false,"allClientScopes":true}]}`))
		case "/v1/mgmt/resourcepolicy/delete":
			_, _ = w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("test", "management-key", server.URL)
	policy := ResourcePolicy{
		ApplicationID:      "AP1",
		ResourceID:         "RS1",
		UserAccessScopes:   []string{"read"},
		ClientAccessScopes: []string{"write"},
		AllClientScopes:    true,
	}
	identity := ResourcePolicyIdentity{ApplicationID: "AP1", ResourceID: "RS1"}

	// When
	created, err := client.CreateResourcePolicy(context.Background(), "P1", policy)
	require.NoError(t, err)
	loaded, err := client.ReadResourcePolicy(context.Background(), "P1", identity)
	require.NoError(t, err)
	updated, err := client.UpdateResourcePolicy(context.Background(), "P1", policy)
	require.NoError(t, err)
	err = client.DeleteResourcePolicy(context.Background(), "P1", identity)

	// Then
	require.NoError(t, err)
	require.Equal(t, policy, *created)
	require.Equal(t, policy, *loaded)
	require.Equal(t, policy, *updated)
	require.Equal(t, []capturedRequest{
		{Method: http.MethodPost, Path: "/v1/mgmt/resourcepolicy/create", Body: map[string]any{
			"thirdPartyApplicationId": "AP1", "resourceId": "RS1",
			"userAccessScopes": []any{"read"}, "clientAccessScopes": []any{"write"},
			"allUserScopes": false, "allClientScopes": true,
		}},
		{Method: http.MethodPost, Path: "/v1/mgmt/resourcepolicy/app/load", Body: map[string]any{"thirdPartyApplicationId": "AP1"}},
		{Method: http.MethodPost, Path: "/v1/mgmt/resourcepolicy/update", Body: map[string]any{
			"thirdPartyApplicationId": "AP1", "resourceId": "RS1",
			"userAccessScopes": []any{"read"}, "clientAccessScopes": []any{"write"},
			"allUserScopes": false, "allClientScopes": true,
		}},
		{Method: http.MethodPost, Path: "/v1/mgmt/resourcepolicy/delete", Body: map[string]any{
			"thirdPartyApplicationId": "AP1", "resourceId": "RS1",
		}},
	}, requests)
}

func TestClient_ReadResourcePolicy_returns_not_found_when_application_has_no_match(t *testing.T) {
	t.Parallel()

	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"resourcePolicies":[]}`))
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	_, err := client.ReadResourcePolicy(context.Background(), "P1", ResourcePolicyIdentity{
		ApplicationID: "AP1",
		ResourceID:    "RS1",
	})

	// Then
	require.ErrorIs(t, err, ErrResourcePolicyNotFound)
}
