package smartcar

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"

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
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int64  `json:"year"`
}

// VehicleSetUnits sets the unit system that information of the vehicle will be returned in.
func VehicleSetUnits(vehicle *Vehicle, unit string) error {
	if !(unit == "metric" || unit == "imperial") {
		return errors.New("unit must either be metric or imperial")
	}
	vehicle.UnitSystem = unit
	return nil
}

func vehicleGETAPIRequest(vehicle Vehicle, endpoint string) (map[string]interface{}, error) {
	requestPath := constants.VehiclePath + vehicle.ID + endpoint
	vehicleURL := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   requestPath,
	}

	authorization := "Bearer " + vehicle.AccessToken

	res, resErr := requests.GET(vehicleURL.String(), authorization)
	if resErr != nil {
		return nil, resErr
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
	response, err := vehicleGETAPIRequest(vehicle, "")
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
	response, err := vehicleGETAPIRequest(vehicle, "/vin")
	if err != nil {
		return "", err
	}

	return response["vin"].(string), nil
}

func VehicleOdometer(vehicle Vehicle) (float64, error) {
	response, err := vehicleGETAPIRequest(vehicle, "/vin")
	if err != nil {
		return 0, err
	}

	return response["distance"].(float64), nil
}
