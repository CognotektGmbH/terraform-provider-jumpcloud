package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceJumpCloudApplication_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceJumpCloudApplicationConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_application.test_application", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_application.test_application", "name", rName),
				),
			},
		},
	})
}

func testAccDataSourceJumpCloudApplicationConfig(applicationName string) string {
	return fmt.Sprintf(`
resource "jumpcloud_application" "test_application" {
  display_name = "%s"
  // ... other required attributes ...
}

data "jumpcloud_application" "test_application" {
  name = jumpcloud_application.test_application.display_name
}`, applicationName)
}
