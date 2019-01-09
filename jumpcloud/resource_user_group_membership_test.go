package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUserGroupMembership(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// The only reasonable step is to check if the useris is in the state
				// It will be deleted from the state in case the membership could not be
				// established
				Config: testAccUserGroupMembership(rName),
				Check: resource.TestCheckResourceAttrSet("jumpcloud_user_group_membership.test_membership_"+rName,
					"userid"),
			},
		},
	})
}

// This needs to be moved to a group acceptance test later
func testAccUserGroupMembership(name string) string {
	return fmt.Sprintf(`

		resource "jumpcloud_user" "test_user_%s" {
		username = "%s"
		email = "%s@testorg.com"
		}

		resource "jumpcloud_user_group" "test_group_%s" {
		name = "testgroup_%s"
		}



	resource "jumpcloud_user_group_membership" "test_membership_%s" {
  userid = "${jumpcloud_user.test_user_%s.id}"
	groupid = "${jumpcloud_user_group.test_group_%s.id}"
  }
`,
		name, name, name, name, name, name, name, name,
	)
}
