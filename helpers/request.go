package smartcar

import (
	"io"
	"net/http"
)

// POSTRequest is a helper for sending POST requests to API.
func POSTRequest(url string, authorization string, data io.Reader) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, err
}
