package helpers

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func HasUnknownValues(values ...any) bool {
	for _, v := range values {
		switch v := v.(type) {
		case interface{ IsUnknown() bool }:
			if v.IsUnknown() {
				return true
			}
		case []types.String:
			if slices.ContainsFunc(v, func(s types.String) bool { return s.IsUnknown() }) {
				return true
			}
		}
	}
	return false
}
