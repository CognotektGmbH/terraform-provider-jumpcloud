package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroupsSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupsSystemCreate,
		Read:   resourceGroupsSystemRead,
		Update: resourceGroupsSystemUpdate,
		Delete: resourceGroupsSystemDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceGroupsSystemCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	body := jcapiv2.SystemGroupData{Name: d.Get("name").(string)}

	req := map[string]interface{}{
		"body": body,
	}
	group, res, err := client.SystemGroupsApi.GroupsSystemPost(context.TODO(),
		"", headerAccept, req)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error creating system group %s: %s - response = %+v",
			(req["body"].(jcapiv2.SystemGroupData)).Name, err, res)
	}

	d.SetId(group.Id)
	return resourceGroupsSystemRead(d, m)
}

func resourceGroupsSystemRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string

	id = d.Id()

	group, res, err := client.SystemGroupsApi.GroupsSystemGet(context.TODO(),
		id, "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error reading system group ID %s: %s - response = %+v",
			d.Get("id").(string), err, res)
	}

	d.SetId(group.Id)
	d.Set("name", group.Name)
	return nil
}

func resourceGroupsSystemUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string

	id = d.Id()

	body := jcapiv2.SystemGroupData{Name: d.Get("name").(string)}

	req := map[string]interface{}{
		"body": body,
	}

	group, res, err := client.SystemGroupsApi.GroupsSystemPut(context.TODO(),
		id, "", headerAccept, req)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error updating system group %s: %s - response = %+v",
			d.Get("name"), err, res)
	}

	d.SetId(group.Id)
	return resourceGroupsSystemRead(d, m)
}

func resourceGroupsSystemDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	res, err := client.SystemGroupsApi.GroupsSystemDelete(context.TODO(),
		d.Id(), "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting system group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}
