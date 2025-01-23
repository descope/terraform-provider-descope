package helpers

import "slices"

func AnySliceToStringSlice(data map[string]any, key string) []string {
	var strs []string
	if objects, ok := data[key].([]any); ok {
		for _, o := range objects {
			if s, ok := o.(string); ok {
				strs = append(strs, s)
			}
		}
	}
	return strs
}

func EqualStringSliceElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !slices.Contains(b, v) {
			return false
		}
	}
	return true
}
