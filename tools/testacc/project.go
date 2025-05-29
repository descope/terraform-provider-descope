package testacc

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
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
	Resource
	Name string
}

func (p *ProjectResource) Config(s ...string) string {
	n := fmt.Sprintf(`name = %q`, p.Name)
	s = append([]string{n}, s...)
	return p.Resource.Config(s...)
}
