package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const DescopeTemplate = "System"

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

// A model whose writes are version checked against a computed "version" attribute. An update sends the version
// from state, since the plan cannot know it.
type VersionChecked interface {
	VersionChecked()
}

// Models without this interface are unprotected when the deletion protection attribute is unset.
type DeletionProtectionDefaulter interface {
	DeletionProtectionDefault(ctx context.Context) bool
}
