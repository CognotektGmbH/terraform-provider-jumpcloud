package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccUserGroupAssociation(t *testing.T) {
	randSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:           func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testUserGroupAssocConfig(randSuffix),
				Check: resource.TestCheckResourceAttrSet("jumpcloud_user_group_association.test_association",
					"group_id"),
			},
		},
	})
}

func testUserGroupAssocConfig(randSuffix string) string {
	return fmt.Sprintf(`
resource "jumpcloud_application" "test_application" {
	display_name        = "test_aws_account_%s"
	sso_url              = "https://sso.jumpcloud.com/saml2/example-application-%s"
}

resource "jumpcloud_user_group" "test_group" {
	name = "testgroup_%s"
}

resource "jumpcloud_user_group_association" "test_association" {
	object_id = jumpcloud_application.test_application.id
	group_id  = jumpcloud_user_group.test_group.id
	type      = "application"
}
`, randSuffix, randSuffix, randSuffix)
}