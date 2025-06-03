package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const RootKey string = ""

func Must[T any](x T, diags diag.Diagnostics) T {
	if errs := diags.Errors(); len(errs) > 0 {
		panic(fmt.Sprintf("%s: %s", errs[0].Summary(), errs[0].Detail()))
	}
	return x
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
		}
	}
	return false
}
