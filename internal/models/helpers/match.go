package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MatchableModel[T any] interface {
	Model[T]
	GetName() types.String
	GetID() types.String
	SetID(id types.String)
}

// A simple heuristic for preserving model object IDs by matching names or the order in the list.
func MatchModels[T any, M MatchableModel[T]](_ context.Context, plan []M, state []M) {
	unmatched := []M{}
	// First look for each existing model object look for a matching one in the plan and
	// give it the ID value, effectively mimicking UseStateForUnknown. This should usually
	// be enough to handle the first common case where a model object is added to a list
	// but the other objects in the list aren't changed.
	for _, existing := range state {
		matched := false
		for _, planned := range plan {
			if planned.GetName().Equal(existing.GetName()) {
				planned.SetID(existing.GetID())
				matched = true
				break
			}
		}
		// keep any unmatched existing model objects in the same order to use below
		if !matched {
			unmatched = append(unmatched, existing)
		}
	}
	// Any model object in the plan that wasn't matched with an existing one will
	// get any leftover IDs from the existing model objects. This heuristic matching
	// should usually be neough to handle second common case where a model object's
	// name is changed but the list structure itself remains the same.
	index := 0
	for _, planned := range plan {
		if planned.GetID().IsUnknown() {
			if index < len(unmatched) {
				planned.SetID(unmatched[index].GetID())
				index += 1
			}
		}
	}
}