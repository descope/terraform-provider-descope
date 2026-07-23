package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestResources_registers_outbound_SCIM_configuration(t *testing.T) {
	// Given
	provider := &descopeProvider{}

	// When
	constructors := provider.Resources(context.Background())

	// Then
	var found bool
	for _, constructor := range constructors {
		var response resource.MetadataResponse
		constructor().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "descope"}, &response)
		if response.TypeName == "descope_outbound_scim_configuration" {
			found = true
			break
		}
	}
	require.True(t, found)
}
