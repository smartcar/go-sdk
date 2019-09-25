package smartcar

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/smartcar/go-sdk/helpers/constants"
	"github.com/smartcar/go-sdk/helpers/requests"
)

// Vehicle contains vehicle request information such as ID and AccessToken and used for Smartcar vehicle functions.
type Vehicle struct {
	ID          string
	AccessToken string
	UnitSystem  string
}

// VehicleInfoResponse contains the vehicle information response returned from the VehicleInfo function.
type VehicleInfoResponse struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

// VehicleLocationResponse contains the vehicle location response returned from the VehicleLocation function.
type VehicleLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// VehicleFuelEVResponse contains the vehicle FuelEV response returned from the VehicleFuel and VehicleBattery functions.
type VehicleFuelEVResponse struct {
	AmountRemaining  float64 `json:"amountRemaining"`
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
}

// VehicleChargeResponse contains the vehicle charging state response returned from the VehicleCharge function.
type VehicleChargeResponse struct {
	IsPluggedIn bool   `json:"isPluggedIn"`
	State       string `json:"state"`
}

// VehicleResponse contains a general vehicle status response returned from the VehicleLock, VehicleUnlock and VehicleDisconnect functions.
type VehicleResponse struct {
	Status string `json:"status"`
}

// SetUnits takes a vehicle and sets the unit system that information for the vehicle will be returned in.
func (vehicle *Vehicle) SetUnits(unit string) error {
	if !(unit == "metric" || unit == "imperial") {
		return errors.New("unit must either be metric or imperial")
	}
	vehicle.UnitSystem = unit
	return nil
}

// Request is an internal functions used to make requests to Smartcar's vehicle API.
func (vehicle *Vehicle) request(path string, method string, data io.Reader) (http.Response, error) {
	// Build url
	requestPath := constants.VehiclePath + vehicle.ID + path
	url := url.URL{
		Scheme: constants.APIScheme,
		Host:   constants.APIHost,
		Path:   requestPath,
	}
	authorization := "Bearer " + vehicle.AccessToken

	// Send request
	return requests.Request(method, url.String(), authorization, data)
}

// Info uses a Vehicle and returns vehicle information from Smartcar in a VehicleInfoResponse.
func (vehicle *Vehicle) Info() VehicleInfoResponse {
	res, err := vehicle.request("/", requests.GET, nil)
	if err != nil {
		return VehicleInfoResponse{}
	}

	formattedResponse := new(VehicleInfoResponse)

	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	return *formattedResponse
}

// VIN uses a Vehicle and returns the vehicle's VIN from Smartcar in a string.
func (vehicle *Vehicle) VIN() string {
	res, err := vehicle.request("/vin", requests.GET, nil)
	if err != nil {
		return ""
	}

	type vinResponse struct {
		Vin string
	}
	formattedResponse := new(vinResponse)

	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	return formattedResponse.Vin
}

// Odometer uses a Vehicle and returns the vehicle's odometer reading from Smartcar in a float64.
func (vehicle *Vehicle) Odometer() float64 {
	res, err := vehicle.request("/odometer", requests.GET, nil)
	if err != nil {
		return -1
	}

	type odometerResponse struct {
		Distance float64
	}
	formattedResponse := new(odometerResponse)

	defer res.Body.Close()
	requests.FormatResponse(res.Body, formattedResponse)

	return formattedResponse.Distance
}

// // VehicleLocation uses a Vehicle and returns the vehicle's location from Smartcar in a VehicleLocationResponse.
// func (vehicle Vehicle) Location() (VehicleLocationResponse, error) {
// 	response, err := vehicle.request("/location", requests.GET, "")
// 	if err != nil {
// 		return VehicleLocationResponse{}, err
// 	}

// 	var vehicleLocation VehicleLocationResponse
// 	err = mapstructure.Decode(response, &vehicleLocation)
// 	if err != nil {
// 		return VehicleLocationResponse{}, err
// 	}

// 	return vehicleLocation, nil
// }

// // VehicleFuel uses a Vehicle and returns the vehicle's fuel level from Smartcar in a VehicleFuelEVResponse.
// func (vehicle Vehicle) Fuel() (VehicleFuelEVResponse, error) {
// 	response, err := vehicle.request("/fuel", requests.GET, "")
// 	if err != nil {
// 		return VehicleFuelEVResponse{}, err
// 	}

// 	var vehicleFuel VehicleFuelEVResponse
// 	err = mapstructure.Decode(response, &vehicleFuel)
// 	if err != nil {
// 		return VehicleFuelEVResponse{}, err
// 	}

// 	return vehicleFuel, nil
// }

// // VehicleBattery uses a Vehicle and returns the vehicle's battery level from Smartcar in a VehicleFuelEVResponse.
// func (vehicle Vehicle) Battery() (VehicleFuelEVResponse, error) {
// 	response, err := vehicle.request("/battery", requests.GET, "")
// 	if err != nil {
// 		return VehicleFuelEVResponse{}, err
// 	}

// 	var vehicleBattery VehicleFuelEVResponse
// 	err = mapstructure.Decode(response, &vehicleBattery)
// 	if err != nil {
// 		return VehicleFuelEVResponse{}, err
// 	}

// 	return vehicleBattery, nil
// }

// // VehicleCharge uses a Vehicle and returns the vehicle's charging status from Smartcar in a VehicleChargeResponse.
// func (vehicle Vehicle) Charge() (VehicleChargeResponse, error) {
// 	response, err := vehicle.request("/charge", requests.GET, "")
// 	if err != nil {
// 		return VehicleChargeResponse{}, err
// 	}

// 	var vehicleCharge VehicleChargeResponse
// 	err = mapstructure.Decode(response, &vehicleCharge)
// 	if err != nil {
// 		return VehicleChargeResponse{}, err
// 	}

// 	return vehicleCharge, nil
// }

// // VehiclePermissions uses a Vehicle and returns the vehicle's authorized permissions in a []string.
// func (vehicle Vehicle) Permissions() ([]string, error) {
// 	response, err := vehicle.request("/permissions", requests.GET, "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var permissions []string
// 	err = mapstructure.Decode(response["permissions"], &permissions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return permissions, nil
// }

// // VehicleHasPermissions uses a Vehicle and a slice of permissions and returns whether the vehicle has the specified permissions.
// func (vehicle Vehicle) HasPermissions(permissions ...string) (bool, error) {
// 	vehiclePermissions, err := vehicle.Permissions()
// 	if err != nil {
// 		return false, err
// 	}

// 	set := make(map[string]bool)
// 	for _, value := range vehiclePermissions {
// 		value = strings.TrimPrefix(value, "required:")
// 		set[value] = true
// 	}

// 	for _, value := range permissions {
// 		value = strings.TrimPrefix(value, "required:")
// 		if hasPermission, found := set[value]; !found {
// 			return false, nil
// 		} else if !hasPermission {
// 			return false, nil
// 		}
// 	}

// 	return true, nil
// }

// // VehicleLock uses a Vehicle to send a vehicle lock request to Smartcar.
// // It returns whether the request was successful in a VehicleResponse.
// func (vehicle Vehicle) Lock() (VehicleResponse, error) {
// 	response, err := vehicle.request("/security", "POST", "LOCK")
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	var vehicleLockResponse VehicleResponse
// 	err = mapstructure.Decode(response, &vehicleLockResponse)
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	return vehicleLockResponse, nil
// }

// // VehicleUnlock uses a Vehicle to send a vehicle unlock request to Smartcar.
// // It returns whether the request was successful in a VehicleResponse.
// func (vehicle Vehicle) Unlock() (VehicleResponse, error) {
// 	response, err := vehicle.request("/security", "POST", "UNLOCK")
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	var vehicleUnlockResponse VehicleResponse
// 	err = mapstructure.Decode(response, &vehicleUnlockResponse)
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	return vehicleUnlockResponse, nil
// }

// // VehicleDisconnect uses a Vehicle to disconnect it from the application.
// // It returns whether the disconnect was successful in a VehicleResponse.
// func (vehicle Vehicle) Disconnect() (VehicleResponse, error) {
// 	response, err := vehicle.request("/application", "DELETE", "")
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	var vehicleDisconnectResponse VehicleResponse
// 	err = mapstructure.Decode(response, &vehicleDisconnectResponse)
// 	if err != nil {
// 		return VehicleResponse{}, err
// 	}

// 	return vehicleDisconnectResponse, nil
// }

// // VehicleIsCompatible checks compatibility for a vehicle VIN with Smartcar for the provided scopes.
// // It takes a VIN and authClient credentials and will return a bool indicating compatibility.
// func (authClient AuthClient) IsCompatible(vin string) (bool, error) {
// 	type CompatibleResponse struct {
// 		Compatible bool `json:"compatible"`
// 	}

// 	compatiblityURL := url.URL{
// 		Scheme: constants.APIScheme,
// 		Host:   constants.APIHost,
// 		Path:   constants.CompatibilityPath,
// 	}

// 	query := compatiblityURL.Query()
// 	query.Set("vin", vin)
// 	query.Set("scope", strings.Join(authClient.Scope, " "))
// 	compatiblityURL.RawQuery = query.Encode()

// 	authString := authClient.ClientId + ":" + authClient.ClientSecret
// 	encodedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authString))

// 	res, resErr := requests.GET(compatiblityURL.String(), encodedAuth)
// 	if resErr != nil {
// 		return false, resErr
// 	}
// 	defer res.Body.Close()
// 	jsonDecoder := json.NewDecoder(res.Body)

// 	if res.StatusCode != 200 {
// 		var err Error
// 		// TODO:
// 		// Put the next lines in a handler, there is repetition of code.
// 		jsonErr := jsonDecoder.Decode(&err)
// 		if jsonErr != nil {
// 			jsonErr = errors.New("Decoding JSON error")
// 			return false, jsonErr
// 		}
// 		return false, &Error{err.Name, err.Message, err.Code}
// 	}

// 	var compatibleResponse CompatibleResponse
// 	// TODO:
// 	// Repetition of code, use handler.
// 	jsonErr := jsonDecoder.Decode(&compatibleResponse)
// 	if jsonErr != nil {
// 		jsonErr = errors.New("Decoding JSON error")
// 		return false, jsonErr
// 	}

// 	return compatibleResponse.Compatible, nil
// }
