package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// GET is a helper for sending GET requests to Smartcar.
func GET(url string, authorization string) (*http.Response, error) {
	client := &http.Client{}
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
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	return res, err
}

// HandleStatusCode returns error type from Smartcar depending on HTTP code.
func HandleStatusCode(code int) string {
	switch code {
	case 400:
		return "validation"
	case 401:
		return "authentication"
	case 403:
		return "permission"
	case 404:
		return "resource not found"
	case 409:
		return "vehicle state"
	case 429:
		return "rate limiting"
	case 430:
		return "monthly limit exceeded"
	case 500:
		return "server error"
	case 501:
		return "smartcar or vehicle not capable"
	default:
		return "gateway timeout"
	}
}

// GetJSONMapFromResponse returns a map[string]string from a []byte.
func GetJSONMapFromResponse(response []byte) (map[string]interface{}, error) {
	jsonResponse := make(map[string]interface{})
	jsonErr := json.Unmarshal(response, &jsonResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return nil, jsonErr
	}
	return jsonResponse, nil
}