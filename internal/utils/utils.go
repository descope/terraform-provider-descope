package utils

// Returns the first non-null pointer value or a pointer to the empty value if none was found
func ZVL[T any](o *T, rest ...*T) *T {
	if o != nil {
		return o
	}
	for _, v := range rest {
		if v != nil {
			return v
		}
	}
	var empty T
	return &empty
}
