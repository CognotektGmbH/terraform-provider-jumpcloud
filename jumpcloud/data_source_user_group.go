package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceJumpCloudUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJumpCloudUserGroupRead,
		Schema: map[string]*schema.Schema{
			"group_name": {
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

func dataSourceJumpCloudUserGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	groupName := d.Get("group_name").(string)

	filter := fmt.Sprintf(`{"name":"%s"}`, groupName)

	page := 0
	limit := 100 // Adjust the limit as per your requirement

	for {
		groups, _, err := client.UserGroupsApi.GroupsUserList(context.Background(), "_id, name", filter, nil)
		if err != nil {
			return err
		}

		for _, group := range groups {
			if group.Name == groupName {
				d.SetId(group.Id)
				return nil
			}
		}

		if len(groups) < limit {
			break
		}

		page++
	}

	return fmt.Errorf("No user group found with name: %s", groupName)
}
