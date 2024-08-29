package infra

import "github.com/descope/go-sdk/descope"

const (
	errCodeValidationError = "E113007"
)

func AsValidationError(err error) (failure string, ok bool) {
	if err, ok := err.(*descope.Error); ok && err.Code == errCodeValidationError && err.Message != "" {
		return err.Message, true
	}
	return
}
