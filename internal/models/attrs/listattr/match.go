package listattr

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

// A simple heuristic for preserving model object IDs by matching names or the order in the list.
func ModifyMatching[T any, M helpers.MatchableModel[T]](h *helpers.Handler, plan *Type[T], state Type[T]) {
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

// Like Set but looks for matching model objects in the list by name.
func SetMatching[T any, M helpers.MatchableModel[T]](l *Type[T], data map[string]any, key string, h *helpers.Handler) {
	// convert the data in the map to a slice of objects
	objects := []map[string]any{}
	values, _ := data[key].([]any)
	for i := range values {
		if v, ok := values[i].(map[string]any); ok {
			objects = append(objects, v)
		}
	}

	// get the current elements in the list, so we can update/delete them
	current, diags := l.ToSlice(h.Ctx)
	h.Diagnostics.Append(diags...)
	if diags.HasError() {
		h.Error("List Conversion Failed", "Could not convert list to slice of elements for setting key '%s'", key)
		return
	}

	// the final list of elements with updated and new ones, and without any deleted ones
	elements := []*T{}

	// for each current element, look for a matching object with the same name
	for _, e := range current {
		var existing M = e
		for i, o := range objects {
			if n, _ := o["name"].(string); n == existing.GetName().ValueString() {
				// if the name matches, we update the existing object
				existing.SetValues(h, o)
				// remove from the list so we know it's not a new model object
				objects = slices.Delete(objects, i, i+1)
				// add the existing object to the matched list as we now know it hasn't been deleted
				elements = append(elements, existing)
				break
			}
		}
	}

	// any objects left here are new model objects that need to be added
	for _, o := range objects {
		var element M = new(T)
		element.SetValues(h, o)
		elements = append(elements, element)
	}

	*l = valueOf(h.Ctx, elements)
}
