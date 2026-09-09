package helpers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DescopeConnector = "Descope"
	DescopeTemplate  = "System"
)

// Provider references come back from the server as `type:id` values, and an empty string means the built-in Descope delivery service.
func SetServiceConnectorID(s *types.String) {
	value := s.ValueString()
	if value == "" {
		*s = types.StringValue(DescopeConnector)
	} else if _, id, found := strings.Cut(value, ":"); found {
		*s = types.StringValue(id)
	}
}

// Pointer receiver interface for model objects.
type Model[T any] interface {
	Values(*Handler) map[string]any
	SetValues(*Handler, map[string]any)
	*T
}

// A model that backs a resource, exposing the ids needed for CRUD operations.
type ResourceModel[T any] interface {
	Model[T]
	GetID() types.String
	SetID(id types.String)
	GetProjectID() types.String
}

// Models without this interface are unprotected when the deletion protection attribute is unset.
type DeletionProtectionDefaulter interface {
	DeletionProtectionDefault(ctx context.Context) bool
}
