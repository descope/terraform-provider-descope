package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

const RootKey string = ""

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
