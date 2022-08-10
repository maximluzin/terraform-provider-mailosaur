package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServerExampleResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailosaur_server.test", "name", "one"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "id"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "password"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "email"),
				),
			},
			// Update and Read testing
			{
				Config: testAccServerExampleResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mailosaur_server.test", "name", "two"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "id"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "password"),
					resource.TestCheckResourceAttrSet("mailosaur_server.test", "email"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccServerExampleResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "mailosaur_server" "test" {
  name = %[1]q
}
`, name)
}
