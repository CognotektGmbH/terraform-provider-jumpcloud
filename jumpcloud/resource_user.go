package jumpcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ldap_binding_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"passwordless_sudo": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password_never_expires": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sudo": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"suspended": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"phone_number": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"mobile",
								"work",
								"work_fax",
							}, false),
						},
					},
				},
			},
			// Currently, only the options necessary for our use case are implemented
			// JumpCloud offers a lot more
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// We receive a v2config from the TF base code but need a v1config to continue. So, we take the only
// preloaded element (the x-api-key) and populate the v1config with it.
func convertV2toV1Config(v2config *jcapiv2.Configuration) *jcapiv1.Configuration {
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", v2config.DefaultHeader["x-api-key"])
	if v2config.DefaultHeader["x-org-id"] != "" {
		configv1.AddDefaultHeader("x-org-id", v2config.DefaultHeader["x-org-id"])
	}
	return configv1
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	var phoneNumbers []jcapiv1.SystemuserputpostPhoneNumbers
	phoneNumbersRaw, _ := json.Marshal(expandPhoneNumbers(d.Get("phone_number").([]interface{})))
	if err := json.Unmarshal(phoneNumbersRaw, &phoneNumbers); err != nil {
		return err
	}

	payload := jcapiv1.Systemuserputpost{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		Password:                    d.Get("password").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
		LdapBindingUser:             d.Get("ldap_binding_user").(bool),
		Sudo:                        d.Get("sudo").(bool),
		Suspended:                   d.Get("suspended").(bool),
		PasswordNeverExpires:        d.Get("password_never_expires").(bool),
		PhoneNumbers:                phoneNumbers,
	}
	req := map[string]interface{}{
		"body": payload,
	}
	returnstruc, _, err := client.SystemusersApi.SystemusersPost(context.TODO(),
		"", "", req)
	if err != nil {
		return err
	}
	d.SetId(returnstruc.Id)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersGet(context.TODO(),
		d.Id(), "", "", nil)

	// If the object does not exist in our infrastructure, we unset the ID
	// Unfortunately, the http request returns 200 even if the resource does not exist
	if err != nil {
		if err.Error() == "EOF" {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(res.Id)

	if err := d.Set("username", res.Username); err != nil {
		return err
	}
	if err := d.Set("email", res.Email); err != nil {
		return err
	}
	if err := d.Set("firstname", res.Firstname); err != nil {
		return err
	}
	if err := d.Set("lastname", res.Lastname); err != nil {
		return err
	}
	if err := d.Set("enable_mfa", res.EnableUserPortalMultifactor); err != nil {
		return err
	}
	if err := d.Set("ldap_binding_user", res.LdapBindingUser); err != nil {
		return err
	}
	if err := d.Set("password_never_expires", res.PasswordNeverExpires); err != nil {
		return err
	}
	if err := d.Set("sudo", res.Sudo); err != nil {
		return err
	}
	if err := d.Set("suspended", res.Suspended); err != nil {
		return err
	}
	if err := d.Set("phone_number", flattenPhoneNumbers(res.PhoneNumbers)); err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	var phoneNumbers []jcapiv1.SystemuserputPhoneNumbers
	phoneNumbersRaw, _ := json.Marshal(expandPhoneNumbers(d.Get("phone_number").([]interface{})))
	if err := json.Unmarshal(phoneNumbersRaw, &phoneNumbers); err != nil {
		return err
	}

	// The code from the create function is almost identical, but the structure is different :
	// jcapiv1.Systemuserput != jcapiv1.Systemuserputpost
	payload := jcapiv1.Systemuserput{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		Password:                    d.Get("password").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
		LdapBindingUser:             d.Get("ldap_binding_user").(bool),
		Sudo:                        d.Get("sudo").(bool),
		Suspended:                   d.Get("suspended").(bool),
		PasswordNeverExpires:        d.Get("password_never_expires").(bool),
		PhoneNumbers:                phoneNumbers,
	}

	req := map[string]interface{}{
		"body": payload,
	}
	_, _, err := client.SystemusersApi.SystemusersPut(context.TODO(),
		d.Id(), "", "", req)
	if err != nil {
		return err
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersDelete(context.TODO(),
		d.Id(), "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting user group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}
