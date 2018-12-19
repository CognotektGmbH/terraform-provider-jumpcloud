# JumpCloud Terraform Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.11 (to build the provider plugin)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/cognotektgmbh/terraform-provider-jumpcloud`

```sh
mkdir -p $GOPATH/src/github.com/cognotektgmbh
cd $GOPATH/src/github.com/cognotektgmbh
git clone git@github.com:cognotektgmbh/terraform-provider-jumpcloud
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/cognotektgmbh/terraform-provider-jumpcloud
make build
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.
