package infra

import (
	"errors"

	"github.com/descope/go-sdk/descope"
)

// ErrNotFound is a sentinel for provider-side lookups (e.g. the role/permission list
// scans) so they participate in the same not-found detection as backend 404 responses.
var ErrNotFound = errors.New("entity not found")

// AsNotFoundError reports whether err means the entity doesn't exist on the backend,
// either as a not-found response from the server or as the provider-side sentinel.
func AsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound) || descope.AsError(err).IsNotFound()
}

// AsValidationError returns the message of backend validation failures, i.e., bad request
// errors that carry a specific message, so resources can present them as configuration
// errors instead of raw wire error strings.
func AsValidationError(err error) (failure string, ok bool) {
	if de := descope.AsError(err); de != nil && de.Message != "" && de.IsBadRequest() {
		return de.Message, true
	}
	return
}
