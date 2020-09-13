build:
	go build -o terraform-provider-jumpcloud
testacc: build
	TF_ACC=true go test -v ./...
