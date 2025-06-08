package helpers

import "github.com/hashicorp/terraform-plugin-framework/types"

const (
	DescopeConnector = "Descope"
	DescopeTemplate  = "System"
)

// Pointer receiver interface for model objects.
type Model[T any] interface {
	Values(*Handler) map[string]any
	SetValues(*Handler, map[string]any)
	*T
}

// A model that can be matched by name and ID, primarily for making more friendly diffs in lists.
type MatchableModel[T any] interface {
	Model[T]
	GetName() types.String
	GetID() types.String
	SetID(id types.String)
}

// A model that can return a list of references to other model objects.
type CollectReferencesModel[T any] interface {
	Model[T]
	CollectReferences(*Handler)
}

// A model that has references that need to be updated after the model is created or updated.
type UpdateReferencesModel[T any] interface {
	Model[T]
	UpdateReferences(*Handler)
}
