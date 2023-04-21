package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceJumpCloudUserGroup_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceJumpCloudUserGroupConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user_group.test_group", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user_group.test_group", "group_name", rName),
				),
			},
		},
	})
}

func testAccDataSourceJumpCloudUserGroupConfig(groupName string) string {
	return fmt.Sprintf(`
resource "jumpcloud_user_group" "test_group" {
  name = "%s"
}

data "jumpcloud_user_group" "test_group" {
  group_name = jumpcloud_user_group.test_group.name
}`, groupName)
}
