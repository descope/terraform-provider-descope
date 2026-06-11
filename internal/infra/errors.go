package infra

import "github.com/descope/go-sdk/descope"

func AsValidationError(err error) (failure string, ok bool) {
	if err, ok := err.(*descope.Error); ok && err.Message != "" {
		if err.Code == "E113007" || err.Code == "E113011" {
			return err.Message, true
		}
	}
	return
}
