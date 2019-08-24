// Package smartcar is the official Go SDK of the Smartcar API.
// Smartcar is the only vehicle API built for developers, by developers.
// Learn more about Smartcar here, https://smartcar.com/
package smartcar

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
)

// AuthClient is used to store your auth credentials when authenticating with Smartcar.
type AuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        []string
	TestMode     bool
}

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

// VehicleIsCompatible checks compatibility for a vehicle VIN with Smartcar for the provided scopes.
// It takes a VIN and auth credentials and will return a bool indicating compatibility.
func VehicleIsCompatible(auth AuthClient, vin string) (bool, error) {
	type CompatibleResponse struct {
		Compatible bool `json:"compatible"`
	}

	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.CompatibilityPath,
	}

	query := compatiblityURL.Query()
	query.Set("vin", vin)
	query.Set("scope", strings.Join(auth.Scope, " "))
	compatiblityURL.RawQuery = query.Encode()

	authString := auth.ClientID + ":" + auth.ClientSecret
	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

	res, resErr := requests.GET(compatiblityURL.String(), encodedAuth)
	if resErr != nil {
		return false, resErr
	}
	defer res.Body.Close()
	jsonDecoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		var err Error
		jsonErr := jsonDecoder.Decode(&err)
		if jsonErr != nil {
			jsonErr = errors.New("Decoding JSON error")
			return false, jsonErr
		}
		return false, &Error{err.Name, err.Message, err.Code}
	}

	var compatibleResponse CompatibleResponse
	jsonErr := jsonDecoder.Decode(&compatibleResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return false, jsonErr
	}

	return compatibleResponse.Compatible, nil
}

// GetUserID returns the ID of the vehicle owner using the access token.
func GetUserID(accessToken string) (string, error) {
	type UserIDResponse struct {
		ID string `json:"id"`
	}

	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.UserPath,
	}

	authorization := "Bearer " + accessToken
	res, resErr := requests.GET(compatiblityURL.String(), authorization)
	if resErr != nil {
		return "", resErr
	}
	defer res.Body.Close()
	jsonDecoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		var err Error
		jsonErr := jsonDecoder.Decode(&err)
		if jsonErr != nil {
			jsonErr = errors.New("Decoding JSON error")
			return "", jsonErr
		}
		return "", &Error{err.Name, err.Message, err.Code}
	}

	var userIDResponse UserIDResponse
	jsonErr := jsonDecoder.Decode(&userIDResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return "", jsonErr
	}

	return userIDResponse.ID, nil
}

// GetVehicleIDs uses an access token and returns the IDs of the vehicles associated with the token.
func GetVehicleIDs(accessToken string) ([]string, error) {
	type VehicleIDResponse struct {
		UUIDs []string `json:"vehicles"`
	}

	vehiclesURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   constants.VehiclePath,
	}

	authorization := "Bearer " + accessToken
	res, resErr := requests.GET(vehiclesURL.String(), authorization)
	if resErr != nil {
		return nil, resErr
	}
	defer res.Body.Close()
	jsonDecoder := json.NewDecoder(res.Body)

	if res.StatusCode != 200 {
		var err Error
		jsonErr := jsonDecoder.Decode(&err)
		if jsonErr != nil {
			jsonErr = errors.New("Decoding JSON error")
			return nil, jsonErr
		}
		return nil, &Error{err.Name, err.Message, err.Code}
	}

	var vehicleIDResponse VehicleIDResponse
	jsonErr := jsonDecoder.Decode(&vehicleIDResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return nil, jsonErr
	}

	return vehicleIDResponse.UUIDs, nil
}
