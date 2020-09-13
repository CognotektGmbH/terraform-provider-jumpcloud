package jumpcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUser(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// This test simply applys a user with the config from testAccUser
				// and checks for the correct username and email in the state
				// The resource is destroyed afterwards via the framework
				Config: testAccUser(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_user.test_user", "username", rName),
					resource.TestCheckResourceAttr("jumpcloud_user.test_user", "email", rName+"@testorg.com"),
				),
			},
		},
	})
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("JUMPCLOUD_API_KEY"); v == "" {
		t.Fatal("JUMPCLOUD_API_KEY= must be set for the acceptance tests")
	}
}

func testAccUser(name string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_user" "test_user" {
  			username = "%s"
			email = "%s@testorg.com"
			firstname = "Firstname"
			lastname = "Lastname"
			enable_mfa = true
		}`, name, name,
	)
}
