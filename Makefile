build:
	GO111MODULE=on go build -o terraform-provider-jumpcloud
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o terraform-provider-jumpcloud_amd64-darwin

testacc: build
	TF_ACC=true go test -v ./...
