install: deps
	mkdir -p ~/.terraform.d/plugins
	go build -o ~/.terraform.d/plugins/jumpcloud-provider-terraform
deps:
	dep ensure
update-deps:
	dep ensure -update
