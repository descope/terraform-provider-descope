package testacc

import (
	"fmt"
	"strconv"
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

func (r *Resource) Check(checks map[string]any, extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	path := r.Path()
	f := []resource.TestCheckFunc{}
	for k, v := range checks {
		if first := strings.TrimSuffix(k, ".=="); first != k {
			second, ok := v.(string)
			if !ok {
				panic(fmt.Sprintf("unexpected non-string argument of type %T in equality check: %v", v, v))
			}
			f = append(f, resource.TestCheckResourceAttrPair(path, first, path, second))
		} else if value, ok := v.(string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, value))
		} else if value, ok := v.(int); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, strconv.Itoa(value)))
		} else if value, ok := v.(bool); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, fmt.Sprintf("%t", value)))
		} else if value, ok := v.(func(string) error); ok {
			f = append(f, resource.TestCheckResourceAttrWith(path, k, value))
		} else if v == AttributeIsSet {
			f = append(f, resource.TestCheckResourceAttrSet(path, k))
		} else if v == AttributeIsNotSet {
			f = append(f, resource.TestCheckNoResourceAttr(path, k))
		} else {
			panic(fmt.Sprintf("unexpected value of type %T in Check(): %v", v, v))
		}
	}
	f = append(f, extras...)
	return resource.ComposeAggregateTestCheckFunc(f...)
}
