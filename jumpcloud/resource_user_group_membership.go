package jumpcloud

import (
	"context"
	//	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupMembershipCreate,
		Read:   resourceUserGroupMembershipRead,
		// We must  not have an update routine as the association cannot be updated
		// Any change in the group ID or user ID forces a recreation of the resource
		Update: nil,
		Delete: resourceUserGroupMembershipDelete,
		Schema: map[string]*schema.Schema{
			"userid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"xorgid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"groupid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	payload := jcapiv2.UserGroupMembersReq{
		Op:    "add",
		Type_: "user",
		Id:    d.Get("userid").(string),
	}

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}

	_, err := client.UserGroupMembersMembershipApi.GraphUserGroupMembersPost(
		context.TODO(), d.Get("groupid").(string), "", "", req)
	if err != nil {
		return err
	}

	return resourceUserGroupMembershipRead(d, m)
}

func resourceUserGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	graphconnect, _, err := client.UserGroupMembersMembershipApi.GraphUserGroupMembersList(
		context.TODO(), d.Get("groupid").(string), "", "", nil)
	if err != nil {
		return err
	}
	// The Userids are hidden in a super-complex construct, see
	// https://github.com/TheJumpCloud/jcapi-go/blob/master/v2/docs/GraphConnection.md
	for _, v := range graphconnect {
		if v.To.Id == d.Get("userid") {
			// Found - As we not have a JC-ID for the membership we simply store
			// the combination of group ID and user ID as our membership ID
			d.SetId(d.Get("groupid").(string) + "/" + d.Get("userid").(string))
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceUserGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	payload := jcapiv2.UserGroupMembersReq{
		Op:    "remove",
		Type_: "user",
		Id:    d.Get("userid").(string),
	}

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}
	client.UserGroupMembersMembershipApi.GraphUserGroupMembersPost(
		context.TODO(), d.Get("groupid").(string), "", "", req)
	return nil
}
