package jumpcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate, //optional
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
