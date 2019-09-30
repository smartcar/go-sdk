// Package smartcar is the official Go SDK of the Smartcar API.
// Smartcar is the only vehicle API built for developers, by developers.
// Learn more about Smartcar here, https://smartcar.com/
package smartcar

import (
	"time"

	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/utils"
)

// GetUserID returns the user ID of the vehicle owner associated with the accessToken.
func GetUserID(accessToken string) (string, error) {

	authorization := utils.BuildBearerAuthorization(accessToken)
	userURL := utils.GetUserURL()

	res, resErr := requests.Request(requests.GET, userURL, authorization, "", nil)
	if resErr != nil {
		return "", resErr
	}

	formattedResponse := new(struct {
		ID string `json:"id"`
	})
	defer res.Body.Close()
	err := requests.FormatResponse(res.Body, formattedResponse)
	if err != nil {
		return "", err
	}

	return formattedResponse.ID, nil
}

// GetVehicleIds returns IDs of the vehicles associated with the token.
func GetVehicleIds(accessToken string) ([]string, error) {
	type vehicleIDResponse struct {
		UUIDs []string `json:"vehicles"`
	}

	vehicleURL := utils.GetVehicleURL()
	authorization := utils.BuildBearerAuthorization(accessToken)

	res, resErr := requests.Request(requests.GET, vehicleURL, authorization, "", nil)
	if resErr != nil {
		return nil, resErr
	}

	formattedResponse := new(vehicleIDResponse)
	defer res.Body.Close()
	err := requests.FormatResponse(res.Body, &formattedResponse)
	if err != nil {
		return nil, err
	}

	return formattedResponse.UUIDs, nil
}

// TokenIsExpired checks if a token is expired by passing in the token expiry time.
func TokenIsExpired(expiration time.Time) bool {
	// Give 10 seconds of
	return time.Now().After(expiration.Add(time.Second * 10))
}
