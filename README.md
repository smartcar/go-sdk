

# Smartcar Go SDK
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/smartcar/go-sdk)

## WARNING: THIS SDK IS IN BETA, INTERFACE IS SUBJECT TO CHANGE

## Overview

The [Smartcar API](https://smartcar.com/docs) lets you read vehicle data (location, odometer, fuel, etc.) and send commands to vehicles (lock, unlock) to connected vehicles using HTTP requests.

## Installation
Install the smartcar package if you are not using Go modules:
```
go get -u github.com/smartcar/go-sdk
```
Import it using:
```
import smartcar "github.com/smartcar/go-sdk"
```

## Getting Started
1. First initialize an `AuthClient` struct with your client credentials, redirect URI, the scopes required and test mode.
```go
// AuthClient is used to store your auth credentials when authenticating with Smartcar.
authClient := smartcar.AuthClient{
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
	RedirectURI:  "REDIRECT_URI",
	Scope:        []string{"read_vehicle_info"},
	TestMode:     false,
}
```
2. Get a connect URL and redirect user to that URL. 
```go
authURL, err := authClient.GetAuthURL(smartcar.AuthURLOptions{
	ForceApproval: false,
	State:  "",
	MakeBypass:   smartcar.MakeBypass{},
	SingleSelect:     smartcar.SingleSelect{},
})
```
3. Setup up redirectURI endpoint to receive authorization code. Exchange auth code for an authorization.
```go
// Exchange initial authorization code
token, err := authClient.ExchangeCode(code)

// Refresh for continued access
token, err := authClient.ExchangeRefreshToken(token.Refresh)
```
4. Get vehicle ids using access token
```go
vehicleIds, err := smartcar.GetVehicleIds(token.Access)
```
5. Construct vehicle
```go
vehicle := smartcar.Vehicle{ID: vehicleIds[0], AccessToken: token.Access}
```
6. Send request to vehicle
```go
// Vehicle Endpoints
vehicleResponse, err := vehicle.Info()
vehicleResponse, err := vehicle.VIN()
vehicleResponse, err := vehicle.Odometer()
vehicleResponse, err := vehicle.Lock()
vehicleResponse, err := vehicle.Unlock()
vehicleResponse, err := vehicle.Location()
vehicleResponse, err := vehicle.Fuel()
vehicleResponse, err := vehicle.Battery()
vehicleResponse, err := vehicle.Charge()
vehicleResponse, err := vehicle.Permissions()
vehicleResponse, err := vehicle.HasPermissions()
vehicleResponse, err := vehicle.Disconnect()
```
