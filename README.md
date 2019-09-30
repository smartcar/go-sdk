# Smartcar Go SDK
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/smartcar/go-sdk)

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
	Auth:     auth,
	ForceApproval: false,
	State:  "",
	MakeBypass:   smartcar.MakeBypass{},
	SingleSelect:     smartcar.SingleSelect{},
})
```
3. Setup up redirectURI endpoint to receive authorization code. Exchange auth code for an authorization.
```go
token, err := authClient.ExchangeCode(auth)
```
4. Get vehicle ids from access token
```go
vehicleIds, err := smartcar.GetVehicleIds(token.Access)
```
5. Construct vehicle
```go
vehicle := smartcar.Vehicle{ID: vehicleIds[0], AccessToken: token.Access}
```
6. Send request to vehicle
```go
vehicleInfo, err := vehicle.Info()
```
