package jumpcloud

import jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"

const (
	Accept      = "application/json"
	ContentType = Accept
)

type Config struct {
	APIKey string
}

func (c *Config) Client() (interface{}, error) {
	config := jcapiv2.NewConfiguration()
	config.AddDefaultHeader("x-api-key", c.APIKey)

	// Instantiate the API client
	return jcapiv2.NewAPIClient(config), nil
}
