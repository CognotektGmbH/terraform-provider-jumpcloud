package jumpcloud

import (
	"context"
	"fmt"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceJumpCloudApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJumpCloudApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceJumpCloudApplicationRead(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)
	applicationName := d.Get("name").(string)
	applicationsResponse, _, err := client.ApplicationsApi.ApplicationsList(context.Background(), "_id, displayName", "", nil)

	if err != nil {
		return err
	}

	applications := applicationsResponse.Results

	for _, application := range applications {
		if application.DisplayName == applicationName {
			d.SetId(application.Id)
			return nil
		}
	}

	return fmt.Errorf("No application found with name: %s", applicationName)
}