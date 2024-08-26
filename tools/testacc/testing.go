package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Run(t *testing.T, steps ...resource.TestStep) {
	t.Parallel()
	resource.Test(t, TestCase(t, steps...))
}

func RunIsolated(t *testing.T, steps ...resource.TestStep) {
	resource.Test(t, TestCase(t, steps...))
}

func TestCase(t *testing.T, steps ...resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps:                    steps,
	}
}
