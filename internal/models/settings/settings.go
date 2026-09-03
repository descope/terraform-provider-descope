package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

// Marks the delivery service as the built-in Descope service when its block is absent, which also resets any custom templates the service had.
func useDescopeService[T any](service objattr.Type[T], data map[string]any, key string) {
	if !service.IsSet() {
		data[key] = helpers.DescopeConnector
	}
}
