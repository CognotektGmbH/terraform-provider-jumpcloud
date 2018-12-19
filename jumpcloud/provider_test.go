package jumpcloud

// see https://www.terraform.io/docs/plugins/provider.html#provider

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
