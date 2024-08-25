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

func (r *Resource) Check(checks map[string]any) resource.TestCheckFunc {
	path := r.Path()
	f := []resource.TestCheckFunc{}
	for k, v := range checks {
		if value, ok := v.(string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, value))
		} else if value, ok := v.(int); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, strconv.Itoa(value)))
		} else if v == AttributeIsSet {
			f = append(f, resource.TestCheckResourceAttrSet(path, k))
		} else if v == AttributeIsNotSet {
			f = append(f, resource.TestCheckNoResourceAttr(path, k))
		} else {
			panic(fmt.Sprintf("unexpected value of type %T in Check(): %v", v, v))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(f...)
}

type attributeSymbol struct {
	v int
}

var AttributeIsSet = attributeSymbol{v: 1}
var AttributeIsNotSet = attributeSymbol{v: -1}
