package main

import (
	"github.com/cognotektgmbh/terraform-provider-jumpcloud/jumpcloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return jumpcloud.Provider()
		},
	})
}
