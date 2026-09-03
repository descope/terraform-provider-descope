package resources

import (
	"net/http"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/stretchr/testify/require"
)

func TestIsNotFound_returns_true_for_HTTP_404(t *testing.T) {
	// Given
	err := (&descope.Error{Code: "E4005"}).WithInfo(descope.ErrorInfoKeys.HTTPResponseStatusCode, http.StatusNotFound)

	// When
	found := isNotFound(err)

	// Then
	require.True(t, found)
}
