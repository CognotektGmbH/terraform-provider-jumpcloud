package jumpcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func dataSourceJumpCloudUserRead(d *schema.ResourceData, m interface{}) error {
    config := m.(*jcapiv2.Configuration)
    userEmail := d.Get("email").(string)

    req, err := http.NewRequest("GET", "https://console.jumpcloud.com/api/systemusers", nil)
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("x-api-key", config.DefaultHeader["x-api-key"])

    resp, err := config.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to fetch users: %s", resp.Status)
    }

    var users []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
        return err
    }

    for _, user := range users {
        if user["email"] == userEmail {
            d.SetId(user["id"].(string))
            d.Set("id", user["id"].(string))
            return nil
        }
    }

    return errors.New("user not found")
}

// func dataSourceJumpCloudUserRead(d *schema.ResourceData, m interface{}) error {
// 	configv2 := m.(*jcapiv2.Configuration)
// 	clientv2 := jcapiv2.NewAPIClient(configv2)

// 	email := d.Get("email").(string)

// 	// Search for the user by email using the v1 API client
// 	queryParams := map[string]interface{}{
// 		"search": email,
// 	}
// 	resp, _, err := client.SystemusersApi.SystemusersGet(context.Background(), "", "", "", queryParams)
// 	if err != nil {
// 		return fmt.Errorf("Error searching for JumpCloud user by email %s: %s", email, err)
// 	}

// 	// Check if a user was found for the specified email
// 	if len(resp.Results) == 0 {
// 		return fmt.Errorf("No JumpCloud user found for email %s", email)
// 	}

// 	// Get the user details using the v2 API client
// 	user, err := client.Users().GetUser(resp.Results[0].ID)
// 	if err != nil {
// 		return fmt.Errorf("Error fetching JumpCloud user by ID %s: %s", resp.Results[0].Id, err)
// 	}

// 	// Set the resource data attributes
// 	d.SetId(user.ID)
// 	d.Set("id", user.ID)
// 	d.Set("username", user.UserName)
// 	d.Set("firstname", user.FirstName)
// 	d.Set("lastname", user.LastName)
// 	d.Set("email", user.Email)
// 	d.Set("enable_mfa", user.MFAEnabled)

// 	log.Printf("[INFO] JumpCloud user %s retrieved successfully", user.UserName)

// 	return nil
// }

