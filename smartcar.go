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

// AuthClient for interacting with Connect and API.
type AuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scope        []string
	TestMode     bool
}

// Error contains error type and message from Smartcar.
type Error struct {
	ErrorType string `json:"error"`
	Message   string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %s, message: %s", e.ErrorType, e.Message)
}

// VehicleIsCompatible checks compatibility for a vin with provided scopes.
func VehicleIsCompatible(vin string, auth AuthClient) (bool, error) {
	type CompatibleResponse struct {
		Compatible bool `json:"compatible"`
	}

	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   "v1.0/compatibility",
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

		return false, &Error{err.ErrorType, err.Message}
	}

	var compatibleResponse CompatibleResponse
	jsonErr := jsonDecoder.Decode(&compatibleResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return false, jsonErr
	}

	return compatibleResponse.Compatible, nil
}

// GetUserID returns the id of the vehicle owner
func GetUserID(accessToken string) (string, error) {
	type UserIDResponse struct {
		ID string `json:"id"`
	}

	compatiblityURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   "v1.0/user",
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
			return "", jsonErr
		}

		return "", &Error{err.ErrorType, err.Message}
	}

	var userIDResponse UserIDResponse
	jsonErr := jsonDecoder.Decode(&userIDResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return "", jsonErr
	}

	return userIDResponse.ID, nil
}

// GetVehicleIDs returns the uuids associated to the access token.
func GetVehicleIDs(accessToken string) ([]string, error) {
	type VehicleIDResponse struct {
		UUIDs []string `json:"vehicles"`
	}

	vehiclesURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   "v1.0/vehicles",
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
			return nil, jsonErr
		}

		return nil, &Error{err.ErrorType, err.Message}
	}

	var vehicleIDResponse VehicleIDResponse
	jsonErr := jsonDecoder.Decode(&vehicleIDResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return nil, jsonErr
	}

	return vehicleIDResponse.UUIDs, nil
}
