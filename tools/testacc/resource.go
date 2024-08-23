package testacc

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type Resource struct {
	Type string
	Name string
}

func (r *Resource) Path() string {
	return fmt.Sprintf(`%s.%s`, r.Type, r.Name)
}

func (r *Resource) Config(s ...string) string {
	return fmt.Sprintf(`
		resource %q %q {
			%s
		}
		`, r.Type, r.Name, strings.Join(s, "\n"))
}

func (r *Resource) Check(checks map[string]any) resource.TestCheckFunc {
	path := r.Path()
	f := []resource.TestCheckFunc{}
	for k, v := range checks {
		if value, ok := v.(string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, value))
		} else {
			f = append(f, resource.TestCheckResourceAttrSet(path, k))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(f...)
}
