package jumpcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"xorgid": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	req := map[string]interface{}{
		"body": jcapiv2.UserGroupPost{
			Name: d.Get("name").(string),
			// Attributes: &jcapiv2.UserGroupPostAttributes{
			// 	//	Note: PosixGroups cannot be edited after group creation, only first member of slice is considered
			// 	PosixGroups: []jcapiv2.UserGroupPostAttributesPosixGroups{
			// 		jcapiv2.UserGroupPostAttributesPosixGroups{Id: int32(posixID), Name: posixName},
			// 	},
			// },
		},
		"xOrgId": d.Get("xorgid").(string),
	}
	group, res, err := client.UserGroupsApi.GroupsUserPost(context.TODO(), "", Accept, req)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error creating user group %s: %s - response = %+v", (req["body"].(jcapiv2.UserGroupPost)).Name, err, res)
	}

	d.SetId(group.Id)
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)

	group, ok, err := trueUserGroupRead(config, d.Id())
	if err != nil {
		return err
	}

	if !ok {
		// not found
		d.SetId("")
		return nil
	}

	d.SetId(group.Id)
	if err := d.Set("name", group.Name); err != nil {
		return err
	}
	// if err := d.Set("attributes", flattenAttributes(&group.Attributes)); err != nil {
	// 	return err
	// }

	return nil
}

func trueUserGroupRead(config *jcapiv2.Configuration, id string) (ug *UserGroup, ok bool, err error) {
	req, err := http.NewRequest(http.MethodGet, config.BasePath+"/usergroups/"+id, nil)
	if err != nil {
		return
	}

	req.Header.Add("x-api-key", config.DefaultHeader["x-api-key"])
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return
	}

	ok = true
	err = json.NewDecoder(res.Body).Decode(&ug)
	return
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	res, err := client.UserGroupsApi.GroupsUserDelete(context.TODO(), d.Id(), "", Accept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting user group: %s - response = %+v", err, res)
	}
	d.SetId("")
	return nil
}
