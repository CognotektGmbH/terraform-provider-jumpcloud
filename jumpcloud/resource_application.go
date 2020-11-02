package jumpcloud

import (
	"context"
	"fmt"
	// "encoding/json"
	// "net/http"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceAppliaction() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,
		Schema: map[string]*schema.Schema{
			"display_label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata_xml": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"acs_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_entity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"idp_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sp_entity_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			
			// Currently, only the options necessary for our use case are implemented
			// JumpCloud offers a lot more
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	payload := jcapiv1.Application{
		Name: 		  "aws", // Api documentation don't declare the application types we can use
		SsoUrl:       d.Get("sso_url").(string),
		DisplayLabel: d.Get("display_label").(string),
		Config: &jcapiv1.ApplicationConfig{
			AcsUrl: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("acs_url").(string)},
			//IdpCertificate: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_certificate").(string)},
			IdpEntityId: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_entity_id").(string)},
			//IdpPrivateKey: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_private_key").(string)},
			SpEntityId: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("sp_entity_id").(string)},
		},
	}

	req := map[string]interface{}{
		"body": payload,
	}

	log.Println(req)

	returnstruc, _, err := client.ApplicationsApi.ApplicationsPost(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(returnstruc.Id)
	return resourceApplicationRead(d, m)
}

func resourceApplicationRead(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.ApplicationsApi.ApplicationsGet(context.TODO(), d.Id(), nil)

	// If the object does not exist, unset the ID
	if err != nil {
		if err.Error() == "EOF" {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(res.Id)

	if err := d.Set("display_label", res.DisplayLabel); err != nil {
		return err
	}
	if err := d.Set("sso_url", res.SsoUrl); err != nil {
		return err
	}
	
	if res.Id != "" {
		log.Println("[INFO] response ID is ", res.Id)
		orgId := configv1.DefaultHeader["x-org-id"]
		apiKey := configv1.DefaultHeader["x-api-key"]

		metadataXml, err := GetApplicationMetadataXml(orgId, res.Id, apiKey)
		if err != nil {
			return err
		}

		if err := d.Set("metadata_xml", metadataXml); err != nil {
			return err
		}
	} else {
		log.Println("[INFO] no ID in response, skipping metadata XML retrieval")
	}

	return nil
}


func resourceApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	payload := jcapiv1.Application{
		Active:       true,
		SsoUrl:       d.Get("sso_url").(string),
		DisplayLabel: d.Get("display_label").(string),
		Config: &jcapiv1.ApplicationConfig{
			AcsUrl: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("acs_url").(string)},
			IdpCertificate: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_certificate").(string)},
			IdpEntityId: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_entity_id").(string)},
			//IdpPrivateKey: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("idp_private_key").(string)},
			SpEntityId: &jcapiv1.ApplicationConfigAcsUrl{Value: d.Get("sp_entity_id").(string)},
		},
	}

	req := map[string]interface{}{
		"body": payload,
	}
	_, _, err := client.ApplicationsApi.ApplicationsPut(context.TODO(),
		d.Id(), req)
	if err != nil {
		return err
	}
	return resourceApplicationRead(d, m)
}

func resourceApplicationDelete(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.ApplicationsApi.ApplicationsDelete(context.TODO(),
		d.Id(), nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting application:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}