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
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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

// Helper to look up a system group by name
func resourceGroupsSystemList_match(d *schema.ResourceData, m interface{}) (jcapiv2.SystemGroup, error) {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	filter := "[name:eq:" + d.Get("name").(string) + "]"

	req := map[string]interface{}{
		"filter": filter,
	}

	result, _, err := client.SystemGroupsApi.GroupsSystemList(context.TODO(),
		"", headerAccept, req)
	if err == nil {
		return result[0], nil
	} else {
		return jcapiv2.SystemGroup{}, err
	}
}

func resourceGroupsSystemRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string

	id = d.Id()

	if d.Id() != "" {
		id_lookup, err := resourceGroupsSystemList_match(d, m)
		if err != nil {
			return fmt.Errorf("unable to locate ID for group %s",
				d.Get("name"))
		}
		d.SetId(id_lookup.Id)
	}

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
