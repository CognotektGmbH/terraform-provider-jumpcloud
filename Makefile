install:
	mkdir -p ~/.terraform.d/plugins
	go build -o ~/.terraform.d/plugins/jumpcloud-provider-terraform
