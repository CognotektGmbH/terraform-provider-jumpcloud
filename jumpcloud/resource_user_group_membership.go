package jumpcloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUserGroupMembership() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a resource for managing user group memberships.",
		Create: resourceUserGroupMembershipCreate,
		Read:   resourceUserGroupMembershipRead,
		// We must not have an update routine as the association cannot be updated.
		// Any change in one of the elements forces a recreation of the resource
		Update:        nil,
		Delete: resourceUserGroupMembershipDelete,
		Schema: map[string]*schema.Schema{
			"userid": {
				Description: "The ID of the `resource_user` object.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"groupid": {
				Description: "The ID of the `resource_user_group` object.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
    _ = d.Set("groupid", s[0])
    _ = d.Set("userid", s[1])

    config := m.(*jcapiv2.Configuration)
    client := jcapiv2.NewAPIClient(config)

    // Check if the user is already a member of the group
    isMember, err := checkUserGroupMembership(client, d.Get("groupid").(string), d.Get("userid").(string))
    if err != nil {
        return nil, err
    }

    if isMember {
        d.SetId(d.Get("groupid").(string) + "/" + d.Get("userid").(string))
        return []*schema.ResourceData{d}, nil
    }

    return nil, fmt.Errorf("User %s is not a member of group %s", d.Get("userid").(string), d.Get("groupid").(string))
}

func checkUserGroupMembership(client *jcapiv2.APIClient, groupID, userID string) (bool, error) {
    for i := 0; ; i++ {
        optionals := map[string]interface{}{
            "groupId": groupID,
            "limit":   int32(100),
            "skip":    int32(i * 100),
        }

        graphconnect, _, err := client.UserGroupMembersMembershipApi.GraphUserGroupMembersList(
            context.TODO(), groupID, "", "", optionals)
        if err != nil {
            return false, err
        }

        for _, v := range graphconnect {
            if v.To.Id == userID {
                return true, nil
            }
        }

        // Break the loop if the number of members in the current batch is less than 100
        if len(graphconnect) < 100 {
            break
        }
    }

    return false, nil
}

func modifyUserGroupMembership(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) error {

	payload := jcapiv2.UserGroupMembersReq{
		Op:    action,
		Type_: "user",
		Id:    d.Get("userid").(string),
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupMembersMembershipApi.GraphUserGroupMembersPost(
		context.TODO(), d.Get("groupid").(string), "", "", req)

		return err
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

	for i := 0; i < 20; i++ { // Prevent infite loop

		optionals := map[string]interface{}{
			"groupId": d.Get("groupid").(string),
			"limit":   int32(100),
			"skip":     int32(i * 100),
		}

		graphconnect, _, err := client.UserGroupMembersMembershipApi.GraphUserGroupMembersList(
			context.TODO(), d.Get("groupid").(string), "", "", optionals)
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

		if len(graphconnect) < 100 {
			break
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Instead of unsetting the ID, return an error to let Terraform retry
	return fmt.Errorf("User ID %s not found in group ID %s", d.Get("userid").(string), d.Get("groupid").(string))
}

func resourceUserGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)
	return modifyUserGroupMembership(client, d, "remove")
}
