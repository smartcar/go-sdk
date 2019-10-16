package smartcar

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

/*
	All of our vehicle endpoints return structs, which are defined in the following lines.
*/

// Disconnect formats response returned from vehicle.Disconnect().
type Disconnect struct {
	Status string `json:"status"`
	smartcarHeaders
}

// Battery formats response returned from vehicle.GetBattery().
type Battery struct {
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
	smartcarHeaders
}

// Charge formats response returned from vehicle.GetCharge().
type Charge struct {
	IsPluggedIn bool   `json:"isPluggedIn"`
	State       string `json:"state"`
	smartcarHeaders
}

// Fuel formats response returned from vehicle.GetFuel().
type Fuel struct {
	AmountRemaining  float64 `json:"amountRemaining"`
	PercentRemaining float64 `json:"percentRemaining"`
	Range            float64 `json:"range"`
	smartcarHeaders
}

// Info formats response returned from vehicle.GetInfo().
type Info struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
	smartcarHeaders
}

// Location formats response returned from vehicle.GetLocation().
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	smartcarHeaders
}

// Odometer formats response returned from vehicle.GetOdometer().
type Odometer struct {
	Distance float64 `json:"distance"`
	smartcarHeaders
}

// Oil formats response returned from vehicle.GetOil().
type Oil struct {
	LifeRemaining float64 `json:"lifeRemaining"`
	smartcarHeaders
}

// Permissions formats response returned from vehicle.GetPermissions().
type Permissions struct {
	Permissions []string `json:"permissions"`
	smartcarHeaders
}

// TirePressure formats response returned from vehicle.GetTirePressure().
type TirePressure struct {
	FrontLeft  float64 `json:"frontLeft"`
	FrontRight float64 `json:"frontRight"`
	BackLeft   float64 `json:"backLeft"`
	BackRight  float64 `json:"backRight"`
	smartcarHeaders
}

// VIN formats response returned from vehicle.GetVIN().
type VIN struct {
	VIN string `json:"vin"`
	smartcarHeaders
}

// Security formats response returned from the vehicle.Lock(), vehicle.Unlock().
type Security struct {
	Status string `json:"status"`
	smartcarHeaders
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
	client      backendClient // nit name
}

// Disconnect sends a request to Smartcar's API vehicle/application endpoint.
func (v *vehicle) Disconnect(ctx context.Context) (*Disconnect, error) {
	disconnect := &Disconnect{}
	return disconnect, v.request(ctx, "/application", http.MethodDelete, v.requestParams, nil, disconnect)
}

// GetBattery sends a request to Smartcar's API vehicle/battery endpoint.
func (v *vehicle) GetBattery(ctx context.Context) (*Battery, error) {
	battery := &Battery{}
	return battery, v.request(ctx, "/battery", http.MethodGet, v.requestParams, nil, battery)
}

// GetCharge sends a request to Smartcar's API vehicle/charge endpoint.
func (v *vehicle) GetCharge(ctx context.Context) (*Charge, error) {
	charge := &Charge{}
	return charge, v.request(ctx, "/charge", http.MethodGet, v.requestParams, nil, charge)
}

// GetFuel sends a request to Smartcar's API vehicle/fuel endpoint.
func (v *vehicle) GetFuel(ctx context.Context) (*Fuel, error) {
	fuel := &Fuel{}
	return fuel, v.request(ctx, "/fuel", http.MethodGet, v.requestParams, nil, fuel)
}

// GetInfo sends a request to Smartcar's API vehicle/ endpoint.
func (v *vehicle) GetInfo(ctx context.Context) (*Info, error) {
	info := &Info{}
	return info, v.request(ctx, "/", http.MethodGet, v.requestParams, nil, info)
}

// GetLocation sends a request to Smartcar's API vehicle/location endpoint.
func (v *vehicle) GetLocation(ctx context.Context) (*Location, error) {
	location := &Location{}
	return location, v.request(ctx, "/location", http.MethodGet, v.requestParams, nil, location)
}

// GetOdometer sends a request to Smartcar's API vehicle/odometer endpoint.
func (v *vehicle) GetOdometer(ctx context.Context) (*Odometer, error) {
	odometer := &Odometer{}
	return odometer, v.request(ctx, "/odometer", http.MethodGet, v.requestParams, nil, odometer)
}

// GetOil sends a request to Smartcar's API vehicle/oil endpoint.
func (v *vehicle) GetOil(ctx context.Context) (*Oil, error) {
	oil := &Oil{}
	return oil, v.request(ctx, "/engine/oil", http.MethodGet, v.requestParams, nil, oil)
}

// GetPermissions sends a request to Smartcar's API vehicle/permissions endpoint.
func (v *vehicle) GetPermissions(ctx context.Context) (*Permissions, error) {
	permissions := &Permissions{}
	return permissions, v.request(ctx, "/permissions", http.MethodGet, v.requestParams, nil, permissions)
}

// GetTiresPressure sends a request to Smartcar's API vehicle/tires/pressure endpoint.
func (v *vehicle) GetTiresPressure(ctx context.Context) (*TirePressure, error) {
	tirePressure := &TirePressure{}
	return tirePressure, v.request(ctx, "/tires/pressure", http.MethodGet, v.requestParams, nil, tirePressure)
}

// GetVIN sends a request to Smartcar's API vehicle/vin endpoint.
func (v *vehicle) GetVIN(ctx context.Context) (*VIN, error) {
	vin := &VIN{}
	return vin, v.request(ctx, "/vin", http.MethodGet, v.requestParams, nil, vin)
}

// Lock sends a request to Smartcar's API vehicle/lock endpoint.
func (v *vehicle) Lock(ctx context.Context) (*Security, error) {
	body := bytes.NewBuffer([]byte(`{"action":"LOCK"}`))
	lock := &Security{}
	return lock, v.request(ctx, "/security", http.MethodPost, v.requestParams, body, lock)
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
	return unlock, v.request(ctx, "/security", http.MethodPost, v.requestParams, body, unlock)
}

/*
  request is an internal function used to make requests to Smartcar's vehicle API. It accepts an interface,
  which is used to format the response.
*/
func (v *vehicle) request(ctx context.Context, path, method string, params requestParams, data io.Reader, target interface{}) error {
	return v.client.Call(backendClientParams{
		ctx:           ctx,
		method:        method,
		url:           BuildVehicleURL(path, v.id),
		authorization: BuildBearerAuthorization(v.accessToken),
		requestParams: params,
		body:          data,
		target:        target,
	})
}
