package listattr

import "github.com/descope/terraform-provider-descope/internal/models/helpers"

// A simple heuristic for preserving model object IDs by matching names or the order in the list.
func Match[T any, M helpers.MatchableModel[T]](h *helpers.Handler, plan *Type[T], state Type[T]) {
	unmatched := []M{}
	// First for each existing model object look for a matching one in the plan and
	// give it the ID value, effectively mimicking UseStateForUnknown. This should usually
	// be enough to handle the first common case where a model object is added to a list
	// but the other objects in the list aren't changed.
	for e := range Iterator(state, h) {
		var existing M = e
		matched := false
		for p := range MutatingIterator(plan, h) {
			var planned M = p
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
	for p := range MutatingIterator(plan, h) {
		var planned M = p
		if planned.GetID().IsUnknown() {
			if index < len(unmatched) {
				planned.SetID(unmatched[index].GetID())
				index += 1
			}
		}
	}
}
