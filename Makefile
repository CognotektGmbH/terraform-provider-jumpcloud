install: deps
	mkdir -p ~/.terraform.d/plugins
	go build -o ~/.terraform.d/plugins/terraform-provider-jumpcloud
deps:
	dep ensure
update-deps:
	dep ensure -update
build: deps
	go build .
testacc: deps
	TF_ACC=true go test -v ./...
