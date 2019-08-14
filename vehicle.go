package smartcar

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
)

// Vehicle is initialized and passed into Smartcar vehicle functions.
type Vehicle struct {
	ID          string
	AccessToken string
	UnitSystem  string
}

type VehicleInfoResponse struct {
	ID    string
	Make  string
	Model string
	Year  int
}

type VehicleLocationResponse struct {
	Latitude  float64
	Longitude float64
}

type VehicleFuelEVResponse struct {
	AmountRemaining  float64 `json:"amountRemaining"`
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
}

type VehicleChargeResponse struct {
	IsPluggedIn bool
	State       string
}

type VehicleResponse struct {
	Status string `json:"status"`
}

// VehicleSetUnits sets the unit system that information of the vehicle will be returned in.
func VehicleSetUnits(vehicle *Vehicle, unit string) error {
	if !(unit == "metric" || unit == "imperial") {
		return errors.New("unit must either be metric or imperial")
	}
	vehicle.UnitSystem = unit
	return nil
}

func vehicleAPIRequest(vehicle Vehicle, endpoint string, httpType string, action string) (map[string]interface{}, error) {
	requestPath := constants.VehiclePath + vehicle.ID + endpoint
	vehicleURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   requestPath,
	}

	authorization := "Bearer " + vehicle.AccessToken

	var res *http.Response

	if httpType == "POST" {
		jsonRequest := map[string]interface{}{
			"action": action,
		}

		request, err := json.Marshal(jsonRequest)
		if err != nil {
			return nil, err
		}

		var resErr error
		res, resErr = requests.POST(vehicleURL.String(), authorization, bytes.NewBuffer(request))
		if resErr != nil {
			return nil, resErr
		}
	} else if httpType == "DELETE" {
		var resErr error
		res, resErr = requests.DELETE(vehicleURL.String(), authorization)
		if resErr != nil {
			return nil, resErr
		}
	} else {
		var resErr error
		res, resErr = requests.GET(vehicleURL.String(), authorization)
		if resErr != nil {
			return nil, resErr
		}
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	if res.StatusCode != 200 {
		var err Error
		jsonErr := json.Unmarshal(body, &err)
		if jsonErr != nil {
			jsonErr = errors.New("Decoding JSON error")
			return nil, jsonErr
		}
		err.Type = requests.HandleStatusCode(res.StatusCode)
		return nil, &Error{err.Type, err.Name, err.Message, err.Code}
	}

	jsonResponse := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &jsonResponse)
	if jsonErr != nil {
		jsonErr = errors.New("Decoding JSON error")
		return nil, jsonErr
	}

	return jsonResponse, nil
}

func VehicleInfo(vehicle Vehicle) (VehicleInfoResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "", "GET", "")
	if err != nil {
		return VehicleInfoResponse{}, err
	}

	var vehicleInfo VehicleInfoResponse
	err = mapstructure.Decode(response, &vehicleInfo)
	if err != nil {
		return VehicleInfoResponse{}, err
	}

	return vehicleInfo, nil
}

func VehicleVIN(vehicle Vehicle) (string, error) {
	response, err := vehicleAPIRequest(vehicle, "/vin", "GET", "")
	if err != nil {
		return "", err
	}

	return response["vin"].(string), nil
}

func VehicleOdometer(vehicle Vehicle) (float64, error) {
	response, err := vehicleAPIRequest(vehicle, "/odometer", "GET", "")
	if err != nil {
		return 0, err
	}

	return response["distance"].(float64), nil
}

func VehicleLocation(vehicle Vehicle) (VehicleLocationResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/location", "GET", "")
	if err != nil {
		return VehicleLocationResponse{}, err
	}

	var vehicleLocation VehicleLocationResponse
	err = mapstructure.Decode(response, &vehicleLocation)
	if err != nil {
		return VehicleLocationResponse{}, err
	}

	return vehicleLocation, nil
}

func VehicleFuel(vehicle Vehicle) (VehicleFuelEVResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/fuel", "GET", "")
	if err != nil {
		return VehicleFuelEVResponse{}, err
	}

	var vehicleFuel VehicleFuelEVResponse
	err = mapstructure.Decode(response, &vehicleFuel)
	if err != nil {
		return VehicleFuelEVResponse{}, err
	}

	return vehicleFuel, nil
}

func VehicleBattery(vehicle Vehicle) (VehicleFuelEVResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/battery", "GET", "")
	if err != nil {
		return VehicleFuelEVResponse{}, err
	}

	var vehicleBattery VehicleFuelEVResponse
	err = mapstructure.Decode(response, &vehicleBattery)
	if err != nil {
		return VehicleFuelEVResponse{}, err
	}

	return vehicleBattery, nil
}

func VehicleCharge(vehicle Vehicle) (VehicleChargeResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/charge", "GET", "")
	if err != nil {
		return VehicleChargeResponse{}, err
	}

	var vehicleCharge VehicleChargeResponse
	err = mapstructure.Decode(response, &vehicleCharge)
	if err != nil {
		return VehicleChargeResponse{}, err
	}

	return vehicleCharge, nil
}

func VehiclePermissions(vehicle Vehicle) ([]string, error) {
	response, err := vehicleAPIRequest(vehicle, "/permissions", "GET", "")
	if err != nil {
		return nil, err
	}

	var permissions []string
	err = mapstructure.Decode(response["permissions"], &permissions)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func VehicleHasPermission(vehicle Vehicle, permission string) (bool, error) {
	vehiclePermissions, err := VehiclePermissions(vehicle)
	if err != nil {
		return false, err
	}

	for _, vehiclePermission := range vehiclePermissions {
		if permission == vehiclePermission {
			return true, nil
		}
	}

	return false, nil
}

func VehicleHasPermissions(vehicle Vehicle, permissions []string) (bool, error) {
	vehiclePermissions, err := VehiclePermissions(vehicle)
	if err != nil {
		return false, err
	}

	set := make(map[string]bool)
	for _, value := range vehiclePermissions {
		value = strings.TrimPrefix(value, "required:")
		set[value] = true
	}

	for _, value := range permissions {
		value = strings.TrimPrefix(value, "required:")
		if hasPermission, found := set[value]; !found {
			return false, nil
		} else if !hasPermission {
			return false, nil
		}
	}

	return true, nil
}

func VehicleLock(vehicle Vehicle) (VehicleResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/security", "POST", "LOCK")
	if err != nil {
		return VehicleResponse{}, err
	}

	var vehicleLockResponse VehicleResponse
	err = mapstructure.Decode(response, &vehicleLockResponse)
	if err != nil {
		return VehicleResponse{}, err
	}

	return vehicleLockResponse, nil
}

func VehicleUnlock(vehicle Vehicle) (VehicleResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/security", "POST", "UNLOCK")
	if err != nil {
		return VehicleResponse{}, err
	}

	var vehicleUnlockResponse VehicleResponse
	err = mapstructure.Decode(response, &vehicleUnlockResponse)
	if err != nil {
		return VehicleResponse{}, err
	}

	return vehicleUnlockResponse, nil
}

func VehicleDisconnect(vehicle Vehicle) (VehicleResponse, error) {
	response, err := vehicleAPIRequest(vehicle, "/application", "DELETE", "")
	if err != nil {
		return VehicleResponse{}, err
	}

	var vehicleDisconnectResponse VehicleResponse
	err = mapstructure.Decode(response, &vehicleDisconnectResponse)
	if err != nil {
		return VehicleResponse{}, err
	}

	return vehicleDisconnectResponse, nil
}
