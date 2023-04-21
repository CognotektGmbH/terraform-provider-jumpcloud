package jumpcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func Test_resourceApplication(t *testing.T) {
	randSuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	fullResourceName := "jumpcloud_application.example_app"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// Create step
			{
				Config: testApplicationConfig(randSuffix, "test_aws_account"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, "display_label", "test_aws_account"),
				),
			},
			userImportStep(fullResourceName),
			// Update Step
			{
				Config: testApplicationConfig(randSuffix, "test_aws_account_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(fullResourceName, "display_label", "test_aws_account_updated"),
				),
			},
			userImportStep(fullResourceName),
		},
	})
}

func testApplicationConfig(randSuffix string, displayLabel string) string {
	return fmt.Sprintf(`
resource "jumpcloud_application" "example_app" {
	display_label        = "%s"
	sso_url              = "https://sso.jumpcloud.com/saml2/example-application_%s"
	saml_role_attribute  = "arn:aws:iam::AWS_ACCOUNT_ID:role/MY_ROLE,arn:aws:iam::AWS_ACCOUNT_ID:saml-provider/MY_SAML_PROVIDER"
	aws_session_duration = 432000
}
`, displayLabel, randSuffix)
}