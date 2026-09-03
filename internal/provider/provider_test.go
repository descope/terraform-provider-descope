package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func Test_Resources_registers_outbound_app(t *testing.T) {
	// Given
	provider := &descopeProvider{}

	// When
	constructors := provider.Resources(context.Background())

	// Then
	for _, constructor := range constructors {
		var response resource.MetadataResponse
		constructor().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "descope"}, &response)
		if response.TypeName == "descope_outbound_app" {
			return
		}
	}
	require.Fail(t, "descope_outbound_app resource is not registered")
}
