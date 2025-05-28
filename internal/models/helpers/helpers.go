package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

const RootKey string = ""

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
		case interface{ Elements() map[string]attr.Value }:
			for _, elem := range v.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		case interface{ Elements() []attr.Value }:
			for _, elem := range v.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		case interface{ Attributes() map[string]attr.Value }:
			for _, elem := range v.Attributes() {
				if elem.IsUnknown() {
					return true
				}
			}
		default:
			panic(fmt.Sprintf("unexpected type in HasUnknownValues: %T", v))
		}
	}
	return false
}
