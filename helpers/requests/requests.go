package requests

import (
	"io"
	"net/http"
	"reflect"
	"time"
)

// GET is a helper for sending GET requests to Smartcar.
func GET(url string, authorization string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 360 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authorization)

	res, err := client.Do(req)
	return res, err
}

// POST is a helper for sending POST requests to Smartcar.
func POST(url string, authorization string, data io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: 360 * time.Second,
	}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authorization)

	if reflect.TypeOf(data).String() == "*bytes.Buffer" {
		req.Header.Add("Content-Type", "application/json")
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := client.Do(req)
	return res, err
}

// DELETE is a helper for sending DELETE requests to Smartcar.
func DELETE(url string, authorization string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 360 * time.Second,
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authorization)

	res, err := client.Do(req)
	return res, err
}
