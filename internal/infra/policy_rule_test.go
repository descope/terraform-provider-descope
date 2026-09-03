package infra

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_CreatePolicyRule_posts_typed_request(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/mgmt/policies/rule/create", r.URL.Path)

		var request createPolicyRuleRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "rule", request.PolicyRule.Name)
		require.Equal(t, []string{"client_access"}, request.PolicyRule.ActionKinds)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"policyRule":{"id":"AR1","version":"1","name":"rule","enabled":true}}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	rule, err := client.CreatePolicyRule(t.Context(), "P1", PolicyRule{
		Name:        "rule",
		Enabled:     true,
		ActionKinds: []string{"client_access"},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "AR1", rule.ID)
	require.EqualValues(t, 1, rule.Version)
}

func TestClient_LoadPolicyRule_loads_by_id(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/mgmt/policies/rule/load", r.URL.Path)

		var request policyRuleIDRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "AR1", request.ID)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"policyRule":{"id":"AR1","version":"2","name":"rule"}}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	rule, err := client.LoadPolicyRule(t.Context(), "P1", "AR1")

	// Then
	require.NoError(t, err)
	require.Equal(t, "AR1", rule.ID)
	require.EqualValues(t, 2, rule.Version)
}

func TestClient_UpdatePolicyRule_sends_version(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/mgmt/policies/rule/update", r.URL.Path)

		var request updatePolicyRuleRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "AR1", request.PolicyRule.ID)
		require.EqualValues(t, 2, request.PolicyRule.Version)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"policyRule":{"id":"AR1","version":"3","name":"updated"}}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	rule, err := client.UpdatePolicyRule(t.Context(), "P1", PolicyRule{ID: "AR1", Version: 2, Name: "updated"})

	// Then
	require.NoError(t, err)
	require.EqualValues(t, 3, rule.Version)
}

func TestClient_DeletePolicyRule_posts_id(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/v1/mgmt/policies/rule/delete", r.URL.Path)

		var request policyRuleIDRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		require.Equal(t, "AR1", request.ID)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"id":"AR1"}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := NewClient("test", "management-key", server.URL)

	// When
	err := client.DeletePolicyRule(t.Context(), "P1", "AR1")

	// Then
	require.NoError(t, err)
}
