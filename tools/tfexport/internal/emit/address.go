package emit

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func ParseProjectAddress(address string) (hcl.Traversal, error) {
	traversal, diags := hclsyntax.ParseTraversalAbs([]byte(address), "", hcl.InitialPos)
	if diags.HasErrors() || !isProjectAddress(traversal) {
		return nil, fmt.Errorf("invalid project address %q: expected a descope_project resource address such as descope_project.main or descope_project.main[0]", address)
	}
	return traversal, nil
}

func ParseModulePrefix(prefix string) (hcl.Traversal, error) {
	traversal, diags := hclsyntax.ParseTraversalAbs([]byte(prefix), "", hcl.InitialPos)
	if diags.HasErrors() || !isModulePath(traversal) {
		return nil, fmt.Errorf("invalid import prefix %q: expected a module path such as module.auth or module.auth[\"prod\"]", prefix)
	}
	return traversal, nil
}

func isProjectAddress(traversal hcl.Traversal) bool {
	if len(traversal) < 2 || len(traversal) > 3 || traversal.RootName() != "descope_project" {
		return false
	}
	if _, ok := traversal[1].(hcl.TraverseAttr); !ok {
		return false
	}
	if len(traversal) == 3 {
		_, ok := traversal[2].(hcl.TraverseIndex)
		return ok
	}
	return true
}

func isModulePath(traversal hcl.Traversal) bool {
	if len(traversal) == 0 || traversal.RootName() != "module" {
		return false
	}
	for i := 1; i < len(traversal); {
		if _, ok := traversal[i].(hcl.TraverseAttr); !ok {
			return false
		}
		i++
		if i < len(traversal) {
			if _, ok := traversal[i].(hcl.TraverseIndex); ok {
				i++
			}
		}
		if i == len(traversal) {
			return true
		}
		if step, ok := traversal[i].(hcl.TraverseAttr); !ok || step.Name != "module" || i == len(traversal)-1 {
			return false
		}
		i++
	}
	return false
}
