package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

// constant request methods
const (
	DELETE = "DELETE"
	GET    = "GET"
	POST   = "POST"
)

// Timeout request set to 5 minutes
const (
	RequestTimeout = time.Duration(5) * time.Minute
)

// Request builds a request and calls execute to send request
func Request(method, url, authorization, unitSystem string, body io.Reader) (http.Response, error) {

	// Build Request
	req, reqErr := http.NewRequest(method, url, body)
	if reqErr != nil {
		return http.Response{}, reqErr
	}

	// Add Headers
	addHeader := addHeader(req)
	addHeader("Authorization", authorization)
	addHeader("User-Agent", getUserAgent())
	if body != nil {
		addHeader("Content-Type", getBodyType(body))
	}
	if unitSystem != "" {
		addHeader("SC-Unit-System", unitSystem)
	}

	// Send request
	return execute(req)
}

//FormatResponse formats a response body using a formatter
func FormatResponse(body io.Reader, formatter interface{}) error {
	return json.NewDecoder(body).Decode(formatter)
}

func execute(req *http.Request) (http.Response, error) {
	client := &http.Client{
		Timeout: RequestTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	if res.StatusCode != 200 {
		return http.Response{}, errors.New(http.StatusText(res.StatusCode))
	}

	return *res, nil
}

func addHeader(req *http.Request) func(string, string) {
	return func(key, value string) {
		req.Header.Add(key, value)
	}
}

func getBodyType(body io.Reader) string {
	switch body.(type) {
	case *bytes.Buffer:
		return "application/json"
	default:
		return "application/x-www-form-urlencoded"
	}
}

// getUserAgent returns the formatted user agent to send to smartcar api
func getUserAgent() string {
	arch := runtime.GOARCH
	os := runtime.GOOS
	version := runtime.Version()

	// NOTE:
	// This is only supported after Go.12
	// if debug.ReadBuildInfo != nil {
	// 	if bi, ok := debug.ReadBuildInfo(); ok {
	// 		return fmt.Sprintf(
	// 			"Smartcar/%s (%s; %s) Go %s",
	// 			bi.Main.Version,
	// 			os,
	// 			arch,
	// 			version,
	// 		)
	// 	}
	// }
	return fmt.Sprintf("Smartcar/unknown (%s;%s) Go %s", os, arch, version)
}
