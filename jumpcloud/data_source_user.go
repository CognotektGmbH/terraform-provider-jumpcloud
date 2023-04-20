package jumpcloud

import (
	"context"
	"errors"
	"fmt"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceJumpCloudUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJumpCloudUserRead,
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getUserDetails(client *jcapiv1.APIClient, userID, email, username string) (*jcapiv1.Systemuserreturn, error) {
	ctx := context.TODO()
	res, _, err := client.SystemusersApi.SystemusersList(ctx, "", "", nil)

	if err != nil {
		return nil, err
	}

	var user *jcapiv1.Systemuserreturn

	if userID != "" {
		for _, u := range res.Results {
			if u.Id == userID {
				user = &u
				break
			}
		}
	} else if email != "" {
		for _, u := range res.Results {
			if u.Email == email {
				user = &u
				break
			}
		}
	} else if username != "" {
		for _, u := range res.Results {
			if u.Username == username {
				user = &u
				break
			}
		}
	}

	if user == nil {
		return nil, fmt.Errorf("no user found with the given query")
	}

	return user, nil
}

func dataSourceJumpCloudUserRead(d *schema.ResourceData, m interface{}) error {
    config := m.(*jcapiv1.Configuration)
    client := jcapiv1.NewAPIClient(config)
    userEmail := d.Get("email").(string)

    // Use the getUserDetails function to query user details using the userEmail
    user, err := getUserDetails(client, "", userEmail, "")

    // If an error occurs or no user is found, return an error
	if err != nil {
		return errors.New("user not found")
	}

    // Set the user ID in the Terraform resource data object
    d.SetId(user.Id)
    d.Set("id", user.Id)

    return nil
}