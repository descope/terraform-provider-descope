package emit

import (
	"context"
	"slices"

	"github.com/hashicorp/hcl/v2"
)

// referenceIDLength guards against substituting short human-chosen ids (flow ids, attribute names); minted ids are long and prefixed.
const referenceIDLength = 20

func references(plan *Plan) map[string]hcl.Traversal {
	refs := map[string]hcl.Traversal{}
	for _, resource := range plan.Resources {
		id := resource.EntityID
		if len(id) < referenceIDLength || id == plan.ProjectID {
			continue
		}
		refs[id] = hcl.Traversal{
			hcl.TraverseRoot{Name: resource.Type},
			hcl.TraverseAttr{Name: resource.Label},
			hcl.TraverseAttr{Name: "id"},
		}
	}
	return refs
}

func projectReference(plan *Plan) hcl.Traversal {
	if plan.ProjectAddress != nil {
		return append(slices.Clone(plan.ProjectAddress), hcl.TraverseAttr{Name: "id"})
	}
	for _, resource := range plan.Resources {
		if resource.Type == "descope_project" {
			return hcl.Traversal{
				hcl.TraverseRoot{Name: resource.Type},
				hcl.TraverseAttr{Name: resource.Label},
				hcl.TraverseAttr{Name: "id"},
			}
		}
	}
	return nil
}

func nameReferences(ctx context.Context, plan *Plan) map[string]hcl.Traversal {
	refs := map[string]hcl.Traversal{}
	for _, resource := range plan.Resources {
		if resource.Type != "descope_permission" {
			continue
		}
		for _, attr := range resource.Attrs {
			if attr.Name != "name" {
				continue
			}
			if name, ok := stringValue(ctx, attr.Value); ok && name != "" {
				refs[name] = hcl.Traversal{
					hcl.TraverseRoot{Name: resource.Type},
					hcl.TraverseAttr{Name: resource.Label},
					hcl.TraverseAttr{Name: "name"},
				}
			}
		}
	}
	return refs
}
