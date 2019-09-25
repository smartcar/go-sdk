// Package smartcar is the official Go SDK of the Smartcar API.
// Smartcar is the only vehicle API built for developers, by developers.
// Learn more about Smartcar here, https://smartcar.com/
package smartcar

import (
	"fmt"
	"net/url"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
)

// Error returns error names, messages and codes from requests to the Smartcar API.
type Error struct {
	Name    string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("error: %s, message: %s, code: %s", e.Name, e.Message, e.Code)
	}
	return fmt.Sprintf("error: %s, message: %s", e.Name, e.Message)
}

// GetUserID returns the ID of the vehicle owner using the access token.
func GetUserID(accessToken string) (string, error) {
	type userIDResponse struct {
		ID string `json:"id"`
	}

	authorization := "Bearer " + accessToken
	URL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}

	res, err := requests.Request(requests.GET, URL.String(), authorization, nil)
	if err != nil {
		return "", err
	}

	formattedResponse := new(userIDResponse)
	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	return formattedResponse.ID, nil
}

// GetVehicleIds uses an access token and returns the IDs of the vehicles associated with the token.
func GetVehicleIds(accessToken string) ([]string, error) {
	type vehicleIDResponse struct {
		UUIDs []string `json:"vehicles"`
	}

	URL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.VehiclePath,
	}
	authorization := "Bearer " + accessToken

	res, resErr := requests.Request("GET", URL.String(), authorization, nil)
	if resErr != nil {
		return nil, resErr
	}

	formattedResponse := new(vehicleIDResponse)

	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	return formattedResponse.UUIDs, nil
}
