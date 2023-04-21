package jumpcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

//Make sure to replace "your-jumpcloud-api-key" with your actual JumpCloud API key.
const testConfig = `
provider "jumpcloud" {
  # Configure your JumpCloud API key
  api_key = "your-jumpcloud-api-key"
}

data "jumpcloud_user" "test" {
  email = "john.doe@example.com"
}
`

func TestAccDataSourceJumpCloudUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.jumpcloud_user.test", "id"),
					resource.TestCheckResourceAttr("data.jumpcloud_user.test", "email", "user@example.com"),
				),
			},
		},
	})
}