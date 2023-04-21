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

// // func TestAccDataSourceJumpCloudUser_basic(t *testing.T) {
// //     userEmail := "<existing-user-email>" // Replace this with an existing user's email address
// //     resourceName := "data.jumpcloud_user.test"

// //     resource.Test(t, resource.TestCase{
// //         PreCheck:     func() { testAccPreCheck(t) },
// //         Providers:    testAccProviders,
// //         Steps: []resource.TestStep{
// //             {
// //                 Config: testAccDataSourceJumpCloudUserConfig(userEmail),
// //                 Check: resource.ComposeTestCheckFunc(
// //                     resource.TestCheckResourceAttrSet(resourceName, "id"),
// //                 ),
// //             },
// //         },
// //     })
// // }

// // func testAccDataSourceJumpCloudUserConfig(userEmail string) string {
// //     return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) + `

// // provider "jumpcloud" {
// //   api_key = "` + os.Getenv("JUMPCLOUD_API_KEY") + `"
// // }

// // data "jumpcloud_user" "test" {
// //   email = "` + userEmail + `"
// // }
// // `
// // }