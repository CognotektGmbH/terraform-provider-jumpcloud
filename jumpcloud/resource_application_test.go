package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccApplication(t *testing.T) {
	rDisplayName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// This test simply applys an application with the config from testAccApplication
				// and checks for the correct data in the state
				// The resource is destroyed afterwards via the framework
				Config: testAccApplication(rDisplayName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_application.test_application", "display_name", rDisplayName),
				),
			},
		},
	})
}

func testAccApplication(display_label string) string {
	return fmt.Sprintf(`
		resource "jumpcloud_application" "test_application" {
  			display_label = "%s"
		}`, display_label,
	)
}
