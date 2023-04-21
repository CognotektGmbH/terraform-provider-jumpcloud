package jumpcloud

import (
	"context"
	"fmt"

	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func resourceUserGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a resource for associating a JumpCloud user group to objects like SSO applications, G Suite, Office 365, LDAP and more.",
		Create:       resourceUserGroupAssociationCreate,
		Read:         resourceUserGroupAssociationRead,
		Update:       nil,
		Delete:       resourceUserGroupAssociationDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the `resource_user_group` resource.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"object_id": {
				Description: "The ID of the object to associate to the group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "The type of the object to associate to the given group. Possible values: `active_directory`, `application`, `command`, `g_suite`, `ldap_server`, `office_365`, `policy`, `radius_server`, `system`, `system_group`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errors []error) {
					allowedValues := []string{
						"active_directory",
						"application",
						"command",
						"g_suite",
						"ldap_server",
						"office_365",
						"policy",
						"radius_server",
						"system",
						"system_group",
					}

					v := val.(string)
					if !stringInSlice(v, allowedValues) {
						errors = append(errors, fmt.Errorf("%q must be one of %q", key, allowedValues))
					}
					return
				},
			},
		},
	}
}


func modifyUserGroupAssociation(client *jcapiv2.APIClient,
	d *schema.ResourceData, action string) diag.Diagnostics {

	payload := jcapiv2.UserGroupGraphManagementReq{
		Op:    action,
		Type_: d.Get("type").(string),
		Id:    d.Get("object_id").(string),
	}

	req := map[string]interface{}{
		"body": payload,
	}

	_, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsPost(
		context.TODO(), d.Get("group_id").(string), "", "", req)

	return diag.FromErr(err)
}

func resourceUserGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	diags := modifyUserGroupAssociation(client, d, "add")
	if diags.HasError() {
		return fmt.Errorf("Error creating user group association: %v", diags)
	}
	return resourceUserGroupAssociationRead(d, meta)
}


func resourceUserGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)

	optionals := map[string]interface{}{
		"groupId": d.Get("group_id").(string),
		"limit":   int32(100),
	}

	graphconnect, _, err := client.UserGroupAssociationsApi.GraphUserGroupAssociationsList(
		context.TODO(), d.Get("group_id").(string), "", "", []string{d.Get("type").(string)}, optionals)
	if err != nil {
		return err
	}

	// the ID of the specified object is buried in a complex construct
	for _, v := range graphconnect {
		if v.To.Id == d.Get("object_id") {
			resourceId := d.Get("group_id").(string) + "/" + d.Get("object_id").(string)
			d.SetId(resourceId)
			return nil
		}
	}

	// element does not exist; unset ID
	d.SetId("")
	return nil
}

func resourceUserGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*jcapiv2.Configuration)
	client := jcapiv2.NewAPIClient(config)
	diags := modifyUserGroupAssociation(client, d, "remove")
	if diags.HasError() {
		return fmt.Errorf("Error deleting user group association: %v", diags)
	}
	return nil
}