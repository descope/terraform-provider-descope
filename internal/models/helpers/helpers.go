package helpers

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
