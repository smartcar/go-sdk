package smartcar

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

/*
	All of our vehicle endpoints return structs, which are defined in the following lines.
*/

// Key is a type to used for API endpoints
type Key string

// Helper types to use in vehicle.Batch()
const (
	BatteryPath      Key = "/battery"
	ChargePath       Key = "/charge"
	FuelPath         Key = "/fuel"
	InfoPath         Key = "/"
	LocationPath     Key = "/location"
	OdometerPath     Key = "/odometer"
	OilPath          Key = "/engine/oil"
	PermissionsPath  Key = "/permissions"
	TirePressurePath Key = "/tires/pressure"
	VINPath          Key = "/vin"

	// DO NOT export the paths that are not supported by Batch.
	securityPath    Key = "/security"
	applicationPath Key = "/application"
	batchPath       Key = "/batch"
)

// Battery formats response returned from vehicle.GetBattery().
type Battery struct {
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
	ResponseHeaders
}

// Charge formats response returned from vehicle.GetCharge().
type Charge struct {
	IsPluggedIn bool   `json:"isPluggedIn"`
	State       string `json:"state"`
	ResponseHeaders
}

// Data formats responses returned from vehicle.Batch().
type Data struct {
	Battery      *Battery      `json:"battery,omitempty"`
	Charge       *Charge       `json:"charge,omitempty"`
	Fuel         *Fuel         `json:"fuel,omitempty"`
	Info         *Info         `json:"info,omitempty"`
	Location     *Location     `json:"location,omitempty"`
	Odometer     *Odometer     `json:"odometer,omitempty"`
	Oil          *Oil          `json:"oil,omitempty"`
	Permissions  *Permissions  `json:"permissions,omitempty"`
	TirePressure *TirePressure `json:"tirePressure,omitempty"`
	VIN          *VIN          `json:"vin,omitempty"`
}

// Disconnect formats response returned from vehicle.Disconnect().
type Disconnect struct {
	Status string `json:"status"`
	ResponseHeaders
}

// Fuel formats response returned from vehicle.GetFuel().
type Fuel struct {
	AmountRemaining  float64 `json:"amountRemaining"`
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
	ResponseHeaders
}

// Info formats response returned from vehicle.GetInfo().
type Info struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
	ResponseHeaders
}

// Location formats response returned from vehicle.GetLocation().
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ResponseHeaders
}

// Odometer formats response returned from vehicle.GetOdometer().
type Odometer struct {
	Distance float64 `json:"distance"`
	ResponseHeaders
}

// Oil formats response returned from vehicle.GetOil().
type Oil struct {
	LifeRemaining float64 `json:"lifeRemaining"`
	ResponseHeaders
}

// Permissions formats response returned from vehicle.GetPermissions().
type Permissions struct {
	Permissions []string `json:"permissions"`
	ResponseHeaders
}

// TirePressure formats response returned from vehicle.GetTirePressure().
type TirePressure struct {
	FrontLeft  float64 `json:"frontLeft"`
	FrontRight float64 `json:"frontRight"`
	BackLeft   float64 `json:"backLeft"`
	BackRight  float64 `json:"backRight"`
	ResponseHeaders
}

// VIN formats response returned from vehicle.GetVIN().
type VIN struct {
	VIN string `json:"vin"`
	ResponseHeaders
}

// Security formats response returned from the vehicle.Lock(), vehicle.Unlock().
type Security struct {
	Status string `json:"status"`
	ResponseHeaders
}

// UnitSystem type that will have either imperic or metric.
type UnitSystem string

/// UnitSystem constants that initialize metric or imperial.
const (
	Metric   UnitSystem = "metric"
	Imperial UnitSystem = "imperial"
)

var unitSystems map[string]UnitSystem = map[string]UnitSystem{
	"metric":   Metric,
	"imperial": Imperial,
}

// UnitsParams struct is
// UnitsParam is a param in vehicle.SetUnitSystem
type UnitsParams struct {
	Units UnitSystem
}

// Vehicle is an interface that contains all public methods available for vehicle. vehicle needs to implement
// this methods to be able to expose them.
type Vehicle interface {
	Batch(context.Context, ...Key) (*Data, error)
	Disconnect(context.Context) (*Disconnect, error)
	GetBattery(context.Context) (*Battery, error)
	GetCharge(context.Context) (*Charge, error)
	GetFuel(context.Context) (*Fuel, error)
	GetInfo(context.Context) (*Info, error)
	GetLocation(context.Context) (*Location, error)
	GetOdometer(context.Context) (*Odometer, error)
	GetOil(context.Context) (*Oil, error)
	GetPermissions(context.Context) (*Permissions, error)
	GetTiresPressure(context.Context) (*TirePressure, error)
	GetVIN(context.Context) (*VIN, error)
	Lock(context.Context) (*Security, error)
	SetUnitSystem(*UnitsParams) error
	Unlock(context.Context) (*Security, error)
}

// vehicle client that implements the Vehicle interface.
type vehicle struct {
	requestParams
	id          string
	accessToken string
	client      backendClient
}

type batchResponse struct {
	Responses []struct {
		Path    string
		Code    int
		Headers struct {
			DataAge    string     `json:"sc-data-age,omitempty"`
			RequestID  string     `json:"sc-request-id,omitempty"`
			UnitSystem UnitSystem `json:"sc-unit-system,omitempty"`
		} `json:"headers,omitempty"`
		Body interface{} `json:"body"`
	} `json:"responses"`
}

// Batch sends a request to Smartcar's API vehicle/batch endpoint.
func (v *vehicle) Batch(ctx context.Context, keys ...Key) (*Data, error) {
	var requests []map[string]string

	for _, path := range keys {
		requests = append(requests, map[string]string{"path": string(path)})
	}
	body := map[string][]map[string]string{
		"requests": requests,
	}
	marshalBody, _ := json.Marshal(body)
	bufferedBody := bytes.NewBuffer([]byte(marshalBody))

	target := new(batchResponse)
	err := v.request(ctx, string(batchPath), http.MethodPost, v.requestParams, bufferedBody, target)
	if err != nil {
		return nil, err
	}

	data := new(Data)
	for _, v := range target.Responses {
		switch v.Path {
		case string(BatteryPath):
			mapstructure.Decode(v.Body, &data.Battery)
			mapstructure.Decode(v.Headers, &data.Battery.ResponseHeaders)
		case string(ChargePath):
			mapstructure.Decode(v.Body, &data.Charge)
			mapstructure.Decode(v.Headers, &data.Charge.ResponseHeaders)
		case string(FuelPath):
			mapstructure.Decode(v.Body, &data.Fuel)
			mapstructure.Decode(v.Headers, &data.Fuel.ResponseHeaders)
		case string(InfoPath):
			mapstructure.Decode(v.Body, &data.Info)
			mapstructure.Decode(v.Headers, &data.Info.ResponseHeaders)
		case string(LocationPath):
			mapstructure.Decode(v.Body, &data.Location)
			mapstructure.Decode(v.Headers, &data.Location.ResponseHeaders)
		case string(OdometerPath):
			mapstructure.Decode(v.Body, &data.Odometer)
			mapstructure.Decode(v.Headers, &data.Odometer.ResponseHeaders)
		case string(OilPath):
			mapstructure.Decode(v.Body, &data.Oil)
			mapstructure.Decode(v.Headers, &data.Oil.ResponseHeaders)
		case string(PermissionsPath):
			mapstructure.Decode(v.Body, &data.Permissions)
			mapstructure.Decode(v.Headers, &data.Permissions.ResponseHeaders)
		case string(TirePressurePath):
			mapstructure.Decode(v.Body, &data.TirePressure)
			mapstructure.Decode(v.Headers, &data.TirePressure.ResponseHeaders)
		case string(VINPath):
			mapstructure.Decode(v.Body, &data.VIN)
			mapstructure.Decode(v.Headers, &data.VIN.ResponseHeaders)
		}
	}

	return data, nil
}

// Disconnect sends a request to Smartcar's API vehicle/application endpoint.
func (v *vehicle) Disconnect(ctx context.Context) (*Disconnect, error) {
	disconnect := &Disconnect{}
	return disconnect, v.request(ctx, string(applicationPath), http.MethodDelete, v.requestParams, nil, disconnect)
}

// GetBattery sends a request to Smartcar's API vehicle/battery endpoint.
func (v *vehicle) GetBattery(ctx context.Context) (*Battery, error) {
	battery := &Battery{}
	return battery, v.request(ctx, string(BatteryPath), http.MethodGet, v.requestParams, nil, battery)
}

// GetCharge sends a request to Smartcar's API vehicle/charge endpoint.
func (v *vehicle) GetCharge(ctx context.Context) (*Charge, error) {
	charge := &Charge{}
	return charge, v.request(ctx, string(ChargePath), http.MethodGet, v.requestParams, nil, charge)
}

// GetFuel sends a request to Smartcar's API vehicle/fuel endpoint.
func (v *vehicle) GetFuel(ctx context.Context) (*Fuel, error) {
	fuel := &Fuel{}
	return fuel, v.request(ctx, string(FuelPath), http.MethodGet, v.requestParams, nil, fuel)
}

// GetInfo sends a request to Smartcar's API vehicle/ endpoint.
func (v *vehicle) GetInfo(ctx context.Context) (*Info, error) {
	info := &Info{}
	return info, v.request(ctx, string(InfoPath), http.MethodGet, v.requestParams, nil, info)
}

// GetLocation sends a request to Smartcar's API vehicle/location endpoint.
func (v *vehicle) GetLocation(ctx context.Context) (*Location, error) {
	location := &Location{}
	return location, v.request(ctx, string(LocationPath), http.MethodGet, v.requestParams, nil, location)
}

// GetOdometer sends a request to Smartcar's API vehicle/odometer endpoint.
func (v *vehicle) GetOdometer(ctx context.Context) (*Odometer, error) {
	odometer := &Odometer{}
	return odometer, v.request(ctx, string(OdometerPath), http.MethodGet, v.requestParams, nil, odometer)
}

// GetOil sends a request to Smartcar's API vehicle/oil endpoint.
func (v *vehicle) GetOil(ctx context.Context) (*Oil, error) {
	oil := &Oil{}
	return oil, v.request(ctx, string(OilPath), http.MethodGet, v.requestParams, nil, oil)
}

// GetPermissions sends a request to Smartcar's API vehicle/permissions endpoint.
func (v *vehicle) GetPermissions(ctx context.Context) (*Permissions, error) {
	permissions := &Permissions{}
	return permissions, v.request(ctx, string(PermissionsPath), http.MethodGet, v.requestParams, nil, permissions)
}

// GetTiresPressure sends a request to Smartcar's API vehicle/tires/pressure endpoint.
func (v *vehicle) GetTiresPressure(ctx context.Context) (*TirePressure, error) {
	tirePressure := &TirePressure{}
	return tirePressure, v.request(ctx, string(TirePressurePath), http.MethodGet, v.requestParams, nil, tirePressure)
}

// GetVIN sends a request to Smartcar's API vehicle/vin endpoint.
func (v *vehicle) GetVIN(ctx context.Context) (*VIN, error) {
	vin := &VIN{}
	return vin, v.request(ctx, string(VINPath), http.MethodGet, v.requestParams, nil, vin)
}

// Lock sends a request to Smartcar's API vehicle/lock endpoint.
func (v *vehicle) Lock(ctx context.Context) (*Security, error) {
	body := bytes.NewBuffer([]byte(`{"action":"LOCK"}`))
	lock := &Security{}
	return lock, v.request(ctx, string(securityPath), http.MethodPost, v.requestParams, body, lock)
}

/*
  SetUnits sets the unit system for a vehicle's instance. (i.e. Setting the unit system to metric, will
		return the odometer in meters).
  Note: Does not send a request to Smartcar's API, it just changes the unitSystem of the vehicle instance.
		Therefore sending a new request after calling this method, the response will return the data using the unitSystem set.
*/
func (v *vehicle) SetUnitSystem(params *UnitsParams) error {
	if !(params.Units == Imperial || params.Units == Metric) { //compare to actuall unit system || params.Units == "imperial") { // check here if its part of the map. go enum
		err := fmt.Sprintf("Unit must be %s or %s", Metric, Imperial)
		return errors.New(err)
	}
	v.requestParams.UnitSystem = params.Units
	return nil
}

// Unlock sends a request to Smartcar's API vehicle/unlock endpoint.
func (v *vehicle) Unlock(ctx context.Context) (*Security, error) {
	body := bytes.NewBuffer([]byte(`{"action":"UNLOCK"}`))
	unlock := &Security{}
	return unlock, v.request(ctx, string(securityPath), http.MethodPost, v.requestParams, body, unlock)
}

/*
  request is an internal function used to make requests to Smartcar's vehicle API. It accepts an interface,
  which is used to format the response.
*/
func (v *vehicle) request(ctx context.Context, path, method string, params requestParams, data io.Reader, target interface{}) error {
	return v.client.Call(backendClientParams{
		ctx:           ctx,
		method:        method,
		url:           buildVehicleURL(path, v.id),
		authorization: buildBearerAuthorization(v.accessToken),
		requestParams: params,
		body:          data,
		target:        target,
	})
}
