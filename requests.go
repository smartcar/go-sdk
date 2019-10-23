package smartcar

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

// HTTPTimeout request set to 5 minutes. This is standard accross our SDKs.
const (
	defaultHTTPTimeout = time.Duration(5) * time.Minute
)

// requestParams is a helper struct to send accross the requests methods.
type requestParams struct {
	UnitSystem UnitSystem
}

type backendClientParams struct {
	ctx                        context.Context
	method, url, authorization string
	requestParams              requestParams
	body                       io.Reader
	target                     interface{}
}

// ResponseHeaders is a struct that has Smartcar's API response headers.
type ResponseHeaders struct {
	// Deprecated: Should use DataAge instead of Age
	Age        string     `json:"age,omitempty"`
	DataAge    string     `json:"dataAge,omitempty"`
	RequestID  string     `json:"requestId,omitempty"`
	UnitSystem UnitSystem `json:"unitSystem,omitempty"`
}

// Call creates a http request and calls the Exectue method with it.
func (c *backend) Call(params backendClientParams) error {
	req, err := c.newRequest(params)
	if err != nil {
		return err
	}

	return c.execute(req, params.target)
}

// execute executes a req and formats response.
func (c *backend) execute(req *http.Request, target interface{}) error {
	client := &http.Client{
		Timeout: defaultHTTPTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New(http.StatusText(res.StatusCode))
	}

	if err := c.formatHeadersResponse(res.Header, target); err != nil {
		return err
	}
	if err := c.formatBodyResponse(res.Body, target); err != nil {
		return err
	}
	return nil
}

// formatBodyResponse formats a response body using an interface.
func (c *backend) formatBodyResponse(body io.Reader, target interface{}) error {
	return json.NewDecoder(body).Decode(target)
}

// formatBodyResponse formats a response body using an interface.
func (c *backend) formatHeadersResponse(headers http.Header, target interface{}) error {
	unitSystem := unitSystems[headers.Get("Sc-Unit-System")]
	h := &ResponseHeaders{
		Age:        headers.Get("Sc-Data-Age"),
		DataAge:    headers.Get("Sc-Data-Age"),
		RequestID:  headers.Get("Sc-Request-Id"),
		UnitSystem: unitSystem,
	}

	b, err := json.Marshal(h)
	if err != nil {
		return errors.New("Unmarshall went wrong")
	}
	return json.Unmarshal(b, target)
}

// newRequest builds a new http request.
func (c *backend) newRequest(params backendClientParams) (*http.Request, error) {
	// Not supported in previous versions og go 1.13
	// req, err := http.NewRequestWithContext(params.ctx, params.method, params.url, params.body)
	req, err := http.NewRequest(params.method, params.url, params.body)
	req = req.WithContext(params.ctx)

	if err != nil {
		return nil, errors.New("Error creating New Request")
	}

	req.Header.Add("Authorization", params.authorization)
	req.Header.Add("User-Agent", getUserAgent())
	if params.body != nil {
		req.Header.Add("Content-Type", getBodyType(params.body))
	}
	if params.requestParams.UnitSystem != "" {
		req.Header.Add("SC-Unit-System", string(params.requestParams.UnitSystem))
	}

	return req, nil
}

func getBodyType(body io.Reader) string {
	switch body.(type) {
	case *bytes.Buffer:
		return "application/json"
	default:
		return "application/x-www-form-urlencoded"
	}
}

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
