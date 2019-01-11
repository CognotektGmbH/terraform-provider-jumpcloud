package jumpcloud

import (
	"context"
	"strings"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserGroupMembershipCreate,
		Read:   resourceUserGroupMembershipRead,
		// We must not have an update routine as the association cannot be updated.
		// Any change in one of the elements forces a recreation of the resource
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
			State: userGroupMembershipImporter,
		},
	}
}

// We cannot use the regular importer as it calls the read function ONLY with the ID field being
// populated.- In our case, we need the group ID and user ID to do the read - But since our
// artificial resource ID is simply the concatenation of user ID group ID seperated by  a '/',
// we can derive both values during our import process
func userGroupMembershipImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), "/")
	d.Set("groupid", s[0])
	d.Set("userid", s[1])
	return []*schema.ResourceData{d}, nil
}

func modifyUserGroupMembership(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) error {
	payload := jcapiv2.UserGroupMembersReq{
		Op:    action,
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
	return nil
}

func resourceUserGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	err := modifyUserGroupMembership(client, d, "add")
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
			// the concatenation of group ID and user ID as our membership ID
			d.SetId(d.Get("groupid").(string) + "/" + d.Get("userid").(string))
			return nil
		}
	}
	// Element does not exist in actual Infrastructure, hence unsetting the ID
	d.SetId("")
	return nil
}

func resourceUserGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	return modifyUserGroupMembership(client, d, "remove")
}
