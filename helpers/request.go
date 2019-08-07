package smartcar

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// POSTRequest is a helper for sending POST requests to API.
func POSTRequest(url string, authorization string, data io.Reader) (*json.Decoder, error) {
	var err error

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body), err
}
