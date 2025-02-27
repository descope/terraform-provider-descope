package helpers

import (
	"context"
	"slices"
)

type contextKey string

const importKey = contextKey("descopeImport")

func SetImport(ctx context.Context) context.Context {
	return context.WithValue(ctx, importKey, true)
}

func IsImport(ctx context.Context) bool {
	return ctx.Value(importKey) == true
}

func InitIfImport[T any](ctx context.Context, o *T) *T {
	if IsImport(ctx) {
		return ZVL(o)
	}
	return o
}

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
