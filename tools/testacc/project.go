package testacc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func Project(t *testing.T) *ProjectResource {
	uuid, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return &ProjectResource{
		Resource: Resource{Type: "descope_project", Name: "test"},
		Name:     fmt.Sprintf("%s-%s", t.Name(), uuid),
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
