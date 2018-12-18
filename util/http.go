package util

import (
	"io"
	"net/http"
)

// RequestHTTP encapsulates boilerplate related to making an HTTP request
func RequestHTTP(method string, headers map[string]string,
	URL string, requestBody io.Reader) (res *http.Response, err error) {

	req, err := http.NewRequest(method, URL, requestBody)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return http.DefaultClient.Do(req)
}
