package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test(t *testing.T, steps ...resource.TestStep) {
	resource.Test(t, TestCase(t, steps...))
}

func TestCase(t *testing.T, steps ...resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps:                    steps,
	}
}

func CleanupStep() resource.TestStep {
	return resource.TestStep{Config: `locals {}`}
}
