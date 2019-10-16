

# Smartcar Go SDK
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/smartcar/go-sdk)

## Overview

The [Smartcar API](https://smartcar.com/docs) lets you read vehicle data (location, odometer, fuel, etc.) and send commands to vehicles (lock, unlock) using HTTP requests.

## Installation
Install the smartcar package if you are not using Go modules:
```
go get -u github.com/smartcar/go-sdk
```
Import it using:
```
import smartcar "github.com/smartcar/go-sdk"
```

## Getting Started Guide

1. Initialize a Smartcar `Client`.

	```go
	// A smartcar Client is needed to talk to any of the methods in the SDK
	smartcarClient := smartcar.NewClient()
	```

1. Initialize an `Auth` Client struct with your `client id`, `client secret`, `redirect URI`, the `scopes` you want, and `test mode`.

	```go
	// An Auth Client is used to generate a smartcar connect url, authenticate with smartcar, and check compatibility
	authClient := smartcarClient.NewAuth(&smartcar.AuthParams{
		ClientID:     "<CLIENT_ID>",
		ClientSecret: "<CLIENT_SECRET>",
		RedirectURI:  "<REDIRECT_URI>",
		Scope:        []string{"read_vehicle_info"},
		TestMode:     true,
	})
	```

1. Get an auth URL and then redirect user to that URL.

	```go
	authURL, err := authClient.GetAuthURL(&smartcar.AuthURLParams{})

	// redirect user here.

	/*
		If using the net/http library, you can use http.Redirect.
		In order for the next line to work you need to have an endpoint that
		has access to a http.ResponseWriter and http.Request.
	*/
	http.Redirect(w, req, authURL, http.StatusSeeOther)
	```

1. Setup up a redirectURI endpoint to receive authorization code. Exchange auth code for an authorization.

	```go
	// Exchange initial authorization code
	token, err := authClient.ExchangeCode(
		context.TODO(),
		&smartcar.ExchangeCodeParams{Code: code},
	)

	// When your token expires, you can exchange it by sending your refresh token for continued access.
	token, err := authClient.ExchangeRefreshToken(
		context.TODO(),
		&smartcar.ExchangeRefreshTokenParams{Token: token.Refresh},
	)

	// You can check the validity of your token by callint the IsTokenExpired method
	isExpired := smartcarClient.IsTokenExpired(&smartcarClient.TokenExpiredParams{
		Expiry: token.AccessExpiry,
	})
	```

1. In order to send a request to a vehicle, you need to create a smartcar.Vehicle, and for that you need a vehicle ID. You can get a vehicleID by sending a request to smartcar.GetVehicleIDs, which will return a list of vehicleIDs associated with an access token.

	```go
	vehicleIDs, err := smartcar.GetVehicleIDs(
		context.TODO(),
		&smartcar.VehicleIDsParams{Access: token.Access},
	)
	```

1. Construct vehicle with an ID, and an AccessToken, a UnitSystem is optional

	```go
	vehicleParams := smartcarClient.VehicleParams{
		ID: (*vehicleIDs)[0],
		AccessToken: token.Access,
	}
	vehicle := smartcarClient.NewVehicle(&vehicleParams)
	```

1. Send request to vehicle. The following endpoints are available. (check the go documentation for the most up to date list.)

	```go
	// Vehicle Endpoints
	battery, err := vehicle.GetBattery(context.TODO())
	charge, err := vehicle.GetCharge(context.TODO())
	disconnect, err := vehicle.Disconnect(context.TODO())
	fuel, err := vehicle.GetFuel(context.TODO())
	info, err := vehicle.GetInfo(context.TODO())
	location, err := vehicle.GetLocation(context.TODO())
	lock, err := vehicle.Lock(context.TODO())
	odometer, err := vehicle.GetOdometer(context.TODO())
	oil, err := vehicle.GetOil(context.TODO())
	permissions, err := vehicle.GetPermissions(context.TODO())
	tirePressure, err := vehicle.GetTiresPressue(context.TODO())
	unlock, err := vehicle.Unlock(context.TODO())
	vin, err := vehicle.GetVIN(context.TODO())

	// You can change the unit systems of a vehicle at any point by doing.
	err := vehicle.SetUnits(smartcar.UnitsParams{Unit: smartcar.UnitSystemMetric})
	```
