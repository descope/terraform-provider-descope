package testacc

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func Project(t *testing.T) *Resource {
	return newResource(t, "project", "test")
}

func Descoper(t *testing.T) *Resource {
	return newResource(t, "descoper", "test")
}

func ManagementKey(t *testing.T) *Resource {
	return newResource(t, "management_key", "test")
}

func newResource(t *testing.T, typ, id string) *Resource {
	return &Resource{
		Type: typ,
		ID:   id,
		Name: GenerateAlias(t),
	}
}

type Resource struct {
	Type string // the resource type without the 'descope_' prefix
	ID   string // the resource name in the Terraform config
	Name string // the value of the 'name' attribute
}

func (r *Resource) Path() string {
	return fmt.Sprintf(`descope_%s.%s`, r.Type, r.ID)
}

func (r *Resource) Variables(s ...string) string {
	return strings.Join(s, "\n")
}

func (r *Resource) Config(s ...string) string {
	n := fmt.Sprintf(`name = %q`, r.Name)
	s = append([]string{n}, s...)
	return fmt.Sprintf(resourceFormat, r.Type, r.ID, strings.Join(s, "\n	"))
}

func (r *Resource) Check(checks map[string]any, extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	path := r.Path()
	f := []resource.TestCheckFunc{}
	checks = flatten(checks, "")
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
		} else if value, ok := v.([]string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k+".#", strconv.Itoa(len(value))))
			for i := range value {
				f = append(f, resource.TestCheckTypeSetElemAttr(path, fmt.Sprintf("%s.*", k), value[i]))
			}
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

const resourceFormat = `
resource "descope_%s" "%s" {
	%s
}
`

func GenerateAlias(t *testing.T) string {
	test := strings.TrimPrefix(t.Name(), "Test")
	ts := time.Now().Format("01021504") // MMddHHmm
	rand, err := uuid.GenerateUUID()
	require.NoError(t, err)
	suffix := rand[len(rand)-8:]
	return fmt.Sprintf("testacc-%s-%s-%s", test, ts, suffix)
}

func flatten(checks map[string]any, keypath string) map[string]any {
	result := map[string]any{}
	for k, v := range checks {
		if keypath != "" {
			k = keypath + "." + k
		}
		if m, ok := v.(map[string]any); ok {
			maps.Copy(result, flatten(m, k))
		} else {
			result[k] = v
		}
	}
	return result
}
