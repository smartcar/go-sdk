package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"time"
)

// constant methods
const (
	GET  = "GET"
	POST = "POST"
)

// Timeout requet
const (
	Timeout = 300
)

// execute TODO: Add decription
func execute(req *http.Request) (http.Response, error) {
	client := &http.Client{
		Timeout: Timeout * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		return http.Response{}, errors.New("Execute returned error")
	}

	if res.StatusCode != 200 {
		return http.Response{}, errors.New(http.StatusText(res.StatusCode))
	}

	return *res, nil
}

// Request creates a request
func Request(method string, url string, authorization string, data io.Reader) (http.Response, error) {

	// Build Request
	req, reqErr := http.NewRequest(method, url, nil)
	if reqErr != nil {
		return http.Response{}, reqErr
	}
	req.Header.Add("Authorization", authorization)

	// Only Add content-type for POST/PUT requests.
	if data != nil {
		if reflect.TypeOf(data).String() == "*bytes.Buffer" {
			req.Header.Add("Content-Type", "application/json")
		} else {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Send request
	return execute(req)
}

//FormatResponse TODO: Add description
func FormatResponse(body io.Reader, formatter interface{}) error {
	return json.NewDecoder(body).Decode(formatter)
}
