package smartcar

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/smartcar/go-sdk/helpers/requests"
	"github.com/smartcar/go-sdk/helpers/utils"
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

// Info sends request to Smartcar API vehicle/ endpoint and returns formatted response.
func (vehicle *Vehicle) Info() (VehicleInfoResponse, error) {
	res, err := vehicle.request("/", requests.GET, nil)
	if err != nil {
		return VehicleInfoResponse{}, err
	}

	formattedResponse := new(VehicleInfoResponse)
	defer res.Body.Close()
	fmtErr := requests.FormatResponse(res.Body, formattedResponse)
	if fmtErr != nil {
		return VehicleInfoResponse{}, fmtErr
	}

	return *formattedResponse, nil
}

// VIN sends request to Smartcar API vehicle/vin endpoint and returns formatted response.
func (vehicle *Vehicle) VIN() (string, error) {
	res, resErr := vehicle.request("/vin", requests.GET, nil)
	if resErr != nil {
		return "", resErr
	}

	formattedResponse := new(struct{ Vin string })
	defer res.Body.Close()
	fmtErr := requests.FormatResponse(res.Body, formattedResponse)
	if fmtErr != nil {
		return "", fmtErr
	}

	return formattedResponse.Vin, nil
}

// Lock sends request to Smartcar API vehicle/odometer endpoint and returns formatted response.
func (vehicle *Vehicle) Lock() (VehicleResponse, error) {
	return vehicle.security("LOCK")
}

// Unlock sends request to Smartcar API vehicle/odometer endpoint and returns formatted response.
func (vehicle *Vehicle) Unlock() (VehicleResponse, error) {
	return vehicle.security("UNLOCK")
}

// Internal function to send lock/unlock request
func (vehicle *Vehicle) security(action string) (VehicleResponse, error) {
	var jsonStr []byte
	if action == "LOCK" {
		jsonStr = []byte(`{"action":"LOCK"}`)
	} else {
		jsonStr = []byte(`{"action":"UNLOCK"}`)
	}

	body := bytes.NewBuffer(jsonStr)

	res, err := vehicle.request("/security", requests.POST, body)
	if err != nil {
		return VehicleResponse{}, err
	}

	formattedResponse := new(VehicleResponse)
	defer res.Body.Close()
	fmtErr := requests.FormatResponse(res.Body, formattedResponse)
	if fmtErr != nil {
		return VehicleResponse{}, fmtErr
	}

	return *formattedResponse, nil
}

// Permissions sends request to Smartcar API vehicle/permissions endpoint and returns formatted response.
func (vehicle *Vehicle) Permissions() ([]string, error) {
	res, err := vehicle.request("/permissions", requests.GET, nil)
	if err != nil {
		return nil, err
	}

	formattedResponse := new(struct{ Permissions []string })
	defer res.Body.Close()
	fmtErr := requests.FormatResponse(res.Body, formattedResponse)
	if fmtErr != nil {
		return nil, fmtErr
	}

	return formattedResponse.Permissions, nil
}

// HasPermissions checks if vehicle has the permissions passed in.
func (vehicle *Vehicle) HasPermissions(permissions ...string) (bool, error) {
	vehiclePermissions, err := vehicle.Permissions()
	if err != nil {
		return false, err
	}

	set := make(map[string]bool)
	for _, value := range vehiclePermissions {
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

// request is an internal functions used to make requests to Smartcar's vehicle API.
func (vehicle *Vehicle) request(path string, method string, data io.Reader) (http.Response, error) {
	// Build url
	url := utils.BuildVehicleURL(path, vehicle.ID)

	// Build authorization
	authorization := utils.BuildBearerAuthorization(vehicle.AccessToken)

	// Send request
	return requests.Request(method, url, authorization, vehicle.UnitSystem, data)
}
