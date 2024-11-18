package testacc

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func Project(t *testing.T) *ProjectResource {
	test := strings.TrimPrefix(t.Name(), "Test")

	time := time.Now().Format("01021504") // MMddHHmm

	uuid, err := uuid.GenerateUUID()
	require.NoError(t, err)
	suffix := uuid[len(uuid)-8:]

	return &ProjectResource{
		Resource: Resource{Type: "descope_project", Name: "testproj"},
		Name:     fmt.Sprintf("testacc-%s-%s-%s", test, time, suffix),
	}
}

// Resource

type ProjectResource struct {
	Resource Resource
	Name     string
}

func (p *ProjectResource) Path() string {
	return p.Resource.Path()
}

func (p *ProjectResource) Config(s ...string) string {
	n := fmt.Sprintf(`name = %q`, p.Name)
	s = append([]string{n}, s...)
	return p.Resource.Config(s...)
}

func (p *ProjectResource) Check(checks map[string]any) resource.TestCheckFunc {
	return p.Resource.Check(checks)
}

// Configs

func (p *ProjectResource) Settings(s ...string) string {
	return fmt.Sprintf(`
		project_settings = {
			%s
		}
		`, strings.Join(s, "\n"))
}
