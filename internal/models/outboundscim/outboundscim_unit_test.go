package outboundscim

import (
	"testing"

	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConfigurationForWrite_merges_write_only_secrets(t *testing.T) {
	// Given
	model := OutboundSCIMConfigurationModel{
		Configuration: types.StringValue(`{"host":"https://scim.example.com","authentication":{"method":"bearerToken"}}`),
		Secrets:       types.StringValue(`{"authentication":{"bearerToken":"secret-value"}}`),
	}

	// When
	configuration, err := model.ConfigurationForWrite()

	// Then
	require.NoError(t, err)
	require.Equal(t, "https://scim.example.com", configuration["host"])
	authentication, ok := configuration["authentication"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "bearerToken", authentication["method"])
	require.Equal(t, "secret-value", authentication["bearerToken"])
}

func TestNormalizeConfigurationForState_removes_internal_and_masked_values(t *testing.T) {
	// Given
	configuration := map[string]any{
		"federatedAppId": "app-1",
		"host":           "https://scim.example.com",
		"authentication": map[string]any{
			"method":      "bearerToken",
			"bearerToken": "MASKED_VALUE",
		},
	}

	// When
	normalized, err := NormalizeConfigurationForState(configuration)

	// Then
	require.NoError(t, err)
	require.JSONEq(t, `{"authentication":{"method":"bearerToken"},"host":"https://scim.example.com"}`, normalized)
}

func TestDecodeConfiguration_rejects_non_object_JSON(t *testing.T) {
	// Given
	raw := `["not", "an", "object"]`

	// When
	_, err := DecodeConfiguration(raw)

	// Then
	require.Error(t, err)
}

func TestDecodeConfiguration_rejects_trailing_JSON(t *testing.T) {
	// Given
	raw := `{} {}`

	// When
	_, err := DecodeConfiguration(raw)

	// Then
	require.Error(t, err)
}

func TestConfigurationForWrite_rejects_secret_fields_in_state_configuration(t *testing.T) {
	// Given
	model := OutboundSCIMConfigurationModel{
		Configuration: types.StringValue(`{"authentication":{"bearerToken":"must-not-enter-state"}}`),
	}

	// When
	_, err := model.ConfigurationForWrite()

	// Then
	require.Error(t, err)
}

func TestConfigurationForWrite_rejects_secret_fields_inside_lists(t *testing.T) {
	// Given
	model := OutboundSCIMConfigurationModel{
		Configuration: types.StringValue(`{"credentials":[{"password":"must-not-enter-state"}]}`),
	}

	// When
	_, err := model.ConfigurationForWrite()

	// Then
	require.Error(t, err)
}

func TestSchema_marks_secrets_write_only_and_sensitive(t *testing.T) {
	// Given
	attribute, ok := Schema.Attributes["secrets"].(resourceschema.StringAttribute)
	require.True(t, ok)

	// When / Then
	require.True(t, attribute.WriteOnly)
	require.True(t, attribute.Sensitive)
}
