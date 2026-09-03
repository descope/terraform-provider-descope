package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func Test_OutboundAppResource_sets_resource_type_name(t *testing.T) {
	// Given
	var response resource.MetadataResponse

	// When
	NewOutboundAppResource().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "descope"}, &response)

	// Then
	require.Equal(t, "descope_outbound_app", response.TypeName)
}

func Test_IsOutboundAppNotFound_returns_true_for_HTTP_404(t *testing.T) {
	// Given
	err := (&descope.Error{Code: "E4005"}).WithInfo(descope.ErrorInfoKeys.HTTPResponseStatusCode, http.StatusNotFound)

	// When
	found := isOutboundAppNotFound(err)

	// Then
	require.True(t, found)
}
