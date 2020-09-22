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
			"jc_id": {
				Type:     schema.TypeString,
				Computed: true,
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

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return resourceGroupsSystemRead(d, m)
}

// Helper to look up a system group by name
func resourceGroupsSystemList_match(d *schema.ResourceData, m interface{}) (jcapiv2.SystemGroup, error) {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var filter []string

	filter = append(filter, "name:eq:"+d.Id())

	optional := map[string]interface{}{
		"filter": filter,
	}

	result, _, err := client.SystemGroupsApi.GroupsSystemList(context.TODO(),
		"", headerAccept, optional)
	if err == nil {
		if len(result) < 1 {
			return jcapiv2.SystemGroup{}, fmt.Errorf("System Group \"%s\" not found.", d.Id())
		} else {
			return result[0], nil
		}
	} else {
		return jcapiv2.SystemGroup{}, err
	}
}

func resourceGroupsSystemRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string

	id = d.Get("jc_id").(string)

	if id == "" {
		id_lookup, err := resourceGroupsSystemList_match(d, m)
		if err != nil {
			return fmt.Errorf("Unable to locate ID for group %s, %+v",
				d.Get("name"), err)
		}
		id = id_lookup.Id
		d.SetId(id_lookup.Name)
		d.Set("name", id_lookup.Name)
		d.Set("jc_id", id_lookup.Id)
	}

	group, res, err := client.SystemGroupsApi.GroupsSystemGet(context.TODO(),
		id, "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error reading system group ID %s: %s - response = %+v",
			d.Id(), err, res)
	}

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return nil
}

func resourceGroupsSystemUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string
	id = d.Get("jc_id").(string)

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

	d.SetId(group.Name)
	d.Set("name", group.Name)
	d.Set("jc_id", group.Id)
	return resourceGroupsSystemRead(d, m)
}

func resourceGroupsSystemDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	var id string
	id = d.Get("jc_id").(string)

	res, err := client.SystemGroupsApi.GroupsSystemDelete(context.TODO(),
		id, "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting system group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}
