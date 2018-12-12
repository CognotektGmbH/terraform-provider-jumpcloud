package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate, //optional
		Delete: resourceUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_samba": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*jcapiv2.APIClient)
	req := map[string]interface{}{
		"body": jcapiv2.UserGroupPost{
			Name: d.Get("name").(string),
			Attributes: &jcapiv2.UserGroupPostAttributes{
				// Note: PosixGroups cannot be edited after group creation
				PosixGroups: []jcapiv2.UserGroupPostAttributesPosixGroups{
					// jcapiv2.UserGroupPostAttributesPosixGroups{Id:0, Name:"foo"}, TODO: add to/read from schema
				},
				SambaEnabled: d.Get("enable_samba").(bool),
			},
		},
		// "xOrgId": "", TODO: add to/read from schema
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
	client := m.(*jcapiv2.APIClient)
	group, _, err := client.UserGroupsApi.GroupsUserGet(context.TODO(), d.Id(), "", Accept, nil)
	if err != nil {
		return err
	}

	d.SetId(group.Id)
	if err = d.Set("name", group.Name); err != nil {
		return nil
	}
	// TODO: attributes?
	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*jcapiv2.APIClient)

	body := jcapiv2.UserGroupPut{
		Name: d.Get("name").(string), // Always set since it is a required value
	}

	if d.HasChange("enable_samba") {
		body.Attributes = &jcapiv2.UserGroupPutAttributes{
			SambaEnabled: d.Get("enable_samba").(bool),
		}
	}

	req := map[string]interface{}{"body": body, "xOrgId": ""} // TODO: xOrgId: add to/read from schema
	_, _, err := client.UserGroupsApi.GroupsUserPut(context.TODO(), d.Id(), "", Accept, req)
	if err != nil {
		return err
	}

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*jcapiv2.APIClient)
	res, err := client.UserGroupsApi.GroupsUserDelete(context.TODO(), d.Id(), "", Accept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting: %s - response = %+v", err, res)
	}
	return nil
}
