terraform {
  required_providers {
    jumpcloud = {
      source = "cheelim1/jumpcloud"
    }
  }
}

provider "jumpcloud" {
  api_key = "MY_API_KEY"
}
