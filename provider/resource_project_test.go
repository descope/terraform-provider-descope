package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("descope_project.test", "id"),
					resource.TestCheckResourceAttr("descope_project.test", "name", "one"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "descope_project.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name", "one"},
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("descope_project.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "descope_project" "test" {
  name = %[1]q
}
`, name)
}
